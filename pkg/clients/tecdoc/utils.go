package tecdoc

import (
	"tec-doc/pkg/clients/model"
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
