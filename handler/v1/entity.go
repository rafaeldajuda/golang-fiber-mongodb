package v1

import "time"

type Animal struct {
	ID       string    `json:"_id,omitempty"`
	Nome     string    `json:"nome"`
	Tipo     string    `json:"tipo"`
	Idade    int       `json:"idade"`
	Dono     string    `json:"dono"`
	Castrado bool      `json:"castrado"`
	Cirurgia time.Time `json:"cirurgia"`
}
