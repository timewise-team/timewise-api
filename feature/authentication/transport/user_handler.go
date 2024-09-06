package transport

import (
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/authentication/dto"
	"github.com/cesc1802/onboarding-and-volunteer-service/feature/authentication/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthenticationHandler struct {
	usecase usecase.UserUsecaseInterface
}

func NewAuthenticationHandler(usecase usecase.UserUsecaseInterface) *AuthenticationHandler {
	return &AuthenticationHandler{usecase: usecase}
}

// Login godoc
// @Summary Login
// @Description Login
// @Produce json
// @Tags authentication
// @Param loginUserRequest body dto.LoginUserRequest true "Login User Request"
// @Success 200 {object} dto.LoginUserResponse{}
// @Router /api/v1/auth/login [post]
func (h *AuthenticationHandler) Login(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, msg := h.usecase.Login(req)
	if msg != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Register godoc
// @Summary Register
// @Description Register
// @Produce json
// @Tags authentication
// @Param registerUserRequest body dto.RegisterUserRequest true "Register User Request"
// @Success 200 {object} dto.RegisterUserResponse{}
// @Router /api/v1/auth/register [post]
func (h *AuthenticationHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, msg := h.usecase.RegisterUser(req)
	if msg != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, resp)
}
