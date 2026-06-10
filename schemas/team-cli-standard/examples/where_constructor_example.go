package where

import "github.com/flyerxp/lib/v2/middleware/gormL"

// DemoListWhere Where查询构造器标准示例
type DemoListWhere struct {
	*gormL.BaseWhere
}

func (w *DemoListWhere) TitleLike(title string) *DemoListWhere {
	w.Where("title LIKE ?", title+"%")
	return w
}
