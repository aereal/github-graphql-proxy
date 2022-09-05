package dto

type Organization struct {
	Login   string `json:"login"`
	Billing *OrganizationBilling
}

type OrganizationBilling struct {
	OrganizationLogin string
}
