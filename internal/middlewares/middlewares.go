package middlewares

import (
	"ecommerce-api/config"
	"ecommerce-api/internal/consts"
	"ecommerce-api/internal/entities"
	"ecommerce-api/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	Cfg *entities.EnvConfig
}

func NewMiddlewares(cfg *entities.EnvConfig) *Middlewares {
	return &Middlewares{
		Cfg: cfg,
	}
}

// Middleware function to check Accept-version from API Header
func (m Middlewares) ApiVersioning() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.Param("version")
		if version == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing version parameter"})
			return
		}

		apiVersion := utils.PrepareVersionName(version)
		apiVersion = strings.ToUpper(apiVersion)

		// set the accepting version in the context
		c.Set(consts.AcceptedVersions, apiVersion)

		// init the system Accepted versions
		// init the env config
		cfg, err := config.LoadConfig(consts.AppName)
		if err != nil {
			panic(err)
		}

		// set the list of system accepting version in the context
		systemAcceptedVersionsList := cfg.AcceptedVersions
		c.Set(consts.ContextSystemAcceptedVersions, systemAcceptedVersionsList)

		// check the version exists in the accepted list
		// find index of version from Accepted versions
		var found bool
		for index, version := range systemAcceptedVersionsList {
			version = strings.ToUpper(version)
			if version == apiVersion {
				found = true
				c.Set(consts.ContextAcceptedVersionIndex, index)
			}

		}
		if !found {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Given version is not supported by the system"})
			return
		}

		c.Next()
	}
}

func (m Middlewares) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set userID from token claims in context (if needed)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["sub"])
		}

		c.Next()
	}
}
