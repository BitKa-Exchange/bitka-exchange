package middleware

import (
	"strings"

	"bitka/pkg/response"
	"bitka/pkg/token"

	"github.com/gofiber/fiber/v2"
)

// Protected returns a middleware that verifies the JWT.
func Protected(v *token.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get Token from Header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid Authorization format")
		}
		tokenStr := parts[1]

		// 2. Validate Token
		// Use UserContext to ensure we respect cancellations
		parsedToken, err := v.Validate(c.UserContext(), tokenStr)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, err.Error())
		}

		// 3. Store in Context
		// We store the Subject (User ID) for the handler to use
		sub, ok := parsedToken.Subject()
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token: missing subject")
		}

		c.Locals("user_id", sub)
		c.Locals("claims", parsedToken) // Store full token if needed

		return c.Next()
	}
}
