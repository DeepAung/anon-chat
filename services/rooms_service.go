package services

import (
	"net/http"
	"time"

	"github.com/DeepAung/anon-chat/utils"
)

type RoomsService struct{}

func NewRoomsService() *RoomsService {
	return &RoomsService{}
}

func (h *RoomsService) SetRoomName(w http.ResponseWriter, roomName string) {
	expires := time.Now().AddDate(0, 0, 1)
	utils.SetCookie(w, "roomName", roomName, expires, 86400)
}

func (h *RoomsService) SetRoomId(w http.ResponseWriter, roomId string) {
	expires := time.Now().AddDate(0, 0, 1)
	utils.SetCookie(w, "roomId", roomId, expires, 86400)
}

func (h *RoomsService) DeleteCookies(w http.ResponseWriter) {
	utils.DeleteCookie(w, "roomName")
	utils.DeleteCookie(w, "roomId")
}
