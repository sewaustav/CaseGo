package models

import (
	"time"
)

type SubscriptionType int

const (
	NoSubscription SubscriptionType = iota
	Basic
	Advanced
)

type SubscriptionInfo struct {
	ID               int64            `db:"id" json:"id"`
	UserID           int64            `db:"user_id" json:"user_id"`
	Subscription     SubscriptionType `db:"subscription" json:"subscription"`
	FirstPaymentDate time.Time        `db:"first_payment_date" json:"first_payment_date"`
	CountOfRenewal   int              `db:"count_of_renewal" json:"count_of_renewal"`
	IsAutoRenew      bool             `db:"is_auto_renew" json:"is_auto_renew"`
	LastPaymentDate  time.Time        `db:"last_payment_date" json:"last_payment_date"`
	CanseledAt       time.Time        `db:"canceled_at" json:"canceled_at"`
	ExpiredAt        time.Time        `db:"expired_at" json:"expired_at"`
}

type PaymentInfo struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"user_id"`
	SubscriptionID *int64    `db:"subscription_id" json:"subscription_id,omitempty"`
	TransactionID  *string   `db:"transaction_id" json:"transaction_id,omitempty"`
	Price          int64     `db:"price" json:"price"`
	Currency       string    `db:"currency" json:"currency"`
	Status         string    `db:"status" json:"status"`
	PaymentDate    time.Time `db:"payment_date" json:"payment_date"`
}

type UserRole int

const (
	Admin UserRole = iota
	User
	Creator
)

type UserIdentity struct {
	UserID int64
	Role   *UserRole
}