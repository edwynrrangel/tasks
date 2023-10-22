package session

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	go_jwt "github.com/golang-jwt/jwt"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/edwynrrangel/tasks/logger"
	"github.com/edwynrrangel/tasks/token"
)

type session struct {
	tokenService token.Service
}

func NewSession(tokenService token.Service) Middleware {
	return &session{
		tokenService: tokenService,
	}
}

func (s *session) Validate(ctx *fiber.Ctx) error {

	// Get JWT from Authorization header
	jwt, err := getJWTFromHeader(ctx)
	if err != nil {
		return err.ToFiber(ctx)
	}

	// Get payload from JWT
	payload, err := getPayloadFromJWT(jwt)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": errMsgExpiredInvalidMissing,
		})
	}

	// Validate JWT
	err = s.tokenService.Validate(jwt)
	if err != nil {
		return err.ToFiber(ctx)
	}

	// Set user in context
	ctx.Locals("user", payload)

	return ctx.Next()
}

func (s *session) ValidateFirstLogin(ctx *fiber.Ctx) error {

	// Get JWT from Authorization header
	jwt, err := getJWTFromHeader(ctx)
	if err != nil {
		return err.ToFiber(ctx)
	}

	// Get payload from JWT
	payload, err := getPayloadFromJWT(jwt)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": errMsgExpiredInvalidMissing,
		})
	}

	// Validate JWT
	err = s.tokenService.ValidateFirstLogin(jwt, payload.ID)
	if err != nil {
		return err.ToFiber(ctx)
	}

	// Set user in context
	ctx.Locals("user", payload)

	return ctx.Next()
}

func getJWTFromHeader(ctx *fiber.Ctx) (string, *errors.Error) {
	// Get JWT from Authorization header
	token := ctx.Get("Authorization")
	if token == "" {
		return "", errors.New(http.StatusUnauthorized, errMsgMissingJwt, nil)
	}

	// Remove Bearer from JWT
	token = token[7:]

	return token, nil
}

func getPayloadFromJWT(jwt string) (*Auth, *errors.Error) {

	// Get claims from JWT wihout validating it
	token, _, err := new(go_jwt.Parser).ParseUnverified(jwt, go_jwt.MapClaims{})
	if err != nil {
		logger.Error(
			"error parsing token",
			"func", "getPayloadFromJWT - jwt.ParseUnverified",
			"error", err,
		)
		return nil, errors.New(http.StatusUnauthorized, errMsgExpiredInvalidMissing, nil)
	}

	// Get payload from claims
	payload := token.Claims.(go_jwt.MapClaims)["payload"].(map[string]interface{})

	// encode payload to json and decode it to User struct
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		logger.Error(
			"error marshaling payload",
			"func", "getPayloadFromJWT - json.Marshal",
			"error", err,
		)
		return nil, errors.New(http.StatusInternalServerError, errors.MsgInternalServerError, nil)
	}

	var auth = new(Auth)
	err = json.Unmarshal(payloadJson, &auth)
	if err != nil {
		logger.Error(
			"error unmarshaling payload",
			"func", "getPayloadFromJWT - json.Unmarshal",
			"error", err,
		)
		return nil, errors.New(http.StatusInternalServerError, errors.MsgInternalServerError, nil)
	}

	return auth, nil
}
