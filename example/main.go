package main

import (
	"fmt"
	"strings"

	fibercasbin "github.com/beluxx/fiber-casbin/v3"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	authz := fibercasbin.New(fibercasbin.Config{
		ModelFilePath: "model.conf",
		PolicyAdapter: fileadapter.NewAdapter("policy.csv"),
		Lookup: func(c *fiber.Ctx) string {
			// get subject from BasicAuth, JWT, Cookie etc in real world
			return "alice"
		},
		FilterRoute: func(whitelist []string) func(c *fiber.Ctx) bool {
			return func(c *fiber.Ctx) bool {
				for _, item := range whitelist {
					return strings.Compare(item, c.Path()) == 0
				}
				return false
			}
		}([]string{"/filter/a"}),
	})

	app.Post("/blog",
		authz.RequiresPermissions([]string{"blog:create"}),
		func(c *fiber.Ctx) error {
			return c.SendString("Blog created")
		},
	)

	app.Put("/blog/:id",
		authz.RequiresRoles([]string{"admin"}),
		func(c *fiber.Ctx) error {
			return c.SendString(fmt.Sprintf("Blog updated with Id: %s", c.Params("id")))
		},
	)

	filter := app.Group("/filter", authz.RoutePermission())
	{
		filter.Get("/a", func(c *fiber.Ctx) error {
			return c.SendString("Your entry ")
		})
		filter.Post("/b", func(c *fiber.Ctx) error {
			return c.SendString("Your entry ")
		})
	}

	app.Listen(":8080")
}
