package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/TonmoyTalukder/go-students-api/internal/storage"
	"github.com/TonmoyTalukder/go-students-api/internal/types"
	"github.com/TonmoyTalukder/go-students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student...")

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "empty body", err)) // err
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Bad Request", err))
			return
		}

		// Request Validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErr(validateErrs))

			return
		}

		lastId, errLastId := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if errLastId != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.ErrorResponse(http.StatusBadRequest, "Error with last id", errLastId))
			return
		}

		response.WriteJson(w, http.StatusCreated, response.SuccessResponse(http.StatusOK, "Student created successfully", map[string]int64{"id": lastId}, nil))

		// w.Write([]byte("Welcome to student api..."))
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Error getting bad request", slog.Any("error", err))
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid student ID", err))
			return
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("Error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve student", err))
			return
		}

		response.WriteJson(w, http.StatusOK, response.SuccessResponse(http.StatusOK, "Student retrieved successfully", student, nil))
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// slog.Info("Getting all students")

		// students, err := storage.GetStudents()
		// if err != nil {
		// 	response.WriteJson(w, http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve students", err))
		// 	return
		// }

		// response.WriteJson(w, http.StatusOK, response.SuccessResponse(http.StatusOK, "Students retrieved successfully", students, nil))

		slog.Info("Getting all students")

		// Get query parameters
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		// Set default values
		page := 1
		limit := 10

		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}

		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}

		// Calculate offset
		offset := (page - 1) * limit

		// Fetch students with limit and offset
		students, err := storage.GetStudentsWithPagination(limit, offset)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve students", err))
			return
		}

		// Optional: total count for frontend pagination UI
		total, err := storage.CountStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to count students", err))
			return
		}

		meta := map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		}

		response.WriteJson(w, http.StatusOK, response.SuccessResponse(http.StatusOK, "Students retrieved successfully", students, meta))
	}
}
