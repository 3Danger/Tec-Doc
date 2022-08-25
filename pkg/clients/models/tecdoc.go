package models

type LinkageTargets struct {
	LinkageTargetId        int    `json:"linkageTargetId"`
	MfrName                string `json:"mfrName"`
	VehicleModelSeriesName string `json:"vehicleModelSeriesName"`
	BeginYearMonth         string `json:"beginYearMonth"`
	EndYearMonth           string `json:"endYearMonth"`
}

type (
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

// Data для записи ответа первого запроса
type Data struct {
	Data struct {
		Array []ArticleLinkages `json:"array"`
	} `json:"data"`
	Status int `json:"status"`
}

type ArticleLinkages struct {
	ArticleLinkages struct {
		LinkingTargetId []struct {
			LinkingTargetId int `json:"linkingTargetId"`
		} `json:"array"`
	} `json:"articleLinkages"`
}

// LinkageTargetsResponse для записи результата
type LinkageTargetsResponse struct {
	Total          int              `json:"total"`
	LinkageTargets []LinkageTargets `json:"linkageTargets"`
	Status         int              `json:"status"`
}
