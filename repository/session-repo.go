package repository

import (
	"api-ddd/entity"
	"errors"
)

type SessionRepository interface {
	Insert(session *entity.ShopperHistory) error
	GetSessionHistory(sessionUUID string) ([]*entity.ShopperHistory, error)
	GetShopperLocation(shopperUUID string) (*entity.ShopperHistory, error)
}

var (
	NotFoundErr = errors.New("Not Found")
)
