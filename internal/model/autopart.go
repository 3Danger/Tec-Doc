package model

//Autopart структура описывающая запчасть с TecDoc
type Autopart struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}
