package interfaces

import (
	"marluxgithub/muehle/pkg/muehle/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Application *application.Application
}

func (client *Client) Start() {
	router := gin.Default()

	// CORS Middleware - erlaubt alle Origins
	router.Use(client.CORS)

	client.generateRouting(router)
	router.Run(":8080")
}

func (client *Client) generateRouting(router *gin.Engine) {
	router.POST("/game", client.CreateGame)
	router.POST("/game/player", client.AddPlayer)
	router.POST("/game/move/put-stone", client.MovePutStone)
	router.POST("/game/move/stone", client.MoveStone)
	router.POST("/game/move/remove-stone", client.RemoveStone)
	router.GET("/game/state", client.GetGameState)
	router.GET("/game/current-player-color", client.GetCurrentPlayerColor)
	router.GET("/game/board", client.GetBoard)
}

func (client *Client) CreateGame(c *gin.Context) {
	client.Application.CreateGame()
	c.JSON(http.StatusOK, gin.H{"message": "Game created"})
}

func (client *Client) AddPlayer(c *gin.Context) {
	playerName := c.PostForm("playerName")

	secret, err := client.Application.AddPlayer(playerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player added", "secret": secret})
}

func (client *Client) MovePutStone(c *gin.Context) {
	fieldIndex := c.PostForm("fieldIndex")

	fieldIndexInt, err := strconv.Atoi(fieldIndex)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	secretCode := c.PostForm("secretCode")

	err = client.Application.MovePutStone(fieldIndexInt, secretCode)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Stone moved"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (client *Client) MoveStone(c *gin.Context) {
	fieldIndex := c.PostForm("fieldIndex")
	toFieldIndex := c.PostForm("toFieldIndex")

	fieldIndexInt, err := strconv.Atoi(fieldIndex)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	toFieldIndexInt, err := strconv.Atoi(toFieldIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	secretCode := c.PostForm("secretCode")

	err = client.Application.MoveStone(fieldIndexInt, toFieldIndexInt, secretCode)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Stone moved"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (client *Client) RemoveStone(c *gin.Context) {
	fieldIndex := c.PostForm("fieldIndex")

	fieldIndexInt, err := strconv.Atoi(fieldIndex)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	secretCode := c.PostForm("secretCode")

	err = client.Application.RemoveStone(fieldIndexInt, secretCode)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Stone moved"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (client *Client) GetGameState(c *gin.Context) {
	state := client.Application.GetGameState()
	c.JSON(http.StatusOK, gin.H{"state": state})
}

func (client *Client) GetCurrentPlayerColor(c *gin.Context) {
	color := client.Application.GetCurrentPlayerColor()
	c.JSON(http.StatusOK, gin.H{"color": color})
}

func (client *Client) GetBoard(c *gin.Context) {
	board := client.Application.GetBoard()
	c.JSON(http.StatusOK, gin.H{"board": board})
}

func (client *Client) CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

func NewClient() *Client {
	return &Client{
		Application: application.NewApplication(),
	}
}
