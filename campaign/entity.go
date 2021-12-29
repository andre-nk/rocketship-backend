package campaign

import (
	"rocketship/user"
	"time"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	Description      string
	ShortDescription string
	Perks            string
	FunderAmount     int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage
	User             user.User
}

type CampaignImage struct {
	ID         int
	CampaignID int
	IsPrimary  int
	FileName   string
	UpdatedAt  time.Time
	CreatedAt  time.Time
}
