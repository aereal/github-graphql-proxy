package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/aereal/github-graphql-proxy/graph/dto"
	"github.com/aereal/github-graphql-proxy/graph/handler"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Actions is the resolver for the actions field.
func (r *organizationBillingResolver) Actions(ctx context.Context, obj *dto.OrganizationBilling) (*dto.ActionBilling, error) {
	billing, _, err := r.githubClient.Billing.GetActionsBillingOrg(ctx, obj.OrganizationLogin)
	if err != nil {
		return nil, gqlerror.Errorf("Billing.GetActionsBillingOrg: %s", err)
	}
	return &dto.ActionBilling{
		TotalMinutesUsed:     billing.TotalMinutesUsed,
		TotalPaidMinutesUsed: billing.TotalPaidMinutesUsed,
		IncludedMinutes:      billing.IncludedMinutes,
		MinutedUsedBreakdown: &dto.ActionBillingBreakdown{
			Ubuntu: &dto.ActionBillingBreakdownUbuntu{
				Total: &billing.MinutesUsedBreakdown.Ubuntu,
			},
			MacOs: &dto.ActionBillingBreakdownMacOs{
				Total: &billing.MinutesUsedBreakdown.MacOS,
			},
			Windows: &dto.ActionBillingBreakdownWindows{
				Total: &billing.MinutesUsedBreakdown.Windows,
			},
		},
	}, nil
}

// Organization is the resolver for the organization field.
func (r *queryResolver) Organization(ctx context.Context, login string) (*dto.Organization, error) {
	return &dto.Organization{Login: login, Billing: &dto.OrganizationBilling{OrganizationLogin: login}}, nil
}

// OrganizationBilling returns handler.OrganizationBillingResolver implementation.
func (r *Resolver) OrganizationBilling() handler.OrganizationBillingResolver {
	return &organizationBillingResolver{r}
}

// Query returns handler.QueryResolver implementation.
func (r *Resolver) Query() handler.QueryResolver { return &queryResolver{r} }

type organizationBillingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
