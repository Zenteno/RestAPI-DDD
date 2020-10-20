package repository

import (
	"api-ddd/entity"
	"time"

	. "github.com/ahmetb/go-linq/v3"
)

type repo struct {
	history []*entity.ShopperHistory
}

func NewMemRepository() SessionRepository {
	return &repo{
		history: []*entity.ShopperHistory{},
	}
}

func (r *repo) GetSessionHistory(sessionUUIDVal string) ([]*entity.ShopperHistory, error) {
	if len(r.history) == 0 {
		return nil, NotFoundErr
	}
	var result []*entity.ShopperHistory
	From(r.history).WhereT(func(i *entity.ShopperHistory) bool {
		return i.SessionUuid == sessionUUIDVal
	}).OrderByDescendingT(
		func(i *entity.ShopperHistory) int64 { return i.ReportedAt.Unix() },
	).ToSlice(&result)
	if len(result) == 0 {
		return nil, NotFoundErr
	}
	return result, nil
}

func (r *repo) GetShopperLocation(shoperUUIDVal string) (*entity.ShopperHistory, error) {
	if len(r.history) == 0 {
		return nil, NotFoundErr
	}
	now := time.Now()
	then := now.Add(-10 * time.Minute)
	result := From(r.history).WhereT(func(i *entity.ShopperHistory) bool {
		return i.ShopperUuid == shoperUUIDVal && i.ReportedAt.UTC().Unix() >= then.UTC().Unix()
	}).OrderByDescendingT(
		func(i *entity.ShopperHistory) int64 { return i.ReportedAt.Unix() },
	).First()
	if result == nil {
		return nil, NotFoundErr
	}
	return result.(*entity.ShopperHistory), nil
}

func (r *repo) Insert(session *entity.ShopperHistory) error {
	r.history = append(r.history, session)
	return nil
}
