package main

import (
	_ "errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var users = []User{
	{ID: 1, Username: "user", Password: "$2a$10$gdQmi2PzjqOpvT.6NKfoO.7sBslirVc8DKW8b9R7iMeaFKkcdMLFW", Role: "user"},
	{ID: 2, Username: "admin", Password: "$2a$10$VPXhVGfVRTxR4I6QLqwD4OAw4W/Py/cLnfI/SK2IYZi.olHR0rLoC", Role: "admin"},
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func main() {
	r := gin.Default()
	// Custom CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	r.POST("/register", Register)
	r.POST("/login", Login)

	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/profile", GetProfile)
		authorized.OPTIONS("/profile", PreflightHandler)
	}

	admin := authorized.Group("/admin")
	admin.Use(AdminOnly())
	{
		admin.GET("/users", GetAllUsers)
		authorized.OPTIONS("/users", PreflightHandler)
		admin.GET("/warehousemanagers", GetAllWarehouseManagers)
		authorized.OPTIONS("/warehousemanagers", PreflightHandler)
	}

	r.Run(":8080")
}

func PreflightHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Status(http.StatusOK)
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	fmt.Print("Password, ", hashedPassword)

	user.Password = string(hashedPassword)
	user.Role = "user"
	user.ID = uint(len(users) + 1)
	users = append(users, user)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, u := range users {
		if u.Username == user.Username {
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
			if err == nil {
				token, err := GenerateToken(u)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"token": token, "role": u.Role})
				return
			}
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

func GenerateToken(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetProfile(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"message": "Profile accessed", "user": username})
}

func GetAllUsers(c *gin.Context) {
	log.Printf("GetAllUsers called. Number of users: %d", len(users))
	c.JSON(http.StatusOK, users)
}

type WareHouseManager struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

var wareHouseManagers = []WareHouseManager{
	{ID: 1, Name: "john doe", Contact: "+44123456789"},
	{ID: 2, Name: "jane doe", Contact: "+44987654321"},
}

func GetAllWarehouseManagers(c *gin.Context) {
	log.Printf("GetAllWarehouseManagers called. Number of warehousemanagers: %d", len(wareHouseManagers))
	c.JSON(http.StatusOK, wareHouseManagers)
}
