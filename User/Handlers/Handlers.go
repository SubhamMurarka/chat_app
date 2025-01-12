package Handler

import (
	"log"
	"net/http"
	"strconv"

	models "github.com/SubhamMurarka/chat_app/User/Models"
	"github.com/SubhamMurarka/chat_app/User/Services"
	util "github.com/SubhamMurarka/chat_app/User/Util"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userSvc Services.UserService
	roomSvc Services.RoomService
	msgSvc  Services.MsgService
}

func NewUserHandler(s Services.UserService, r Services.RoomService, m Services.MsgService) *Handler {
	return &Handler{
		userSvc: s,
		roomSvc: r,
		msgSvc:  m,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u models.CreateUserReq

	if err := c.ShouldBindJSON(&u); err != nil {
		log.Printf("Error binding : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": util.ErrInvalidRequest})
		return
	}

	res, err := h.userSvc.CreateUser(c.Request.Context(), &u)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
	var user models.LoginUserReq

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": util.ErrInvalidRequest})
		return
	}

	u, err := h.userSvc.Login(c.Request.Context(), &user)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, u)
}

func (h *Handler) JoinRoom(c *gin.Context) {
	var room models.Room

	if err := c.ShouldBindJSON(&room); err != nil {
		log.Printf("Error binding : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": util.ErrInvalidRequest})
		return
	}

	// channelid, err := h.roomSvc.IsRoom(c, room.RoomName)
	// if err != nil {
	// 	respondWithError(c, err)
	// 	return
	// }

	userid := c.GetString("userid")
	id, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Printf("Error converting to int64 : %v", err)
		respondWithError(c, util.ErrInternal)
		return
	}

	err = h.roomSvc.AddUserToRoom(c, id, room.RoomID)
	if err != nil && err != util.ErrAlreadyJoined {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": "room joined"})
}

func (h *Handler) GetAllUserRoom(c *gin.Context) {
	userid := c.GetString("userid")
	id, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		log.Printf("Error converting to int64 : %v", err)
		respondWithError(c, util.ErrInternal)
		return
	}

	rooms, err := h.roomSvc.GetAllUserRoom(c, id)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var room models.Room

	if err := c.ShouldBindJSON(&room); err != nil {
		log.Printf("Error binding : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": util.ErrInvalidRequest})
		return
	}
	id, err := h.roomSvc.AddRoom(c, room.RoomName, room.Typ)
	if err != nil {
		respondWithError(c, err)
		return
	}

	roomres := &models.Room{
		RoomName: room.RoomName,
		RoomID:   id,
		Typ:      room.Typ,
	}

	if err != nil {
		respondWithError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"New Room": roomres})
}

func (h *Handler) GetAllMembers(c *gin.Context) {
	room := c.Query("room_id")

	id, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		log.Printf("Error converting to int64 : %v", err)
		respondWithError(c, util.ErrInternal)
		return
	}

	users, err := h.roomSvc.GetAllMembers(c, id)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetMessages(c *gin.Context) {

	lastid := c.DefaultQuery("lastid", "")

	var member models.Members

	var msg []models.Message

	if err := c.ShouldBindJSON(&member); err != nil {
		log.Printf("Error binding : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": util.ErrInvalidRequest})
		return
	}

	is := h.roomSvc.IsUSerMember(c, member.UserID, member.ChannelID)
	if !is {
		err := util.ErrInvalidRequest
		respondWithError(c, err)
		return
	}

	if lastid == "" {
		var err error
		msg, err = h.msgSvc.GetNewMessages(c, member.ChannelID)
		if err != nil {
			respondWithError(c, err)
			return
		}
	} else {
		var err error
		lid, _ := strconv.ParseUint(lastid, 10, 64)
		msg, err = h.msgSvc.PaginationMessages(c, member.ChannelID, lid)
		if err != nil {
			respondWithError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, msg)
}

func respondWithError(c *gin.Context, err error) {
	switch err {
	case util.ErrEmailExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case util.ErrUserEmailData, util.ErrRoomExists, util.ErrRoomExist:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case util.ErrBcrypt:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case util.ErrCredentials:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case util.ErrInvalidToken:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case util.ErrExpiredToken:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case util.ErrInvalidRequest:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		log.Printf("Default Error : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": util.ErrInternal.Error()})
	}
}
