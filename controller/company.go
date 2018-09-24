package controller

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/transport/config"
	"github.com/paulantezana/transport/models"
	"github.com/paulantezana/transport/utilities"
	"net/http"
)

func GetCompanies(c echo.Context) error {
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
	companies := make([]models.Company, 0)

	// Find companies
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&companies).
		Offset(-1).Limit(-1).Count(&total).
		Error; err != nil {
		return err
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customCompanies := make([]models.Company, 0)
		for _, company := range companies {
			customCompanies = append(customCompanies, models.Company{
				ID:   company.ID,
				Name: company.Name,
			})
		}
		return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
			Success:     true,
			Data:        customCompanies,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        companies,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetCompanyByID(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&company, company.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    company,
	})
}

func CreateCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// begin a transaction
	tx := db.Begin()

	// Insert company in database
	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// ==========================================================================
	// Create new user with company id
	user := models.User{}

	// Hash password
	cc := sha256.Sum256([]byte(company.Ruc))
	pwd := fmt.Sprintf("%x", cc)

	// Fill user data
	user.Password = pwd
	user.UserName = company.Ruc
	user.Profile = "company"
	user.CompanyID = company.ID

	// Insert user in database
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}
	tx.Commit()
	// ==========================================================================

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("La empresa %s se registro exitosamente", company.Name),
	})
}

func UpdateCompany(c echo.Context) error {
	// Get data request
	newCompany := models.Company{}
	if err := c.Bind(&newCompany); err != nil {
		return err
	}
	oldCompany := models.Company{
		ID: newCompany.ID,
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation company exist
	if db.First(&oldCompany).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", oldCompany.ID),
		})
	}

	// Update company in database
	if err := db.Model(&newCompany).Update(newCompany).Error; err != nil {
		return err
	}
	if !newCompany.State {
		if err := db.Model(newCompany).UpdateColumn("state", false).Error; err != nil {
			return err
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    newCompany.ID,
		Message: fmt.Sprintf("Los datos del usuario %s, se actualizarón correctamente", oldCompany.Name),
	})
}

func DeleteCompany(c echo.Context) error {
	// Get data request
	company := models.Company{}
	if err := c.Bind(&company); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation company exist
	if db.First(&company).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", company.ID),
		})
	}

	// Delete company in database
	if err := db.Delete(&company).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    company.ID,
		Message: fmt.Sprintf("El usuario %s, se elimino correctamente", company.Name),
	})
}
