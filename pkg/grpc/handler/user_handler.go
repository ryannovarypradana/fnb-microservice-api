// File: ryannovarypradana/fnb-microservice-api/fnb-microservice-api-dd6285232082f71efc6950ba298fd97bc68fbcc3/pkg/grpc/handler/user_handler.go

package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(client pb.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	req := &pb.GetUserRequest{Id: id}

	res, err := h.client.GetUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// HANDLER BARU UNTUK GET ALL USERS
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	storeID := c.Query("store_id")
	companyID := c.Query("company_id")
	search := c.Query("search")

	req := &pb.GetAllUsersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	}
	if storeID != "" {
		req.StoreId = &storeID
	}
	if companyID != "" {
		req.CompanyId = &companyID
	}
	if search != "" {
		req.Search = &search
	}

	res, err := h.client.GetAllUsers(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req pb.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	req.Id = id

	res, err := h.client.UpdateUser(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	req := &pb.DeleteUserRequest{Id: id}

	res, err := h.client.DeleteUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *UserHandler) CreateCompanyWithRep(c *fiber.Ctx) error {
	var req pb.CreateCompanyWithRepRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	res, err := h.client.CreateCompanyWithRep(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}
