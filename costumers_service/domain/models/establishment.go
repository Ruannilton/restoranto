package models

type Establishment struct {
	Id               int
	CompanyName      string
	FantasyName      string
	Cnpj             string
	EstablishmentKey string
	Costumer         Costumer
}
