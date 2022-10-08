package main

import (
	"image-storage/internal/router"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/api/images/upload", router.UploadImage)
	app.Get("/:filename", router.ServeImage)

	port := getPort(":5000")
	log.Fatal(app.Listen(port))
}

func getPort(fallbackPort string) string {
	port := os.Getenv("PORT")

	if port == "" {
		port = fallbackPort
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return port
}
