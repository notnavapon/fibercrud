package handlerUser

import (
	dtoUser "clean/internal/dto/user"
	usecaseUser "clean/internal/usecase/user"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase usecaseUser.UserUsecase
}

func NewUserHandler(usecase usecaseUser.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: usecase}
}

func (handler *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dtoUser.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := handler.userUsecase.CreateUser(&req)

	if err != nil {
		if err.Error() == "email alreay exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})

}

func (handler *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req dtoUser.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	res, err := handler.userUsecase.LoginUser(&req)

	if err != nil {
		if err.Error() == "email not found" || err.Error() == "password don't match" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Jwt_Token",
		Value:    res.Token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Login successful",
		"data":    res,
	})

}

func (handler *UserHandler) LogoutUser(c *fiber.Ctx) error {
	err := handler.userUsecase.LogoutUser()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "Jwt_Token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	return c.JSON(fiber.Map{
		"message": "logout success",
	})
}

func (handler *UserHandler) GetAllUser(c *fiber.Ctx) error {
	email := c.Locals("email")
	userEmail := email.(string)
	users, err := handler.userUsecase.GetAllUser()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"getByuser": userEmail,
		"data":      users,
	})

}

func (handler *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	user, err := handler.userUsecase.GetByID(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Find user success",
		"data":    user,
	})

}

func (handler *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	err = handler.userUsecase.DeleteUser(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Delete User success",
	})

}

func (handler *UserHandler) Update(c *fiber.Ctx) error {
	var req dtoUser.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if req.Email == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	//check cookie.email = req.email
	cookieEmail := c.Locals("email").(string)

	if cookieEmail != req.CurrentEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email in the cookie does not match the logged-in email.",
		})
	}

	res, err := handler.userUsecase.UpdateUser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Jwt_Token",
		Value:    res["token"].(string),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.Status(fiber.StatusAccepted).JSON(res["user"])

}

func (handler *UserHandler) ChangePassword(c *fiber.Ctx) error {
	var req dtoUser.ChangePasswordRequest
	cookieEmail := c.Locals("email").(string)
	fmt.Println("Cookie Email:", cookieEmail)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	res, err := handler.userUsecase.ChangePasswordUser(&req, cookieEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "change password successfully",
		"data":    res,
	})
}
