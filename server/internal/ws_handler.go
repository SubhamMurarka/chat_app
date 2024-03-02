package internal

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"

	"github.com/SubhamMurarka/chat_app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	WsService
}

func NewWsHandler(s WsService) *WsHandler {
	return &WsHandler{
		WsService: s,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WsHandler) WebSocketHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error upgrading")
		return
	}

	h.WsService.handleWebSocket(conn, c)
}

func (h *WsHandler) ImageUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors.New("resend image")})
		return
	}

	file := form.File["files"]

	// TODO can't interrupt process for single if array of images
	var urls []string

	username := c.Request.Header.Get("Username")
	roomid := c.Request.Header.Get("RoomID")

	for _, fileHeader := range file {
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()

		//checking extension of file
		extension := filepath.Ext(fileHeader.Filename)
		if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
			urls = append(urls, "Error: invalid file format")
			continue
		}

		//saving file to file storage temporarily and then to s3
		url, err := SaveFile(f, fileHeader)
		if err != nil {
			urls = append(urls, "Error:"+err.Error())
			continue
		}

		urls = append(urls, url)
	}

	m := &models.Message{
		Username: username,
		RoomID:   roomid,
		ImageUrl: urls,
	}

	var cl *models.Client
	h.WsService.SendMessage(c.Request.Context(), cl, *m)
}
