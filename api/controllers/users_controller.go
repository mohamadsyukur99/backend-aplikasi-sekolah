package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohamadsyukur99/fullstack/api/auth"
	"github.com/mohamadsyukur99/fullstack/api/models"
	"github.com/mohamadsyukur99/fullstack/api/security"
	"github.com/mohamadsyukur99/fullstack/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser ...
func (server *Server) CreateUser(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	user.Prepare()
	errorMessages := user.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": userCreated,
	})
}

// GetUsers ...
func (server *Server) GetUsers(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		errList["No_user"] = "No User Found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": users,
	})
}

// GetUser ...
func (server *Server) GetUser(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	userID := c.Param("id")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	user := models.User{}

	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		errList["No_user"] = "No User Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userGotten,
	})
}

// GetUserByName ...
func (server *Server) GetUserByName(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	userNama := c.Param("nama")

	user := models.User{}

	userGotten, err := user.GetUserByNames(server.DB, userNama)
	if err != nil {
		errList["code"] = "03"
		errList["message"] = "No User Found"
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusNotFound,
			"response": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userGotten,
	})
}

// GetUserUsername ...
func (server *Server) GetUserUsername(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	userNama := c.Param("nama")

	user := models.User{}

	userGotten, err := user.GetUserUsername(server.DB, userNama)
	if err != nil {
		errList["code"] = "03"
		errList["message"] = "No User Found"
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusNotFound,
			"response": errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userGotten,
	})
}

// UpdateUser ...
func (server *Server) UpdateUser(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	userID := c.Param("id")
	// Check if the user id is valid
	uid, err := strconv.ParseUint(userID, 10, 32)
	fmt.Println(uid)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Get user id from the token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If the id is not the authenticated user id
	if tokenID == 0 {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// Start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// Check for previous details
	formerUser := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&formerUser).Error
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	newUser := models.User{}
	//update both the password and the email
	newUser.Username = formerUser.Username //remember, you cannot update the username
	newUser.Email = requestBody["email"]
	newUser.Name = requestBody["name"]
	newUser.Level = requestBody["level"]
	newUser.Password = requestBody["password"]
	newUser.Prepare()
	errorMessages := newUser.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	if requestBody["password"] == "" {
		updatedUser, err := newUser.UpdateAUser(server.DB, uint32(uid))
		if err != nil {
			errList := formaterror.FormatError(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": updatedUser,
		})
	} else {
		updatedUser, err := newUser.UpdateAUserPassword(server.DB, uint32(uid))
		if err != nil {
			errList := formaterror.FormatError(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"response": updatedUser,
		})
	}
}

// UpdateUserOld ...
func (server *Server) UpdateUserOld(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	userID := c.Param("id")
	// Check if the user id is valid
	uid, err := strconv.ParseUint(userID, 10, 32)
	fmt.Println(uid)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Get user id from the token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If the id is not the authenticated user id
	if tokenID == 0 {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// Start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// Check for previous details
	formerUser := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&formerUser).Error
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	newUser := models.User{}

	//When current password has content.
	if requestBody["current_password"] == "" && requestBody["new_password"] != "" {
		errList["Empty_current"] = "Please Provide current password"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	if requestBody["current_password"] != "" && requestBody["new_password"] == "" {
		errList["Empty_new"] = "Please Provide new password"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	if requestBody["current_password"] != "" && requestBody["new_password"] != "" {
		//Also check if the new password
		if len(requestBody["new_password"]) < 6 {
			errList["Invalid_password"] = "Password should be atleast 6 characters"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}

		err = security.VerifyPassword(formerUser.Password, requestBody["current_password"])
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			errList["Password_mismatch"] = "The password not correct"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}
		//update both the password and the email
		newUser.Username = formerUser.Username //remember, you cannot update the username
		newUser.Email = requestBody["email"]
		newUser.Password = requestBody["new_password"]
	}
	//The password fields not entered, so update only the email
	newUser.Username = formerUser.Username
	newUser.Email = requestBody["email"]
	newUser.Name = requestBody["name"]

	newUser.Prepare()
	errorMessages := newUser.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	updatedUser, err := newUser.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedUser,
	})
}

// DeleteUser ...
func (server *Server) DeleteUser(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	var tokenID uint32
	userID := c.Param("id")

	// Check if the user id is valid
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Get user id from the token for valid tokens
	tokenID, err = auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If the id is not the authenticated user id
	// if tokenID != 0 && tokenID != uint32(uid) {
	// 	errList["Unauthorized"] = "Unauthorized"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }
	if tokenID == 0 {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	user := models.User{}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	// Also delete the posts, likes and the comments that this user created if any:
	post := models.Post{}
	_, err = post.DeleteUserPosts(server.DB, uint32(uid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "User deleted",
	})
}

func (Server *Server) Cek(c *gin.Context) {
	errList := map[string]string{}

	var tokenID uint32
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	if tokenID == 0 {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "oke",
	})
}
