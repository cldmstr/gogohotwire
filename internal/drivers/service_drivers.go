package drivers

type DriversService interface {
	Drivers() ([]*Driver, error)
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
