package adapters

import (
	"fmt"
	"os"
	"tg-crypto-tracker/internal/infrastructure/telegram"

	"github.com/gin-gonic/gin"
)

func AuthUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	botToken := os.Getenv("BOT_TOKEN")

	var body struct {
		InitData string `json:"initData"`
	}
	if err := c.BindJSON(&body); err != nil || body.InitData == "" {
		c.JSON(400, gin.H{"error": "initData required"})
		return
	}

	data, ok := telegram.VerifyInitData(body.InitData, botToken)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid initData"})
		return
	}

	fmt.Println(data)

	c.JSON(200, gin.H{
		"user_id": data["user"],
		"auth":    true,
	})
}
