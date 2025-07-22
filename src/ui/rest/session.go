package rest

import (
	domainApp "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/app"
	domainSession "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/session"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type Session struct {
	SessionManager domainSession.ISessionManager
	AppUsecase     domainApp.IAppUsecase
}

func InitRestSession(app *fiber.App, sessionManager domainSession.ISessionManager, appUsecase domainApp.IAppUsecase) {
	rest := Session{SessionManager: sessionManager, AppUsecase: appUsecase}
	app.Post("/sessions/:sessionId", rest.CreateSession)
	app.Get("/sessions", rest.GetAllSessions)
	app.Delete("/sessions/:sessionId", rest.DeleteSession)
	app.Get("/sessions/:sessionId/login", rest.Login)
}

func (h *Session) CreateSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")
	_, err := h.SessionManager.CreateSession(sessionId)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "Session created successfully",
		Results: sessionId,
	})
}

func (h *Session) GetAllSessions(c *fiber.Ctx) error {
	sessions := h.SessionManager.GetAllSessions()
	// We only need to return the session IDs
	sessionIds := make([]string, 0, len(sessions))
	for id := range sessions {
		sessionIds = append(sessionIds, id)
	}

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "Fetched all sessions",
		Results: sessionIds,
	})
}

func (h *Session) DeleteSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")
	err := h.SessionManager.DeleteSession(sessionId)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "Session deleted successfully",
		Results: sessionId,
	})
}

func (h *Session) Login(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")
	// The login usecase now needs the session ID to know which client to use
	response, err := h.AppUsecase.Login(c.UserContext(), sessionId)
	utils.PanicIfNeeded(err)

	return c.JSON(utils.ResponseData{
		Status:  200,
		Code:    "SUCCESS",
		Message: "Login success",
		Results: map[string]any{
			"qr_base64":   response.QRBase64,
			"qr_duration": response.Duration,
		},
	})
} 