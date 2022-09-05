package dto

type Organization struct {
	Login   string `json:"login"`
	Billing *OrganizationBilling
}

type OrganizationBilling struct {
	OrganizationLogin string
}

type Repository struct {
	Owner string `json:"-"`
	Name  string `json:"name"`
}

type RepositoryArtifactConnection struct {
	TotalCount       int         `json:"totalCount"`
	TotalSizeInBytes int64       `json:"totalSizeInBytes"`
	Nodes            []*Artifact `json:"nodes"`
}
