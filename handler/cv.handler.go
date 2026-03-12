package handler

import (
	"bytes"
	"context"
	"fmt"
	"go-cariir-api/database"
	"go-cariir-api/model/request"
	"go-cariir-api/model/response"
	"io"
	"log"
	"mime/multipart"
	"net/http"
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
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericResponse{
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

func UploadCV(ctx *fiber.Ctx) error {

	fileHeader, err := ctx.FormFile("CV_File")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.GenericResponse{
			Status:  "error",
			Message: "File is required",
			Errors:  nil,
			Code:    fiber.StatusBadRequest,
		})
	}

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "failed to open uploaded file",
			Errors:  nil,
			Code:    fiber.StatusInternalServerError,
		})
	}
	defer src.Close()

	// Prepare multipart body for n8n
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("CV_File", fileHeader.Filename)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "failed to create multipart form",
			Errors:  nil,
			Code:    fiber.StatusInternalServerError,
		})
	}

	_, err = io.Copy(part, src)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "failed to copy file content",
			Errors:  nil,
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Optional: add extra fields for n8n
	_ = writer.WriteField("Email Address", "test@demo.com")
	_ = writer.WriteField("Preferred Location", "Surabaya")
	_ = writer.WriteField("Max Jobs to Analyze", "2")

	if err := writer.Close(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "failed to finalize multipart form",
			Errors:  nil,
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Send request to n8n webhook
	n8nURL := "https://n8n.alvianardhiansyah.my.id/webhook/2696e245-1268-4fe0-a5af-3984ea157abd"

	req, err := http.NewRequest(http.MethodPost, n8nURL, &body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericResponse{
			Status:  "error",
			Message: "failed to create request to n8n",
			Errors:  nil,
			Code:    fiber.StatusInternalServerError,
		})
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(response.GenericResponse{
			Status:  "error",
			Message: fmt.Sprintf("failed to send request to n8n: %v", err),
			Errors:  nil,
			Code:    fiber.StatusInternalServerError,
		})
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	return ctx.Status(resp.StatusCode).Send(respBody)

}
