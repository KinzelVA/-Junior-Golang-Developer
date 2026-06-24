package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/KinzelVA/-Junior-Golang-Developer/internal/model"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		pool: pool,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	query := `
INSERT INTO subscriptions (
service_name,
price,
user_id,
start_date,
end_date
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at
`

	return r.pool.QueryRow(
		ctx,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	).Scan(
		&subscription.ID,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*model.Subscription, error) {
	query := `
SELECT
id,
service_name,
price,
user_id,
start_date,
end_date,
created_at,
updated_at
FROM subscriptions
WHERE id = $1
`

	subscription, err := scanSubscription(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}

		return nil, err
	}

	return subscription, nil
}

func (r *SubscriptionRepository) List(ctx context.Context, request model.ListSubscriptionsRequest) ([]model.Subscription, error) {
	query := `
SELECT
id,
service_name,
price,
user_id,
start_date,
end_date,
created_at,
updated_at
FROM subscriptions
WHERE
($1::uuid IS NULL OR user_id = $1::uuid)
AND ($2::text IS NULL OR service_name = $2::text)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4
`

	var userID any
	if request.UserID != nil && *request.UserID != "" {
		userID = *request.UserID
	}

	var serviceName any
	if request.ServiceName != nil && *request.ServiceName != "" {
		serviceName = *request.ServiceName
	}

	rows, err := r.pool.Query(ctx, query, userID, serviceName, request.Limit, request.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]model.Subscription, 0)

	for rows.Next() {
		subscription, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, *subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	query := `
UPDATE subscriptions
SET
service_name = $2,
price = $3,
user_id = $4,
start_date = $5,
end_date = $6,
updated_at = NOW()
WHERE id = $1
RETURNING created_at, updated_at
`

	err := r.pool.QueryRow(
		ctx,
		query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	).Scan(
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSubscriptionNotFound
		}

		return err
	}

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	commandTag, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}

type subscriptionScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(scanner subscriptionScanner) (*model.Subscription, error) {
	var subscription model.Subscription
	var endDate sql.NullTime

	if err := scanner.Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&endDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if endDate.Valid {
		subscription.EndDate = &endDate.Time
	}

	return &subscription, nil
}
