package router

import (
	"fmt"
	"image-storage/internal/db"
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

	log.Print("file size: ", imageFile.Size)

	fileName := getNewFileName(imageFile.Filename)
	err = c.SaveFile(imageFile, filepath.Join("public/images", fileName))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "error writing file to fs",
		})
	}

	newImage := db.Image{
		Filename: fileName,
	}
	err = db.Db.Create(&newImage).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "error adding image to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":    true,
		"image": newImage,
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

func GetImageById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": err.Error(),
		})
	}

	var image db.Image
	err = db.Db.Find(&image, "id", id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":      false,
			"message": err.Error(),
		})
	}

	fullUrl := fmt.Sprintf("%s/%s", c.BaseURL(), image.Filename)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok": true,
		"data": fiber.Map{
			"id":         image.Id,
			"filename":   image.Filename,
			"url":        fullUrl,
			"created_at": image.CreatedAt,
			"updated_at": image.UpdatedAt,
		},
	})
}

func DeleteImageById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": err.Error(),
		})
	}

	var image db.Image
	err = db.Db.Find(&image, "id", id).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": err.Error(),
		})
	}

	err = os.Remove(filepath.Join("public/images", image.Filename))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": err.Error(),
		})
	}

	err = db.Db.Delete(&image, "id", id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ok":      false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"message": "image deleted successfully",
	})
}

func getNewFileName(filename string) string {
	basePath := "public/images"

	newFileName := ""
	for newFileName == "" {
		possibleName := fmt.Sprintf("%d-%s", time.Now().Unix(), filename)
		if _, err := os.Stat(filepath.Join(basePath, possibleName)); os.IsExist(err) {
			continue
		}

		newFileName = possibleName
	}

	return newFileName
}
