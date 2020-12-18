package model

// Owner model is used as part of breeder struct
type Owner struct {
	ID       int    `json:"id,omitempty"`
	Forename string `json:"forename"`
	Surname  string `json:"surname"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

// Breeder model is the data structure used to store interface with database data
type Breeder struct {
	ID         int     `json:"id,omitempty"`
	Affix      string  `json:"affix"`
	ShortAffix string  `json:"shortAffix"`
	Website    string  `json:"website"`
	Owners     []Owner `json:"owners"`
}
