package cars

import (
	"time"

	"github.com/codepnw/go-car-management/modules/engines"
	"github.com/google/uuid"
)

type Car struct {
	CarID     uuid.UUID       `json:"carId" db:"car_id"`
	Name      string          `json:"name" db:"name"`
	Year      uint16          `json:"year" db:"year"`
	Brand     string          `json:"brand" db:"brand"`
	FuelType  string          `json:"fuelType" db:"fuel_type"`
	Engine    *engines.Engine `json:"engineId" db:"engine_id"`
	Price     float64         `json:"price" db:"price"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time       `json:"updatedAt" db:"updated_at"`
}

type CarRequest struct {
	Name     string          `json:"name" validate:"required"`
	Year     uint16          `json:"year" validate:"required,gte=1886"`
	Brand    string          `json:"brand" validate:"required"`
	FuelType string          `json:"fuelType" validate:"oneof=Petrol Diesel Electric Hybrid"`
	Engine   *engines.Engine `json:"engine" validate:"required"`
	Price    float64         `json:"price" validate:"required,gte=1"`
}
