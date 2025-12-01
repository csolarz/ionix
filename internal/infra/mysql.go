package infra

import (
	"context"
	"fmt"

	utils "github.com/csolarz/ionix/internal/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define la interfaz del repositorio de base de datos, sirve para desacoplar la lógica de negocio
// de la implementación concreta del acceso a datos.
// Se usa mockery para generar mocks de esta interfaz para pruebas unitarias.

//go:generate mockery --name=DBRepository --output=./mock --outpkg=mock --case=snake
type DBRepository interface {
	GetByID(ctx context.Context, id int64, data any) error
	GetAll(ctx context.Context, data any) error
	Create(ctx context.Context, data any) error
	Update(ctx context.Context, data any) error
	Delete(ctx context.Context, data any) error
}

type DBGorm struct {
	db *gorm.DB
}

func NewDBGorm(dsn string) (*DBGorm, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %w", err)
	}

	return &DBGorm{db: db}, nil
}

// TODO: Gracefully close the database connection
func (r *DBGorm) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (r *DBGorm) GetByID(ctx context.Context, id int64, data any) error {
	result := r.db.WithContext(ctx).First(data, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
		}
		return fmt.Errorf("error obteniendo registro: %w", result.Error)
	}

	return nil
}

func (r *DBGorm) GetAll(ctx context.Context, data any) error {
	result := r.db.WithContext(ctx).Find(data)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return utils.ErrNotFound
		}
		return fmt.Errorf("error obteniendo registros: %w", result.Error)
	}

	return nil
}

func (r *DBGorm) Create(ctx context.Context, data any) error {
	result := r.db.WithContext(ctx).Create(data)
	if result.Error != nil {
		return fmt.Errorf("error creando registro: %w", result.Error)
	}

	return nil
}

func (r *DBGorm) Update(ctx context.Context, data any) error {
	result := r.db.WithContext(ctx).Save(data)
	if result.Error != nil {
		return fmt.Errorf("error actualizando registro: %w", result.Error)
	}
	return nil
}

func (r *DBGorm) Delete(ctx context.Context, data any) error {
	result := r.db.WithContext(ctx).Delete(data)
	if result.Error != nil {
		return fmt.Errorf("error eliminando registro: %w", result.Error)
	}
	return nil
}
