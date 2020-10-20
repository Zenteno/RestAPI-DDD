package controller

import (
	"api-ddd/entity"
	"api-ddd/repository"
	"api-ddd/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SessionController interface {
	HistorySession(c *fiber.Ctx) error
	CurrentLocation(c *fiber.Ctx) error
	AddLocation(c *fiber.Ctx) error
	NotFound(c *fiber.Ctx) error
}

type sessionController struct{}

var sessionService service.SessionService

func NewSessionController(service service.SessionService) SessionController {
	sessionService = service
	return &sessionController{}
}

// GetHistorySessions godoc
// @Summary Get location history of a session
// @Description Get location history of a session
// @Tags History Session
// @Param session_uuid path string true "Session UUID"
// @Produce  json
// @Success 200 {array} entity.ShopperHistory
// @Router /api/v1/session_location_history/{session_uuid} [get]
func (*sessionController) HistorySession(c *fiber.Ctx) error {
	sessionUUID := c.Params("session_uuid")
	data, err := sessionService.GetSessionHistory(sessionUUID)
	if err != nil {
		if err == repository.NotFoundErr {
			c.SendStatus(404)
		} else {
			c.SendStatus(400)
		}
		return nil
	}
	return c.JSON(data)
}

// GetCurrentLocation godoc
// @Summary Get current position of a shopper
// @Description Get current position of a shopper
// @Tags Shopper Position
// @Param shopper_uuid path string true "Shopper UUID"
// @Produce  json
// @Success 200 {object} entity.ShopperHistory
// @Router /api/v1/current_shopper_location/{shopper_uuid} [get]
func (*sessionController) CurrentLocation(c *fiber.Ctx) error {
	shopperUUID := c.Params("shopper_uuid")
	data, err := sessionService.GetShopperLocation(shopperUUID)

	if err != nil {
		if err == repository.NotFoundErr {
			c.SendStatus(404)
		} else {
			c.SendStatus(400)
		}
		return nil

	}
	return c.JSON(data)
}

// AddLocation godoc
// @Summary Insert location data
// @Description Insert location data
// @Param account body entity.ShopperHistory true "Add Location"
// @Tags Insert Location
// @Accept json
// @Produce  json
// @Success 200
// @Router /api/v1/location [post]
func (*sessionController) AddLocation(c *fiber.Ctx) error {
	history := new(entity.ShopperHistory)
	if err := c.BodyParser(history); err != nil {
		fmt.Println(err)
		c.SendStatus(400)
		return nil
	}

	err := sessionService.Validate(history)
	if err != nil {
		fmt.Println(err)
		c.Status(400).JSON(fiber.Map{
			"message": "missing parameters",
		})
		return nil
	}
	err = sessionService.Insert(history)
	if err == nil {
		return c.SendStatus(201)
	}
	return c.SendStatus(400)
}

func (*sessionController) NotFound(c *fiber.Ctx) error {
	return c.SendStatus(404)
}
