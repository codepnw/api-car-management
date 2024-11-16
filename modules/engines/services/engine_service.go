package engservices

import (
	"context"

	"github.com/codepnw/go-car-management/modules/engines"
	engrepositories "github.com/codepnw/go-car-management/modules/engines/repositories"
	"github.com/go-playground/validator/v10"
)

type IEngineService interface {
	GetEngineByID(ctx context.Context, id string) (*engines.Engine, error)
	CreateEngine(ctx context.Context, req *engines.EngineRequest) (*engines.Engine, error)
	UpdateEngine(ctx context.Context, id string, req *engines.EngineRequest) (*engines.Engine, error)
	DeleteEngine(ctx context.Context, id string) (*engines.Engine, error)
}

type engineService struct {
	repo engrepositories.IEngineRepository
}

var validate *validator.Validate

func NewEngineService(repo engrepositories.IEngineRepository) IEngineService {
	return &engineService{repo: repo}
}

func (s *engineService) GetEngineByID(ctx context.Context, id string) (*engines.Engine, error) {
	engine, err := s.repo.GetEngineByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func (s *engineService) CreateEngine(ctx context.Context, req *engines.EngineRequest) (*engines.Engine, error) {
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}
		return nil, err
	}

	createdEngine, err := s.repo.CreateEngine(ctx, req)
	if err != nil {
		return nil, err
	}
	return createdEngine, nil
}

func (s *engineService) UpdateEngine(ctx context.Context, id string, req *engines.EngineRequest) (*engines.Engine, error) {
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}
		return nil, err
	}

	updatedEngine, err := s.repo.UpdateEngine(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return updatedEngine, nil
}

func (s *engineService) DeleteEngine(ctx context.Context, id string) (*engines.Engine, error) {
	deletedEngine, err := s.repo.DeleteEngine(ctx, id)
	if err != nil {
		return nil, err
	}
	return deletedEngine, nil
}
