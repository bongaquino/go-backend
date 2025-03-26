package tokens

import (
	"net/http"

	"koneksi/server/app/helper"
	"koneksi/server/app/service"

	"github.com/gin-gonic/gin"
)

// RequestTokenController handles user authentication and token generation
type RequestTokenController struct {
	tokenService *service.TokenService
	userService  *service.UserService
	mfaService   *service.MFAService
}

// NewRequestTokenController initializes a new RequestTokenController
func NewRequestTokenController(tokenService *service.TokenService, userService *service.UserService, mfaService *service.MFAService) *RequestTokenController {
	return &RequestTokenController{
		tokenService: tokenService,
		userService:  userService,
		mfaService:   mfaService,
	}
}

// Handle processes the login request and returns an access & refresh token
func (rc *RequestTokenController) Handle(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	// Validate the payload
	if err := rc.validatePayload(ctx, &request); err != nil {
		return
	}

	// Get user details
	user, _, err := rc.userService.GetUserProfileByEmail(ctx, request.Email)
	if err != nil {
		helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to retrieve user profile", nil, nil)
		return
	}

	// Check if MFA is enabled
	if user.IsMFAEnabled {
		// Generate login code
		loginCode, err := rc.mfaService.GenerateLoginCode(ctx.Request.Context(), user.ID.Hex())
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusInternalServerError, "failed to generate login code", nil, nil)
		}

		// Respond with boolean flag indicating MFA is enabled
		helper.FormatResponse(ctx, "success", http.StatusOK, "request token successful", gin.H{
			"is_mfa_enabled": true,
			"login_code":     loginCode,
		}, nil)
	} else {
		// Authenticate user and generate tokens
		accessToken, refreshToken, err := rc.tokenService.AuthenticateUser(ctx.Request.Context(), request.Email, request.Password)
		if err != nil {
			helper.FormatResponse(ctx, "error", http.StatusUnauthorized, err.Error(), nil, nil)
			return
		}

		// Respond with tokens
		helper.FormatResponse(ctx, "success", http.StatusOK, "request token successful", gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}, nil)
	}

}

// validatePayload validates the incoming request payload
func (rc *RequestTokenController) validatePayload(ctx *gin.Context, request any) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		helper.FormatResponse(ctx, "error", http.StatusBadRequest, "invalid input", nil, nil)
		return err
	}
	return nil
}
