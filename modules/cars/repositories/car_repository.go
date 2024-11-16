package carRepositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/codepnw/go-car-management/modules/cars"
	"github.com/codepnw/go-car-management/modules/engines"
	"github.com/google/uuid"
)

type ICarRepository interface {
	GetCarById(ctx context.Context, id string) (*cars.Car, error)
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]*cars.Car, error)
	CreateCar(ctx context.Context, req *cars.CarRequest) (*cars.Car, error)
	UpdateCar(ctx context.Context, id string, req *cars.CarRequest) (*cars.Car, error)
	DeleteCar(ctx context.Context, id string) (*cars.Car, error)
}

type carRepository struct {
	db *sql.DB
}

func NewCarRepository(db *sql.DB) ICarRepository {
	return &carRepository{db: db}
}

func (r *carRepository) GetCarById(ctx context.Context, id string) (*cars.Car, error) {
	var response cars.Car

	query := `
		SELECT c.car_id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at,
			e.engine_id, e.displacement, e.no_of_cylinders, e.car_range 
		FROM cars c 
		LEFT JOIN engines e ON c.engine_id = e.engine_id 
		WHERE c.car_id = $1;
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&response.CarID,
		&response.Name,
		&response.Year,
		&response.Brand,
		&response.FuelType,
		&response.Engine.EngineID,
		&response.Price,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.Engine.EngineID,
		&response.Engine.Displacement,
		&response.Engine.NoOfCylinders,
		&response.Engine.CarRange,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &response, nil
		}
		return &response, err
	}

	return &response, nil
}

func (r *carRepository) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]*cars.Car, error) {
	var response []*cars.Car
	var query string

	if isEngine {
		query = `
			SELECT c.car_id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.Price, c.created_at, c.updated_at
				e.engine_id, e.displacement, e.no_of_cylinders, e.car_range
			FROM cars c
			LEFT JOIN engines e ON c.engine_id = e.engine_id
			WHERE c.brand = $1;
		`
	} else {
		query = `
			SELECT car_id, name, year, brand, fuel_type, engine_id, Price, created_at, updated_at
			FROM cars WHERE brand = $1;
		`
	}

	rows, err := r.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car cars.Car
		if isEngine {
			var engine engines.Engine
			err := rows.Scan(
				&car.CarID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Engine.EngineID,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
				&car.Engine.EngineID,
				&car.Engine.Displacement,
				&car.Engine.NoOfCylinders,
				&car.Engine.CarRange,
			)
			if err != nil {
				return nil, err
			}

			car.Engine = &engine
		} else {
			err := rows.Scan(
				&car.CarID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Engine.EngineID,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}

		response = append(response, &car)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

func (r *carRepository) CreateCar(ctx context.Context, req *cars.CarRequest) (*cars.Car, error) {
	var createCar cars.Car
	var engineID uuid.UUID

	err := r.db.QueryRowContext(ctx, "SELECT id FROM engines WHERE id = $1;", req.Engine.EngineID).Scan(&engineID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &createCar, errors.New("engine_id does not exists in the engine table")
		}
		return &createCar, err
	}

	carId := uuid.New()
	createdAt := time.Now().Local()
	updatedAt := createdAt

	newCar := cars.Car{
		CarID:     carId,
		Name:      req.Name,
		Year:      req.Year,
		Brand:     req.Brand,
		FuelType:  req.FuelType,
		Engine:    req.Engine,
		Price:     req.Price,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &createCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		INSERT INTO cars (car_id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING car_id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at;
	`
	err = tx.QueryRowContext(
		ctx,
		query,
		newCar.CarID,
		newCar.Name,
		newCar.Year,
		newCar.Brand,
		newCar.FuelType,
		newCar.Engine.EngineID,
		newCar.Price,
		newCar.CreatedAt,
		newCar.UpdatedAt,
	).Scan(
		&createCar.CarID,
		&createCar.Name,
		&createCar.Year,
		&createCar.Brand,
		&createCar.FuelType,
		&createCar.Engine.EngineID,
		&createCar.Price,
		&createCar.CreatedAt,
		&createCar.UpdatedAt,
	)

	if err != nil {
		return &createCar, err
	}
	return &createCar, nil
}

func (r *carRepository) UpdateCar(ctx context.Context, id string, req *cars.CarRequest) (*cars.Car, error) {
	var updatedCar cars.Car

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &updatedCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		UPDATE cars
		SET name=$2, year=$3, brand=$4, fuel_type=$5, engine_id=$6, price=$7, updated_at=$8
		WHERE car_id = $1
		RETURNING car_id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		id,
		req.Name,
		req.Year,
		req.Brand,
		req.FuelType,
		req.Engine.EngineID,
		req.Price,
		time.Now().Local(),
	).Scan(
		&updatedCar.CarID,
		&updatedCar.Name,
		&updatedCar.Year,
		&updatedCar.Brand,
		&updatedCar.FuelType,
		&updatedCar.Engine.EngineID,
		&updatedCar.Price,
		&updatedCar.UpdatedAt,
	)

	if err != nil {
		return &updatedCar, err
	}
	return &updatedCar, nil
}

func (r *carRepository) DeleteCar(ctx context.Context, id string) (*cars.Car, error) {
	var deletedCar cars.Car

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &deletedCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = tx.QueryRowContext(
		ctx,
		`SELECT car_id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
		FROM cars WHERE car_id = $1;`,
		id,
	).Scan(
		&deletedCar.CarID,
		&deletedCar.Name,
		&deletedCar.Year,
		&deletedCar.Brand,
		&deletedCar.FuelType,
		&deletedCar.Engine.EngineID,
		&deletedCar.Price,
		&deletedCar.CreatedAt,
		&deletedCar.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &cars.Car{}, errors.New("car not found")
		}
		return &cars.Car{}, err
	}

	result, err := tx.ExecContext(ctx, "DELETE FROM cars WHERE car_id = $1;", id)
	if err != nil {
		return &cars.Car{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &cars.Car{}, err
	}

	if rowsAffected == 0 {
		return &cars.Car{}, errors.New("no rows delete")
	}

	return &deletedCar, nil
}
