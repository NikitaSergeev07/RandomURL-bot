package redis

import (
	"RandomURL/storage"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"time"
)

var ctx = context.Background()

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisStorage{
		client: rdb,
	}
}

func (r *RedisStorage) Save(p *storage.Page) error {
	hash, err := p.Hash()
	if err != nil {
		return err
	}

	log.Printf("Saving page with hash %s and URL %s for user %s", hash, p.URL, p.UserName)

	err = r.client.HSet(ctx, p.UserName, hash, p.URL).Err()
	if err != nil {
		return err
	}

	// Устанавливаем время жизни (TTL) на 1 день
	err = r.client.Expire(ctx, p.UserName, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	log.Printf("TTL set for user %s to 24 hours", p.UserName)

	return nil
}

func (r *RedisStorage) PickRandom(userName string) (*storage.Page, error) {
	urls, err := r.client.HGetAll(ctx, userName).Result()
	if err != nil || len(urls) == 0 {
		log.Printf("No pages found for user %s", userName)

		return nil, storage.ErrNoSavedPages
	}
	log.Printf("Found %d pages for user %s", len(urls), userName)

	rand.Seed(time.Now().UnixNano())
	keys := make([]string, 0, len(urls))
	for key := range urls {
		keys = append(keys, key)
	}
	randomKey := keys[rand.Intn(len(keys))]

	page := &storage.Page{
		URL:      urls[randomKey],
		UserName: userName,
	}

	return page, nil
}

func (r *RedisStorage) Remove(p *storage.Page) error {
	hash, err := p.Hash()
	if err != nil {
		return err
	}

	return r.client.HDel(ctx, p.UserName, hash).Err()
}

func (r *RedisStorage) IsExists(p *storage.Page) (bool, error) {
	hash, err := p.Hash()
	if err != nil {
		return false, err
	}

	exists, err := r.client.HExists(ctx, p.UserName, hash).Result()
	return exists, err
}
