package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize migration database
	migration()

	// COR
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
	}))

	// Static Files =========================================================================
	static := e.Group("/static")
	static.Static("", "static")

	// API
	api.PublicApi(e)
	api.ProtectedApi(e)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = config.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}


// migration Init migration database
func migration() {
	db := config.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.User{},
		&models.Mobile{},
	)

	// -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	usr := models.User{}
	db.First(&usr)
	// hash password
	cc := sha256.Sum256([]byte("admin"))
	pwd := fmt.Sprintf("%x", cc)
	// create model
	user := models.User{
		UserName: "admin",
		Password: pwd,
		Profile:  "admin",
		Email:    "yoel.antezana@gmail.com",
	}
	// insert database
	if usr.ID == 0 {
		db.Create(&user)
	}

	// First Setting
	cg := models.Setting{}
	db.First(&cg)
	co := models.Setting{
		Item:       10,
		Company:    "REQUIREMENT WEB",
		Quotations: 3,
		Logo:       "static/logo.png",
	}
	// Insert database
	if cg.ID == 0 {
		db.Create(&co)
	}
}

