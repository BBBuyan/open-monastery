package models

import "time"

type Monastery struct {
	Id        int64
	Name      string
	CreatedBy *string
	CreatedAt time.Time
}

func NewMonastery(
	name string,
	createdBy *string,
) *Monastery {
	return &Monastery{
		Name:      name,
		CreatedBy: createdBy,
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
	}
}
