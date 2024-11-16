package engines

import "github.com/google/uuid"

type Engine struct {
	EngineID      uuid.UUID `json:"engineId" db:"engine_id"`
	Displacement  uint16    `json:"displacement" db:"displacement"`
	NoOfCylinders uint16    `json:"noOfCylinders" db:"no_of_cylinders"`
	CarRange      uint16    `json:"carRange" db:"car_range"`
}

type EngineRequest struct {
	Displacement  uint16 `json:"displacement" validate:"required"`
	NoOfCylinders uint16 `json:"noOfCylinders" validate:"required"`
	CarRange      uint16 `json:"carRange" validate:"required"`
}
