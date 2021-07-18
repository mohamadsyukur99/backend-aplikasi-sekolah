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
func (server *Server) CreateSiswa(c *gin.Context) {

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

	siswa := models.Siswa{}

	err = json.Unmarshal(body, &siswa)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	siswaCreated, err := siswa.SaveSiswa(server.DB)
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
		"response": siswaCreated,
	})
}

// GetUSiswa ...
func (server *Server) GetSiswaAll(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	siswa := models.Siswa{}

	datasiswa, err := siswa.FindAllSiswa(server.DB)
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
		"response": datasiswa,
	})
}

// GetSiswa ...
func (server *Server) GetSiswa(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	siswaID := c.Param("id")

	uid, err := strconv.ParseUint(siswaID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	siswa := models.Siswa{}

	siswaGotten, err := siswa.FindSiswaByID(server.DB, uint32(uid))
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
		"response": siswaGotten,
	})
}

// GetSiswaByName ...
func (server *Server) GetSiswaByName(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	userNama := c.Param("nama")

	siswa := models.Siswa{}

	siswaGotten, err := siswa.GetSiswaByNames(server.DB, userNama)
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

// GetSiswaNoInduk ...
func (server *Server) GetSiswaNoInduk(c *gin.Context) {

	//clear previous error if any
	errList := map[string]string{}

	userNama := c.Param("nama")

	siswa := models.Siswa{}

	siswaGotten, err := siswa.GetSiswaNoInduk(server.DB, userNama)
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

// UpdateSiswa ...
func (server *Server) UpdateSiswa(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	siswaID := c.Param("id")
	// Check if the siswa id is valid
	uid, err := strconv.ParseUint(siswaID, 10, 32)
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
	formerSiswa := models.Siswa{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&formerSiswa).Error
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	newSiswa := models.Siswa{}
	//update both the password and the email
	newSiswa.No_induk = formerSiswa.No_induk //remember, you cannot update the username
	newSiswa.Nama = requestBody["nama"]
	newSiswa.Tempat_lahir = requestBody["tempat_lahir"]
	newSiswa.Tanggal_lahir = requestBody["tanggal_lahir"]
	newSiswa.Nama_wali = requestBody["nama_wali"]
	newSiswa.Alamat = requestBody["alamat"]
	newSiswa.Jenis_kelamin = requestBody["jenis_kelamin"]
	newSiswa.Kelas = requestBody["kelas"]
	newSiswa.Agama = requestBody["agama"]

	updatedUser, err := newSiswa.UpdateASiswa(server.DB, uint32(uid))
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

// DeleteSiswa ...
func (server *Server) DeleteSiswa(c *gin.Context) {
	//clear previous error if any
	errList := map[string]string{}

	var tokenID uint32
	siswaID := c.Param("id")

	// Check if the user id is valid
	uid, err := strconv.ParseUint(siswaID, 10, 32)
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

	siswa := models.Siswa{}
	_, err = siswa.DeleteASiswa(server.DB, uint32(uid))
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
