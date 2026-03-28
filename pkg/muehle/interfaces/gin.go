package interfaces

import (
	"io/fs"
	"log"
	"marluxgithub/muehle/pkg/muehle/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTP-API (öffentliche Routen): Kanonische Beschreibung in
// specs/001-rest-routing-cleanup/contracts/http-api.md
type Client struct {
	Registry *application.GameRegistry
}

func (client *Client) Start() {
	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("muehle: SetTrustedProxies: %v", err)
	}

	router.Use(client.CORS)

	client.generateRouting(router)
	listen := "0.0.0.0:40000"
	log.Printf("muehle: lausche auf http://%s — Diagnose: GET /health, GET /openapi.yaml, GET /swagger", listen)
	if err := router.Run(listen); err != nil {
		log.Fatalf("muehle: %v", err)
	}
}

func (client *Client) generateRouting(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml", openAPISpec)
	})

	swaggerFS, err := fs.Sub(swaggerUIStatic, "swaggerui")
	if err != nil {
		log.Printf("muehle: Swagger-UI nicht gemountet (embed/swaggerui): %v", err)
	} else {
		indexHTML, err := fs.ReadFile(swaggerFS, "index.html")
		if err != nil {
			log.Printf("muehle: Swagger index.html: %v", err)
		} else {
			// Kein Redirect (VPN/Proxy): HTML direkt unter /swagger.
			// Assets unter /swagger-static — vermeidet Gin-1.11-Konflikt /swagger + /swagger/*filepath.
			router.GET("/swagger", func(c *gin.Context) {
				c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
			})
		}
		router.StaticFS("/swagger-static", http.FS(swaggerFS))
	}

	router.POST("/games", client.postGames)

	games := router.Group("/games/:gameId")
	games.GET("/players", client.listGamePlayers)
	games.POST("/players", client.postGamePlayers)
	games.POST("/moves", client.postGameMoves)
	games.GET("/state", client.getGameState)
	games.GET("/current-player", client.getGameCurrentPlayer)
	games.GET("/board", client.getGameBoard)
}

func (client *Client) resolveGame(c *gin.Context) (*application.Application, bool) {
	app, ok := client.Registry.Get(c.Param("gameId"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return nil, false
	}
	return app, true
}

func (client *Client) postGames(c *gin.Context) {
	id, err := client.Registry.CreateGame()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Game created", "id": id})
}

func (client *Client) listGamePlayers(c *gin.Context) {
	app, ok := client.resolveGame(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"players": app.GetPlayersPublic()})
}

func (client *Client) postGamePlayers(c *gin.Context) {
	app, ok := client.resolveGame(c)
	if !ok {
		return
	}
	playerName := c.PostForm("playerName")
	secret, color, err := app.AddPlayer(playerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Player added", "secret": secret, "color": color})
}

func (client *Client) postGameMoves(c *gin.Context) {
	app, ok := client.resolveGame(c)
	if !ok {
		return
	}
	action := c.PostForm("action")
	secretCode := c.PostForm("secretCode")

	switch action {
	case "place":
		fieldIndex := c.PostForm("fieldIndex")
		if fieldIndex == "" || secretCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fieldIndex and secretCode required"})
			return
		}
		fieldIndexInt, err := strconv.Atoi(fieldIndex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fieldIndex"})
			return
		}
		if err := app.MovePutStone(fieldIndexInt, secretCode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Stone moved"})
	case "move":
		fieldIndex := c.PostForm("fieldIndex")
		toFieldIndex := c.PostForm("toFieldIndex")
		if fieldIndex == "" || toFieldIndex == "" || secretCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fieldIndex, toFieldIndex and secretCode required"})
			return
		}
		from, err := strconv.Atoi(fieldIndex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fieldIndex"})
			return
		}
		to, err := strconv.Atoi(toFieldIndex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid toFieldIndex"})
			return
		}
		if err := app.MoveStone(from, to, secretCode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Stone moved"})
	case "remove":
		fieldIndex := c.PostForm("fieldIndex")
		if fieldIndex == "" || secretCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fieldIndex and secretCode required"})
			return
		}
		idx, err := strconv.Atoi(fieldIndex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fieldIndex"})
			return
		}
		if err := app.RemoveStone(idx, secretCode); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Stone moved"})
	default:
		if action == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "action required (place, move, remove)"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action"})
		}
	}
}

func (client *Client) getGameState(c *gin.Context) {
	app, ok := client.resolveGame(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"state": app.GetGameState()})
}

func (client *Client) getGameCurrentPlayer(c *gin.Context) {
	app, ok := client.resolveGame(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"color": app.GetCurrentPlayerColor()})
}

func (client *Client) getGameBoard(c *gin.Context) {
	app, ok := client.resolveGame(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"board": app.GetBoard()})
}

func (client *Client) CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func NewClient() *Client {
	return &Client{
		Registry: application.NewGameRegistry(),
	}
}
