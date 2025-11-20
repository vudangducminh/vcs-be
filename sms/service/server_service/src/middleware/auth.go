package middleware

import (
	"log"
	"net/http"
	"server_service/src/algorithm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthAddServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Invalid JWT token:", err)
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT token is not valid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT claims are not of type jwt.MapClaims")
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Username claim not found or not a string")
			c.Abort()
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Role claim not found or not a string")
			c.Abort()
			return
		}
		str := ""
		admin := false
		addServer := false
		for _, r := range role {
			if r == ',' {
				if str == "admin" {
					admin = true
				}
				if str == "add-server" {
					addServer = true
				}
				str = ""
			} else {
				str += string(r)
			}
		}
		if str == "admin" {
			admin = true
		}
		if str == "add-server" {
			addServer = true
		}
		if !admin && !addServer {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		// Store username and role in context for further handlers
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

func AuthDeleteServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Invalid JWT token:", err)
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT token is not valid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT claims are not of type jwt.MapClaims")
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Username claim not found or not a string")
			c.Abort()
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Role claim not found or not a string")
			c.Abort()
			return
		}
		str := ""
		admin := false
		deleteServer := false
		for _, r := range role {
			if r == ',' {
				if str == "admin" {
					admin = true
				}
				if str == "delete-server" {
					deleteServer = true
				}
				str = ""
			} else {
				str += string(r)
			}
		}
		if str == "admin" {
			admin = true
		}
		if str == "delete-server" {
			deleteServer = true
		}
		if !admin && !deleteServer {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		// Store username and role in context for further handlers
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

func AuthUpdateServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Invalid JWT token:", err)
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT token is not valid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT claims are not of type jwt.MapClaims")
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Username claim not found or not a string")
			c.Abort()
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Role claim not found or not a string")
			c.Abort()
			return
		}
		str := ""
		admin := false
		updateServer := false
		for _, r := range role {
			if r == ',' {
				if str == "admin" {
					admin = true
				}
				if str == "update-server" {
					updateServer = true
				}
				str = ""
			} else {
				str += string(r)
			}
		}
		if str == "admin" {
			admin = true
		}
		if str == "update-server" {
			updateServer = true
		}
		if !admin && !updateServer {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		// Store username and role in context for further handlers
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

func AuthViewServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Invalid JWT token:", err)
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT token is not valid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT claims are not of type jwt.MapClaims")
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Username claim not found or not a string")
			c.Abort()
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Role claim not found or not a string")
			c.Abort()
			return
		}
		str := ""
		admin := false
		viewServer := false
		for _, r := range role {
			if r == ',' {
				if str == "admin" {
					admin = true
				}
				if str == "view-server" {
					viewServer = true
				}
				str = ""
			} else {
				str += string(r)
			}
		}
		if str == "admin" {
			admin = true
		}
		if str == "view-server" {
			viewServer = true
		}
		if !admin && !viewServer {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		// Store username and role in context for further handlers
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

func AuthImportExcel() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Invalid JWT token:", err)
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT token is not valid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT claims are not of type jwt.MapClaims")
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Username claim not found or not a string")
			c.Abort()
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Role claim not found or not a string")
			c.Abort()
			return
		}
		str := ""
		admin := false
		importExcel := false
		for _, r := range role {
			if r == ',' {
				if str == "admin" {
					admin = true
				}
				if str == "import-excel" {
					importExcel = true
				}
				str = ""
			} else {
				str += string(r)
			}
		}
		if str == "admin" {
			admin = true
		}
		if str == "import-excel" {
			importExcel = true
		}
		if !admin && !importExcel {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		// Store username and role in context for further handlers
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}

func AuthExportExcel() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Invalid JWT token:", err)
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT token is not valid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT claims are not of type jwt.MapClaims")
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Username claim not found or not a string")
			c.Abort()
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Role claim not found or not a string")
			c.Abort()
			return
		}
		str := ""
		admin := false
		exportExcel := false
		for _, r := range role {
			if r == ',' {
				if str == "admin" {
					admin = true
				}
				if str == "export-excel" {
					exportExcel = true
				}
				str = ""
			} else {
				str += string(r)
			}
		}
		if str == "admin" {
			admin = true
		}
		if str == "export-excel" {
			exportExcel = true
		}
		if !admin && !exportExcel {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		// Store username and role in context for further handlers
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}
