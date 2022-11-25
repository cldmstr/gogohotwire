package drivers

import (
	"github.com/pkg/errors"
)

type DriversService interface {
	Drivers() ([]*Driver, error)
	Driver(id int) (*Driver, error)
}

type DriversAdapter interface {
	Drivers() ([]*Driver, error)
}

type driverService struct {
	store DriversAdapter
}

var _ DriversService = &driverService{}

func NewDriversService(drivers DriversAdapter) DriversService {
	s := driverService{
		store: drivers,
	}

	return &s
}

func (s *driverService) Drivers() ([]*Driver, error) {
	return s.store.Drivers()
}

func (s *driverService) Driver(id int) (*Driver, error) {
	drivers, err := s.store.Drivers()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load driver with id %d", id)
	}
	for _, d := range drivers {
		if id == d.Id {
			return d, nil
		}
	}

	return nil, errors.Errorf("failed to find driver %d", id)
}
