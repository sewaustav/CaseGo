package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/sewaustav/Payment/internal/payment/dto"
	"github.com/sewaustav/Payment/internal/payment/models"
)

type PostgresPaymentRepo struct {
	db DBTX
	tx Tx
}

func NewPostgresPaymentRepo(db *sql.DB) *PostgresPaymentRepo {
	return &PostgresPaymentRepo{
		db: db,
	}
}

func (r *PostgresPaymentRepo) WithTx(tx Tx) PaymentRepo {
	return &PostgresPaymentRepo{
		db: r.db,
		tx: tx,
	}
}

func (r *PostgresPaymentRepo) Begin(ctx context.Context) (Tx, error) {
	if db, ok := r.db.(*sql.DB); ok {
		return db.BeginTx(ctx, nil)
	}
	return nil, fmt.Errorf("could not begin transaction: db is not *sql.DB")
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (r *PostgresPaymentRepo) InitSubscription(ctx context.Context, sub *models.SubscriptionInfo) (*models.SubscriptionInfo, error) {
	query := psql.Insert("subscription_info").
		Columns("user_id", "subscription", "first_payment_date", "count_of_renewal", "is_auto_renew", "last_payment_date", "expired_at").
		Values(sub.UserID, sub.Subscription, sub.FirstPaymentDate, sub.CountOfRenewal, sub.IsAutoRenew, sub.LastPaymentDate, sub.ExpiredAt).
		Suffix("RETURNING id")

	err := query.RunWith(r.db).QueryRowContext(ctx).Scan(&sub.ID)
	return sub, err
}

func (r *PostgresPaymentRepo) CreatePayment(ctx context.Context, p *models.PaymentInfo) (*models.PaymentInfo, error) {
	query := psql.Insert("payment_info").
		Columns("user_id", "subscription_id", "transaction_id", "price", "currency", "status", "payment_date").
		Values(p.UserID, p.SubscriptionID, p.TransactionID, p.Price, p.Currency, p.Status, p.PaymentDate).
		Suffix("RETURNING id")

	err := query.RunWith(r.db).QueryRowContext(ctx).Scan(&p.ID)
	return p, err
}

func (r *PostgresPaymentRepo) GetSubscriptionStatus(ctx context.Context, userID int64) (dto.SubscriptionStatusDto, error) {
	var res dto.SubscriptionStatusDto
	query, args, _ := psql.Select("subscription", "expired_at").
		From("subscription_info").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&res.Status, &res.ExpiredAt)
	return res, err
}

func (r *PostgresPaymentRepo) GetUserSubscriptionInfo(ctx context.Context, userID int64) (*models.SubscriptionInfo, error) {
	var s models.SubscriptionInfo
	query, args, _ := psql.Select("id", "user_id", "subscription", "first_payment_date", "count_of_renewal", "is_auto_renew", "last_payment_date", "expired_at").
		From("subscription_info").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&s.ID, &s.UserID, &s.Subscription, &s.FirstPaymentDate,
		&s.CountOfRenewal, &s.IsAutoRenew, &s.LastPaymentDate, &s.ExpiredAt,
	)
	return &s, err
}

func (r *PostgresPaymentRepo) GetUserPayments(ctx context.Context, userID int64, limit, offset int) ([]models.PaymentInfo, error) {
	query, args, _ := psql.Select("id", "user_id", "subscription_id", "transaction_id", "price", "currency", "status", "payment_date").
		From("payment_info").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("payment_date DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PaymentInfo
	for rows.Next() {
		var p models.PaymentInfo
		if err := rows.Scan(&p.ID, &p.UserID, &p.SubscriptionID, &p.TransactionID, &p.Price, &p.Currency, &p.Status, &p.PaymentDate); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PostgresPaymentRepo) UpdateSubscription(ctx context.Context, userID int64, sub *dto.UpadateSubcriptionInfoDto) error {
	builder := psql.Update("subscription_info").Where(sq.Eq{"userID": userID})

	if sub.Subscription != nil {
		builder = builder.Set("subscription", *sub.Subscription)
	}
	if sub.IsAutoRenew != nil {
		builder = builder.Set("is_auto_renew", *sub.IsAutoRenew)
	}
	if sub.IsRenew {
		builder = builder.Set("count_of_renewal", sq.Expr("count_of_renewal + 1")).
			Set("last_payment_date", time.Now())
	}

	_, err := builder.RunWith(r.db).ExecContext(ctx)
	return err
}

func (r *PostgresPaymentRepo) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*models.PaymentInfo, error) {
	var p models.PaymentInfo
	query, args, _ := psql.Select("id", "user_id", "subscription_id", "transaction_id", "price", "currency", "status", "payment_date").
		From("payment_info").
		Where(sq.Eq{"transaction_id": transactionID}).
		ToSql()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&p.ID, &p.UserID, &p.SubscriptionID, &p.TransactionID, &p.Price, &p.Currency, &p.Status, &p.PaymentDate,
	)
	return &p, err
}

func (r *PostgresPaymentRepo) GetPaymentByID(ctx context.Context, id int64) (*models.PaymentInfo, error) {
	var p models.PaymentInfo
	query, args, _ := psql.Select("id", "user_id", "subscription_id", "transaction_id", "price", "currency", "status", "payment_date").
		From("payment_info").
		Where(sq.Eq{"id": id}).
		ToSql()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&p.ID, &p.UserID, &p.SubscriptionID, &p.TransactionID, &p.Price, &p.Currency, &p.Status, &p.PaymentDate,
	)
	return &p, err
}

func (r *PostgresPaymentRepo) DeleteUser(ctx context.Context, userID int64) error {
	_, err := psql.Delete("payment_info").Where(sq.Eq{"user_id": userID}).RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return err
	}
	_, err = psql.Delete("subscription_info").Where(sq.Eq{"user_id": userID}).RunWith(r.db).ExecContext(ctx)
	return err
}