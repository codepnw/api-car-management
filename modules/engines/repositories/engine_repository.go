package engineRepositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/go-car-management/modules/engines"
	"github.com/google/uuid"
)

type IEngineRepository interface {
	GetEngineByID(ctx context.Context, id string) (*engines.Engine, error)
	CreateEngine(ctx context.Context, req *engines.EngineRequest) (*engines.Engine, error)
	UpdateEngine(ctx context.Context, id string, req *engines.EngineRequest) (*engines.Engine, error)
	DeleteEngine(ctx context.Context, id string) (*engines.Engine, error)
}

type enginRepository struct {
	db *sql.DB
}

func NewEngineRepository(db *sql.DB) IEngineRepository {
	return &enginRepository{db: db}
}

func (r *enginRepository) GetEngineByID(ctx context.Context, id string) (*engines.Engine, error) {
	var engine engines.Engine

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &engine, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	err = tx.QueryRowContext(
		ctx,
		"SELECT engine_id, displacement, no_of_cylinders, car_range FROM engines WHERE engine_id = $1;",
		id,
	).Scan(
		&engine.EngineID,
		&engine.Displacement,
		&engine.NoOfCylinders,
		&engine.CarRange,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &engine, nil
		}
		return &engine, err
	}

	return &engine, err
}

func (r *enginRepository) CreateEngine(ctx context.Context, req *engines.EngineRequest) (*engines.Engine, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &engines.Engine{}, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	engineID := uuid.New()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO engines (engine_id, displacement, no_of_cylinders, car_range)
		VALUES ($1, $2, $3, $4);`,
		engineID,
		req.Displacement,
		req.NoOfCylinders,
		req.CarRange,
	)
	if err != nil {
		return &engines.Engine{}, err
	}

	engine := &engines.Engine{
		EngineID:      engineID,
		Displacement:  req.Displacement,
		NoOfCylinders: req.NoOfCylinders,
		CarRange:      req.CarRange,
	}

	return engine, nil
}

func (r *enginRepository) UpdateEngine(ctx context.Context, id string, req *engines.EngineRequest) (*engines.Engine, error) {
	engineID, err := uuid.Parse(id)
	if err != nil {
		return &engines.Engine{}, fmt.Errorf("invalid engine id: %w", err)
	}

	// Transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &engines.Engine{}, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("transaction commit error: %v\n", cmErr)
			}
		}
	}()

	results, err := tx.ExecContext(
		ctx,
		"UPDATE engines SET displacement = $1, no_of_cylinders = $2, car_range = $3 WHERE engine_id = $4",
		req.Displacement,
		req.NoOfCylinders,
		req.CarRange,
		engineID,
	)
	if err != nil {
		return &engines.Engine{}, err
	}

	rowAffected, err := results.RowsAffected()
	if err != nil {
		return &engines.Engine{}, err
	}

	if rowAffected == 0 {
		return &engines.Engine{}, errors.New("no rows updated")
	}

	engine := &engines.Engine{
		EngineID:      engineID,
		Displacement:  req.Displacement,
		NoOfCylinders: req.NoOfCylinders,
		CarRange:      req.CarRange,
	}

	return engine, nil
}

func (r *enginRepository) DeleteEngine(ctx context.Context, id string) (*engines.Engine, error) {
	var engine engines.Engine

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return &engines.Engine{}, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	err = tx.QueryRowContext(
		ctx,
		"SELECT engine_id, displacement, no_of_cylinders, car_range FROM engines WHERE engine_id = $1;",
		id,
	).Scan(
		&engine.EngineID,
		&engine.Displacement,
		&engine.NoOfCylinders,
		&engine.CarRange,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &engine, nil
		}
		return &engine, err
	}

	result, err := tx.ExecContext(ctx, "DELETE FROM engines WHERE engine_id = $1;", id)
	if err != nil {
		return &engines.Engine{}, err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return &engines.Engine{}, err
	}

	if rowAffected == 0 {
		return &engines.Engine{}, errors.New("no rows delete")
	}

	return &engine, nil
}
