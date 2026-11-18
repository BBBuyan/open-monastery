package models

import "time"

type Monastery struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
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
