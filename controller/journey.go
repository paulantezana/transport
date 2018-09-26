package controller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/transport/config"
    "github.com/paulantezana/transport/models"
    "github.com/paulantezana/transport/utilities"
    "net/http"
)

type getJourneyRequest struct {
    CompanyID uint `json:"company_id"`
}

func GetJourneys(c echo.Context) error {
	// Get data request
	request := getJourneyRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Check the number of matches
	journeys := make([]models.Journey, 0)

	// Find journeys
	if err := db.Debug().Where("company_id = ?" ,request.CompanyID).
		Order("id desc").Find(&journeys).Error; err != nil {
		    return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success:     true,
		Data:        journeys,
	})
}

func GetJourneyByID(c echo.Context) error {
	// Get data request
	journey := models.Journey{}
	if err := c.Bind(&journey); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&journey, journey.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    journey,
	})
}

func CreateJourney(c echo.Context) error {
	// Get data request
	journey := models.Journey{}
	if err := c.Bind(&journey); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()


	// Insert journey in database
	if err := db.Create(&journey).Error; err != nil {
		db.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    journey.ID,
		Message: fmt.Sprintf("La empresa %s se registro exitosamente", journey.Name),
	})
}

func UpdateJourney(c echo.Context) error {
	// Get data request
	newJourney := models.Journey{}
	if err := c.Bind(&newJourney); err != nil {
		return err
	}
	oldJourney := models.Journey{
		ID: newJourney.ID,
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation journey exist
	if db.First(&oldJourney).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", oldJourney.ID),
		})
	}

	// Update journey in database
	if err := db.Model(&newJourney).Update(newJourney).Error; err != nil {
		return err
	}
	if !newJourney.State {
		if err := db.Model(newJourney).UpdateColumn("state", false).Error; err != nil {
			return err
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    newJourney.ID,
		Message: fmt.Sprintf("Los datos del usuario %s, se actualizarón correctamente", oldJourney.Name),
	})
}

func DeleteJourney(c echo.Context) error {
	// Get data request
	journey := models.Journey{}
	if err := c.Bind(&journey); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation journey exist
	if db.First(&journey).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", journey.ID),
		})
	}

	// Delete journey in database
	if err := db.Delete(&journey).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    journey.ID,
		Message: fmt.Sprintf("El usuario %s, se elimino correctamente", journey.Name),
	})
}
