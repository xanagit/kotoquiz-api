package services

import (
	"fmt"
	"gorm.io/gorm"
)

type ApiHealthService interface {
	Check() error
}

type ApiHealthServiceImpl struct {
	DB *gorm.DB
}

// Make sure that ApiHealthServiceImpl implements ApiHealthService
var _ ApiHealthService = (*ApiHealthServiceImpl)(nil)

func (hs *ApiHealthServiceImpl) Check() error {
	sqlDB, err := hs.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}

	return nil
}
