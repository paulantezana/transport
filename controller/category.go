package controller

import (
    "fmt"
    "github.com/labstack/echo"
    "github.com/paulantezana/transport/config"
    "github.com/paulantezana/transport/models"
    "github.com/paulantezana/transport/utilities"
    "net/http"
)

func GetCategoriesPaginate(c echo.Context) error {
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
    categories := make([]models.Category, 0)

    // Find categories
    if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
        Order("id desc").
        Offset(offset).Limit(request.Limit).Find(&categories).
        Offset(-1).Limit(-1).Count(&total).
        Error; err != nil {
        return err
    }

    // Type response
    // 0 = all data
    // 1 = minimal data
    if request.Type == 1 {
        customCategories := make([]models.Category, 0)
        for _, category := range categories {
            customCategories = append(customCategories, models.Category{
                ID:   category.ID,
                Name: category.Name,
            })
        }
        return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
            Success:     true,
            Data:        customCategories,
            Total:       total,
            CurrentPage: request.CurrentPage,
        })
    }
    // Return response
    return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
        Success:     true,
        Data:        categories,
        Total:       total,
        CurrentPage: request.CurrentPage,
    })
}

func GetCategoriesAll(c echo.Context) error {
    // Get connection
    db := config.GetConnection()
    defer db.Close()

    // Check the number of matches
    categories := make([]models.Category, 0)

    // Find categories
    if err := db.Find(&categories).Error; err != nil  {
        return err
    }

    // Return response
    return c.JSON(http.StatusCreated, utilities.Response{
        Success:     true,
        Data:        categories,
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
    if !newCategory.State {
        if err := db.Model(newCategory).UpdateColumn("state", false).Error; err != nil {
            return err
        }
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
