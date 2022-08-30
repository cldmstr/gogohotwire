package drivers

type DriversStore []*Driver

var _ DriversAdapter = DriversStore{}

func NewDriversStore() DriversStore {
	store := make([]*Driver, 0, 3)
	store = append(store,
		&Driver{
			FirstName:  "Max",
			LastName:   "Kowalski",
			NickName:   "Lightning",
			Attributes: []string{"Top Speed", "Acceleration"},
			Motto:      "You break, you loose!",
		},
		&Driver{
			FirstName:  "Kimiko",
			LastName:   "Takahashi",
			NickName:   "ByeBye",
			Attributes: []string{"Ruthless Passing", "Slipstream"},
			Motto:      "Oh, sorry, didn't see you there!",
		},
		&Driver{
			FirstName:  "Ueli",
			LastName:   "Anderegg",
			NickName:   "Chugublitz",
			Attributes: []string{"Hill Climbing", "Precision"},
			Motto:      "Driving like clockwork!",
		},
	)

	return store
}

func (s DriversStore) Drivers() ([]*Driver, error) {
	return []*Driver(s), nil
}
