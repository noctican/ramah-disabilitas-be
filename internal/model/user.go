package model

import (
	"time"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleAdmin   UserRole = "admin"
)

type RankTier string

const (
	RankBronze   RankTier = "Bronze"
	RankSilver   RankTier = "Silver"
	RankGold     RankTier = "Gold"
	RankPlatinum RankTier = "Platinum"
	RankDiamond  RankTier = "Diamond"
)

type User struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string    `gorm:"type:varchar(255)" json:"name"`
	Email            string    `gorm:"type:varchar(255);unique" json:"email"`
	Password         string    `gorm:"type:varchar(255)" json:"-"`
	Role             UserRole  `gorm:"type:varchar(20)" json:"role"`
	Avatar           string    `gorm:"type:varchar(255)" json:"avatar"`
	Points           int       `gorm:"default:0" json:"points"`
	RankTier         RankTier  `gorm:"type:varchar(20)" json:"rank_tier"`
	MMR              int       `json:"mmr"`
	CurrentStreak    int       `gorm:"default:0" json:"current_streak"`
	LastActivityDate time.Time `gorm:"type:date" json:"last_activity_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type FriendshipStatus string

const (
	StatusPending  FriendshipStatus = "pending"
	StatusAccepted FriendshipStatus = "accepted"
	StatusBlocked  FriendshipStatus = "blocked"
)

type Friendship struct {
	ID          uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	RequesterID uint64           `json:"requester_id"`
	AddresseeID uint64           `json:"addressee_id"`
	Status      FriendshipStatus `gorm:"type:varchar(20)" json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	Requester   User             `gorm:"foreignKey:RequesterID" json:"requester"`
	Addressee   User             `gorm:"foreignKey:AddresseeID" json:"addressee"`
}
