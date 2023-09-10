package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

func main() {
	engine := django.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine})
	app.Static("/static", "./static")
	app.Static("/static", "./styles")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/contacts")
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		q := c.Queries()["q"]

		db := Db()
		users := []User{}
		if q == "" {
			db.Find(&users)
		} else {
			db.Where("name LIKE ?", "%"+q+"%").Find(&users)
		}

		return c.Render("search", fiber.Map{"q": q, "users": users})
	})

	app.Post("/contacts", func(c *fiber.Ctx) error {
		q := c.FormValue("q")
		db := Db()
		users := []User{}
		if q == "" {
			db.Find(&users)
		} else {
			db.Where("name LIKE ?", "%"+q+"%").Find(&users)
		}

		return c.Render("search", fiber.Map{"q": q, "users": users})
	})

	app.Get("/contacts/add", func(c *fiber.Ctx) error {
		return c.Render("add", fiber.Map{})
	})

	app.Post("/contacts/add", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		phone := c.FormValue("phone")
		fmt.Println(name, phone)

		db := Db()
		db.Create(&User{Name: name, Phone: phone})

		return c.Redirect("/contacts")
	})

	app.Get("/contacts/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		db := Db()
		user := User{}
		db.Where("id = ?", id).First(&user)
		fmt.Println(user)
		return c.Render("view", fiber.Map{"user": user})
	})

	app.Post("contacts/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		name := c.FormValue("name")
		phone := c.FormValue("phone")
		db := Db()
		user := User{}
		db.First(&user, id)
		user.Name = name
		user.Phone = phone
		db.Save(&user)

		return c.Redirect("/contacts")
	})

	app.Delete("/contacts/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user := User{}
		db := Db()
		db.Delete(&user, id)
		return c.Redirect("/contacts", 303)
	})

	app.Get("/contacts/:id/edit", func(c *fiber.Ctx) error {
		id := c.Params("id")
		db := Db()
		user := User{}
		db.Where("id = ?", id).First(&user)
		return c.Render("edit", fiber.Map{"user": user})

	})

	app.Get("/getcontacts", func(c *fiber.Ctx) error {
		return c.SendString(`<ul>
		<li><a href="mailto:joe@example.com">Joe</a></li>
		<li><a href="mailto:sarah@example.com">Sarah</a></li>
		<li><a href="mailto:fred@example.com">Fred</a></li>
	  </ul>
	  <div><button hx-target="change" hx-swap="delete">hide</button></div>`)
	})

	app.Get("/spinner", func(c *fiber.Ctx) error {
		time.Sleep(time.Second)
		fmt.Println("connected")
		return c.SendString("All data loaded")
	})

	app.Listen(":3000")

}

// GORM EXAMPLES

// MIGRATE
// db.AutoMigrate(&User{})

// CREATE
// db.Create(&User{Name: "Janis", Phone: "1232-2323-5342"})

// FIND WHERE
// db.Where("name LIKE ?", "%"+q+"%").Find(&users)

// DELETE
// db.Delete(&user, id)

// https://hypermedia.systems/htmx-in-action/
