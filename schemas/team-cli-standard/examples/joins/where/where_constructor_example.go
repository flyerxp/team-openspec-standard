package where

import (
	"github.com/flyerxp/lib/v2/middleware/gormL"
	"gorm.io/gorm"
)

// DemoJoinWhere 连表专属查询条件
// 【必须指针内嵌】原因：BaseWhere内部维护Wheres切片，值内嵌会发生切片拷贝，条件丢失
type DemoJoinWhere struct {
	*gormL.BaseWhere
	// 无冗余入参结构体字段，全部改为链式入参方法
}

// NewDemoJoinWhere 构造方法，统一初始化基类
func NewDemoJoinWhere() *DemoJoinWhere {
	return &DemoJoinWhere{
		BaseWhere: gormL.GetBaseWhere(),
	}
}

// TitleLike 【你提出的标准链式写法，替代原有入参结构体】
// 调用基类Where方法，返回自身实现连续链式调用
func (w *DemoJoinWhere) TitleLike(title string) *DemoJoinWhere {
	w.Where("`demo_info`.`title` LIKE ?", "%"+title+"%")
	return w
}

// MinPriceGte 关联表价格筛选链式方法
func (w *DemoJoinWhere) MinPriceGte(price float64) *DemoJoinWhere {
	w.Where("`demo_info_price`.`price_amount` >= ?", price)
	return w
}

// Build 仅做透传，对齐全局where规范，外部统一调用Build
func (w *DemoJoinWhere) Build(db *gorm.DB) *gorm.DB {
	if w == nil {
		return db
	}
	return w.BaseWhere.Build(db)
}
