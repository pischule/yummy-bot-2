package main

import (
	"encoding/json"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
		r.Use(cors.Default())
	} else {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"https://yummy.pischule.xyz", "https://y.pischule.xyz"}
		r.Use(cors.New(config))
	}
	r.GET("/menu", func(c *gin.Context) {
		menu, err := GetMenu()
		if err != nil || menu.Items == "" {
			log.Println("get menu failed", err)
			c.JSON(404, gin.H{"error": "today's menu not found"})
			return
		}
		var items = make([]string, 0)
		if err := json.Unmarshal([]byte(menu.Items), &items); err != nil {
			log.Println("menu items unmarshall failed")
			c.JSON(500, gin.H{"error": err})
			return
		}

		ruWeekday := weekdayToRussian[menu.DeliveryDate.Weekday()]
		title := "Меню на " + ruWeekday
		c.JSON(200, gin.H{
			"title": title,
			"items": items,
		})
	})

	r.POST("/order", func(c *gin.Context) {
		var order OrderRequest
		if err := c.ShouldBindJSON(&order); err != nil {
			log.Println("order request unmarshall failed", err)
			c.JSON(400, gin.H{"error": err})
			return
		}
		log.Println("post order", order)
		err := PostOrderInChat(order)
		if err != nil {
			log.Println("post order in chat failed", err)
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "ok"})
	})

	if err := r.Run(); err != nil {
		log.Println("webserver failed to start", err)
		panic(err)
	}
}
