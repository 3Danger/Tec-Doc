// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package content

type requestABACCheckAccess struct {
	Scope      string                   `json:"scope"`
	FeatureKey string                   `json:"featureKey"`
	UserID     *uint64                  `json:"userID"`
	Key        [16]byte                 `json:"key"`
	Attributes []map[string]interface{} `json:"attributes"` // This field was defined with ellipsis (...).
}

type responseABACCheckAccess struct {
	Decision bool `json:"decision"`
}
