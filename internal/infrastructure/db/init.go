package db

import (
	"context"
	"fmt"
	"time"

	"course-platform/internal/configs"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMySQL 初始化 MySQL 資料庫連接
func InitMySQL(config configs.MySQLConfig) (*gorm.DB, error) {
	// 構建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 設定日誌級別
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 嘗試連接資料庫
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("連接 MySQL 失敗: %v", err)
	}

	// 獲取底層 sql.DB 來設定連接池參數
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("獲取 SQL DB 實例失敗: %v", err)
	}

	// 設定連接池參數
	sqlDB.SetMaxIdleConns(10)           // 設定最大空閒連接數
	sqlDB.SetMaxOpenConns(100)          // 設定最大打開連接數
	sqlDB.SetConnMaxLifetime(time.Hour) // 設定連接的最大生命週期

	// 測試連接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("資料庫連接測試失敗: %v", err)
	}

	return db, nil
}

// InitRedis 初始化 Redis 連接
func InitRedis(config configs.RedisConfig) (*redis.Client, error) {
	fmt.Printf("🔍 嘗試連接 Redis: %s\n", config.Addr)

	// 建立 Redis 客戶端，增加超時設置
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  15 * time.Second, // 連接超時
		ReadTimeout:  5 * time.Second,  // 讀取超時
		WriteTimeout: 5 * time.Second,  // 寫入超時
	})

	// 測試連接，設置10秒超時
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("🔍 發送 PING 測試...\n")
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("連接 Redis 失敗: %v", err)
	}

	fmt.Printf("✅ Redis 連接成功\n")
	return rdb, nil
}
