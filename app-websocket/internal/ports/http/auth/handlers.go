package auth

import (
	"app-websocket/internal/domain/auth"
	response "app-websocket/internal/ports/http"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthManager interface {
	Create(ctx context.Context, req auth.CreateUserReq) error
	Login(ctx context.Context, req auth.LoginReq) (auth.LoginResp, error)
	Validate(ctx context.Context, req auth.ValidateTokenReq) (auth.ValidateTokenResp, error)
}

type Handler struct {
	auth AuthManager
}

func NewHandler(auth AuthManager) *Handler {
	return &Handler{
		auth: auth,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	var input auth.CreateUserReq
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := c.BindJSON(&input); err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	err := h.auth.Create(ctx, input)
	if err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusInternalServerError, "couldn't create an account, try again")
		return
	}
	c.JSON(http.StatusOK, response.NewStatusResponse("you succesfully signed up!"))
}

func (h *Handler) SignIn(c *gin.Context) {
	var input auth.LoginReq
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	token, err := h.auth.Login(ctx, input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
