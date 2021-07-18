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
	"github.com/mohamadsyukur99/fullstack/api/utils/formaterror"
)

// CreateSiswa ...
func (server *Server) CreateGuru(c *gin.Context) {

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

	guru := models.Guru{}

	err = json.Unmarshal(body, &guru)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	guruCreated, err := guru.SaveGuru(server.DB)
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
		"response": guruCreated,
	})
}

// GetUGuru ...
func (server *Server) GetGuruAll(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	guru := models.Guru{}

	dataguru, err := guru.FindAllGuru(server.DB)
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
		"response": dataguru,
	})
}

// GetGuru ...
func (server *Server) GetGuru(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	guruID := c.Param("id")

	uid, err := strconv.ParseUint(guruID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	guru := models.Guru{}

	guruGotten, err := guru.FindGuruByID(server.DB, uint32(uid))
	if err != nil {
		errList["No_siswa"] = "No siswa Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": guruGotten,
	})
}

// GetGuruByName ...
func (server *Server) GetGuruByName(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	userNama := c.Param("nama")

	guru := models.Guru{}

	guruGotten, err := guru.GetGuruByNames(server.DB, userNama)
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
		"response": guruGotten,
	})
}

// GetGuruNip ...
func (server *Server) GetGuruNip(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	userNama := c.Param("nama")

	guru := models.Guru{}

	siswaGotten, err := guru.GetGuruNip(server.DB, userNama)
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
		"response": siswaGotten,
	})
}

// UpdateGuru ...
func (server *Server) UpdateGuru(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	guruID := c.Param("id")
	// Check if the guru id is valid
	uid, err := strconv.ParseUint(guruID, 10, 32)
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
	formerGuru := models.Guru{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&formerGuru).Error
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	newGuru := models.Guru{}
	//update both the password and the email
	newGuru.Nip = formerGuru.Nip //remember, you cannot update the username
	newGuru.Nama = requestBody["nama"]
	newGuru.Tempat_lahir = requestBody["tempat_lahir"]
	newGuru.Tanggal_lahir = requestBody["tanggal_lahir"]
	newGuru.Alamat = requestBody["alamat"]
	newGuru.Agama = requestBody["agama"]
	newGuru.Jenis_kelamin = requestBody["jenis_kelamin"]

	updatedGuru, err := newGuru.UpdateAGuru(server.DB, uint32(uid))
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
		"response": updatedGuru,
	})
}

// DeleteSiswa ...
func (server *Server) DeleteGuru(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	var tokenID uint32
	guruID := c.Param("id")

	// Check if the user id is valid
	uid, err := strconv.ParseUint(guruID, 10, 32)
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

	if tokenID == 0 {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	guru := models.Guru{}
	_, err = guru.DeleteAGuru(server.DB, uint32(uid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "User deleted",
	})
}
