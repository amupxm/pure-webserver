package controller

import (
	"log"

	"github.com/amupxm/pure-webserver/config"
	"github.com/amupxm/pure-webserver/logic"
	httpEngine "github.com/amupxm/pure-webserver/pkg/httpEngine"
)

type (
	engine struct {
		ProductLogic logic.ProductLogic
	}
	Engine interface {
		// GetOne returns one product
		GetOne(c *httpEngine.ServerContext)
		// GetAll returns all products
		GetAll(c *httpEngine.ServerContext)
		// UpdateOne updates one product
		UpdateOne(c *httpEngine.ServerContext)
		// CreateProduct creates one product
		CreateProduct(c *httpEngine.ServerContext)
		// DeleteOne deletes one product
		DeleteOne(c *httpEngine.ServerContext)
	}
)

func NewEngine(pl logic.ProductLogic) Engine {
	return &engine{
		ProductLogic: pl,
	}
}

func InitNewEngine(pl logic.ProductLogic) {
	en := NewEngine(pl)
	server := httpEngine.NewServer()

	server.AddHandler("/v1/toys", "GET", en.GetAll)
	server.AddHandler("/v1/toys/:iid", "GET", en.GetOne)
	server.AddHandler("/v1/toys/:iid", "DELETE", en.DeleteOne)
	server.AddHandler("/v1/toys/:iid", "PATCH", en.UpdateOne)
	server.AddHandler("/v1/toys/:iid", "PUT", en.CreateProduct)

	// listen on port 8080 , You can change this port from config.json
	server.StartServer(config.AppConf.Http.Port)
	log.Printf("server started on port %s\n", config.AppConf.Http.Port)
}
