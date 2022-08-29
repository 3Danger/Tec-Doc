package tecdoc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tec-doc/pkg/model"
)

func convertArticleCriteriaRaw(cr model.ArticleCriteriaRaw) model.ArticleCriteria {
	return model.ArticleCriteria{
		CriteriaDescription:     cr.CriteriaDescription,
		CriteriaAbbrDescription: cr.CriteriaAbbrDescription,
		CriteriaUnitDescription: cr.CriteriaUnitDescription,
		CriteriaType:            cr.CriteriaType,
		RawValue:                cr.RawValue,
		FormattedValue:          cr.FormattedValue,
	}
}

func contains(arr []string, s string) bool {
	for _, elem := range arr {
		if elem == s {
			return true
		}
	}
	return false
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
