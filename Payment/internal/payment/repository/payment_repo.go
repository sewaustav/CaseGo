package repository

import (
	"context"
	"database/sql"

	"github.com/sewaustav/Payment/internal/payment/dto"
	"github.com/sewaustav/Payment/internal/payment/models"
)

type PaymentRepo interface {
	Begin(ctx context.Context) (Tx, error)
	WithTx(tx Tx) PaymentRepo

	InitSubscription(ctx context.Context, sub *models.SubscriptionInfo) (*models.SubscriptionInfo, error)
	CreatePayment(ctx context.Context, payment *models.PaymentInfo) (*models.PaymentInfo, error)

	GetSubscriptionStatus(ctx context.Context, userID int64) (dto.SubscriptionStatusDto, error)
	GetUserSubscriptionInfo(ctx context.Context, userID int64) (*models.SubscriptionInfo, error)
	GetUserPayments(ctx context.Context, userID int64, limit, offset int) ([]models.PaymentInfo, error)
	GetPaymentByID(ctx context.Context, id int64) (*models.PaymentInfo, error)
	GetPaymentByTransactionID(ctx context.Context, id string) (*models.PaymentInfo, error)

	UpdateSubscription(ctx context.Context, id int64, sub *dto.UpadateSubcriptionInfoDto) error

	DeleteUser(ctx context.Context, userID int64) error
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row

	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

type Tx interface {
	DBTX
	Commit() error
	Rollback() error
}