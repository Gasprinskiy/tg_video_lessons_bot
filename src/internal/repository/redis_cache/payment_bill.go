package redis_cache

import (
	"context"
	"encoding/json"
	"tg_video_lessons_bot/internal/entity/payment"
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

func (r *paymentBillCache) SetPaymentBill(ctx context.Context, id string, bill payment.Bill) error {
	byteData, err := json.Marshal(bill)
	if err != nil {
		return err
	}

	return r.db.Set(ctx, id, byteData, r.ttl).Err()
}
