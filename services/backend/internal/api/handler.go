package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/service/auth"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/image"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	ll2datasyncer "github.com/vamosdalian/launchdate-backend/internal/service/ll2_data_syncer"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type Handler struct {
	IsProduction bool
	logger       *logrus.Logger
	core         *core.MainService
	ll2Server    *ll2.LL2Service
	ll2syncer    *ll2datasyncer.LL2DataSyncer
	authHandler  *AuthHandler
	jwtM         *util.JWTManager
	imageService *image.ImageService
}

func NewHandler(logger *logrus.Logger, cfg *config.Config, db *db.MongoDB, ll2server *ll2.LL2Service,
	ll2syncer *ll2datasyncer.LL2DataSyncer, core *core.MainService, image *image.ImageService) *Handler {
	isProduction := cfg.Server.Env == config.ENV_PRODUCTION
	jwtManager := util.NewJWTManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTokenExpireMin,
		cfg.Auth.RefreshTokenExpireDays,
		cfg.Auth.Issuer,
	)
	authService := auth.NewAuthService(db.Database, jwtManager, logger)
	authHandler := NewAuthHandler(authService, logger, isProduction)

	return &Handler{
		IsProduction: isProduction,
		logger:       logger,
		core:         core,
		ll2Server:    ll2server,
		ll2syncer:    ll2syncer,
		authHandler:  authHandler,
		jwtM:         jwtManager,
		imageService: image,
	}
}

func (h *Handler) Health(c *gin.Context) {
	h.Json(c, "ok")
}
