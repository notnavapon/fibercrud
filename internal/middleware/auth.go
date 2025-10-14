package middleware

import (
	jwtpkg "clean/internal/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected(jwtSecret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("Jwt_Token")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "please login",
			})
		}

		claims, err := jwtpkg.ValidateToken(token, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("email", claims.Email)

		return c.Next()
	}
}
