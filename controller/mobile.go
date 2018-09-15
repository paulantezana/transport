package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/transport/config"
	"github.com/paulantezana/transport/models"
	"github.com/paulantezana/transport/utilities"
	"net/http"
)

type loginMobileResponse struct {
	Mobile models.Mobile `json:"mobile"`
	Token  interface{}   `json:"token"`
}

func MobileLogin(c echo.Context) error {
	// Get data request
	mobile := models.Mobile{}
	if err := c.Bind(&mobile); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validate mobile and email
	if db.Where("name = ? and key = ?", mobile.Name, mobile.Key).First(&mobile).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("El ID o Clave es incorecta"),
		})
	}

	// Check state mobile
	if !mobile.State {
		return c.NoContent(http.StatusForbidden)
	}

	// Prepare response data
	mobile.Key = ""

	// get token key
	token := utilities.GenerateJWT(mobile)

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Bienvenido al sistema %s", mobile.Name),
		Data: loginMobileResponse{
			Mobile: mobile,
			Token:  token,
		},
	})
}

func GetMobiles(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Pagination calculate
	if request.CurrentPage == 0 {
		request.CurrentPage = 1
	}
	offset := request.Limit*request.CurrentPage - request.Limit

	// Check the number of matches
	var total uint
	mobiles := make([]models.Mobile, 0)

	// Find mobiles
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&mobiles).
		Offset(-1).Limit(-1).Count(&total).
		Error; err != nil {
		return err
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customMobiles := make([]models.Mobile, 0)
		for _, mobile := range mobiles {
			customMobiles = append(customMobiles, models.Mobile{
				ID:   mobile.ID,
				Name: mobile.Name,
			})
		}
		return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
			Success:     true,
			Data:        customMobiles,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        mobiles,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetMobileByID(c echo.Context) error {
	// Get data request
	mobile := models.Mobile{}
	if err := c.Bind(&mobile); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&mobile, mobile.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    mobile,
	})
}

func CreateMobile(c echo.Context) error {
	// Get data request
	mobile := models.Mobile{}
	if err := c.Bind(&mobile); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert mobile in database
	if err := db.Create(&mobile).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    mobile.ID,
		Message: fmt.Sprintf("El usuario %s se registro exitosamente", mobile.Name),
	})
}

func UpdateMobile(c echo.Context) error {
	// Get data request
	newMobile := models.Mobile{}
	if err := c.Bind(&newMobile); err != nil {
		return err
	}
	oldMobile := models.Mobile{
		ID: newMobile.ID,
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation mobile exist
	if db.First(&oldMobile).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", oldMobile.ID),
		})
	}

	// Update mobile in database
	if err := db.Model(&newMobile).Update(newMobile).Error; err != nil {
		return err
	}
	if !newMobile.State {
		if err := db.Model(newMobile).UpdateColumn("state", false).Error; err != nil {
			return err
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    newMobile.ID,
		Message: fmt.Sprintf("Los datos del usuario %s, se actualizarón correctamente", oldMobile.Name),
	})
}

func DeleteMobile(c echo.Context) error {
	// Get data request
	mobile := models.Mobile{}
	if err := c.Bind(&mobile); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation mobile exist
	if db.First(&mobile).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", mobile.ID),
		})
	}

	// Delete mobile in database
	if err := db.Delete(&mobile).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    mobile.ID,
		Message: fmt.Sprintf("El usuario %s, se elimino correctamente", mobile.Name),
	})
}
