package files

import "time"

type File struct {
	Id        uint16    `json:"id" sql:"id"`
	Path      string    `json:"path" sql:"path"`
	Name      string    `json:"name" sql:"name"`
	CreatedAt time.Time `json:"createdAt" sql:"created_at"`
}
