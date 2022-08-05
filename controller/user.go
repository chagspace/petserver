package controller

import (
	"net/http"

	"github.com/chagspace/petserver/common"
	"github.com/chagspace/petserver/database"
	"github.com/chagspace/petserver/model"
	"github.com/chagspace/petserver/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"message": "get_users",
	})
}
func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"message": "get_user",
	})
}

func CreateUser(c *gin.Context) {
	user := &model.UserModel{}
	c.BindJSON(&user)

	// check if username exists
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   1,
			"msg":    "username is required",
			"status": "error",
		})
		return
	}

	// check if username is exists in database
	database_user, exist_user := service.GetUser(user.Username)
	if exist_user && database_user.Username == user.Username {
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"msg":    "username already exists",
			"status": "error",
		})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(password)

	service.CreateUser(user)

	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"msg":      "success",
		"username": user.Username,
		"userId":   user.ID,
		"email":    user.Email,
		"uid":      user.UID,
	})
}

func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"message": "update_user",
	})
}
func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"message": "delete_user",
	})
}

// subscribe a user
func SubscribeUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"message": "subscribe_user",
	})
}
func UnsubscribeUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"message": "unsubscribe_user",
	})
}

// notify  a user
func NotifyUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "success",
		"message": "notify_user",
	})
}

// Login user
func Login(c *gin.Context) {
	var user model.UserModel

	// try parser to json
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	// note:
	// because of the uniqueness of the user name, the password of the first user found is hashed with the current password
	database_user, allowed_user := service.GetUser(user.Username)
	if !allowed_user {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "unauthorized"})
		return
	}
	password_error := bcrypt.CompareHashAndPassword([]byte(database_user.Password), []byte(user.Password))
	if password_error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "password error"})
		return
	}

	// generate tokens (JWT)
	token, err := common.CreateToken(uint(database_user.UID), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	// set token to cookies and enable httpOnly
	common.SetHttpOnlyCookie(c, "access_token", token, 3600)
	// set token to redis
	database.GlobalRedis.Set(token, database_user.UID, 3600)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "username": user.Username, "uid": database_user.UID})
}

func Logout(c *gin.Context) {}
