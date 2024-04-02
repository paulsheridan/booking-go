package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/paulsheridan/booking-go/model"
)

type RedisDatabase struct {
	Client *redis.Client
}

func clientIDKey(id uuid.UUID) string {
	return fmt.Sprintf("client:%d", id)
}

func (r *RedisDatabase) Insert(ctx context.Context, client model.Client) error {
	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("failed to encode client: %w", err)
	}

	key := clientIDKey(client.ClientID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "clients", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add orders to set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisDatabase) FindByID(ctx context.Context, id uuid.UUID) (model.Client, error) {
	key := clientIDKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Client{}, ErrNotExist
	} else if err != nil {
		return model.Client{}, fmt.Errorf("failed to find client: %w", err)
	}

	var client model.Client
	err = json.Unmarshal([]byte(value), &client)
	if err != nil {
		return model.Client{}, fmt.Errorf("failed to decode client json: %w", err)
	}

	return client, nil
}

func (r *RedisDatabase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	key := clientIDKey(id)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to find client: %w", err)
	}

	if err := txn.SRem(ctx, "clients", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove clients set: %w", err)
	}
	return nil
}

func (r *RedisDatabase) Update(ctx context.Context, client model.Client) error {
	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("failed to encode client: %w", err)
	}

	key := clientIDKey(client.ClientID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("failed to set client: %w", err)
	}
	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Clients []model.Client
	Cursor  uint64
}

func (r *RedisDatabase) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "clients", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get client ids: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Clients: []model.Client{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get clients: %w", err)
	}

	clients := make([]model.Client, len(xs))

	for i, x := range xs {
		x := x.(string)
		var client model.Client

		err := json.Unmarshal([]byte(x), &client)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode clients json: %w", err)
		}

		clients[i] = client
	}

	return FindResult{
		Clients: clients,
		Cursor:  cursor,
	}, nil
}

var ErrNotExist = errors.New("order does not exist")
