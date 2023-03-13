package services

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"greenbone-task/logger"
	db "greenbone-task/models/db"
	"log"
	"sync"
	"time"
)

var DbConnection *gorm.DB

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		Config.DBUserName, Config.DBUserPassword, Config.DBHost, Config.DBPort, Config.DBName)

	DbConnection, err = gorm.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to the Database", zap.Error(err))
	}

	DbConnection.AutoMigrate(&db.Computer{})
	DbConnection.AutoMigrate(&db.Employee{})
	DbConnection.AutoMigrate(&db.Token{})
}

var redisDefaultClient *redis.Client
var redisDefaultOnce sync.Once

var redisCache *cache.Cache
var redisCacheOnce sync.Once

func GetRedisDefaultClient() *redis.Client {
	redisDefaultOnce.Do(func() {
		redisDefaultClient = redis.NewClient(&redis.Options{
			Addr:     Config.RedisDefaultAddr,
			Password: Config.RedisPassword,
		})
	})

	return redisDefaultClient
}

func GetRedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		redisCache = cache.New(&cache.Options{
			Redis:      GetRedisDefaultClient(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	})

	return redisCache
}

func CheckRedisConnection() {
	redisClient := GetRedisDefaultClient()
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Failed to connected redis:", err.Error())
	}

	logger.Info("Connected to Redis!")
}
