package where

import "github.com/flyerxp/lib/v2/middleware/gormL"

// DemoListWhere Where查询构造器标准示例
type DemoListWhere struct {
	*gormL.BaseWhere
}

// GetDemoListWhere Where对象唯一获取入口
func GetDemoListWhere() *DemoListWhere {
	return &DemoListWhere{BaseWhere: &gormL.BaseWhere{}}
}

// TitleLike 标题模糊查询条件
func (w *DemoListWhere) TitleLike(title string) *DemoListWhere {
	w.Where("title LIKE ?", title+"%")
	return w
}

// StatusEq 状态精准匹配条件
func (w *DemoListWhere) StatusEq(status int) *DemoListWhere {
	w.Where("status = ?", status)
	return w
}
