package router

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UploadImageBody struct {
	Image []byte `form:"image"`
}

func UploadImage(c *fiber.Ctx) error {
	imageFile, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": "request body should have file in field 'image'",
		})
	}

	file, err := imageFile.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "error opening file",
		})
	}
	defer file.Close()

	// var data []byte
	// n, err := file.Read(data)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"ok":      false,
	// 		"message": "error reading file",
	// 	})
	// }

	log.Print("file size: ", imageFile.Size)

	fileName := getNewFileName(imageFile.Filename)
	// err = ioutil.WriteFile(fileName, data, 0644)
	err = c.SaveFile(imageFile, fileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "error writing file to fs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":    true,
		"image": fileName,
	})
}

func ServeImage(c *fiber.Ctx) error {
	filename := c.Params("filename")

	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": "no filename in request path",
		})
	}

	filePath := filepath.Join("public/images", filename)
	data, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":      false,
			"message": "file not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "error reading file",
		})
	}

	return c.Status(fiber.StatusOK).Type(filepath.Ext(filename)).Send(data)
}

func getNewFileName(filename string) string {
	basePath := "public/images"

	newFileName := ""
	for newFileName == "" {
		possibleName := filepath.Join(basePath, fmt.Sprintf("%d-%s", time.Now().Unix(), filename))
		if _, err := os.Stat(possibleName); os.IsExist(err) {
			continue
		}

		newFileName = possibleName
	}

	return newFileName
}
