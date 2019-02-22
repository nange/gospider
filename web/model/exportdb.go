package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/web/core"
	"github.com/pkg/errors"
)

//go:generate goqueryset -in exportdb.go
// gen:qs
type ExportDB struct {
	ID        uint64    `json:"id" gorm:"column:id;type:bigint unsigned AUTO_INCREMENT;primary_key"`
	ShowName  string    `json:"show_name" gorm:"column:show_name;type:varchar(64);not null;unique_index:uk_show_name"`
	Host      string    `json:"host" gorm:"column:host;type:varchar(128);not null"`
	Port      int       `json:"port" gorm:"column:port;type:int;not null"`
	User      string    `json:"user" gorm:"column:user;type:varchar(32);not null"`
	Password  string    `json:"password" gorm:"column:password;type:varchar(32);not null;default:''"`
	DBName    string    `json:"db_name" gorm:"column:db_name;type:varchar(64);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (o *ExportDB) TableName() string {
	return "gospider_exportdb"
}

func init() {
	core.Register(&ExportDB{})
}

func GetExportDBList(db *gorm.DB, size, offset int) ([]ExportDB, int, error) {
	queryset := NewExportDBQuerySet(db)
	count, err := queryset.Count()
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	queryset = NewExportDBQuerySet(db.Limit(size).Offset(offset))
	ret := make([]ExportDB, 0)
	if err := queryset.OrderDescByCreatedAt().All(&ret); err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return ret, count, nil
}
