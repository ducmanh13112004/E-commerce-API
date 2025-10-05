package utils_call

import (
	"ecom_promotion_v2/internal"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func GetDataRedis(key string, funcName string) string {
	key = internal.ServiceName + ":" + key
	_, err := internal.RedisCache.Ping().Result()
	if err != nil {
		internal.Log.Error("RedisCache.Ping", zap.Any("funcName", funcName), zap.Error(err))
		return ""
	}
	value, err := internal.RedisCache.Get(key).Result()
	if err != nil {
		internal.Log.Error("RedisCache.Get", zap.Any("funcName", funcName), zap.Any("key", key), zap.Error(err))
		return ""
	}
	return value
}

func SaveCacheRedis(keyCache string, value string, timeCache time.Duration, funcName string) {
	keyCache = internal.ServiceName + ":" + keyCache
	err := internal.RedisCache.Set(keyCache, value, timeCache).Err()
	if err != nil {
		internal.Log.Error("SaveCacheRedis FAIL", zap.Any("funcName", funcName), zap.Any("key", keyCache), zap.Any("value", value), zap.Error(err))
	}
}

func RedisKeyHash(redisKey, input, funcName string) string {
	redisKey = internal.ServiceName + ":" + redisKey
	if _, err := internal.RedisCache.Ping().Result(); err != nil {
		internal.Log.Error("RedisCache.Ping failed", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.Error(err))
		return ""
	}
	value, err := internal.RedisCache.HGet(redisKey, input).Result()
	if err != nil {
		if err == redis.Nil {
			// Trường hợp không có field hoặc key -> không trả error, đọc bool để biết có key hay không
			internal.Log.Info("Redis field not found", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.String("field", input))
			return ""
		} else {
			internal.Log.Error("RedisCache.Exists failed", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.Error(err))
			return ""
		}
	}
	return value
}

func RedisKeyHashSet(redisKey, field, value string, activeTTL bool, timeCache time.Duration, funcName string) error {
	redisKey = internal.ServiceName + ":" + redisKey
	// Ping Redis
	if _, err := internal.RedisCache.Ping().Result(); err != nil {
		internal.Log.Error("RedisCache.Ping failed", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.Error(err))
		return err
	}
	// Set field vào hash
	err := internal.RedisCache.HSet(redisKey, field, value).Err()
	if err != nil {
		internal.Log.Error("RedisCache.HSet failed", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.String("field", field), zap.Error(err))
		return err
	}
	// Set TTL
	if activeTTL && timeCache != 0 {
		// ttl := utils.GetTTLUntilMidnight(time.Second * timeCache)
		err = internal.RedisCache.Expire(redisKey, timeCache).Err()
		if err != nil {
			internal.Log.Error("RedisCache.Expire failed, deleting hash", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.Any("activeTTL", activeTTL), zap.Duration("ttl", timeCache), zap.Error(err))
			return fmt.Errorf("failed to set TTL, cache cleared")
		}
		internal.Log.Info("RedisCache.HSet success", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.String("field", field), zap.Any("activeTTL", activeTTL), zap.String("value", value), zap.Duration("ttl", timeCache))
	}
	return nil
}

func RedisDeleteHashField(redisKey, input, funcName string) (int64, error) {
	redisKey = internal.ServiceName + ":" + redisKey
	// Ping Redis
	if _, err := internal.RedisCache.Ping().Result(); err != nil {
		internal.Log.Error("RedisCache.Ping failed", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.Error(err))
		return 0, err
	}
	// Xóa field khỏi hash
	deleted, err := internal.RedisCache.HDel(redisKey, input).Result()
	if err != nil {
		internal.Log.Error("RedisCache.HDel failed", zap.String("funcName", funcName), zap.String("redisKey", redisKey), zap.String("field", input), zap.Error(err))
		return 0, err
	}
	return deleted, err
}
