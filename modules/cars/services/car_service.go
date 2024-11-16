package carservices

import (
	"context"

	"github.com/codepnw/go-car-management/modules/cars"
	"github.com/codepnw/go-car-management/modules/cars/carrepositories"
	"github.com/go-playground/validator/v10"
)

type ICarService interface {
	GetCarById(ctx context.Context, id string) (*cars.Car, error) 
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]*cars.Car, error) 
	CreateCar(ctx context.Context, req *cars.CarRequest) (*cars.Car, error) 
	UpdateCar(ctx context.Context, id string, req *cars.CarRequest) (*cars.Car, error)
	DeleteCar(ctx context.Context, id string) (*cars.Car, error)
}

type carService struct {
	repo carrepositories.ICarRepository
}

var validate *validator.Validate

func NewCarService(repo carrepositories.ICarRepository) ICarService {
	return &carService{repo: repo}
}

func (s *carService) GetCarById(ctx context.Context, id string) (*cars.Car, error) {
	car, err := s.repo.GetCarById(ctx, id)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (s *carService) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]*cars.Car, error) {
	results, err := s.repo.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *carService) CreateCar(ctx context.Context, req *cars.CarRequest) (*cars.Car, error) {
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}
		return nil, err
	}

	createdCar, err := s.repo.CreateCar(ctx, req)
	if err != nil {
		return nil, err
	}
	return createdCar, nil
}

func (s *carService) UpdateCar(ctx context.Context, id string, req *cars.CarRequest) (*cars.Car, error) {
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}
		return nil, err
	}

	updatedCar, err := s.repo.UpdateCar(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return updatedCar, nil
}

func (s *carService) DeleteCar(ctx context.Context, id string) (*cars.Car, error) {
	deletedCar, err := s.repo.DeleteCar(ctx, id)
	if err != nil {
		return nil, err
	}
	return deletedCar, nil
}