package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"log"
	
	"github.com/go-redis/redis/v8"
	"github.com/ajitirto/ms/product-service/internal/domain"
	"github.com/ajitirto/ms/product-service/internal/repository"
)

// Pastikan ProductCacheRedis mengimplementasikan repository.ProductCache
var _ repository.ProductCache = (*ProductCacheRedis)(nil)

// ProductCacheRedis mengimplementasikan antarmuka repository.ProductCache menggunakan Redis.
type ProductCacheRedis struct {
	Client *redis.Client
	TTL    time.Duration // Time to Live untuk entri cache
}

// NewProductCacheRedis membuat implementasi cache Redis yang baru.
func NewProductCacheRedis(client *redis.Client, ttl time.Duration) *ProductCacheRedis {
	return &ProductCacheRedis{
		Client: client,
		TTL:    ttl,
	}
}

// key menghasilkan kunci Redis untuk ID produk
func (r *ProductCacheRedis) key(id int) string {
	return fmt.Sprintf("product:%d", id)
}

// Get mengambil produk dari Redis.
func (r *ProductCacheRedis) Get(ctx context.Context, id int) (*domain.Product, error) {
	val, err := r.Client.Get(ctx, r.key(id)).Result()
	if err == redis.Nil {
		// Cache Miss, kembalikan nil, nil
		return nil, nil 
	}
	if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var product domain.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product from cache: %w", err)
	}

	return &product, nil
}

// Set menyimpan produk di Redis dengan TTL yang dikonfigurasi.
func (r *ProductCacheRedis) Set(ctx context.Context, product *domain.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product for cache: %w", err)
	}

	log.Printf("Attempting to SET key: %s with TTL: %s", r.key(product.ProductID), r.TTL.String())

	// Simpan data dengan TTL
	if err := r.Client.Set(ctx, r.key(product.ProductID), data, r.TTL).Err(); err != nil {
		log.Printf("Redis SET failed for key %s: %v", r.key(product.ProductID), err)
		return fmt.Errorf("redis set error: %w", err)
	}

	log.Printf("Successfully SET key: %s", r.key(product.ProductID))

	return nil
}

// Delete menghapus produk dari Redis (invalidasi cache).
func (r *ProductCacheRedis) Delete(ctx context.Context, id int) error {
	if err := r.Client.Del(ctx, r.key(id)).Err(); err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}

func (r *ProductCacheRedis) GetAll(ctx context.Context, key string) ([]domain.Product, error) {
	val, err := r.Client.Get(ctx, key).Result()
	
	if err == redis.Nil {
		return nil, nil // Cache Miss
	}
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	err = json.Unmarshal([]byte(val), &products) // Unmarshal sebagai slice
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ALL products JSON from Redis: %w", err)
	}

	return products, nil // Cache Hit
}

func (r *ProductCacheRedis) SetAll(ctx context.Context, key string, products []domain.Product) error {
	data, err := json.Marshal(products) // Marshal slice
	if err != nil {
		return fmt.Errorf("failed to marshal ALL products to JSON: %w", err)
	}

	// Simpan data dengan TTL
	return r.Client.Set(ctx, key, data, r.TTL).Err()
}