package main

import (
	"github.com/SohailHaqyar/friday/data"
	"github.com/SohailHaqyar/friday/user"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"xorm.io/xorm"
)

func setupRoutes(app *fiber.App, engine *xorm.Engine) {
	user.SetupRoutes(app, engine)
}

func main() {
	app := fiber.New()
	engine := data.SetupDatabase()
	setupRoutes(app, engine)

	private := app.Group("/api")
	private.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		},
	}))

	private.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "You can only learn to love yourself from other people"})
	})

	private.Get("/users", func(c *fiber.Ctx) error {
		users := new([]data.User)
		if err := engine.Find(users); err != nil {
			return err
		}
		return c.JSON(users)
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

}
