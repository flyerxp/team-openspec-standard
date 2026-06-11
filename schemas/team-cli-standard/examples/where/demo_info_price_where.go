package where

import "github.com/flyerxp/lib/v2/middleware/gormL"

// DemoInfoPriceWhere 价格表查询构造器标准示例
type DemoInfoPriceWhere struct {
	*gormL.BaseWhere
}

// StrategyId 按策略ID精确筛选
func (w *DemoInfoPriceWhere) StrategyId(id int) *DemoInfoPriceWhere {
	w.Where("strategy_id = ?", id)
	return w
}
