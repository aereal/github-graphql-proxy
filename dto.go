package githubgraphqlproxy

type Organization struct {
	Login   string `json:"login"`
	Billing *OrganizationBilling
}

func (Organization) IsEntity() {}

type OrganizationBilling struct {
	OrganizationLogin string
}

type Repository struct {
	Owner         string `json:"-"`
	Name          string `json:"name"`
	NameWithOwner string
}

func (Repository) IsEntity() {}

type RepositoryArtifactConnection struct {
	TotalCount       int         `json:"totalCount"`
	TotalSizeInBytes int64       `json:"totalSizeInBytes"`
	Nodes            []*Artifact `json:"nodes"`
}
