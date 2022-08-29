package tecdoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"tec-doc/pkg/clients/model"
)

// LIMIT Лимит на единицу запроса
const LIMIT = 100

// Applicability получение списка применимости
func (c *tecDocClient) Applicability(legacyArticleId int) ([]model.LinkageTargets, error) {
	linkageTargetsBodies, err := c.getLinkageTargetsResponse(legacyArticleId)
	if err != nil {
		return nil, err
	}

	length := len(linkageTargetsBodies)
	targetCh := make(chan []model.LinkageTargets, length)
	errChan := make(chan error, length+1)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer close(targetCh)
		defer close(errChan)
		for i := range linkageTargetsBodies {
			wg.Add(1)
			go func(wg *sync.WaitGroup, LinkageTargets model.GetLinkageTargetsResponse) {
				defer wg.Done()
				targets, err := c.getLinkageTargets(LinkageTargets)
				if err != nil {
					errChan <- err
					return
				}
				targetCh <- targets
			}(wg, linkageTargetsBodies[i])
		}
		wg.Done()
		wg.Wait()
		errChan <- nil
	}(wg)
	linkageTargets := make([]model.LinkageTargets, 0, length*LIMIT)
	for target := range targetCh {
		linkageTargets = append(linkageTargets, target...)
	}
	if err = <-errChan; err != nil {
		return nil, err
	}
	return linkageTargets, nil
}

func (c *tecDocClient) getLinkageTargetsResponse(legacyArticleId int) ([]model.GetLinkageTargetsResponse, error) {
	var (
		LinkageTargetsBody []model.GetLinkageTargetsResponse
		linkages           []model.ArticleLinkages
		err                error
	)
	if linkages, err = c.getArticleLinkedAllLinkingTarget4(legacyArticleId); err != nil {
		return nil, err
	}
	var linkageTargetIds []map[string]any
	for _, array := range linkages {
		for _, data := range array.ArticleLinkages.LinkingTargetId {
			linkageTargetIds = append(linkageTargetIds, map[string]any{
				"type": "P",
				"id":   data.LinkingTargetId,
			})
		}
	}
	steps := len(linkageTargetIds) / LIMIT
	for i := 0; i < steps; i++ {
		start, end := i*LIMIT, (i+1)*LIMIT
		LinkageTargetsBody = append(LinkageTargetsBody, model.GetLinkageTargetsResponse{
			GetLinkageTargets: model.GetLinkageTargets{
				PerPage: LIMIT, Page: 1,
				LinkageTargetCountry: "RU", Lang: "ru",
				LinkageTargetIds: linkageTargetIds[start:end]},
		})
	}
	LinkageTargetsBody = append(LinkageTargetsBody, model.GetLinkageTargetsResponse{
		GetLinkageTargets: model.GetLinkageTargets{
			PerPage: LIMIT, Page: 1,
			LinkageTargetCountry: "RU", Lang: "ru",
			LinkageTargetIds: linkageTargetIds[steps*LIMIT:]},
	})
	return LinkageTargetsBody, nil
}

func (c *tecDocClient) getArticleLinkedAllLinkingTarget4(legacyArticleId int) ([]model.ArticleLinkages, error) {
	var responseFirst model.Data
	{
		requestBody := bytes.NewReader([]byte(fmt.Sprintf(
			`{
						"getArticleLinkedAllLinkingTarget4": {
							"articleCountry": "RU",
							"articleId": %d,
							"country": "RU",
							"lang": "ru",
							"linkingTargetType": "P"
						}
					}`, legacyArticleId)))

		if err := c.doRequest(http.MethodPost, requestBody, &responseFirst); err != nil {
			return nil, err
		}
		if responseFirst.Status != 200 {
			return nil, fmt.Errorf("bad status code %d", responseFirst.Status)
		}
	}
	return responseFirst.Data.Array, nil
}

// getLinkageTargets
// максимальный размер массива внутри структуры LinkageTargetIds = 100 (указал в константе LIMIT)
func (c *tecDocClient) getLinkageTargets(LinkageTargetsBody model.GetLinkageTargetsResponse) (linkageTargets []model.LinkageTargets, err error) {
	var length = len(LinkageTargetsBody.GetLinkageTargets.LinkageTargetIds)
	linkageTargets = make([]model.LinkageTargets, 0, length)
	LinkageTargetsBody.GetLinkageTargets.Page = 1
	var linkageTargetsResponse = model.LinkageTargetsResponse{Total: length}
	for arrivedCount := 0; arrivedCount < linkageTargetsResponse.Total; {
		var requestByte []byte
		if requestByte, err = json.Marshal(LinkageTargetsBody); err != nil {
			return nil, err
		}
		if err = c.doRequest(http.MethodPost, bytes.NewReader(requestByte), &linkageTargetsResponse); err != nil {
			return nil, err
		}
		if linkageTargetsResponse.Status != http.StatusOK {
			return nil, fmt.Errorf("bad status code %d", linkageTargetsResponse.Status)
		}
		arrivedCount += len(linkageTargetsResponse.LinkageTargets)
		linkageTargets = append(linkageTargets, linkageTargetsResponse.LinkageTargets...)
		LinkageTargetsBody.GetLinkageTargets.Page++
	}
	return linkageTargets, nil
}
