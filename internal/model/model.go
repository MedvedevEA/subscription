package model

import (
	"subscription/internal/pkg/types"

	"github.com/google/uuid"
)

type Service struct {
	ServiceId int    `json:"service_id"`
	Name      string `json:"name"`
}

type Subscription struct {
	SubscriptionId int               `json:"subscription_id"`
	ServiceId      int               `json:"service_id"`
	Price          int               `json:"price"`
	UserId         uuid.UUID         `json:"user_id"`
	StartDate      types.CustomDate  `json:"start_date"`
	StopDate       *types.CustomDate `json:"stop_date"`
}
