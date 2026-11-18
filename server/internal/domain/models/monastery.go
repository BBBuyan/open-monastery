package models

import "time"

type Monastery struct {
	Id        int64
	Name      string
	CreatedBy int64
	CreatedAt time.Time
}

func NewMonastery(
	name string,
	userId int64,
) *Monastery {
	return &Monastery{
		Name:      name,
		CreatedBy: userId,
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
	}
}
