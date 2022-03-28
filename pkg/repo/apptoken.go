package repo

import (
	"github.com/girishg4t/app_invite_service/pkg/logging"
	"github.com/girishg4t/app_invite_service/pkg/model"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// IAppToken methods to perform crud operation on apptoken
//counterfeiter:generate . IAppToken
type IAppToken interface {
	SaveAppToken(t *model.AppToken) error
	GetAppToken(at *model.AppToken) (model.AppToken, error)
	UpdateAppToken(at *model.AppToken) (*model.AppToken, error)
	GetAllAppToken() ([]model.AppToken, error)
}

// AppTokenProcessor to initialize the app token instance
type AppTokenProcessor struct {
	Conn   *gorm.DB
	logger *zap.Logger
}

// NewAppTokenProcessor send the new instance of app token repo
func NewAppTokenProcessor(conn *gorm.DB) *AppTokenProcessor {
	log := logging.GetLogger().Named("apptoken_repo")
	return &AppTokenProcessor{
		conn,
		log,
	}
}

// SaveAppToken save's the app token
func (p *AppTokenProcessor) SaveAppToken(t *model.AppToken) error {
	t.ID = model.UUID()
	return p.Conn.Save(&t).Error
}

// GetAppToken returns the apptoken object by token
func (p *AppTokenProcessor) GetAppToken(at *model.AppToken) (apptoken model.AppToken, err error) {
	err = p.Conn.Where("token = ?", at.Token).First(&apptoken).Error
	return apptoken, err
}

// UpdateAppToken update the apptoken object status
func (p *AppTokenProcessor) UpdateAppToken(at *model.AppToken) (user *model.AppToken, err error) {
	result, err := p.GetAppToken(at)
	if err != nil {
		return nil, err
	}
	at.ID = result.ID
	err = p.Conn.Save(&at).Error
	return &result, err
}

// GetAllAppToken retun's all the apptoken
func (p *AppTokenProcessor) GetAllAppToken() (allAppToken []model.AppToken, err error) {
	err = p.Conn.Find(&allAppToken).Error
	return allAppToken, err
}
