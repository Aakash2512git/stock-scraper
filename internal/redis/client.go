package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/models"
	"main/pkg/config"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Initialize Redis client
var Rdb *redis.Client

func InitRedis(ctx context.Context) {
	cfg := config.GetConfig()

	// Parse redis://host:port from cfg.RedisURL
	redisAddr := strings.TrimPrefix(cfg.RedisURL, "redis://")

	fmt.Println("Connecting to Redis at:", redisAddr)

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // add if you set one
		DB:       0,
	})

	if err := Rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("❌ Could not connect to Redis:", err)
	} else {
		fmt.Println("✅ Connected to Redis successfully")
	}
}

// Save a single company
func SaveCompany(ctx context.Context, company models.CompanyInfo) error {
	data, err := json.Marshal(company)
	if err != nil {
		return err
	}

	key := "company:" + company.Company
	return Rdb.Set(ctx, key, data, 0).Err()
}

// Save multiple companies
func SaveCompanies(ctx context.Context, companies []models.CompanyInfo) error {
	for _, company := range companies {
		if err := SaveCompany(ctx, company); err != nil {
			return err
		}
	}
	return nil
}

// Get a single company
func GetCompanyInfo(ctx context.Context, name string) (models.CompanyInfo, error) {
	var info models.CompanyInfo
	key := "company:" + name

	val, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		return info, err
	}

	if err := json.Unmarshal([]byte(val), &info); err != nil {
		return info, err
	}

	return info, nil
}

// Get all companies
func GetAllCompanies(ctx context.Context) ([]models.CompanyInfo, error) {
	var companies []models.CompanyInfo

	keys, err := Rdb.Keys(ctx, "company:*").Result()
	if err != nil {
		return companies, err
	}

	for _, key := range keys {
		val, err := Rdb.Get(ctx, key).Result()
		if err != nil {
			return companies, err
		}

		var company models.CompanyInfo
		if err := json.Unmarshal([]byte(val), &company); err != nil {
			return companies, err
		}

		companies = append(companies, company)
	}

	return companies, nil
}
