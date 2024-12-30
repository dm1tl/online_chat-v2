package client

type AddClientReq struct {
	RoomID   int64
	ClientID int64
	Username string
	Password string
}
