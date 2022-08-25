package tecdoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tec-doc/pkg/clients/model"
)

// Applicability получение списка применимости
func (c *tecDocClient) Applicability(legacyArticleId int) (linkageTargets []model.LinkageTargets, err error) {
	type (
		// DataFirstResponse для записи ответа первого запроса
		DataFirstResponse struct {
			Data struct {
				Array []struct {
					ArticleLinkages struct {
						LinkingTargetId []struct {
							LinkingTargetId int `json:"linkingTargetId"`
						} `json:"array"`
					} `json:"articleLinkages"`
				} `json:"array"`
			} `json:"data"`
			Status int `json:"status"`
		}
		// LinkageTargetsResponse для записи результата
		LinkageTargetsResponse struct {
			Total          int                    `json:"total"`
			LinkageTargets []model.LinkageTargets `json:"linkageTargets"`
			Status         int                    `json:"status"`
		}

		// GetLinkageTargetsResponse для запроса
		GetLinkageTargets struct {
			PerPage              int              `json:"perPage"`
			Page                 int              `json:"page"`
			LinkageTargetCountry string           `json:"linkageTargetCountry"`
			Lang                 string           `json:"lang"`
			LinkageTargetIds     []map[string]any `json:"linkageTargetIds"`
		}
		GetLinkageTargetsResponse struct {
			GetLinkageTargets GetLinkageTargets `json:"getLinkageTargets"`
		}
	)

	const LIMIT = 100
	var responseFirst DataFirstResponse
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

		if err = c.doRequest(http.MethodPost, requestBody, &responseFirst); err != nil {
			return nil, err
		}
		if responseFirst.Status != 200 {
			return nil, fmt.Errorf("bad status code %d", responseFirst.Status)
		}
	}

	var length int
	var LinkageTargetsBody []GetLinkageTargetsResponse
	{
		var linkageTargetIds []map[string]any
		for _, array := range responseFirst.Data.Array {
			for _, data := range array.ArticleLinkages.LinkingTargetId {
				linkageTargetIds = append(linkageTargetIds, map[string]any{"type": "P", "id": data.LinkingTargetId})
			}
		}
		length = len(linkageTargetIds)
		steps := length / LIMIT
		for i := 0; i < steps; i++ {
			start, end := i*LIMIT, (i+1)*LIMIT
			LinkageTargetsBody = append(LinkageTargetsBody, GetLinkageTargetsResponse{GetLinkageTargets{
				LIMIT, 1,
				"RU", "ru",
				linkageTargetIds[start:end]},
			})
		}
		LinkageTargetsBody = append(LinkageTargetsBody, GetLinkageTargetsResponse{GetLinkageTargets{
			LIMIT, 1,
			"RU", "ru",
			linkageTargetIds[steps*LIMIT:]},
		})
	}

	linkageTargets = make([]model.LinkageTargets, 0, length)
	for _, targets := range LinkageTargetsBody {
		var requestByte []byte
		targets.GetLinkageTargets.Page = 1

		var linkageTargetsResponse = LinkageTargetsResponse{Total: LIMIT << 1}
		for arrivedCount := 0; arrivedCount < linkageTargetsResponse.Total; {
			if requestByte, err = json.Marshal(targets); err != nil {
				return nil, err
			}
			if err = c.doRequest(http.MethodPost, bytes.NewReader(requestByte), &linkageTargetsResponse); err != nil {
				return nil, err
			}
			if linkageTargetsResponse.Status != http.StatusOK {
				return nil, fmt.Errorf("bad status code %d", responseFirst.Status)
			}
			linkageTargets = append(linkageTargets, linkageTargetsResponse.LinkageTargets...)
			arrivedCount += len(linkageTargetsResponse.LinkageTargets)
			targets.GetLinkageTargets.Page++
		}
	}
	return linkageTargets, nil
}

// doRequest делает запрос и заполняет данными JSON структуру outStructPtr. аналог BindJSON() из gin
func (c *tecDocClient) doRequest(method string, body io.Reader, outStructPtr interface{}) (err error) {
	var (
		response *http.Response
		request  *http.Request
	)
	if request, err = http.NewRequest(method, c.tecDocCfg.URL, body); err != nil {
		return fmt.Errorf("can't create new request: %w", err)
	}
	request.Header = http.Header{"Content-Type": {"application/json"}, "X-Api-Key": {c.tecDocCfg.XApiKey}}
	if response, err = c.Do(request); err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code %d", response.StatusCode)
	}
	if err = json.NewDecoder(response.Body).Decode(&outStructPtr); err != nil {
		return err
	}
	return nil
}
