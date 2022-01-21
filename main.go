package main

import (
	"github.com/SohailHaqyar/friday/data"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	app := fiber.New()
	engine := data.SetupDatabase()

	private := app.Group("/private")
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

	public := app.Group("/public")

	public.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello World"})
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

}
