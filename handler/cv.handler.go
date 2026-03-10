package handler

import (
	"context"
	"fmt"
	"go-cariir-api/database"
	"go-cariir-api/model/request"
	"go-cariir-api/model/response"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func InsertManyCV(data []request.JobSearchRequest) error {
	collection := database.MongoDB.Collection("cv")

	var docs []interface{}
	for _, item := range data {
		item.CreatedAt = time.Now()
		docs = append(docs, item)
	}

	_, err := collection.InsertMany(context.Background(), docs)
	return err
}

func CreateCVHandler(ctx *fiber.Ctx) error {
	var reqCVs []request.JobSearchRequest

	if ctx.Get("API_KEY") != os.Getenv("API_KEY") {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Invalid api key",
			Errors:  nil,
			Code:    fiber.StatusUnauthorized,
		})
	}

	log.Println("body data:", string(ctx.Body()))

	if err := ctx.BodyParser(&reqCVs); err != nil {
		log.Println("body parser error:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Invalid request body",
			Errors:  err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	validate := validator.New()
	for i, reqCV := range reqCVs {
		if err := validate.Struct(reqCV); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
				Status:  "error",
				Message: fmt.Sprintf("Validation failed at index %d", i),
				Errors:  err.Error(),
				Code:    fiber.StatusBadRequest,
			})
		}
	}

	if err := InsertManyCV(reqCVs); err != nil {
		log.Println("insert cv error:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "Failed to create cv",
			Errors:  err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.GenericResponse{
		Status:  "success",
		Message: "CV created successfully",
		Data:    reqCVs,
		Code:    fiber.StatusCreated,
	})
}
