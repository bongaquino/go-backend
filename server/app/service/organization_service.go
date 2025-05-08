package service

import (
	"context"
	"errors"
	"fmt"
	"koneksi/server/app/dto"
	"koneksi/server/app/model"
	"koneksi/server/app/repository"
	"koneksi/server/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationService struct {
	orgRepo         *repository.OrganizationRepository
	policyRepo      *repository.PolicyRepository
	permissionRepo  *repository.PermissionRepository
	orgUserRoleRepo *repository.OrganizationUserRoleRepository
	userRepo        *repository.UserRepository
	roleRepo        *repository.RoleRepository
}

func NewOrganizationService(orgRepo *repository.OrganizationRepository,
	policyRepo *repository.PolicyRepository,
	permissionRepo *repository.PermissionRepository,
	orgUserRoleRepo *repository.OrganizationUserRoleRepository,
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:         orgRepo,
		policyRepo:      policyRepo,
		permissionRepo:  permissionRepo,
		orgUserRoleRepo: orgUserRoleRepo,
		userRepo:        userRepo,
		roleRepo:        roleRepo,
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

func (os *OrganizationService) ReadOrg(ctx context.Context, orgID string) (*model.Organization, []map[string]string, error) {
	// Fetch the organization from the repository
	org, err := os.orgRepo.Read(ctx, orgID)
	if err != nil {
		logger.Log.Error("error fetching organization", logger.Error(err))
		return nil, nil, errors.New("error fetching organization")
	}
	if org == nil {
		return nil, nil, errors.New("organization not found")
	}

	// Fetch organization members
	orgMembers, err := os.orgUserRoleRepo.ReadByOrganizationID(ctx, orgID)
	if err != nil {
		logger.Log.Error("error fetching organization members", logger.Error(err))
		return nil, nil, errors.New("error fetching organization members")
	}

	// Fetch roles and map to ID => Name
	roles, err := os.roleRepo.List(ctx)
	if err != nil {
		logger.Log.Error("error fetching roles", logger.Error(err))
		return nil, nil, errors.New("error fetching roles")
	}
	roleMap := make(map[string]string, len(roles))
	for _, role := range roles {
		roleMap[role.ID.Hex()] = role.Name
	}

	// Pre-allocate members slice
	members := make([]map[string]string, 0, len(orgMembers))

	// Loop through members and populate user info and role name
	for _, member := range orgMembers {
		user, err := os.userRepo.Read(ctx, member.UserID.Hex())
		if err != nil {
			logger.Log.Error("error fetching user details", logger.Error(err))
			return nil, nil, errors.New("error fetching user details")
		}
		if user != nil {
			members = append(members, map[string]string{
				"email": user.Email,
				"role":  roleMap[member.RoleID.Hex()],
			})
		}
	}

	return org, members, nil
}

func (os *OrganizationService) UpdateOrg(ctx context.Context, orgID string, dto *dto.UpdateOrgDTO) (*model.Organization, error) {
	// Check if the organization exists
	org, err := os.orgRepo.Read(ctx, orgID)
	if err != nil {
		logger.Log.Error("error fetching organization", logger.Error(err))
		return nil, errors.New("error fetching organization")
	}
	if org == nil {
		return nil, errors.New("organization not found")
	}

	// Precompute policyID and parentID
	policyID, err := primitive.ObjectIDFromHex(dto.PolicyID)
	if err != nil {
		logger.Log.Error("invalid policy ID", logger.Error(err))
		return nil, errors.New("invalid policy ID")
	}

	var parentID primitive.ObjectID
	if dto.ParentID != nil {
		parentID, err = primitive.ObjectIDFromHex(*dto.ParentID)
		if err != nil {
			logger.Log.Error("invalid parent ID", logger.Error(err))
			return nil, errors.New("invalid parent ID")
		}
	}

	// Update the organization fields
	orgUpdate := bson.M{
		"name":      dto.Name,
		"domain":    dto.Domain,
		"contact":   dto.Contact,
		"policy_id": policyID,
		"parent_id": parentID,
	}

	// Update the organization in the repository
	err = os.orgRepo.Update(ctx, orgID, orgUpdate)
	if err != nil {
		logger.Log.Error("error updating organization", logger.Error(err))
		return nil, errors.New("error updating organization")
	}

	// Fetch the updated organization
	updatedOrg, err := os.orgRepo.Read(ctx, orgID)
	if err != nil {
		logger.Log.Error("error fetching updated organization", logger.Error(err))
		return nil, errors.New("error fetching updated organization")
	}
	if updatedOrg == nil {
		logger.Log.Error("organization not found after update", logger.Error(err))
		return nil, errors.New("organization not found after update")
	}

	// Return the updated organization
	return updatedOrg, nil
}

func (os *OrganizationService) AddMember(ctx context.Context, orgID string, userID string, roleID string) error {
	// Check if the organization exists
	org, err := os.orgRepo.Read(ctx, orgID)
	if err != nil {
		logger.Log.Error("error fetching organization", logger.Error(err))
		return errors.New("error fetching organization")
	}
	if org == nil {
		return errors.New("organization not found")
	}

	// Check if the user exists
	user, err := os.userRepo.Read(ctx, userID)
	if err != nil {
		logger.Log.Error("error fetching user", logger.Error(err))
		return errors.New("error fetching user")
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Check if the role exists
	role, err := os.roleRepo.Read(ctx, roleID)
	if err != nil {
		logger.Log.Error("error fetching role", logger.Error(err))
		return errors.New("error fetching role")
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Check if the user is already a member of the organization with the same role
	existingMember, err := os.orgUserRoleRepo.ReadByUserIDOrganizationID(ctx, userID, orgID)
	if err != nil {
		logger.Log.Error("error checking existing member", logger.Error(err))
		return errors.New("error checking existing member")
	}
	fmt.Println("existingMember", existingMember)
	if existingMember != nil {
		return errors.New("user is already a member")
	}

	// Check if the user is already a member of another organization
	orgUserRoles, err := os.orgUserRoleRepo.ReadByUserID(ctx, userID)
	if err != nil {
		logger.Log.Error("error checking existing member", logger.Error(err))
		return errors.New("error checking existing member")
	}
	for _, orgUserRole := range orgUserRoles {
		if orgUserRole.OrganizationID != org.ID {
			return errors.New("user is already a member of another organization")
		}
	}

	// Add the member to the organization
	err = os.orgUserRoleRepo.Create(ctx, &model.OrganizationUserRole{
		OrganizationID: org.ID,
		UserID:         user.ID,
		RoleID:         role.ID,
	})
	if err != nil {
		logger.Log.Error("error adding member to organization", logger.Error(err))
		return errors.New("error adding member to organization")
	}

	return nil
}
