package resolvers

import "errors"

var (
	ErrOrganizationPlanIsNil = errors.New("organization.plan in the response from GitHub is nil")
)
