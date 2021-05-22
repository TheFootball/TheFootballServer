package history

import (
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	service *service
}

func (c *controller) getAllHistory(ctx *fiber.Ctx) error {
	res := c.service.getAllHistory()

	return ctx.JSON(fiber.Map{"message": res})
}

func (c* controller) createGameHistory(ctx *fiber.Ctx) error {
	dto := new(createGameHistoryDTO)
	err := ctx.BodyParser(dto)
	if err != nil {
		panic(err)
	}

	err = c.service.createGameHistory(dto)
	if err != nil {
		panic(err)
	}

	return ctx.JSON(fiber.Map{"message": "ok"})
}

func initController(r fiber.Router, ser *service) {
	c := &controller{ser}
	r.Get("", c.getAllHistory)
	r.Post("", c.createGameHistory)
}