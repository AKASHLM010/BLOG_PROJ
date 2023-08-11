package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

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

	// Check if email is already registered
	var existingUser models.User
	err = database.DB.QueryRow("SELECT id FROM users WHERE email = $1", user.Email).Scan(&existingUser.ID)
	if err == nil {
		return c.Status(http.StatusConflict).SendString("Email already registered")
	} else if err != sql.ErrNoRows {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Check if phone number is already registered
	err = database.DB.QueryRow("SELECT id FROM users WHERE phone = $1", user.Phone).Scan(&existingUser.ID)
	if err == nil {
		return c.Status(http.StatusConflict).SendString("Phone number already registered")
	} else if err != sql.ErrNoRows {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	query := "INSERT INTO users (first_name, last_name, email, password, phone) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = database.DB.QueryRow(query, user.FirstName, user.LastName, user.Email, hashedPassword, user.Phone).Scan(&user.ID)
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var storedPassword string
	query := "SELECT id, password FROM users WHERE email = $1"
	err = database.DB.QueryRow(query, user.Email).Scan(&user.ID, &storedPassword)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Email or Password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Email or Password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  user.Email,
		"userID": user.ID,
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the JWT token as a cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24), // Set the cookie expiration time
		HTTPOnly: true,                           // Ensure the cookie is not accessible via JavaScript
		Secure:   true,                           // Set the cookie to be secure (HTTPS only)
		SameSite: "Strict",                       // Set the SameSite attribute to Strict
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"token": tokenString})
}

func Logout(c *fiber.Ctx) error {
    // Clear the JWT token cookie
    c.ClearCookie("jwt")

    // Redirect the user to the login page or any other desired page
    return c.Redirect("/login")
}

func UserProfile(c *fiber.Ctx) error {
	// Get the JWT token from the cookie
	cookie := c.Cookies("jwt")

	// Parse the JWT token
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}

	// Extract user information from the JWT token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
	}

	// Retrieve user data based on the claims (email and userID)
	email := claims["email"].(string)
	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return c.Status(http.StatusInternalServerError).SendString("Invalid user ID")
	}
	userID := int(userIDFloat)

	// Retrieve the user's profile data from the database using the email or userID
	var user models.User
	query := "SELECT first_name FROM users WHERE email = $1 AND id = $2"
	err = database.DB.QueryRow(query, email, userID).Scan(&user.FirstName)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Retrieve the user's blog data from the database using the user's ID
	blogs, err := GetUserBlogs(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Create a data structure to hold the user's profile data and blog IDs
	data := struct {
		FirstName string
		BlogIDs   []int
	}{
		FirstName: user.FirstName,
	}

	// Extract the blog IDs from the retrieved blogs
	var blogIDs []int
	for _, blog := range blogs {
		blogIDs = append(blogIDs, blog.ID)
	}

	data.BlogIDs = blogIDs

	// Render the profile.html template with the data
	return c.Render("public/profile.html", data)
}

func CheckAuthentication(c *fiber.Ctx) error {
	// Get the JWT token from the request cookie
	cookie := c.Cookies("jwt")
	if cookie == "" {
		// No JWT token found, user is not authenticated
		return c.JSON(fiber.Map{"isAuthenticated": false})
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// Provide the same signing key used in the AuthenticateUser function
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		// Invalid token or token verification failed, user is not authenticated
		return c.JSON(fiber.Map{"isAuthenticated": false})
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Token is valid, user is authenticated
		return c.JSON(fiber.Map{"isAuthenticated": true})
	}

	// Token verification failed, user is not authenticated
	return c.JSON(fiber.Map{"isAuthenticated": false})
}
