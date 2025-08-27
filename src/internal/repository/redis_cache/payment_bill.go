package redis_cache

import (
	"context"
	"tg_video_lessons_bot/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type paymentBillCache struct {
	db  *redis.Client
	ttl time.Duration
}

func NewPaymentBillCache(db *redis.Client, ttl time.Duration) repository.PaymentBillCache {
	return &paymentBillCache{db, ttl}
}

func (r *paymentBillCache) SetPaymentBill(ctx context.Context, ID string) error {
	return r.db.Set(ctx, ID, true, r.ttl).Err()
}
