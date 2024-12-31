package chathandle

import (
	"app-websocket/internal/domain"
	"app-websocket/internal/domain/room"
	response "app-websocket/internal/ports/http"
	"app-websocket/internal/ports/ws"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type RoomManager interface {
	CreateRoom(ctx context.Context, req room.CreateRoomReq) (int64, error)
	GetRooms(ctx context.Context) ([]domain.Room, error)
}

type ServiceChatPusher interface {
	Subscribe(ctx context.Context, client *ws.ClientConnection) error
	PushMessage(ctx context.Context, msg *domain.Message) error
	Unsubscribe(ctx context.Context, client *ws.ClientConnection) error
}

type Handler struct {
	room   RoomManager
	pusher ServiceChatPusher
}

func NewHandler(room RoomManager, pusher ServiceChatPusher) *Handler {
	return &Handler{
		room:   room,
		pusher: pusher,
	}
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

	id, err := h.room.CreateRoom(ctx, *input)

	if err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "couldn't create room")
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("you succesfully created room № %d", id))
}

func (h *Handler) GetRooms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	rooms, err := h.room.GetRooms(ctx)
	if err != nil {
		logrus.Error(err)
		response.NewErrorResponse(c, http.StatusBadRequest, "couldn't get all rooms, try again")
		return
	}

	c.JSON(http.StatusOK, dtoRoomsResp(rooms))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cl, err := dtoAddClientReq(c)
	if err != nil {
		logrus.Error(err)
		return
	}

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error(err)
		wsConn.WriteJSON(response.ErrorResponse{
			Message: "couldn't estalish connection",
		})
		return
	}
	defer wsConn.Close()

	client := ws.NewClientConnection(wsConn, cl.Username, cl.Password, cl.ClientID, cl.RoomID, h.pusher)

	if err := h.pusher.Subscribe(ctx, client); err != nil {
		logrus.Error(err)
		wsConn.WriteJSON(response.ErrorResponse{
			Message: "couldn't join chat",
		})
	}

	go client.ReadMessage(ctx)
	client.WriteMessage()

}
