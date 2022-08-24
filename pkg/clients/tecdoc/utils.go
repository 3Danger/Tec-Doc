package tecdoc

import "tec-doc/internal/tec-doc/model"

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
