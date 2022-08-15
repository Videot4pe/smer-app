package smer

import "time"

type Smer struct {
	Id        uint16    `json:"id" sql:"id"`
	UserId    uint16    `json:"userId" sql:"user_id"`
	Situation string    `json:"situation" sql:"situation"`
	Thoughts  []string  `json:"thoughts" sql:"thoughts"`
	Emotions  []string  `json:"emotions" sql:"emotions"`
	Reactions []string  `json:"reactions" sql:"reactions"`
	CreatedAt time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" sql:"updated_at"`
}

type NewSmerDto struct {
	Situation string   `json:"situation" sql:"situation"`
	Thoughts  []string `json:"thoughts" sql:"thoughts"`
	Emotions  []string `json:"emotions" sql:"emotions"`
	Reactions []string `json:"reactions" sql:"reactions"`
}
