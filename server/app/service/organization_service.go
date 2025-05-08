package service

import (
	"context"
	"errors"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"
)

type OrganizationService struct {
	orgRepo      *repository.OrganizationRepository
	policiesRepo *repository.PolicyRepository
}

func NewOrganizationService(orgRepo *repository.OrganizationRepository, policiesRepo *repository.PolicyRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:      orgRepo,
		policiesRepo: policiesRepo,
	}
}

func (os *OrganizationService) ListPolicies(ctx context.Context) ([]*model.Policy, error) {
	// Fetch policies from the repository
	policies, err := os.policiesRepo.List(ctx)
	if err != nil {
		logger.Log.Error("error fetching policies", logger.Error(err))
		return nil, errors.New("internal server error")
	}
	// Convert []model.Policy to []*model.Policy
	policyPointers := make([]*model.Policy, len(policies))
	copy(policyPointers, policies)
	return policyPointers, nil
}

func (os *OrganizationService) ListOrgs(ctx context.Context, page, limit int) ([]*model.Organization, error) {
	// Fetch orgs from the repository
	orgs, err := os.orgRepo.List(ctx, page, limit)
	if err != nil {
		logger.Log.Error("error fetching orgs", logger.Error(err))
		return nil, errors.New("internal server error")
	}
	// Convert []model.Organization to []*model.Organization
	orgPointers := make([]*model.Organization, len(orgs))
	for i := range orgs {
		orgPointers[i] = &orgs[i]
	}
	return orgPointers, nil
}
