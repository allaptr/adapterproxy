package model

type CompanyInfo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Active_until string `json:"active_until,omitempty"`
}

type CompanyLegacyV1 struct {
	Cn         string `json:"cn"`
	Created_on string `json:"created_on"`
	Closed_on  string `json:"closed_on,omitempty"`
}

type CompanyLegacyV2 struct {
	Company_name string `json:"company_name"`
	Tin          string `json:"tin"`
	Dissolved_on string `json:"dissolved_on,omitempty"`
}
