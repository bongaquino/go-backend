package service

import (
	"context"
	"errors"
	"koneksi/server/app/dto"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationService struct {
	orgRepo        *repository.OrganizationRepository
	policyRepo     *repository.PolicyRepository
	permissionRepo *repository.PermissionRepository
}

func NewOrganizationService(orgRepo *repository.OrganizationRepository,
	policyRepo *repository.PolicyRepository,
	permissionRepo *repository.PermissionRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:        orgRepo,
		policyRepo:     policyRepo,
		permissionRepo: permissionRepo,
	}
}

func (os *OrganizationService) ListPermissions(ctx context.Context) ([]*model.Permission, error) {
	// Fetch permissions from the repository
	permissions, err := os.permissionRepo.List(ctx)
	if err != nil {
		logger.Log.Error("error fetching permissions", logger.Error(err))
		return nil, errors.New("error fetching permissions")
	}
	// Convert []model.Permission to []*model.Permission
	permissionPointers := make([]*model.Permission, len(permissions))
	copy(permissionPointers, permissions)
	return permissionPointers, nil
}

func (os *OrganizationService) ListPolicies(ctx context.Context) ([]*model.Policy, error) {
	// Fetch policies from the repository
	policies, err := os.policyRepo.List(ctx)
	if err != nil {
		logger.Log.Error("error fetching policies", logger.Error(err))
		return nil, errors.New("error fetching policies")
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
		return nil, errors.New("error fetching orgs")
	}
	// Convert []model.Organization to []*model.Organization
	orgPointers := make([]*model.Organization, len(orgs))
	for i := range orgs {
		orgPointers[i] = &orgs[i]
	}
	return orgPointers, nil
}

func (os *OrganizationService) CreateOrg(ctx context.Context, request *dto.CreateOrgDTO) (*model.Organization, error) {
	// Map the request to the organization model
	org := &model.Organization{
		Name:    request.Name,
		Domain:  request.Domain,
		Contact: request.Contact,
		PolicyID: func() primitive.ObjectID {
			policyID, err := primitive.ObjectIDFromHex(request.PolicyID)
			if err != nil {
				logger.Log.Error("invalid policy ID", logger.Error(err))
				return primitive.NilObjectID
			}
			return policyID
		}(),
		SubscriptionPlanID:   primitive.NilObjectID,
		SubscriptionStatusID: primitive.NilObjectID,
		ParentID: func() primitive.ObjectID {
			if request.ParentID != nil {
				parentID, err := primitive.ObjectIDFromHex(*request.ParentID)
				if err == nil {
					return parentID
				}
				logger.Log.Error("invalid parent ID", logger.Error(err))
			}
			return primitive.NilObjectID
		}(),
	}

	// Create the organization
	err := os.orgRepo.Create(ctx, org)
	if err != nil {
		logger.Log.Error("error creating organization", logger.Error(err))
		return nil, errors.New("internal server error")
	}

	return org, nil
}
