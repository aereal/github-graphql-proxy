package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/aereal/github-graphql-proxy/graph/dto"
	"github.com/aereal/github-graphql-proxy/graph/handler"
	"github.com/google/go-github/v47/github"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Plan is the resolver for the plan field.
func (r *organizationResolver) Plan(ctx context.Context, obj *dto.Organization) (*dto.Plan, error) {
	org, _, err := r.githubClient.Organizations.Get(ctx, obj.Login)
	if err != nil {
		return nil, gqlerror.Errorf("Organizations.Get: %s", err)
	}
	if org.Plan == nil {
		return nil, nil
	}
	return &dto.Plan{
		Name:          org.Plan.Name,
		Space:         org.Plan.Space,
		Collaborators: org.Plan.Collaborators,
		PrivateRepos:  org.Plan.PrivateRepos,
		FilledSeats:   org.Plan.FilledSeats,
		Seats:         org.Plan.Seats,
	}, nil
}

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

// Storage is the resolver for the storage field.
func (r *organizationBillingResolver) Storage(ctx context.Context, obj *dto.OrganizationBilling) (*dto.StorageBilling, error) {
	billing, _, err := r.githubClient.Billing.GetStorageBillingOrg(ctx, obj.OrganizationLogin)
	if err != nil {
		return nil, gqlerror.Errorf("Billing.GetStorageBillingOrg: %s", err)
	}
	return &dto.StorageBilling{
		DaysLeftInBillingCycle:       billing.DaysLeftInBillingCycle,
		EstimatedPaidStorageForMonth: billing.EstimatedPaidStorageForMonth,
		EstimatedStorageForMonth:     billing.EstimatedStorageForMonth,
	}, nil
}

// Organization is the resolver for the organization field.
func (r *queryResolver) Organization(ctx context.Context, login string) (*dto.Organization, error) {
	return &dto.Organization{Login: login, Billing: &dto.OrganizationBilling{OrganizationLogin: login}}, nil
}

// Repository is the resolver for the repository field.
func (r *queryResolver) Repository(ctx context.Context, owner string, name string) (*dto.Repository, error) {
	return &dto.Repository{Owner: owner, Name: name}, nil
}

// Artifacts is the resolver for the artifacts field.
func (r *repositoryResolver) Artifacts(ctx context.Context, obj *dto.Repository, first *int, page *int) (*dto.RepositoryArtifactConnection, error) {
	listOpts := &github.ListOptions{}
	if first != nil {
		listOpts.PerPage = *first
	}
	if page != nil {
		listOpts.Page = *page
	}
	artifacts, _, err := r.githubClient.Actions.ListArtifacts(ctx, obj.Owner, obj.Name, listOpts)
	if err != nil {
		return nil, gqlerror.Errorf("Actions.ListArtifacts: %v", err)
	}
	out := &dto.RepositoryArtifactConnection{Nodes: make([]*dto.Artifact, len(artifacts.Artifacts))}
	out.TotalCount = len(artifacts.Artifacts)
	for i, artifact := range artifacts.Artifacts {
		size := artifact.GetSizeInBytes()
		out.TotalSizeInBytes += size
		a := &dto.Artifact{
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

// Organization returns handler.OrganizationResolver implementation.
func (r *Resolver) Organization() handler.OrganizationResolver { return &organizationResolver{r} }

// OrganizationBilling returns handler.OrganizationBillingResolver implementation.
func (r *Resolver) OrganizationBilling() handler.OrganizationBillingResolver {
	return &organizationBillingResolver{r}
}

// Query returns handler.QueryResolver implementation.
func (r *Resolver) Query() handler.QueryResolver { return &queryResolver{r} }

// Repository returns handler.RepositoryResolver implementation.
func (r *Resolver) Repository() handler.RepositoryResolver { return &repositoryResolver{r} }

type organizationResolver struct{ *Resolver }
type organizationBillingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type repositoryResolver struct{ *Resolver }
