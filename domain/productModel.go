package domain

import "github.com/amupxm/pure-webserver/pkg/database"

type Product struct {
	database.DbModel
	Name    string `json:"name"`
	Brand   string `json:"brand"`
	Company string `json:"company"`
	Iid     string `json:"iid"`
	Id      string `json:"id"`
}
