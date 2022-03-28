package repo

import (
	"github.com/girishg4t/app_invite_service/pkg/logging"
	"github.com/girishg4t/app_invite_service/pkg/model"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// IUser contains the user specific methods
//counterfeiter:generate . IUser
type IUser interface {
	GetUser(u *model.User) (model.User, error)
}

// UserProcessor require to initialize the user instance
type UserProcessor struct {
	Conn   *gorm.DB
	logger *zap.Logger
}

// NewUserProcessor retun's the user repo instance
func NewUserProcessor(conn *gorm.DB, logger *zap.Logger) *UserProcessor {
	log := logging.GetLogger().Named("user_repo")
	return &UserProcessor{
		conn,
		log,
	}
}

// GetUser get's the user by username
func (p *UserProcessor) GetUser(u *model.User) (user model.User, err error) {
	err = p.Conn.Where("username = ?", u.Username).First(&user).Error
	return user, err
}
