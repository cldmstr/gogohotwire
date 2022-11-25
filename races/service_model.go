package races

import (
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	Setup State = iota
	Ready
	Running
	Finished
)

type State int

type RaceDriver struct {
	Id       int
	Name     string
	rank     int
	Position int
}

func (r RaceDriver) Rank() string {
	switch r.rank {
	case 1:
		return "one"
	case 2:
		return "two"
	case 3:
		return "three"
	default:
		return "none"
	}
}

type Race struct {
	Id       uuid.UUID
	Name     string
	Drivers  []RaceDriver
	state    State
	watchers []chan RaceUpdate
}

func NewRace(name string) *Race {
	r := Race{}
	r.Name = name
	r.Id = uuid.New()
	r.state = Ready

	return &r
}

func (r *Race) StatusColor() string {
	switch r.state {
	case Ready:
		return "red"
	case Running:
		return "green"
	case Finished:
		return "checkered"
	default:
		return ""
	}
}

func (r *Race) Watch(watcher chan RaceUpdate) {
	if r.state == Finished {
		return
	}

	r.watchers = append(r.watchers, watcher)
}

func (r *Race) Start() {
	if r.state != Ready {
		return
	}
	r.state = Running

	go func() {
		for i := 0; i < 120; i++ {
			time.Sleep(time.Millisecond * 100)
			r.Drivers[0].Position += 2
			r.Drivers[1].Position += 3
			r.Drivers[2].Position += 5
			update := RaceUpdate{
				One:   r.Drivers[0].Position,
				Two:   r.Drivers[1].Position,
				Three: r.Drivers[2].Position,
			}
			for _, c := range r.watchers {
				c <- update
			}
		}
		r.Drivers[0].rank = 3
		r.Drivers[1].rank = 1
		r.Drivers[2].rank = 2
		r.state = Finished
	}()
}

type Races []*Race

var _ sort.Interface = &Races{}

func (r Races) Len() int {
	return len(r)
}

func (r Races) Less(i, j int) bool {
	a, b := r[i], r[j]
	if a.state == b.state {
		return strings.Compare(a.Name, b.Name) < 0
	}

	return a.state < b.state
}

func (r Races) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
