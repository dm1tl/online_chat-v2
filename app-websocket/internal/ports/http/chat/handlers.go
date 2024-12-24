package chat

import (
	"app-websocket/internal/domain/room"
	response "app-websocket/internal/ports/http"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RoomManager interface {
	CreateRoom(ctx context.Context, req room.CreateRoomReq) error
}

type Handler struct {
	room RoomManager
}

func (h *Handler) CreateRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	input, err := dtoCreateRoom(c)
	if err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "incorrect input data")
		return
	}
	if err := h.room.CreateRoom(ctx, *input); err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "couldn't create room")
		return
	}
	c.JSON(http.StatusOK, response.NewStatusResponse("you succesfully created a room"))
}
