package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"database/sql"

	"github.com/AKASHLM010/BLOG_PROJ/config"
	"github.com/AKASHLM010/BLOG_PROJ/database"
	"github.com/AKASHLM010/BLOG_PROJ/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetAllBlogs(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT * FROM blogs ORDER BY created_at DESC")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	blogs := []models.Blog{}
	for rows.Next() {
		var blog models.Blog
		var userID sql.NullInt64
		err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt, &userID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		if userID.Valid {
			blog.UserID = int(userID.Int64)
		}
		blog.CreatedAt = blog.CreatedAt.Local() // Convert to local time (if needed)
		blog.UpdatedAt = blog.UpdatedAt.Local() // Convert to local time (if needed)
		blogs = append(blogs, blog)
	}

	return c.JSON(blogs)
}

func CreateBlog(c *fiber.Ctx) error {
	var blog models.Blog
	err := json.Unmarshal(c.Body(), &blog)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	currentTime := time.Now()
	blog.CreatedAt = currentTime
	blog.UpdatedAt = currentTime

	// Get the logged-in user's full name
	fullName, err := getUserDetails(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	// Set the author field of the blog to the user's full name
	blog.Author = fullName

	// Get the logged-in user's ID
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	blog.UserID = userID

	query := "INSERT INTO blogs (title, content, author, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	row := database.DB.QueryRow(query, blog.Title, blog.Content, blog.Author, blog.CreatedAt, blog.UpdatedAt, blog.UserID)
	if err := row.Scan(&blog.ID); err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Create a response map with the blogID
	response := map[string]interface{}{
		"blogID": blog.ID,
	}

	return c.JSON(response)
}

// getUserDetails retrieves the concatenated first name and last name of the logged-in user
func getUserDetails(c *fiber.Ctx) (string, error) {
	// Retrieve the JWT token from the cookie
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return "", errors.New("missing JWT token")
	}

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and secret key
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	// Verify and extract the claims from the token
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid JWT token")
	}

	// Extract the user details from the claims
	email, ok := (*claims)["email"].(string)
	if !ok || email == "" {
		return "", errors.New("invalid user email")
	}

	userID, ok := (*claims)["userID"].(float64)
	if !ok || userID == 0 {
		return "", errors.New("invalid user ID")
	}

	// Fetch the user details from the database using the extracted email and user ID
	// Replace this with your actual database query to retrieve the user details
	var user models.User
	err = database.DB.QueryRow("SELECT first_name, last_name FROM users WHERE email = $1 AND id = $2", email, int(userID)).Scan(&user.FirstName, &user.LastName)
	if err != nil {
		return "", err
	}

	// Concatenate the first name and last name
	fullName := user.FirstName + " " + user.LastName

	return fullName, nil
}

func getUserID(c *fiber.Ctx) (int, error) {
	// Retrieve the JWT token from the cookie
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return 0, errors.New("missing JWT token")
	}

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and secret key
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}

	// Verify and extract the claims from the token
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid JWT token")
	}

	// Extract the userID from the claims
	userIDFloat, ok := (*claims)["userID"].(float64)
	if !ok || userIDFloat == 0 {
		return 0, errors.New("invalid user ID")
	}

	userID := int(userIDFloat)

	return userID, nil
}

func GetBlogForEdit(c *fiber.Ctx) error {
	// Retrieve the logged-in user's ID
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	// Extract the blog ID from the request parameters
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Fetch the blog from the database using the blog ID and the logged-in user's ID
	blog, err := getBlogByIDAndUserID(id, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(blog)
}

func GetBlogsForEdit(c *fiber.Ctx) error {
	// Retrieve the logged-in user's ID
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	// Retrieve the user's blogs from the database using the user's ID
	blogs, err := GetUserBlogs(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(blogs)
}

func UpdateBlog(c *fiber.Ctx) error {
	// Retrieve the logged-in user's ID
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	// Extract the blog ID from the request parameters
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Fetch the existing blog from the database using the blog ID and the logged-in user's ID
	existingBlog, err := getBlogByIDAndUserID(id, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Parse the updated blog data from the request body
	var updatedBlog models.Blog
	err = c.BodyParser(&updatedBlog)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Update the fields of the existing blog with the updated values
	existingBlog.Title = updatedBlog.Title
	existingBlog.Content = updatedBlog.Content
	existingBlog.UpdatedAt = time.Now()

	// Update the blog in the database
	_, err = database.DB.Exec("UPDATE blogs SET title = $1, content = $2, updated_at = $3 WHERE id = $4 AND user_id = $5",
		existingBlog.Title, existingBlog.Content, existingBlog.UpdatedAt, existingBlog.ID, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("Blog updated successfully")
}

func getBlogByIDAndUserID(id, userID int) (models.Blog, error) {
	// Query the database to fetch the blog by its ID and user ID
	row := database.DB.QueryRow("SELECT * FROM blogs WHERE id = $1 AND user_id = $2", id, userID)

	var blog models.Blog
	var userIDNullable sql.NullInt64
	err := row.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt, &userIDNullable)
	if err != nil {
		return blog, err
	}

	if userIDNullable.Valid {
		blog.UserID = int(userIDNullable.Int64)
	}

	return blog, nil
}

func GetUserBlogs(userID int) ([]models.Blog, error) {
	// Fetch the user's blogs from the database using the user's ID
	query := "SELECT id, title, content, author, created_at, updated_at FROM blogs WHERE user_id = $1"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			return nil, err
		}
		blog.CreatedAt = blog.CreatedAt.Local() // Convert to local time (if needed)
		blog.UpdatedAt = blog.UpdatedAt.Local() // Convert to local time (if needed)
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func GetBlogsForDelete(c *fiber.Ctx) error {
	// Retrieve the logged-in user's ID
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	// Retrieve the user's blogs from the database using the user's ID
	blogs, err := GetUserBlogs(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Return the blogs data as JSON
	return c.JSON(blogs)
}

func DeleteBlog(c *fiber.Ctx) error {
	// Retrieve the logged-in user's ID
	userID, err := getUserID(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).SendString(err.Error())
	}

	// Extract the blog ID from the request parameters
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Delete the blog from the database using the blog ID and the logged-in user's ID
	err = deleteBlogByIDAndUserID(id, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Redirect to the delete page to show the updated list of blogs
	return c.Redirect("/delete")
}

func deleteBlogByIDAndUserID(id, userID int) error {
	// Perform the deletion operation on the database
	_, err := database.DB.Exec("DELETE FROM blogs WHERE id = $1 AND user_id = $2", id, userID)
	return err
}

func GetUserBlogsByUserID(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		// Handle the error, such as returning an error response
		return err
	}

	blogs, err := GetUserBlogs(userID)
	if err != nil {
		// Handle the error, such as returning an error response
		return err
	}

	// Return the blogs data as JSON
	return c.JSON(blogs)
}

func ViewBlog(c *fiber.Ctx) error {
	blogID := c.Params("id")

	blog, err := GetBlogByID(blogID)
	if err != nil {
		// Handle the error, such as returning an error response or redirecting to an error page
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Blog not found",
		})
	}

	// Return the blog data as JSON response
	return c.JSON(blog)
}

func GetBlogByID(blogID string) (*models.Blog, error) {
	// Perform a database query to fetch the blog by its ID
	query := "SELECT id, title, content, author, created_at, updated_at FROM blogs WHERE id = $1"
	row := database.DB.QueryRow(query, blogID)

	// Create a new Blog struct to hold the retrieved blog data
	blog := &models.Blog{}

	// Scan the row's values into the blog struct
	err := row.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where the blog doesn't exist
			return nil, errors.New("blog not found")
		}
		// Handle any other database query error
		return nil, err
	}

	return blog, nil
}
