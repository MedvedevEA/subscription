package dto

import (
	"subscription/internal/pkg/types"

	"github.com/google/uuid"
)

type UpdateService struct {
	ServiceId int
	Name      string
}

type AddSubscription struct {
	ServiceName string
	Price       int
	UserId      uuid.UUID
	StartDate   types.CustomDate
	StopDate    *types.CustomDate
}

type GetSubscriptions struct {
	Offset *int
	Limit  *int
}
type GetSubscriptionTotal struct {
	StartDate   types.CustomDate
	StopDate    types.CustomDate
	UserId      *uuid.UUID
	ServiceName *string
}

type UpdateSubscription struct {
	SubscriptionId int
	ServiceName    string
	Price          int
	UserId         uuid.UUID
	StartDate      types.CustomDate
	StopDate       *types.CustomDate
}
