package dto

import (
	"subscription/internal/pkg/types"

	"github.com/google/uuid"
)

type AddService struct {
	Name *string `json:"name" validate:"required,min=1"`
}
type GetService struct {
	ServiceId *int `params:"id" validate:"required,gte=1"`
}
type UpdateService struct {
	ServiceId *int    `params:"id" validate:"required,gte=1"`
	Name      *string `json:"name" validate:"required,min=1"`
}
type RemoveService struct {
	ServiceId *int `params:"id" validate:"required,gte=1"`
}

type AddSubscription struct {
	ServiceName *string           `json:"service_name" validate:"required,min=1"`
	Price       *int              `json:"price" validate:"required,gte=0"`
	UserId      *uuid.UUID        `json:"user_id" validate:"required"`
	StartDate   *types.CustomDate `json:"start_date" validate:"required"`
	StopDate    *types.CustomDate `json:"stop_date" validate:"omitempty"`
}

type GetSubscription struct {
	SubscriptionId *int `params:"id" validate:"required,gte=1"`
}
type GetSubscriptions struct {
	Offset *int `query:"offset" validate:"omitempty,gte=0"`
	Limit  *int `query:"limit" validate:"omitempty,gte=0"`
}
type GetSubscriptionTotal struct {
	StartDate   *types.CustomDate `json:"start_date" validate:"required"`
	StopDate    *types.CustomDate `json:"stop_date" validate:"required"`
	UserId      *uuid.UUID        `json:"user_id" validate:"omitempty"`
	ServiceName *string           `json:"service_name" validate:"omitempty,min=1"`
}

type UpdateSubscription struct {
	SubscriptionId *int              `params:"id" validate:"required,gte=1"`
	ServiceName    *string           `json:"service_name" validate:"required,min=1"`
	Price          *int              `json:"price" validate:"required,gte=0"`
	UserId         *uuid.UUID        `json:"user_id" validate:"required"`
	StartDate      *types.CustomDate `json:"start_date" validate:"required"`
	StopDate       *types.CustomDate `json:"stop_date,omitempty"`
}
type RemoveSubscription struct {
	SubscriptionId *int `params:"id" validate:"required,gte=1"`
}
