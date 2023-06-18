package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AKASHLM010/BLOG_PROJ/database"
	"github.com/AKASHLM010/BLOG_PROJ/models"
	"github.com/gofiber/fiber/v2"
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
		err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		blogs = append(blogs, blog)
	}

	return c.JSON(blogs)
}

func GetBlogByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	var blog models.Blog
	err = database.DB.QueryRow("SELECT * FROM blogs WHERE id = $1", id).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(blog)
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

	result, err := database.DB.Exec("INSERT INTO blogs (title, content, author, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		blog.Title, blog.Content, blog.Author, blog.CreatedAt, blog.UpdatedAt)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	blog.ID = int(id)

	return c.JSON(blog)
}

func UpdateBlog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	var blog models.Blog
	err = json.Unmarshal(c.Body(), &blog)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	blog.ID = id
	blog.UpdatedAt = time.Now()

	_, err = database.DB.Exec("UPDATE blogs SET title = $1, content = $2, author = $3, updated_at = $4 WHERE id = $5",
		blog.Title, blog.Content, blog.Author, blog.UpdatedAt, blog.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(blog)
}

func DeleteBlog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	_, err = database.DB.Exec("DELETE FROM blogs WHERE id = $1", id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(http.StatusOK)
}