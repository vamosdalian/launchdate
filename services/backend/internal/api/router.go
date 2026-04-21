package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/middleware"
)

// SetupRouter sets up the API routes
func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(handler.IsProduction))
	router.Use(middleware.Logger(handler.logger))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "https://launch-date.com")
	})

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/health", handler.Health)
		apiV1.GET("/page-backgrounds", handler.GetPublicPageBackgrounds)
		apiV1.GET("/launch", handler.GetPublicLaunches)
		apiV1.GET("/launch/:id", handler.GetPublicLaunchByID)
		apiV1.GET("/rocket", handler.GetPublicRockets)
		apiV1.GET("/rocket/:id", handler.GetPublicRocketByID)
		apiV1.GET("/companies", handler.GetPublicCompanies)
		apiV1.GET("/companies/:id", handler.GetPublicCompanyByID)
		apiV1.GET("/launch-bases", handler.GetPublicLaunchBases)
		apiV1.GET("/launch-bases/:id", handler.GetPublicLaunchBaseByID)
		apiV1.POST("/subscriptions", handler.Subscribe)
		apiV1.GET("/subscriptions/unsubscribe", handler.Unsubscribe)

		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", handler.authHandler.Login)
			auth.POST("/refresh", handler.authHandler.Refresh)
			auth.POST("/logout", handler.authHandler.Logout)
			auth.POST("/bootstrap", handler.authHandler.Bootstrap)
			auth.GET("/me", handler.authHandler.Me)
		}

		data := apiV1.Group("/data")
		data.Use(middleware.AuthMiddleware(handler.jwtM, handler.logger))
		{
			data.GET("/stats", handler.GetStats)
			data.GET("/page-backgrounds", handler.GetPageBackgrounds)
			data.PUT("/page-backgrounds/:pageKey", handler.UpdatePageBackground)
			data.GET("/rockets", handler.GetRockets)
			data.GET("/rockets/:id", handler.GetRocketByID)
			data.PUT("/rockets/:id", handler.UpdateRocket)
			data.GET("/agencies", handler.GetAgencies)
			data.GET("/agencies/:id", handler.GetAgencyByID)
			data.PUT("/agencies/:id", handler.UpdateAgency)
			data.GET("/launchbases", handler.GetLaunchBases)
			data.GET("/launchbases/:id", handler.GetLaunchBaseByID)
			data.GET("/launches", handler.GetLaunches)
			data.GET("/launches/:id", handler.GetLaunchByID)
			data.PUT("/launches/:id", handler.UpdateLaunch)
		}

		images := apiV1.Group("/images")
		images.Use(middleware.AuthMiddleware(handler.jwtM, handler.logger))
		{
			images.POST("", handler.UploadImage)
			images.POST("/thumb", handler.GenerateThumb)
			images.GET("", handler.ListImages)
			images.DELETE("/:key", handler.DeleteImage)
		}

		task := apiV1.Group("/task")
		task.Use(middleware.AuthMiddleware(handler.jwtM, handler.logger))
		{
			task.POST("", handler.StartTask)
			task.GET("", handler.GetTask)
			task.GET("/history", handler.GetTaskHistory)
			task.POST("/action", handler.TaskAction)
		}

		ll2 := apiV1.Group("/ll2")
		ll2.Use(middleware.AuthMiddleware(handler.jwtM, handler.logger))
		{
			ll2.GET("/launches", handler.GetLL2Launches)
			ll2.POST("/launches/update", handler.StartLL2LaunchUpdate)
			ll2.GET("/angecies", handler.GetLL2Angecy)
			ll2.POST("/angecies/update", handler.StartLL2AngecyUpdate)
			ll2.GET("/launcher-families", handler.GetLL2LauncherFamilies)
			ll2.GET("/launchers", handler.GetLL2Launchers)
			ll2.POST("/launchers/update", handler.StartLL2LauncherUpdate)
			ll2.POST("/launcher-families/update", handler.StartLL2LauncherFamilyUpdate)
			ll2.GET("/locations", handler.GetLL2Locations)
			ll2.POST("/locations/update", handler.StartLL2LocationUpdate)
			ll2.GET("/pads", handler.GetLL2Pads)
			ll2.POST("/pads/update", handler.StartLL2PadUpdate)
		}
	}

	return router
}
