package races

import (
	"fmt"
	"math/rand"
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

type watcher struct {
	id     uuid.UUID
	update chan RaceUpdate
	finish chan bool
	closed bool
}

type Race struct {
	Id       uuid.UUID
	Name     string
	Drivers  []*RaceDriver
	state    State
	watchers []*watcher
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

func (r *Race) Watch() (uuid.UUID, chan RaceUpdate, chan bool, bool) {
	if r.state == Finished {
		return uuid.Nil, nil, nil, true
	}

	w := make(chan RaceUpdate, 1)
	d := make(chan bool, 1)
	watcher := &watcher{
		id:     uuid.New(),
		update: w,
		finish: d,
	}
	r.watchers = append(r.watchers, watcher)

	return watcher.id, w, d, false
}

func (r *Race) UnWatch(id uuid.UUID) {
	for _, w := range r.watchers {
		if w.id == id {
			fmt.Printf("Unwatch %v\n", w.id)
			w.closed = true
		}
	}
}

func (r *Race) Start() {
	if r.state != Ready {
		return
	}
	r.state = Running

	go func() {
		defer r.Finish()
	CheckeredFlag:
		for {
			time.Sleep(time.Millisecond * 100)
			for _, d := range r.Drivers {
				d.Position += rand.Intn(12)
			}
			update := RaceUpdate{
				One:   r.Drivers[0].Position,
				Two:   r.Drivers[1].Position,
				Three: r.Drivers[2].Position,
			}
			for _, w := range r.watchers {
				if w.closed {
					continue
				}
				w.update <- update
			}
			for _, d := range r.Drivers {
				if d.Position > 750 {
					break CheckeredFlag
				}
			}
		}
	}()
}

func (r *Race) Finish() {
	sort.Slice(r.Drivers, func(i, j int) bool {
		return r.Drivers[i].Position < r.Drivers[j].Position
	})
	for pos, d := range r.Drivers {
		d.rank = pos + 1
	}
	r.state = Finished
	for _, w := range r.watchers {
		w.finish <- true
	}
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
