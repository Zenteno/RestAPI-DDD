package main

import (
	"api-ddd/controller"
	_ "api-ddd/docs"
	"api-ddd/repository"
	"api-ddd/service"
	"flag"
	"log"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port      = flag.String("port", "3000", "Port to listen on")
	prod      = flag.Bool("prod", false, "Enable prefork in Production")
	mongoHost = os.Getenv("MONGO_HOST")
	mongoPort = 27017
	mongoDB   = "my-db"
)

// @title Test DDD Golang API
// @version 1.0
// @description This is a sample service for managing sessions.
// @termsOfService http://swagger.io/terms/
// @contact.name Alberto Zenteno
// @contact.email x.zenteno.a@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	var repo repository.SessionRepository = repository.NewMongoRepository(mongoHost, mongoPort, mongoDB, 10)
	var sessionService service.SessionService = service.NewSessionService(repo)
	var sessionController controller.SessionController = controller.NewSessionController(sessionService)

	flag.Parse()
	app := fiber.New(fiber.Config{
		Prefork: *prod,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	v1 := app.Group("/api/v1")

	v1.Get("/current_shopper_location/:shopper_uuid", sessionController.CurrentLocation)
	v1.Get("/session_location_history/:session_uuid", sessionController.HistorySession)
	v1.Post("/location", sessionController.AddLocation)

	app.Use("/docs", swagger.Handler)
	app.Use(sessionController.NotFound)
	log.Fatal(app.Listen(":" + *port))
}
