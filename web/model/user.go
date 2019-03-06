package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/web/core"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//go:generate goqueryset -in user.go
// gen:qs
type User struct {
	ID           uint64    `json:"id,string" gorm:"column:id;type:bigint unsigned AUTO_INCREMENT;primary_key"`
	UserName     string    `json:"user_name" gorm:"column:user_name;type:varchar(32) not null;unique_index:uk_uname"`
	Password     string    `json:"-" gorm:"column:password;type:varchar(128) not null"`
	Email        string    `json:"email" gorm:"column:email;type:varchar(32) not null default '';index:idx_email"`
	Roles        string    `json:"roles" gorm:"column:roles;type:varchar(128) not null default ''"`
	Introduction string    `json:"introduction" gorm:"column:introduction;type:varchar(128) not null default ''"`
	Avatar       string    `json:"avatar" gorm:"column:avatar;type:varchar(256) not null default ''"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:idx_updated_at"`
}

func (o *User) TableName() string {
	return "gospider_user"
}

func init() {
	core.Register(&User{})
}

func IsValidUser(db *gorm.DB, username, password string) (bool, *User, error) {
	user := &User{}
	query := NewUserQuerySet(db)
	if err := query.UserNameEq(username).One(user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, nil
		}
		log.Errorf("fetch user by name err [%+v]", err)
		return false, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, nil, nil
	}

	return true, user, nil
}

func GenUserHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func InitAdminUserIfNeeded(db *gorm.DB) error {
	user := &User{}
	query := NewUserQuerySet(db)
	err := query.UserNameEq("admin").One(user)
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		user.UserName = "admin"
		pw, err := GenUserHashPassword("admin")
		if err != nil {
			return err
		}
		user.Password = pw
		user.Avatar = "/admin/gopher.png"
		user.Introduction = "admin user"
		user.Roles = "admin"

		return user.Create(db)
	}

	return err
}
