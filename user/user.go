package user

import (
	"github.com/SohailHaqyar/friday/data"
	"github.com/SohailHaqyar/friday/jwt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

func SetupRoutes(app *fiber.App, dbEngine *xorm.Engine) {
	app.Post("/signup", func(c *fiber.Ctx) error {
		return signup(c, dbEngine)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return login(c, dbEngine)
	})
}

func signup(c *fiber.Ctx, engine *xorm.Engine) error {
	req := new(SignupRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	// check that the request body is valid
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid signup credentials",
		})
	}
	// save the user to the database
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save user to the database",
		})
	}
	// create
	user := &data.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}
	_, err = engine.Insert(user)
	if err != nil {
		return err
	}
	token, exp, err := jwt.CreateJWTToken(*user)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
}

func login(c *fiber.Ctx, engine *xorm.Engine) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}
	// check that the request body is valid
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid login credentials",
		})
	}

	user := new(data.User)
	has, err := engine.Where("email = ?", req.Email).Get(user)
	if err != nil {
		return err
	}

	if !has {
		fiber.NewError(fiber.StatusBadRequest, "invalid login credentials ")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid login credentials",
		})
	}

	token, exp, err := jwt.CreateJWTToken(*user)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
}
