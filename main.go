package main

import (
	"database/sql"
	"log"

	"github.com/AKASHLM010/BLOG_PROJ/database"
	"github.com/AKASHLM010/BLOG_PROJ/handler"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	app := fiber.New()

	app.Static("/", "./public")
	app.Static("/api/uploads", "./uploads")
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
	app.Get("/logout", handler.Logout)

	app.Get("/profile", handler.UserProfile)

	app.Get("/blogs", handler.GetAllBlogs)
	app.Get("/view", handler.GetUserBlogsByUserID)

	
	app.Post("/api/create-blog-with-image", handler.CreateBlog)
	app.Get("/edit", handler.GetBlogsForEdit)
    app.Get("/blogs/:id", handler.ViewBlog)
	

	app.Get("/edit/:id", handler.GetBlogForEdit)

	app.Patch("/edit/:id", handler.UpdateBlog)
	app.Get("/delete", handler.GetBlogsForDelete)
    app.Post("/delete/:id", handler.DeleteBlog)
    app.Get("/check-authentication", handler.CheckAuthentication)
  
	// Route to post a comment on a blog
	app.Post("/api/blogs/:blog_id/comments", handler.PostComment)

	// Route to get comments for a specific blog
	app.Get("/api/blogs/:blog_id/comments", handler.GetCommentsForBlog)

	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running")
	})

	log.Println("Server is running on http://localhost:8000")
	log.Fatal(app.Listen(":8000"))
}