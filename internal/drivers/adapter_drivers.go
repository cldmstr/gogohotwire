package drivers

type DriversStore []*Driver

var _ DriversAdapter = DriversStore{}

func NewDriversStore() DriversStore {
	store := make([]*Driver, 0, 3)
	store = append(store,
		&Driver{
			Id:         1,
			FirstName:  "Max",
			LastName:   "Kowalski",
			NickName:   "Cheetah",
			Attributes: []string{"Top Speed", "Acceleration"},
			Motto:      "You break, you loose!",
			Color:      "blue",
		},
		&Driver{
			Id:         2,
			FirstName:  "Kimiko",
			LastName:   "Takahashi",
			NickName:   "ByeBye",
			Attributes: []string{"Ruthless Passing", "Slipstream"},
			Motto:      "Oh, sorry, didn't see you there!",
			Color:      "black",
		},
		&Driver{
			Id:         3,
			FirstName:  "Ueli",
			LastName:   "Anderegg",
			NickName:   "Chugublitz",
			Attributes: []string{"Hill Climbing", "Precision"},
			Motto:      "Drive like clockwork!",
			Color:      "beige",
		},
	)

	return store
}

func (s DriversStore) Drivers() ([]*Driver, error) {
	return []*Driver(s), nil
}
