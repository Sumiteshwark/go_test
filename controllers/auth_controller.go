package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignupHandler(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	auth0Payload := map[string]string{
		"client_id":     os.Getenv("AUTH0_CLIENT_ID"),
		"audience":      os.Getenv("AUTH0_AUDIENCE"),
		"client_secret": os.Getenv("AUTH0_CLIENT_SECRET"),
		"email":         req.Email,
		"password":      req.Password,
		"connection":    "Username-Password-Authentication",
	}

	payloadBytes, _ := json.Marshal(auth0Payload)
	auth0URL := fmt.Sprintf("https://%s/dbconnections/signup", os.Getenv("AUTH0_DOMAIN"))

	resp, err := http.Post(auth0URL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach Auth0"})
		return
	}
	defer resp.Body.Close()

	var auth0Response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&auth0Response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Auth0 response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": auth0Response["error_description"]})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful! Please log in."})
}

func LoginHandler(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	auth0Payload := map[string]string{
		"grant_type":    "password",
		"client_id":     os.Getenv("AUTH0_CLIENT_ID"),
		"client_secret": os.Getenv("AUTH0_CLIENT_SECRET"),
		"username":      req.Email,
		"password":      req.Password,
		"audience":      os.Getenv("AUTH0_AUDIENCE"),
		"scope":         "openid profile email",
		"connection":    "Username-Password-Authentication",
	}

	payloadBytes, _ := json.Marshal(auth0Payload)
	auth0URL := fmt.Sprintf("https://%s/oauth/token", os.Getenv("AUTH0_DOMAIN"))

	resp, err := http.Post(auth0URL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach Auth0"})
		return
	}
	defer resp.Body.Close()

	var auth0Response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&auth0Response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Auth0 response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": auth0Response["error_description"]})
		return
	}

	c.JSON(http.StatusOK, auth0Response)
}

func GetValidateUserInfoHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	auth0UserInfoURL := fmt.Sprintf("https://%s/userinfo", os.Getenv("AUTH0_DOMAIN"))

	req, err := http.NewRequest("GET", auth0UserInfoURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach Auth0"})
		return
	}
	defer resp.Body.Close()
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Auth0 response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": userInfo["error_description"]})
		return
	}

	c.JSON(resp.StatusCode, userInfo)
}
