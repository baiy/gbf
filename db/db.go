package db

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/baiy/gbf/log"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type GormCreate func(db *sql.DB) (gorm.Dialector, *Option)
type Db struct {
	db     *sql.DB
	gorm   *gorm.DB
	logger log.Logger
	mu     sync.Mutex

	gormCreate GormCreate
}

// Gorm gorm.DB
func (d *Db) Gorm() *gorm.DB {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.gorm != nil {
		return d.gorm
	}
	dialector, option := d.gormCreate(d.db)
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   option.tablePrefix,
			SingularTable: true,
		},
		Logger: &gormLogger{logger: d.logger},
	})
	if err != nil {
		d.logger.Errorf("get gorm error:%s", err)
		panic(err)
	}
	d.gorm = gormDB
	return d.gorm
}

// Model 获取模型操作对象
func (d *Db) Model[T any](ctxs ...context.Context) *Model[T] {
	ctx := context.Background()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	return newModel(ctx, gorm.G[T](d.Gorm()).Where(""), d.logger, d)
}

// ManyToMany 获取模型操作对象
func (d *Db) ManyToMany[R any, E any, T any](ctxs ...context.Context) *ManyToMany[R, E, T] {
	ctx := context.Background()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	return &ManyToMany[R, E, T]{
		db:     d,
		ctx:    ctx,
		logger: d.logger,
	}
}

// Transaction 事务
func (d *Db) Transaction(fc func(db *Db) error) error {
	return d.Gorm().Transaction(func(tx *gorm.DB) error {
		return fc(&Db{db: d.db, logger: d.logger, gorm: tx})
	})
}

// SqlDb 底层 sql.DB
func (d *Db) SqlDb() *sql.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *Db) Close() error {
	return d.db.Close()
}

// CurrentTime 获取数据库当前时间
func (d *Db) CurrentTime() time.Time {
	u, err := gorm.G[int64](d.Gorm()).Raw("SELECT REPLACE(UNIX_TIMESTAMP(NOW(3)), '.', '')").Find(context.Background())
	if err != nil {
		d.logger.Errorf("get current time error:%s", err)
		return time.Now()
	}
	return time.UnixMilli(u[0])
}

func NewDb(db *sql.DB, logger log.Logger, gormCreate GormCreate) *Db {
	return &Db{
		db:         db,
		logger:     logger,
		gormCreate: gormCreate,
	}
}
