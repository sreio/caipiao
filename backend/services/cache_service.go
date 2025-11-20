package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// CacheService Redis缓存服务
type CacheService struct {
	client *redis.Client
}

// NewCacheService 创建缓存服务实例
func NewCacheService() *CacheService {
	// 从环境变量获取Redis地址，默认为localhost:6379
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // 默认无密码
		DB:       0,  // 使用默认DB
	})

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("警告: Redis连接失败，将不使用缓存: %v\n", err)
		return &CacheService{client: nil}
	}

	fmt.Println("Redis缓存服务初始化成功")
	return &CacheService{client: client}
}

// Set 设置缓存（带过期时间）
func (s *CacheService) Set(key string, value interface{}, expiration time.Duration) error {
	if s.client == nil {
		return fmt.Errorf("Redis未连接")
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, jsonData, expiration).Err()
}

// Get 获取缓存
func (s *CacheService) Get(key string, dest interface{}) error {
	if s.client == nil {
		return fmt.Errorf("Redis未连接")
	}

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete 删除缓存
func (s *CacheService) Delete(key string) error {
	if s.client == nil {
		return nil
	}
	return s.client.Del(ctx, key).Err()
}

// DeletePattern 删除匹配模式的所有缓存
func (s *CacheService) DeletePattern(pattern string) error {
	if s.client == nil {
		return nil
	}

	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists 检查缓存是否存在
func (s *CacheService) Exists(key string) bool {
	if s.client == nil {
		return false
	}

	count, err := s.client.Exists(ctx, key).Result()
	return err == nil && count > 0
}
