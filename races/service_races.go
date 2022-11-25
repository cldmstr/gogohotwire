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

func (r *raceService) Race(id uuid.UUID) (*Race, error) {
	race, ok := r.races[id]
	if !ok {
		return nil, errors.Errorf("no race found with id %q", id.String())
	}

	return race, nil
}

func (r *raceService) Races() (Races, error) {
	races := make(Races, 0, len(r.races))
	for _, r := range r.races {
		races = append(races, r)
	}
	sort.Sort(races)

	return races, nil
}

func (r *raceService) Create(name string) (*Race, error) {
	race := &Race{
		Id:   uuid.New(),
		Name: name,
		Drivers: []RaceDriver{
			{
				Id:       1,
				rank:     0,
				Position: 0,
			},
			{
				Id:       2,
				rank:     0,
				Position: 0,
			},
			{
				Id:       3,
				rank:     0,
				Position: 0,
			},
		},
		state: Ready,
	}

	r.races[race.Id] = race
	return race, nil
}

func (r *raceService) Start(id uuid.UUID) (*Race, error) {
	race, err := r.Race(id)
	if err != nil {
		return nil, err
	}
	race.Start()

	return race, nil
}

func (r *raceService) prefillRaces() {
	r1 := &Race{
		Id:   uuid.New(),
		Name: "Prairie Circuit",
		Drivers: []RaceDriver{
			{
				Id:       1,
				rank:     2,
				Position: 180,
			},
			{
				Id:       2,
				rank:     3,
				Position: 180,
			},
			{
				Id:       3,
				rank:     1,
				Position: 180,
			},
		},
		state: Finished,
	}
	r2 := &Race{
		Id:   uuid.New(),
		Name: "Gophtona 500",
		Drivers: []RaceDriver{
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
		},
		state: Finished,
	}
	r3 := &Race{
		Id:   uuid.New(),
		Name: "South Rodent Ring",
		Drivers: []RaceDriver{
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
		},
		state: Finished,
	}
	r.races[r1.Id] = r1
	r.races[r2.Id] = r2
	r.races[r3.Id] = r3
}
