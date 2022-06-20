package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

type OrderRequest struct {
	UserId          string      `json:"userId"`
	Name            string      `json:"name"`
	Items           []OrderItem `json:"items"`
	DataCheckString string      `json:"dataCheckString"`
	Hash            string      `json:"hash"`
}

type OrderItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func RunWeb(dev bool) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	if dev {
		r.Use(CORSMiddleware())
	}
	r.GET("/menu", func(c *gin.Context) {
		menu, err := GetMenu(Today())
		if err != nil || menu.Items == "" {
			c.JSON(404, gin.H{"error": "today's menu not found"})
			return
		}
		var items = make([]string, 0)
		if err := json.Unmarshal([]byte(menu.Items), &items); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		weekday := menu.DeliveryDate.Weekday()
		ruWeekday := LocalizeWeekday(weekday)
		title := "Меню на " + ruWeekday
		c.JSON(200, gin.H{
			"title": title,
			"items": items,
		})
	})

	r.POST("/order", func(c *gin.Context) {
		var order OrderRequest
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		log.Println(order)
		err := PostOrderInChat(order)
		if err != nil {
			log.Println(order)
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	r.Static("/static", "./frontend/build/static")
	r.StaticFile("/", "./frontend/build/index.html")
	r.StaticFile("/favicon.ico", "./frontend/build/favicon.ico")
	r.StaticFile("/manifest.json", "./frontend/build/manifest.json")
	r.StaticFile("/logo192.png", "./frontend/build/logo192.png")
	r.StaticFile("/logo512.png", "./frontend/build/logo512.png")

	r.Static("/rects-tool", "./rects-tool")

	if err := r.Run(); err != nil {
		panic(err)
	}
}
