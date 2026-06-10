package demo

import (
	"context"

	"github.com/flyerxp/lib/v2/logger"
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Demo1 手动事务管理（灵活控制提交/回滚时机）
func Demo1(ctx context.Context) error {
	ctx := logger.GetContext(ctx, "test")
	// 初始化事务
	tx, e := gormL.NewTx(gormL.GetDB(ctx))
	if e != nil {
		logger.AddError(ctx, zap.Error(e))
		return e
	}
	defer tx.Close()

	// 此处编写业务逻辑...

	// 提交事务
	e = tx.Commit()
	if e != nil {
		logger.AddError(ctx, zap.Error(e))
		return e
	}

	// 若失败则回滚
	// e = tx.Rollback()
	// if e != nil {
	// 	logger.AddError(ctx, zap.Error(e))
	// 	return e
	// }

	return nil
}

// Demo2 自动事务管理（GORM内置方法，自动提交/回滚）
func Demo2(ctx context.Context) error {
	ctx := logger.GetContext(ctx, "test")
	// 自动事务：返回nil自动提交，返回error自动回滚
	e := gormL.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		/**
		  此处编写业务逻辑...
		*/

		// 返回 nil 会自动提交
		return nil
	})
	if e != nil {
		logger.AddError(ctx, zap.Error(e))
		return e
	}
	return nil
}
