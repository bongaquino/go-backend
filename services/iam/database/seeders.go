package database

import (
	"context"
	"fmt"
	"koneksi/services/iam/app/models"
	"koneksi/services/iam/app/services/mongo"
	"koneksi/services/iam/core/logger"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// SeedCollections seeds initial data into MongoDB collections
func SeedCollections(mongoService *mongo.MongoService) {
	db := mongoService.GetDB()
	ctx := context.Background()

	seeders := []struct {
		Name   string
		Seeder func(*mongoDriver.Database, context.Context) error
	}{
		{
			Name:   "permissions",
			Seeder: seedPermissions,
		},
		{
			Name:   "roles",
			Seeder: seedRoles,
		},
		{
			Name:   "role_permissions",
			Seeder: seedRolePermissions,
		},
	}

	for _, seeder := range seeders {
		if err := seeder.Seeder(db, ctx); err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to seed collection: %s", seeder.Name), logger.Error(err))
		} else {
			logger.Log.Info(fmt.Sprintf("Seeded collection: %s", seeder.Name))
		}
	}
}

// seedPermissions inserts initial permissions
func seedPermissions(db *mongoDriver.Database, ctx context.Context) error {
	collection := db.Collection("permissions")

	existing, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if existing > 0 {
		logger.Log.Info("Skipping permissions seeding: Data already exists")
		return nil
	}

	permissions := []any{
		models.Permission{Name: "upload_files"},
		models.Permission{Name: "download_files"},
		models.Permission{Name: "list_files"},
	}

	_, err = collection.InsertMany(ctx, permissions)
	return err
}

// seedRoles inserts initial roles
func seedRoles(db *mongoDriver.Database, ctx context.Context) error {
	collection := db.Collection("roles")

	existing, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if existing > 0 {
		logger.Log.Info("Skipping roles seeding: Data already exists")
		return nil
	}

	roles := []any{
		models.Role{Name: "user"},
	}

	_, err = collection.InsertMany(ctx, roles)
	return err
}

// seedRolePermissions assigns all permissions to the "user" role
func seedRolePermissions(db *mongoDriver.Database, ctx context.Context) error {
	roleCollection := db.Collection("roles")
	permissionCollection := db.Collection("permissions")
	rolePermissionCollection := db.Collection("role_permissions")

	// Find the "user" role
	var userRole models.Role
	err := roleCollection.FindOne(ctx, bson.M{"name": "user"}).Decode(&userRole)
	if err != nil {
		return fmt.Errorf("failed to find user role: %w", err)
	}

	// Get all permissions
	cursor, err := permissionCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to retrieve permissions: %w", err)
	}
	defer cursor.Close(ctx)

	var permissions []models.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		return fmt.Errorf("failed to decode permissions: %w", err)
	}

	// Check if data already exists
	existing, err := rolePermissionCollection.CountDocuments(ctx, bson.M{"role_id": userRole.ID})
	if err != nil {
		return err
	}
	if existing > 0 {
		logger.Log.Info("Skipping role_permissions seeding: Data already exists")
		return nil
	}

	// Create role-permission entries
	var rolePermissions []any
	for _, perm := range permissions {
		rolePermissions = append(rolePermissions, models.RolePermission{
			RoleID:       userRole.ID,
			PermissionID: perm.ID,
		})
	}

	_, err = rolePermissionCollection.InsertMany(ctx, rolePermissions)
	return err
}
