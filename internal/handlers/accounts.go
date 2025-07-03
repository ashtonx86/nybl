package handlers

import (
	"context"
	"fmt"
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
	Mail         *dependencies.MailSingleton
}

func NewAccountsHandler(su *supervisor.Supervisor) *AccountsHandler {
	sqlite := dependencies.MustGetSQLite(su)

	return &AccountsHandler{
		Supervisor:   su,
		AccountStore: store.NewSQLiteAccountStore(sqlite.DB),
		OTPStore:     store.NewSQLiteOTPStore(sqlite.DB),
		Mail:         dependencies.MustGetMail(su),
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

	go func(acc *schemas.Account, otp *schemas.OTP) {
		err := h.Mail.SendMail(acc.EmailID, fmt.Sprintf("Your OTP for Nybl is %s", otp.Code), "", fmt.Sprintf("<p>Hello <b>%s</b>! Please verify your email, <b>%s</b> using the One Time Password below.</p> <br><h3><b>%s</b></h3><br>DO NOT SHARE THIS CODE WITH ANYONE", acc.Name, acc.EmailID, otp.Code))

		if err != nil {
			slog.Error("Failed to send OTP", "account_id", acc.ID, "error", err)
		}
	}(acc, otp)

	res := schemas.NewOkResponse(map[string]interface{}{
		"acc":  acc,
		"next": "Verify email",
	})
	return c.JSON(res)
}
