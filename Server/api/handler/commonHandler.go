package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/mskKandula/oes/api/middleware"
	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/util/websock"
)

func (h *Handler) Login(c *gin.Context) {

	userLogin := model.UserLogin{}

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, userType, err := h.CommonService.UserLogin(userLogin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenString, expiriesIn, err := middleware.GenerateJWT(userLogin, id, userType)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Path:    "/",
		Expires: expiriesIn,
	})

	c.JSON(http.StatusOK, gin.H{"userType": userType})
}

func (h *Handler) GetAllRoutes(c *gin.Context) {
	userId := c.GetInt("userId")

	routes, err := h.CommonService.GetRoutes(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (h *Handler) GetAllVideos(c *gin.Context) {
	videos, err := h.CommonService.GetVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

func (h *Handler) ServeWs(pool *websock.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket Endpoint Hit")

	var details websock.Details

	decoder := schema.NewDecoder()

	decoder.Decode(&details, r.URL.Query())
	// if err != nil {
	//     log.Fprintf(w, "%+v\n", err)
	// }

	conn, err := websock.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websock.Client{
		Conn:    conn,
		Pool:    pool,
		Details: &details,
	}

	go client.Read()

	pool.Register <- client
}

func (h *Handler) Logout(c *gin.Context) {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
		Path:   "/",
	})

}
