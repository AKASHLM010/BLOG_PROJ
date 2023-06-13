package handler

import (
	"encoding/json"
	"net/http"
    
	"github.com/AKASHLM010/BLOG_PROJ/config"
	"github.com/AKASHLM010/BLOG_PROJ/database"
	"github.com/AKASHLM010/BLOG_PROJ/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *fiber.Ctx) error {
	var user models.User
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err = database.DB.QueryRow(query, user.Username, hashedPassword).Scan(&user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	user.Password = ""

	return c.JSON(user)
}

func AuthenticateUser(c *fiber.Ctx) error {
	var user models.User
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	var storedPassword string
	query := "SELECT id, password FROM users WHERE username = $1"
	err = database.DB.QueryRow(query, user.Username).Scan(&user.ID, &storedPassword)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"userID":   user.ID,
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"token": tokenString})
}
