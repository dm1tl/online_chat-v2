package chat

import (
	"app-websocket/internal/domain/room"
	response "app-websocket/internal/ports/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func dtoCreateRoom(c *gin.Context) (*room.CreateRoomReq, error) {
	var input *room.CreateRoomReq
	if err := c.BindJSON(&input); err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return nil, err
	}
	return input, nil
}
