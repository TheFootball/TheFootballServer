package history

import (
	"github.com/gofiber/fiber/v2"
	"onair/src/database"
)

func InitModule(router fiber.Router, DB *database.DB) {
	initController(router, newService(DB))
}