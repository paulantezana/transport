package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/transport/config"
	"github.com/paulantezana/transport/models"
	"github.com/paulantezana/transport/utilities"
	"net/http"
)

func GetVehicles(c echo.Context) error {
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
	vehicles := make([]models.Vehicle, 0)

	// Find vehicles
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&vehicles).
		Offset(-1).Limit(-1).Count(&total).
		Error; err != nil {
		return err
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customVehicles := make([]models.Vehicle, 0)
		for _, vehicle := range vehicles {
			customVehicles = append(customVehicles, models.Vehicle{
				ID: vehicle.ID,
				//Name: vehicle.Name,
			})
		}
		return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
			Success:     true,
			Data:        customVehicles,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        vehicles,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetVehicleByID(c echo.Context) error {
	// Get data request
	vehicle := models.Vehicle{}
	if err := c.Bind(&vehicle); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&vehicle, vehicle.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    vehicle,
	})
}

func CreateVehicle(c echo.Context) error {
	// Get data request
	vehicle := models.Vehicle{}
	if err := c.Bind(&vehicle); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert vehicle in database
	if err := db.Create(&vehicle).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    vehicle.ID,
		Message: fmt.Sprintf("El usuario %s se registro exitosamente", vehicle.Name),
	})
}

func UpdateVehicle(c echo.Context) error {
	// Get data request
	newVehicle := models.Vehicle{}
	if err := c.Bind(&newVehicle); err != nil {
		return err
	}
	oldVehicle := models.Vehicle{
		ID: newVehicle.ID,
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation vehicle exist
	if db.First(&oldVehicle).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", oldVehicle.ID),
		})
	}

	// Update vehicle in database
	if err := db.Model(&newVehicle).Update(newVehicle).Error; err != nil {
		return err
	}
	if !newVehicle.State {
		if err := db.Model(newVehicle).UpdateColumn("state", false).Error; err != nil {
			return err
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    newVehicle.ID,
		Message: fmt.Sprintf("Los datos del usuario %s, se actualizarón correctamente", oldVehicle.Name),
	})
}

func DeleteVehicle(c echo.Context) error {
	// Get data request
	vehicle := models.Vehicle{}
	if err := c.Bind(&vehicle); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation vehicle exist
	if db.First(&vehicle).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", vehicle.ID),
		})
	}

	// Delete vehicle in database
	if err := db.Delete(&vehicle).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    vehicle.ID,
		Message: fmt.Sprintf("El usuario %s, se elimino correctamente", vehicle.Name),
	})
}
