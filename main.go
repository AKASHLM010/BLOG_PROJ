package main

import (
	"database/sql"
	"log"

	"github.com/AKASHLM010/BLOG_PROJ/database"
	"github.com/AKASHLM010/BLOG_PROJ/handler"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const secretKey = "jwxruguggkgghiihkg"

var db *sql.DB

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	err := database.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/home.html") // Serve the home.html file as the home page
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.SendFile("./public/register.html")
	})
	app.Post("/register", handler.RegisterUser)
	app.Post("/login", handler.AuthenticateUser)
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendFile("./public/login.html")
	})
	app.Get("/blogs", handler.GetAllBlogs)
	app.Get("/blogs/:id", handler.GetBlogByID)
	app.Post("/blogs", handler.CreateBlog)
	app.Put("/blogs/:id", handler.UpdateBlog)
	app.Delete("/blogs/:id", handler.DeleteBlog)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running")
	})

	log.Println("Server is running on http://localhost:8000")
	log.Fatal(app.Listen(":8000"))
}