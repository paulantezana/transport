package controller

import (
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/labstack/echo"
    "github.com/paulantezana/transport/config"
    "github.com/paulantezana/transport/models"
    "github.com/paulantezana/transport/utilities"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

func GetJourneyDetails(c echo.Context) error {
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

	// Check the number of matches
	journeyDetails := make([]models.JourneyDetail, 0)

	// Find journeyDetails
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").Find(&journeyDetails).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success:     true,
		Data:        journeyDetails,
	})
}

func GetJourneyDetailByID(c echo.Context) error {
	// Get data request
	journeyDetail := models.JourneyDetail{}
	if err := c.Bind(&journeyDetail); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&journeyDetail, journeyDetail.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    journeyDetail,
	})
}

func CreateJourneyDetail(c echo.Context) error {
	// Get data request
	journeyDetail := models.JourneyDetail{}
	if err := c.Bind(&journeyDetail); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert journeyDetail in database
	if err := db.Create(&journeyDetail).Error; err != nil {
		db.Rollback()
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    journeyDetail.ID,
		Message: fmt.Sprintf("La empresa %s se registro exitosamente", journeyDetail.Name),
	})
}

func UpdateJourneyDetail(c echo.Context) error {
	// Get data request
	newJourneyDetail := models.JourneyDetail{}
	if err := c.Bind(&newJourneyDetail); err != nil {
		return err
	}
	oldJourneyDetail := models.JourneyDetail{
		ID: newJourneyDetail.ID,
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation journeyDetail exist
	if db.First(&oldJourneyDetail).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", oldJourneyDetail.ID),
		})
	}

	// Update journeyDetail in database
	if err := db.Model(&newJourneyDetail).Update(newJourneyDetail).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    newJourneyDetail.ID,
		Message: fmt.Sprintf("Los datos del usuario %s, se actualizarón correctamente", oldJourneyDetail.Name),
	})
}

func DeleteJourneyDetail(c echo.Context) error {
	// Get data request
	journeyDetail := models.JourneyDetail{}
	if err := c.Bind(&journeyDetail); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation journeyDetail exist
	if db.First(&journeyDetail).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("No se encontró el registro con id %d", journeyDetail.ID),
		})
	}

	// Delete journeyDetail in database
	if err := db.Delete(&journeyDetail).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    journeyDetail.ID,
		Message: fmt.Sprintf("El usuario %s, se elimino correctamente", journeyDetail.Name),
	})
}


func SetTempUploadJourneyDetail(c echo.Context) error {
    // Source
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    // Destination
    auxDir := "temp/journeyDetail" + filepath.Ext(file.Filename)
    dst, err := os.Create(auxDir)
    if err != nil {
        return err
    }
    defer dst.Close()

    // Copy
    if _, err = io.Copy(dst, src); err != nil {
        return err
    }

    // ---------------------
    // Read File whit Excel
    // ---------------------
    xlsx, err := excelize.OpenFile(auxDir)
    if err != nil {
        return err
    }

    // Prepare
    journeyDetails := make([]models.JourneyDetail, 0)
    ignoreCols := 1

    // Get all the rows in the
    rows := xlsx.GetRows("Sheet1")
    for k, row := range rows {
        if k >= ignoreCols {
            sequence, _ := strconv.ParseUint(row[2],10,64)
            latitude, _ := strconv.ParseFloat(row[3], 64)
            longitude, _ := strconv.ParseFloat(row[4], 64)
            journeyID, _ := strconv.ParseUint(row[5],10, 64)
            journeyDetails = append(journeyDetails, models.JourneyDetail{
                Name:    strings.TrimSpace(row[1]),
                Sequence: uint(sequence),
                Latitude:  latitude,
                Longitude:   longitude,
                JourneyID: uint(journeyID),
            })
        }
    }

    // get connection
    db := config.GetConnection()
    defer db.Close()

    // Insert providers in database
    tr := db.Begin()
    for _, journeyDs := range journeyDetails {
        if err := tr.Create(&journeyDs).Error; err != nil {
            tr.Rollback()
            return c.JSON(http.StatusOK, utilities.Response{
                Success: false,
                Message: fmt.Sprintf("Ocurrió un error al insertar las rutas %s con "+
                    " es posible que este proveedor ya este en la base de datos o los datos son incorrectos, "+
                    "Error: %s, no se realizo ninguna cambio en la base de datos", journeyDs.Name, err),
            })
        }
    }
    tr.Commit()

    // Response success
    return c.JSON(http.StatusOK, utilities.Response{
        Success: true,
        Message: fmt.Sprintf("Se guardo %d registros den la base de datos", len(journeyDetails)),
    })
}


