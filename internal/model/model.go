package model

import (
	"time"
)

type Products struct {
	Article     int     `json:"article,omitempty"`
	Catalog     string  `json:"catalog,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	PhotoUrl    []byte  `json:"photo_url,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Length      int     `json:"length"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Weight      int     `json:"weight"`
}

//	type Products struct {
//		Article     int
//		Catalog     string
//		Name        string
//		Description string
//		PhotoUrl    string
//		Price       float32
//	}
type Corzine struct {
	ID        int       `json:"id,omitempty"`
	TgId      int64     `json:"tg_id,omitempty"`
	Article   int       `json:"article,omitempty"` //В наличии
	Quantity  int       `json:"quantity,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type Orders struct {
	ID          int       `json:"id,omitempty"`
	TgID        int64     `json:"tg_id,omitempty"`
	StatusOrder int       `json:"status_order,omitempty"`
	Pvz         string    `json:"pvz,omitempty"`
	Order       string    `json:"order,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Users struct {
	id         int   `json:"id,omitempty"`
	TgID       int64 `json:"tg_id,omitempty"`
	StatusUser int   `json:"status_user,omitempty"`
	StateUser  int   `json:"state_user,omitempty"`
	//Corzine    []int     `json:"corzine,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
