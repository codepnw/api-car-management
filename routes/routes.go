package routes

import (
	"database/sql"

	"github.com/codepnw/go-car-management/modules/cars/carrepositories"
	carhandlers "github.com/codepnw/go-car-management/modules/cars/handlers"
	carservices "github.com/codepnw/go-car-management/modules/cars/services"
	enghandlers "github.com/codepnw/go-car-management/modules/engines/handlers"
	engrepositories "github.com/codepnw/go-car-management/modules/engines/repositories"
	engservices "github.com/codepnw/go-car-management/modules/engines/services"
	"github.com/gin-gonic/gin"
)

func NewRoutes(db *sql.DB, r *gin.Engine, version string) {
	carRoutes(db, r, version)
	engineRoutes(db, r, version)
}

func carRoutes(db *sql.DB, r *gin.Engine, version string) {
	g := r.Group(version + "/cars")

	repo := carrepositories.NewCarRepository(db)
	service := carservices.NewCarService(repo)
	handler := carhandlers.NewCarHandler(service)

	idParam := "/:id"

	g.GET(idParam, handler.GetCarByID)
	g.GET("/", handler.GetCarByBrand)
	g.POST("/", handler.CreateCar)
	g.PATCH(idParam, handler.UpdateCar)
	g.DELETE(idParam, handler.DeleteCar)
}

func engineRoutes(db *sql.DB, r *gin.Engine, version string) {
	g := r.Group(version + "/engines")

	repo := engrepositories.NewEngineRepository(db)
	service := engservices.NewEngineService(repo)
	handler := enghandlers.NewEngineHandler(service)

	idParam := "/:id"

	g.GET(idParam, handler.GetEngineByID)
	g.POST("/", handler.CreateEngine)
	g.PATCH(idParam, handler.UpdateEngine)
	g.DELETE(idParam, handler.DeleteEngine)
}