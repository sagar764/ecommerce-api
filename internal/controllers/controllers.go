package controllers

import (
	"ecommerce-api/internal/middlewares"
	"ecommerce-api/internal/usecases"
	"ecommerce-api/version"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EcommerceController struct {
	router   *gin.RouterGroup
	m        *middlewares.Middlewares
	useCases usecases.EcommerceUseCaseImply
}

// EcommerceController
func NewEcommerceController(router *gin.RouterGroup, m *middlewares.Middlewares,
	ecommerceUseCase usecases.EcommerceUseCaseImply) *EcommerceController {
	return &EcommerceController{
		router:   router,
		m:        m,
		useCases: ecommerceUseCase,
	}
}

// InitRoutes
func (ecommerce *EcommerceController) InitRoutes() {

	ecommerce.router.GET("/:version/health", func(ctx *gin.Context) {
		version.RenderHandler(ctx, ecommerce, "HealthHandler")
	})

	// Create a protected group with JWT authentication
	protectedRoutes := ecommerce.router.Group("/:version")
	protectedRoutes.Use(ecommerce.m.JWTAuth())

	// Define category routes
	categoryRoutes := protectedRoutes.Group("/category")
	{
		categoryRoutes.POST("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "CreateCategory")
		})

		categoryRoutes.PATCH("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "UpdateCategory")
		})

		categoryRoutes.DELETE("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "DeleteCategory")
		})

		categoryRoutes.GET("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetCategoryByID")
		})

		categoryRoutes.GET("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetAllCategories")
		})
	}

	productRoutes := protectedRoutes.Group("/:version/products")
	{
		productRoutes.POST("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "CreateProduct")
		})

		productRoutes.PATCH("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "UpdateProduct")
		})

		productRoutes.DELETE("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "DeleteProduct")
		})

		productRoutes.GET("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetProductByID")
		})

		productRoutes.GET("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetAllProducts")
		})
	}

	variantRoutes := protectedRoutes.Group("/:version/variants")
	{
		variantRoutes.POST("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "CreateVariant")
		})

		variantRoutes.PATCH("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "UpdateVariant")
		})

		variantRoutes.DELETE("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "DeleteVariant")
		})

		variantRoutes.GET("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetVariantByID")
		})

		variantRoutes.GET("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetAllVariants")
		})
	}

	orderRoutes := protectedRoutes.Group("/:version/orders")
	{
		orderRoutes.POST("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "CreateOrder")
		})

		orderRoutes.GET("/:id", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetOrderByID")
		})

		orderRoutes.GET("/", func(ctx *gin.Context) {
			version.RenderHandler(ctx, ecommerce, "GetAllOrders")
		})
	}
}

// HealthHandler
func (ecommerce *EcommerceController) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "server run with base version",
	})

}
