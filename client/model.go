package client

import "github.com/parnurzeal/gorequest"

type Client struct {
	Host		string
	Log			bool
	superAgent *gorequest.SuperAgent
}

type Base struct {
	ID		int64	`json:"id"`
	Name	string	`json:"nome"`
}

type State struct {
	Base
	Acronym	string	`json:"sigla"`

	Region 	`json:"regiao"`
}

type Region struct {
	Base
	Acronym	string	`json:"sigla"`
}

type MicroRegion struct {
	Base
	MesorRegion 	`json:"mesorregiao"`
}

type MesorRegion struct {
	Base
	State	`json:"UF"`
}

type ImmediateRegion struct {
	Base
	State	`json:"UF"`
}

type County struct {
	Base
	MicroRegion `json:"microrregiao"`
	ImmediateRegion `json:"regiao-intermediaria"`
}
