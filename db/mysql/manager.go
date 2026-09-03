package mysql

import (
	"cmp"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/baiy/gbf/closer"
	"github.com/baiy/gbf/db"
	"github.com/baiy/gbf/manager"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func create(mo manager.Option[Option]) (*db.Db, error) {
	option := mo.Value

	tls := "skip-verify"
	if option.Tls != "" {
		tls = option.Tls
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s&time_zone=%s",
		option.Username,
		option.Password,
		option.Host,
		cmp.Or(option.Port, "3306"),
		option.Database,
		tls,
		url.QueryEscape("'"+time.Local.String()+"'"),
	)

	sqlDB, err := sql.Open("mysql", dsn)
	logStr := fmt.Sprintf("database connect %s@%s:%s/%s", mo.Value.Username, mo.Value.Host, mo.Value.Port, mo.Value.Database)
	if err != nil {
		mo.Logger.Errorf("%s err:%s", logStr, err)
		return nil, err
	}
	mo.Logger.Info(logStr)
	sqlDB.SetMaxOpenConns(cmp.Or(int(option.MaxOpenConns), 200))
	sqlDB.SetMaxIdleConns(cmp.Or(int(option.MaxIdleConns), 5))
	sqlDB.SetConnMaxLifetime(cmp.Or(time.Duration(option.ConnMaxLifetime)*time.Second, time.Hour))
	sqlDB.SetConnMaxIdleTime(cmp.Or(time.Duration(option.ConnMaxIdleTime)*time.Second, time.Minute*10))
	DB := db.NewDb(sqlDB, mo.Logger, func(d *sql.DB) (gorm.Dialector, *db.Option) {
		return mysql.New(mysql.Config{Conn: d}), db.NewOption(db.WithTablePrefix(option.TablePrefix))
	})
	closer.Register(fmt.Sprintf("mysql:%s", mo.Key), DB)
	return DB, nil
}

var Manager = manager.New[*db.Db, Option](
	"mysql",
	create,
	nil,
	func() Option {
		return Option{
			Host:     "127.0.0.1",
			Username: "root",
			Password: "root",
			Database: "app",
		}
	},
)
