package races

import (
	"sort"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RacesService interface {
	Create(name string) (*Race, error)
	Race(id uuid.UUID) (*Race, error)
	Races() (Races, error)
	Start(id uuid.UUID) (*Race, error)
}

var _ RacesService = &raceService{}

func New() RacesService {
	s := raceService{}
	s.races = make(map[uuid.UUID]*Race, 3)

	s.prefillRaces()

	return &s
}

type raceService struct {
	races map[uuid.UUID]*Race
}

func (s *raceService) Race(id uuid.UUID) (*Race, error) {
	race, ok := s.races[id]
	if !ok {
		return nil, errors.Errorf("no race found with id %q", id.String())
	}

	return race, nil
}

func (s *raceService) Races() (Races, error) {
	races := make(Races, 0, len(s.races))
	for _, r := range s.races {
		races = append(races, r)
	}
	sort.Sort(races)

	return races, nil
}

func (s *raceService) Create(name string) (*Race, error) {
	race := NewRace(name)
	race.Drivers = []*RaceDriver{
		{
			Id: 1,
		},
		{
			Id: 2,
		},
		{
			Id: 3,
		},
	}

	s.races[race.Id] = race
	return race, nil
}

func (s *raceService) Start(id uuid.UUID) (*Race, error) {
	race, err := s.Race(id)
	if err != nil {
		return nil, err
	}
	race.Start()

	return race, nil
}

func (s *raceService) prefillRaces() {
	names := []string{
		"Prairie Circuit",
		"Gophtona 500",
		"South Rodent Ring",
	}

	for index, drivers := range initialDrivers() {
		race := NewRace(names[index])
		race.Drivers = drivers
		s.races[race.Id] = race
	}
}

func initialDrivers() [][]*RaceDriver {
	drivers := make([][]*RaceDriver, 3)
	drivers[0] = []*RaceDriver{
		{
			Id: 1,
		},
		{
			Id: 2,
		},
		{
			Id: 3,
		},
	}
	drivers[1] = []*RaceDriver{
		{
			Id:       1,
			rank:     1,
			Position: 180,
		},
		{
			Id:       2,
			rank:     3,
			Position: 180,
		},
		{
			Id:       3,
			rank:     2,
			Position: 180,
		},
	}
	drivers[2] = []*RaceDriver{
		{
			Id:       1,
			rank:     1,
			Position: 180,
		},
		{
			Id:       2,
			rank:     2,
			Position: 180,
		},
		{
			Id:       3,
			rank:     3,
			Position: 180,
		},
	}

	return drivers
}
