package gormL

import (
	"context"

	"github.com/flyerxp/lib/v2/logger"
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"gorm.io/gorm"
)

// DBClient 多数据库客户端管理结构体
type DBClient struct {
	Shop *gorm.DB
}

var dbClient *DBClient

// Init 数据库全局初始化入口
func Init(ctx context.Context) error {
	if dbClient != nil {
		return nil
	}

	// 初始化Shop库
	shopDB, err := gormL.GetEngine(ctx, "shop")
	if err != nil {
		return err
	}

	dbClient = &DBClient{
		Shop: shopDB,
	}
	return nil
}

// GetShopDB 获取Shop库的上下文绑定DB实例，自动兜底初始化
func GetShopDB(ctx context.Context) *gorm.DB {
	if dbClient == nil {
		initCtx := logger.GetContext(context.Background(), "gormInit")
		if err := Init(initCtx); err != nil {
			panic(err)
		}
	}
	return dbClient.Shop.WithContext(ctx)
}
