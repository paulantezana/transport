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
    pb := e.Group("/api/v1")

    pb.POST("/user/login", controller.Login)
    pb.POST("/user/forgot/search", controller.ForgotSearch)
    pb.POST("/user/forgot/validate", controller.ForgotValidate)
    pb.POST("/user/forgot/change", controller.ForgotChange)
}

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
    ar.POST("/user/byid", controller.GetUserByID)
    ar.POST("/user", controller.CreateUser)
    ar.PUT("/user", controller.UpdateUser)
    ar.DELETE("/user", controller.DeleteUser)
    ar.POST("/user/upload/avatar", controller.UploadAvatarUser)
    ar.POST("/user/reset/password", controller.ResetPasswordUser)
    ar.POST("/user/change/password", controller.ChangePasswordUser)

    // Global settings
    ar.POST("/setting/global", controller.GetGlobalSettings)
    ar.GET("/setting", controller.GetSetting)
    ar.PUT("/setting", controller.UpdateSetting)
    ar.POST("/setting/upload/logo", controller.UploadLogoSetting)
    ar.GET("/setting/download/logo", controller.DownloadLogoSetting)
}