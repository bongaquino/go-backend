package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	ID                         primitive.ObjectID `bson:"_id,omitempty"`
	UserID                     primitive.ObjectID `bson:"user_id"`
	BackupCycle                string             `bson:"backup_cycle"`
	BackupCustomDay            int                `bson:"backup_custom_day"`
	NotificationsFrequency     string             `bson:"notifications_frequency"`
	RecoveryPriorityOrder      string             `bson:"recovery_priority_order"`
	RecoveryCustomOrder        string             `bson:"recovery_custom_order"`
	IsMFAEnabled               bool               `bson:"is_mfa_enabled"`
	IsRealtimeBackupEnabled    bool               `bson:"is_realtime_backup_enabled"`
	IsEmailNotificationEnabled bool               `bson:"is_email_notification_enabled"`
	IsSMSNotificationEnabled   bool               `bson:"is_sms_notification_enabled"`
	IsVersionHistoryEnabled    bool               `bson:"is_version_history_enabled"`
	CreatedAt                  time.Time          `bson:"created_at"`
	UpdatedAt                  time.Time          `bson:"updated_at"`
}

func (Setting) GetIndexes() []bson.D {
	return []bson.D{
		{{Key: "user_id", Value: 1}},
	}
}
