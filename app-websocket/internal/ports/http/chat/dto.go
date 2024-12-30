package chat

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/domain/room"
	response "app-websocket/internal/ports/http"
	"net/http"
	"strconv"

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

type RoomsResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func dtoRoomsResp(rooms []domain.Room) []RoomsResp {
	var resp []RoomsResp
	for _, room := range rooms {
		respRoom := &RoomsResp{
			ID:   room.ID,
			Name: room.Name,
		}
		resp = append(resp, *respRoom)
	}
	return resp
}

func dtoAddClientReq(c *gin.Context) (*domain.Client, error) {
	roomID, err := strconv.ParseInt(c.Param("roomID"), 10, 64)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "incorrect room id")
		return nil, err
	}

	clientID, err := getUserId(c)
	if err != nil {
		logrus.Error("userId", err)
		response.NewErrorResponse(c, http.StatusBadRequest, "incorrect client id")
		return nil, err
	}

	username := c.Query("username")
	if username == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "username is required")
		return nil, err
	}

	password := c.Query("password")
	return &domain.Client{
		RoomID:   roomID,
		Username: username,
		ClientID: clientID,
		Password: password,
	}, nil
}
