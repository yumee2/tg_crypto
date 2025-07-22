package adapters

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
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
	fmt.Println("Raw initData:", body.InitData)

	expIn := 24 * time.Hour

	err := initdata.Validate(body.InitData, botToken, expIn)

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid initData"})
		return
	}

	// data, ok := telegram.VerifyInitData(body.InitData, botToken)
	// if !ok {
	// 	c.JSON(401, gin.H{"error": "Invalid initData"})
	// 	return
	// }

	// fmt.Println(data)

	c.JSON(200, gin.H{
		"success": true,
	})
}
