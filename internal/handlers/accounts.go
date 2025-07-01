package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/ashtonx86/nybl/internal/dependencies"
	"github.com/ashtonx86/nybl/internal/schemas"
	"github.com/ashtonx86/nybl/internal/store"
	"github.com/ashtonx86/nybl/internal/supervisor"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

type AccountsHandler struct {
	Supervisor   *supervisor.Supervisor
	AccountStore store.AccountStore
	OTPStore     store.OTPStore
}

func NewAccountsHandler(su *supervisor.Supervisor) *AccountsHandler {
	sqlite := dependencies.MustGetSQLite(su)

	return &AccountsHandler{
		Supervisor:   su,
		AccountStore: store.NewSQLiteAccountStore(sqlite.DB),
		OTPStore:     store.NewSQLiteOTPStore(sqlite.DB),
	}
}

func (handler *AccountsHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Hello world")
}

func (h *AccountsHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	body := new(schemas.RequestAccountCreate)
	c.BodyParser(&body)

	err := Validator.StructCtx(ctx, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewBadRequestError(err))
	}

	acc, err := h.AccountStore.Create(ctx, body.Name, body.EmailID)
	if err != nil {
		slog.Error("Failed to create new account", "err", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed")
	}

	otp, err := h.OTPStore.Create(ctx, acc.EmailID, acc.ID)
	if err != nil {
		slog.Error("Failed to create OTP for account", "accountID", acc.ID, "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to initate OTP")
	}

	return c.JSON(otp)
}
