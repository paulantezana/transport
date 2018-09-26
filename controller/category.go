package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/transport/config"
	"github.com/paulantezana/transport/models"
	"github.com/paulantezana/transport/utilities"
	"net/http"
)

func GetCategoriesAll(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Check the number of matches
	categories := make([]models.Category, 0)

	// Find categories
	if err := db.Order("id desc").Find(&categories).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    categories,
	})
}

func GetCategoryByID(c echo.Context) error {
	// Get data request
	category := models.Category{}
	if err := c.Bind(&category); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&category, category.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    category,
	})
}

func CreateCategory(c echo.Context) error {
	// Get data request
	category := models.Category{}
	if err := c.Bind(&category); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert category in database
	if err := db.Create(&category).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    category.ID,
		Message: fmt.Sprintf("El usuario %s se registro exitosamente", category.Name),
	})
}

func UpdateCategory(c echo.Context) error {
	// Get data request
	newCategory := models.Category{}
	if err := c.Bind(&newCategory); err != nil {
		return err
	}
	oldCategory := models.Category{
		ID: newCategory.ID,
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation category exist
	if db.First(&oldCategory).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", oldCategory.ID),
		})
	}

	// Update category in database
	if err := db.Model(&newCategory).Update(newCategory).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    newCategory.ID,
		Message: fmt.Sprintf("Los datos del usuario %s, se actualizarón correctamente", oldCategory.Name),
	})
}

func DeleteCategory(c echo.Context) error {
	// Get data request
	category := models.Category{}
	if err := c.Bind(&category); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation category exist
	if db.First(&category).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", category.ID),
		})
	}

	// Delete category in database
	if err := db.Delete(&category).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    category.ID,
		Message: fmt.Sprintf("El usuario %s, se elimino correctamente", category.Name),
	})
}
