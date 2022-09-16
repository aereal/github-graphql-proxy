package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	githubgraphqlproxy "github.com/aereal/github-graphql-proxy"
	"github.com/google/go-github/v47/github"
)

// Plan is the resolver for the plan field.
func (r *organizationResolver) Plan(ctx context.Context, obj *githubgraphqlproxy.Organization) (*githubgraphqlproxy.Plan, error) {
	org, _, err := r.githubClient.Organizations.Get(ctx, obj.Login)
	if err != nil {
		return nil, fmt.Errorf("Organizations.Get: %w", err)
	}
	if org.Plan == nil {
		return nil, ErrOrganizationPlanIsNil
	}
	return &githubgraphqlproxy.Plan{
		Name:          org.Plan.Name,
		Space:         org.Plan.Space,
		Collaborators: org.Plan.Collaborators,
		PrivateRepos:  org.Plan.PrivateRepos,
		FilledSeats:   org.Plan.FilledSeats,
		Seats:         org.Plan.Seats,
	}, nil
}

// Actions is the resolver for the actions field.
func (r *organizationBillingResolver) Actions(ctx context.Context, obj *githubgraphqlproxy.OrganizationBilling) (*githubgraphqlproxy.ActionBilling, error) {
	billing, _, err := r.githubClient.Billing.GetActionsBillingOrg(ctx, obj.OrganizationLogin)
	if err != nil {
		return nil, fmt.Errorf("Billing.GetActionsBillingOrg: %w", err)
	}
	return &githubgraphqlproxy.ActionBilling{
		TotalMinutesUsed:     billing.TotalMinutesUsed,
		TotalPaidMinutesUsed: billing.TotalPaidMinutesUsed,
		IncludedMinutes:      billing.IncludedMinutes,
		MinutedUsedBreakdown: &githubgraphqlproxy.ActionBillingBreakdown{
			Ubuntu: &githubgraphqlproxy.ActionBillingBreakdownUbuntu{
				Total: &billing.MinutesUsedBreakdown.Ubuntu,
			},
			MacOs: &githubgraphqlproxy.ActionBillingBreakdownMacOs{
				Total: &billing.MinutesUsedBreakdown.MacOS,
			},
			Windows: &githubgraphqlproxy.ActionBillingBreakdownWindows{
				Total: &billing.MinutesUsedBreakdown.Windows,
			},
		},
	}, nil
}

// Storage is the resolver for the storage field.
func (r *organizationBillingResolver) Storage(ctx context.Context, obj *githubgraphqlproxy.OrganizationBilling) (*githubgraphqlproxy.StorageBilling, error) {
	billing, _, err := r.githubClient.Billing.GetStorageBillingOrg(ctx, obj.OrganizationLogin)
	if err != nil {
		return nil, fmt.Errorf("Billing.GetStorageBillingOrg: %w", err)
	}
	return &githubgraphqlproxy.StorageBilling{
		DaysLeftInBillingCycle:       billing.DaysLeftInBillingCycle,
		EstimatedPaidStorageForMonth: billing.EstimatedPaidStorageForMonth,
		EstimatedStorageForMonth:     billing.EstimatedStorageForMonth,
	}, nil
}

// Organization is the resolver for the organization field.
func (r *queryResolver) Organization(ctx context.Context, login string) (*githubgraphqlproxy.Organization, error) {
	return &githubgraphqlproxy.Organization{Login: login, Billing: &githubgraphqlproxy.OrganizationBilling{OrganizationLogin: login}}, nil
}

// Repository is the resolver for the repository field.
func (r *queryResolver) Repository(ctx context.Context, owner string, name string, followRenames *bool) (*githubgraphqlproxy.Repository, error) {
	return &githubgraphqlproxy.Repository{Owner: owner, Name: name}, nil
}

// Artifacts is the resolver for the artifacts field.
func (r *repositoryResolver) Artifacts(ctx context.Context, obj *githubgraphqlproxy.Repository, first *int, page *int) (*githubgraphqlproxy.RepositoryArtifactConnection, error) {
	listOpts := &github.ListOptions{}
	if first != nil {
		listOpts.PerPage = *first
	}
	if page != nil {
		listOpts.Page = *page
	}
	artifacts, _, err := r.githubClient.Actions.ListArtifacts(ctx, obj.Owner, obj.Name, listOpts)
	if err != nil {
		return nil, fmt.Errorf("Actions.ListArtifacts: %w", err)
	}
	out := &githubgraphqlproxy.RepositoryArtifactConnection{Nodes: make([]*githubgraphqlproxy.Artifact, len(artifacts.Artifacts))}
	out.TotalCount = len(artifacts.Artifacts)
	for i, artifact := range artifacts.Artifacts {
		size := artifact.GetSizeInBytes()
		out.TotalSizeInBytes += size
		a := &githubgraphqlproxy.Artifact{
			ID:                 int(artifact.GetID()),
			Name:               artifact.GetName(),
			SizeInBytes:        int(size),
			ArchiveDownloadURL: artifact.GetArchiveDownloadURL(),
			Expired:            artifact.GetExpired(),
		}
		if artifact.CreatedAt != nil {
			a.CreatedAt = artifact.CreatedAt.Time
		}
		if artifact.ExpiresAt != nil {
			a.ExpiresAt = artifact.ExpiresAt.Time
		}
		out.Nodes[i] = a
	}
	return out, nil
}

// Organization returns githubgraphqlproxy.OrganizationResolver implementation.
func (r *Resolver) Organization() githubgraphqlproxy.OrganizationResolver {
	return &organizationResolver{r}
}

// OrganizationBilling returns githubgraphqlproxy.OrganizationBillingResolver implementation.
func (r *Resolver) OrganizationBilling() githubgraphqlproxy.OrganizationBillingResolver {
	return &organizationBillingResolver{r}
}

// Query returns githubgraphqlproxy.QueryResolver implementation.
func (r *Resolver) Query() githubgraphqlproxy.QueryResolver { return &queryResolver{r} }

// Repository returns githubgraphqlproxy.RepositoryResolver implementation.
func (r *Resolver) Repository() githubgraphqlproxy.RepositoryResolver { return &repositoryResolver{r} }

type organizationResolver struct{ *Resolver }
type organizationBillingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type repositoryResolver struct{ *Resolver }
