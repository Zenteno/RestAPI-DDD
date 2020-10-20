package entity

import "time"

type ShopperHistory struct {
	ShopperUuid string    `json:"shopper_uuid" form:"shopper_uuid" validate:"required" bson:"shopper_uuid" faker:"{uuid}"`
	SessionUuid string    `json:"session_uuid" form:"session_uuid" validate:"required" bson:"session_uuid" faker:"{uuid}"`
	Lat         float64   `json:"lat" form:"lat" validate:"required" bson:"lat" faker:"{latitude}"`
	Lng         float64   `json:"lng" form:"lng" validate:"required" bson:"lng" faker:"{longitude}"`
	Precision   float64   `json:"precision" form:"precision" validate:"required" bson:"precision" faker:"{floatrange:1,30}"`
	ReportedAt  time.Time `json:"reported_at" form:"reported_at" validate:"required" bson:"reported_at" faker:"skip"`
	InsertedAt  time.Time `json:"-"  bson:"inserted_at" faker:"skip"`
}
