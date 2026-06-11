package webExamples

import (
	"context"
	"webExamples/examples/where"

	"github.com/flyerxp/lib/v2/middleware/gormL"
	"gorm.io/gorm"
)

// DemoInfoPrice 关联表源模型
type DemoInfoPrice struct {
	Id          int     `gorm:"column:id;primaryKey"`
	StrategyId  int     `gorm:"column:strategy_id;comment:关联主表root_id"`
	PriceAmount float64 `gorm:"column:price_amount;comment:对应价格"`
}

func (DemoInfoPrice) TableName() string {
	return "demo_info_price"
}

// DemoInfoPriceRepo 数据仓储层标准示例
type DemoInfoPriceRepo struct{}

// NewDemoInfoPriceRepo 构造方法：每次新建实例，禁止全局单例
func NewDemoInfoPriceRepo() *DemoInfoPriceRepo {
	return &DemoInfoPriceRepo{}
}

// GetWhere 强制实现：Where对象唯一获取入口
func (r *DemoInfoPriceRepo) GetWhere() *where.DemoInfoPriceWhere {
	return &where.DemoInfoPriceWhere{BaseWhere: &gormL.BaseWhere{}}
}

// getGormModel 获取绑定上下文的DB实例，支持事务传递
func (r *DemoInfoPriceRepo) getGormModel(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.Session(&gorm.Session{}).Model(&DemoInfoPrice{})
	}
	return GetShopDB(ctx).Model(&DemoInfoPrice{})
}
