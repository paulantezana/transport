package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paulantezana/transport/config"
	"github.com/paulantezana/transport/controller"
	"github.com/paulantezana/transport/utilities"
	"net/http"
)

// PublicApi public routes
func PublicApi(e *echo.Echo) {
	e.GET("/", func(context echo.Context) error {
		return context.NoContent(http.StatusOK)
	})
	pb := e.Group("/api/v1/public")

	// Company user
	pb.POST("/user/login", controller.Login)
	pb.POST("/user/forgot/search", controller.ForgotSearch)
	pb.POST("/user/forgot/validate", controller.ForgotValidate)
	pb.POST("/user/forgot/change", controller.ForgotChange)

	// Conductor user
	pb.POST("/mobile/login", controller.MobileLogin)
}

// ProtectedApi protected routes
func ProtectedApi(e *echo.Echo) {
	ar := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	con := middleware.JWTConfig{
		Claims:     &utilities.Claim{},
		SigningKey: []byte(config.GetConfig().Server.Key),
	}
	ar.Use(middleware.JWTWithConfig(con))

	// Crud user
	ar.POST("/user/all", controller.GetUsers)
	ar.POST("/user/by/id", controller.GetUserByID)
	ar.POST("/user/create", controller.CreateUser)
	ar.PUT("/user/update", controller.UpdateUser)
	ar.DELETE("/user/delete", controller.DeleteUser)
	ar.POST("/user/upload/avatar", controller.UploadAvatarUser)
	ar.POST("/user/reset/password", controller.ResetPasswordUser)
	ar.POST("/user/change/password", controller.ChangePasswordUser)

	// Crud mobile
	ar.POST("/mobile/all", controller.GetMobiles)
	ar.POST("/mobile/by/id", controller.GetMobileByID)
	ar.POST("/mobile/create", controller.CreateMobile)
	ar.PUT("/mobile/update", controller.UpdateMobile)
	ar.DELETE("/mobile/delete", controller.DeleteMobile)

	// Crud mobile
	ar.POST("/category/all", controller.GetCategoriesAll)
	ar.POST("/category/by/id", controller.GetCategoryByID)
	ar.POST("/category/create", controller.CreateCategory)
	ar.PUT("/category/update", controller.UpdateCategory)
	ar.DELETE("/category/delete", controller.DeleteCategory)

	// Crud company
	ar.POST("/company/all", controller.GetCompanies)
	ar.POST("/company/by/id", controller.GetCompanyByID)
	ar.POST("/company/create", controller.CreateCompany)
	ar.PUT("/company/update", controller.UpdateCompany)
	ar.DELETE("/company/delete", controller.DeleteCompany)

	// Crud vehicle
	ar.POST("/vehicle/paginate/by/company/id", controller.GetVehiclesPaginateByCompanyID)
	ar.POST("/vehicle/by/id", controller.GetVehicleByID)
	ar.POST("/vehicle/create", controller.CreateVehicle)
	ar.PUT("/vehicle/update", controller.UpdateVehicle)
	ar.DELETE("/vehicle/delete", controller.DeleteVehicle)

    // Crud vehicle
    ar.POST("/journey/all/by/company/id", controller.GetJourneys)
    ar.POST("/journey/by/id", controller.GetJourneyByID)
    ar.POST("/journey/create", controller.CreateJourney)
    ar.PUT("/journey/update", controller.UpdateJourney)
    ar.DELETE("/journey/delete", controller.DeleteJourney)

    // Crud vehicle
    ar.POST("/journeyDetail/all", controller.GetJourneyDetails)
    ar.POST("/journeyDetail/by/id", controller.GetJourneyDetailByID)
    ar.POST("/journeyDetail/create", controller.CreateJourneyDetail)
    ar.PUT("/journeyDetail/update", controller.UpdateJourneyDetail)
    ar.DELETE("/journeyDetail/delete", controller.DeleteJourneyDetail)

	// Global settings
	ar.POST("/setting/global", controller.GetGlobalSettings)
	ar.POST("/setting/global/mobile", controller.GetGlobalSettingsMobile)
	ar.GET("/setting", controller.GetSetting)
	ar.PUT("/setting", controller.UpdateSetting)
	ar.POST("/setting/upload/logo", controller.UploadLogoSetting)
	ar.GET("/setting/download/logo", controller.DownloadLogoSetting)
}
