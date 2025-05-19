package controllers

import (
	"backend/services"
	"backend/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GetSystemHealthUpdate(ctx *gin.Context) {
	healthyCount, unhealthyCount, err := services.GetSystemHealthUpdate()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"healthy_count":   healthyCount,
		"unhealthy_count": unhealthyCount,
	})
}

func GetSystemUsage(ctx *gin.Context) {
	storageUsage, memoryUsage, averageCPUUsage, err := services.GetSystemUsage()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"storage_usage": storageUsage,
		"memory_usage":  memoryUsage,
		"cpu_usage":     averageCPUUsage,
	})
}

// WebSocket upgrader with proper CORS handling
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler for system usage monitoring
func GetSystemUsageWS(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("WebSocket connection established")

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("WebSocket closed with code %d: %s\n", code, text)
		return nil
	})

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			storageUsage, memoryUsage, averageCPUUsage, err := services.GetSystemUsage()
			if err != nil {
				fmt.Println("Error getting system usage:", err)
				continue
			}
			data := gin.H{
				"storage_usage": storageUsage,
				"memory_usage":  memoryUsage,
				"cpu_usage":     averageCPUUsage,
			}

			if err := conn.WriteJSON(data); err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					fmt.Println("Client disconnected")
					return
				}
				fmt.Println("WebSocket send error:", err)
				return
			}
		}
	}
}
