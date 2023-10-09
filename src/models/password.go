package models

type Password struct {
	New     string `json:"new"`
	Current string `json:"current"`
}
