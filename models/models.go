package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	FirstName    string `gorm:"not null" validate:"required,max=20"`
	MiddleName   string
	LastName     string `gorm:"not null" validate:"required,max=20"`
	Mobile       string `gorm:"unique;not null" validate:"required,min=5,max=20"`
	Email        string `gorm:"unique;not null" validate:"required,min=5,max=45"`
	PasswordHash string `gorm:"not null" validate:"required,min=5,max=85"`
	RegisteredAt time.Time
	LastLogin    time.Time
	Intro        string
	Friends      []UserFriend   `gorm:"foreignKey:SourceID;references:ID"`
	Followers    []UserFollower `gorm:"foreignKey:SourceID;references:ID"`
	Messages     []Message      `gorm:"foreignKey:MessageSenderID;references:ID"`
	Posts        []Post         `gorm:"foreignKey:UserID;references:ID"`
	Groups       []Group        `gorm:"foreignKey:CreatedByID;references:ID"`
	Token        Token          `gorm:"foreignKey:UserID;references:ID"`
}

type UserFriend struct {
	gorm.Model

	SourceID uint `gorm:"not null;type:bigint;index;uniqueIndex:idx_source_target_followers"`
	Source   User `gorm:"foreignKey:SourceID;references:ID"`
	TargetID uint `gorm:"not null;type:bigint;index;uniqueIndex:idx_source_target_followers"`
	Target   User `gorm:"foreignKey:TargetID;references:ID"`
	Type     int
	Status   int
	Notes    string
}

type UserFollower struct {
	gorm.Model

	SourceID uint `gorm:"not null;type:bigint;uniqueIndex:idx_source_target"`
	Source   User `gorm:"foreignKey:SourceID"`
	TargetID uint `gorm:"not null;type:bigint;uniqueIndex:idx_source_target"`
	Target   User `gorm:"foreignKey:TargetID"`
	Type     int
}

type Message struct {
	gorm.Model

	MessageSenderID    uint `gorm:"not null" gorm:"type:bigint;index"`
	MessageSender      User `gorm:"foreignKey:MessageSenderID"`
	MessageRecipientID uint `gorm:"not null" gorm:"type:bigint;index"`
	MessageRecipient   User `gorm:"foreignKey:MessageRecipientID"`
	Message            string
	
}

type Post struct {
	gorm.Model

	UserID   uint `gorm:"type:bigint;index"`
	User     User `gorm:"foreignKey:UserID"`
	SenderID uint `gorm:"type:bigint;index"`
	Sender   User `gorm:"foreignKey:SenderID"`
	Message  string
}

type Group struct {
	gorm.Model

	CreatedByID uint `gorm:"type:bigint;index"`
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`
	UpdatedByID uint `gorm:"type:bigint;index"`
	UpdatedBy   User `gorm:"foreignKey:UpdatedByID"`
	Title       string
	MetaTitle   string
	Slug        string `gorm:"unique"`
	Summary     string
	Status      int
	Profile     string
	Content     string
}

type GroupMeta struct {
	gorm.Model

	GroupID uint
	Group   Group `gorm:"foreignKey:GroupID"`
	Key     string
	Content string
}

type GroupMember struct {
	gorm.Model

	GroupID uint
	Group   Group `gorm:"foreignKey:GroupID"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	Status  int
	Notes   string
}

type GroupMessage struct {
	gorm.Model

	GroupID uint
	Group   Group `gorm:"foreignKey:GroupID"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	Message string
}

type Token struct {
	gorm.Model

	UserID    uint
	Token     string
	ExpiresAt int64
}

func AutoMigrate(db *gorm.DB) {
	// AutoMigrate will create the necessary tables in the database
	db.AutoMigrate(&User{}, &Message{}, &UserFriend{}, &UserFollower{}, &Message{}, &Post{}, &Group{}, &GroupMeta{}, &GroupMember{}, &GroupMessage{}, &Token{})

}
