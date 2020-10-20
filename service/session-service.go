package service

import (
	"api-ddd/entity"
	"api-ddd/repository"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

type SessionService interface {
	Validate(*entity.ShopperHistory) error
	Insert(session *entity.ShopperHistory) error
	GetSessionHistory(sessionUUID string) ([]*entity.ShopperHistory, error)
	GetShopperLocation(shopperUUID string) (*entity.ShopperHistory, error)
}

type service struct{}

var repo repository.SessionRepository

func NewSessionService(repository repository.SessionRepository) SessionService {
	repo = repository
	return &service{}
}

func (s *service) Validate(shopper *entity.ShopperHistory) error {
	validate = validator.New()
	return validate.Struct(shopper)
}
func (s *service) Insert(session *entity.ShopperHistory) error {
	return repo.Insert(session)
}
func (s *service) GetSessionHistory(sessionUUID string) ([]*entity.ShopperHistory, error) {
	return repo.GetSessionHistory(sessionUUID)
}
func (s *service) GetShopperLocation(shopperUUID string) (*entity.ShopperHistory, error) {
	return repo.GetShopperLocation(shopperUUID)
}
