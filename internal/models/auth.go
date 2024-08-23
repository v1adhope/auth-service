package models

type TokenPair struct {
	Id      string `json:"-"`
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}
