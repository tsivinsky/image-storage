package main

import (
	"image-storage/internal/db"
	"image-storage/internal/router"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/:filename", router.ServeImage)
	app.Post("/api/images/upload", router.UploadImage)
	app.Get("/api/images/:id", router.GetImageById)
	app.Delete("/api/images/:id", router.DeleteImageById)

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
