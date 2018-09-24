package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/transport/api"
	"github.com/paulantezana/transport/config"
	"github.com/paulantezana/transport/models"
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
	api.SocketApi(e)

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
		&models.Category{},
		&models.Company{},
		&models.Pink{},
		&models.User{},
		&models.Mobile{},
		&models.Setting{},
		&models.Vehicle{},
		&models.VehicleAuthorized{},
		&models.Rute{},
	)

	db.Model(&models.Mobile{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.VehicleAuthorized{}).AddForeignKey("company_id", "companies(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.VehicleAuthorized{}).AddForeignKey("vehicle_id", "vehicles(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Pink{}).AddForeignKey("vehicle_id", "vehicles(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Company{}).AddForeignKey("category_id", "categories(id)", "RESTRICT", "RESTRICT")

	//db.Model(&models.Category{}).AddForeignKey("category_parent_id", "categories(id)", "RESTRICT", "RESTRICT")

	// -------------------------------------------------------------
	// INSERT FIST DATA --------------------------------------------
	// -------------------------------------------------------------
	usr := models.User{}
	db.First(&usr)
	// Validate
	if usr.ID == 0 {
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
		db.Create(&user)
	}

	// =======================
	// First Setting
	// =======================
	cg := models.Setting{}
	db.First(&cg)

	// Validate
	if cg.ID == 0 {
		co := models.Setting{
			Item:    10,
			Company: "TRANSPORT WEB",
			Logo:    "static/logo.png",
		}
		db.Create(&co)
	}
}
