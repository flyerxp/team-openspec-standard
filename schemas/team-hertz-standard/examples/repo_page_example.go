package shop

import (
	"context"

	"github.com/flyerxp/lib/v2/middleware/gormL"
	"github.com/flyerxp/manage/v2/golang/biz/dal/gormL/shop/where"
	"github.com/flyerxp/globalStruct/widget"
	"gorm.io/gorm"
)

// DemoInfo 示例表GORM模型
type DemoInfo struct {
	Id     int    `gorm:"column:id;primaryKey;autoIncrement"`
	Title  string `gorm:"column:title;size:255"`
	Status int    `gorm:"column:status"`
	Path   string `gorm:"column:path;size:255"`
	RootId int    `gorm:"column:root_id"`
}

// TableName 指定数据库表名
func (DemoInfo) TableName() string {
	return "demo_info"
}

// DemoRepo 数据仓储层标准示例
type DemoRepo struct{}

// NewDemoRepo 构造方法：每次新建实例，禁止全局单例
func NewDemoRepo() *DemoRepo {
	return &DemoRepo{}
}

// GetWhere 强制实现：Where对象唯一获取入口
func (r *DemoRepo) GetWhere() *where.DemoListWhere {
	return &where.DemoListWhere{BaseWhere: &gormL.BaseWhere{}}
}

// getGormModel 获取绑定上下文的DB实例，支持事务传递
func (r *DemoRepo) getGormModel(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		// 有事务：创建新会话避免污染外部事务
		return tx.Session(&gorm.Session{}).Model(&DemoInfo{})
	}
	// 无事务：使用默认连接
	return gormL.GetDB(ctx).Model(&DemoInfo{})
}

// DemoListPage 统一分页结构体标准示例
type DemoListPage struct {
	List []DemoInfo
	Page widget.Page
}

// DoPage 统一分页裁剪、自动计算是否有下一页
func (p *DemoListPage) DoPage() *DemoListPage {
	p.Page.HasMore = len(p.List) > p.Page.Size
	if p.Page.HasMore {
		p.List = p.List[:p.Page.Size]
	}
	return p
}

// ListPage 分页查询标准示例方法
func (r *DemoRepo) ListPage(ctx context.Context, w *where.DemoListWhere, page, limit int) (*DemoListPage, error) {
	var list []DemoInfo
	pageObj := DemoListPage{
		List: list,
		Page: widget.Page{Page: page, Size: limit},
	}

	// 基础查询
	db := r.getGormModel(ctx, nil)
	// 拼接where条件
	if w != nil {
		db = w.Build(db)
	}
	// 排序
	db = db.Order("id desc")
	// limit+1 查询，精准判断是否有下一页
	offset := (page - 1) * limit
	db = db.Offset(offset).Limit(limit + 1)

	if err := db.Find(&list).Error; err != nil {
		return pageObj.DoPage(), err
	}
	pageObj.List = list
	return pageObj.DoPage(), nil
}

// UpdatePathById 普通更新：仅返回DB错误，不校验更新行数
func (r *DemoRepo) UpdatePathById(ctx context.Context, id int, path string, rootId int, tx *gorm.DB) error {
	db := r.getGormModel(ctx, tx)
	return db.Where("id = ?", id).Updates(map[string]interface{}{
		"path":    path,
		"root_id": rootId,
	}).Error
}

// UpdatePathByIdMust 强制更新：额外校验更新行数，确保数据存在
func (r *DemoRepo) UpdatePathByIdMust(ctx context.Context, id int, path string, rootId int, tx *gorm.DB) error {
	db := r.getGormModel(ctx, tx)
	result := db.Where("id = ?", id).Updates(map[string]interface{}{
		"path":    path,
		"root_id": rootId,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Save 保存数据标准示例
func (r *DemoRepo) Save(ctx context.Context, info *DemoInfo, tx *gorm.DB) error {
	return r.getGormModel(ctx, tx).Save(info).Error
}
