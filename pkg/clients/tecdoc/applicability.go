package tecdoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DataFirstResponse для записи ответа первого запроса
type DataFirstResponse struct {
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

// GetLinkageTargetsBody для записи ответа втрого запроса
type GetLinkageTargetsBody struct {
	GetLinkageTargets GetLinkageTargets `json:"getLinkageTargets"`
}

type GetLinkageTargets struct {
	LinkageTargetCountry string             `json:"linkageTargetCountry"`
	Lang                 string             `json:"lang"`
	LinkageTargetIds     []LinkageTargetIds `json:"linkageTargetIds"`
}

type LinkageTargetIds struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
}

// LinkageTargetsResponse для записи результата
type LinkageTargetsResponse struct {
	LinkageTargets []LinkageTargets `json:"linkageTargets"`
	Status         int              `json:"status"`
}

type LinkageTargets struct {
	LinkageTargetId        int    `json:"linkageTargetId"`
	MfrName                string `json:"mfrName"`
	VehicleModelSeriesName string `json:"vehicleModelSeriesName"`
	BeginYearMonth         string `json:"beginYearMonth"`
	EndYearMonth           string `json:"endYearMonth"`
}

// Applicability получение списка применимости
func (c *tecDocClient) Applicability(legacyArticleId int) (linkageTargets map[int]LinkageTargets, err error) {
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

	var LinkageTargetsBody []GetLinkageTargetsBody
	{
		var linkageTargetIds []LinkageTargetIds
		for _, array := range responseFirst.Data.Array {
			for _, data := range array.ArticleLinkages.LinkingTargetId {
				linkageTargetIds = append(linkageTargetIds, LinkageTargetIds{"P", data.LinkingTargetId})
			}
		}
		steps := len(linkageTargetIds) / LIMIT
		for i := 0; i < steps; i++ {
			start, end := i*LIMIT, (i+1)*LIMIT
			LinkageTargetsBody = append(LinkageTargetsBody, GetLinkageTargetsBody{GetLinkageTargets{
				"RU", "ru",
				linkageTargetIds[start:end]},
			})
		}
		LinkageTargetsBody = append(LinkageTargetsBody, GetLinkageTargetsBody{GetLinkageTargets{
			"RU", "ru",
			linkageTargetIds[steps*LIMIT:]},
		})
	}

	linkageTargets = make(map[int]LinkageTargets, len(LinkageTargetsBody)*LIMIT)
	for _, targets := range LinkageTargetsBody {
		var requestByte []byte
		if requestByte, err = json.Marshal(targets); err != nil {
			return nil, err
		}

		var linkageTargetsResponse LinkageTargetsResponse
		if err = c.doRequest(http.MethodPost, bytes.NewReader(requestByte), &linkageTargetsResponse); err != nil {
			return nil, err
		}
		if linkageTargetsResponse.Status != http.StatusOK {
			return nil, fmt.Errorf("bad status code %d", responseFirst.Status)
		}
		for _, target := range linkageTargetsResponse.LinkageTargets {
			linkageTargets[target.LinkageTargetId] = target
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
