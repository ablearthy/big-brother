// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

type VkActivity string

const (
	VkActivityOnline  VkActivity = "online"
	VkActivityOffline VkActivity = "offline"
)

func (e *VkActivity) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VkActivity(s)
	case string:
		*e = VkActivity(s)
	default:
		return fmt.Errorf("unsupported scan type for VkActivity: %T", src)
	}
	return nil
}

type VkMessageEventType string

const (
	VkMessageEventTypeNew    VkMessageEventType = "new"
	VkMessageEventTypeEdit   VkMessageEventType = "edit"
	VkMessageEventTypeDelete VkMessageEventType = "delete"
)

func (e *VkMessageEventType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VkMessageEventType(s)
	case string:
		*e = VkMessageEventType(s)
	default:
		return fmt.Errorf("unsupported scan type for VkMessageEventType: %T", src)
	}
	return nil
}

type VkPlatform string

const (
	VkPlatformMobile  VkPlatform = "mobile"
	VkPlatformIphone  VkPlatform = "iphone"
	VkPlatformIpad    VkPlatform = "ipad"
	VkPlatformAndroid VkPlatform = "android"
	VkPlatformWphone  VkPlatform = "wphone"
	VkPlatformWindows VkPlatform = "windows"
	VkPlatformWeb     VkPlatform = "web"
)

func (e *VkPlatform) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = VkPlatform(s)
	case string:
		*e = VkPlatform(s)
	default:
		return fmt.Errorf("unsupported scan type for VkPlatform: %T", src)
	}
	return nil
}

type InviteCode struct {
	UserID     int32
	InviteCode string
}

type User struct {
	ID        int32
	Username  string
	Password  string
	InviterID int32
}

type UserToken struct {
	UserID      int32
	AccessToken sql.NullString
}

type VkActivityEvent struct {
	ID              int32
	VkOwnerID       int32
	TargetID        int32
	Activity        VkActivity
	Platform        VkPlatform
	KickedByTimeout sql.NullBool
	CreatedAt       time.Time
}

type VkMessage struct {
	ID        int32
	VkOwnerID int32
	MessageID int32
	Message   pgtype.JSONB
}

type VkMessageEvent struct {
	ID                int32
	InternalMessageID int32
	MType             VkMessageEventType
	CreatedAt         time.Time
}

type VkToken struct {
	AccessToken string
	VkUserID    int32
}
