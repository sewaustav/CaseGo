package dto

import (
	"time"

	"github.com/sewaustav/Payment/internal/payment/models"
)

type UpadateSubcriptionInfoDto struct {
	Subscription *models.SubscriptionType `db:"subscription" json:"subscription"`
	IsAutoRenew  *bool                    `db:"is_auto_renew" json:"is_auto_renew"`
	IsRenew      bool
}

type SubscriptionStatusDto struct {
	Status    int       `db:"subscription" json:"subscription"`
	ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
}
