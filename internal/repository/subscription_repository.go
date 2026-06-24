package repository

import (
"context"

"github.com/jackc/pgx/v5/pgxpool"

"github.com/KinzelVA/-Junior-Golang-Developer/internal/model"
)

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
