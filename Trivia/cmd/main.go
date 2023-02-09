package main

import "github.com/gofiber/fiber/v2"

// https://www.youtube.com/watch?v=p08c0-99SyU

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Aarti Hitesh Chhasiya!")
	})

	app.Listen(":3000")
}
