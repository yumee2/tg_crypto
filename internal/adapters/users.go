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

	initData := c.Query("initData")
	if initData == "" {
		c.JSON(400, gin.H{"error": "initData required"})
		return
	}

	data, ok := telegram.VerifyInitData(initData, botToken)
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
