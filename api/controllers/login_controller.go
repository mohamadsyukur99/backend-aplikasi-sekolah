package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamadsyukur99/fullstack/api/auth"
	"github.com/mohamadsyukur99/fullstack/api/models"
	"github.com/mohamadsyukur99/fullstack/api/security"
	"github.com/mohamadsyukur99/fullstack/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

// Login ...
func (server *Server) Login(c *gin.Context) {
	//clear previous error if any
	// errList = map[string]string{}
	// var data []string
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":      http.StatusUnprocessableEntity,
			"first error": "Unable to get request",
		})
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Cannot unmarshal body",
		})
		return
	}
	user.Prepare()
	errorMessages := user.Validate("login")
	if len(errorMessages) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errorMessages,
		})
		return
	}

	userData, err := server.SignIn(user.Username, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  formattedError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userData,
	})
}

// SignIn ...
func (server *Server) SignIn(username, password string) (map[string]interface{}, error) {
	var err error

	userData := make(map[string]interface{})

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		userData["code"] = "01"
		userData["message"] = "Email Tidak ditemukan"
		return userData, nil
	}

	err = security.VerifyPassword(user.Password, password)
	fmt.Println(err)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		userData["code"] = "02"
		userData["message"] = "Password Salah"
		return userData, nil
	}

	token, err := auth.CreateToken(user.ID)
	if err != nil {
		fmt.Println("this is the error creating the token: ", err)
		return nil, err
	}
	userData["code"] = "00"
	userData["token"] = token
	userData["id"] = user.ID
	userData["email"] = user.Email
	userData["username"] = user.Username
	userData["name"] = user.Name

	return userData, nil
}
