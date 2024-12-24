package room

type CreateRoomReq struct {
	ID       int64  `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateRoomResp struct {
	Status string `json:"status"`
}

type GetRoomResp struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
