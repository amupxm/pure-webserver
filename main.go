package main

import (
	"github.com/amupxm/pure-webserver/config"
	"github.com/amupxm/pure-webserver/controller"
	"github.com/amupxm/pure-webserver/logic"
	"github.com/amupxm/pure-webserver/pkg/database"
	"github.com/amupxm/pure-webserver/repository"
)

func main() {
	config.Init()
	database := database.NewDatabase()
	productsRepository := repository.NewProductRepository(database)
	productLogic := logic.NewProductLogic(productsRepository)
	controller.InitNewEngine(productLogic)
}
