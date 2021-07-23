package controller

import (
	"errors"

	"github.com/amupxm/pure-webserver/constants"
	"github.com/amupxm/pure-webserver/domain"
	httpEngine "github.com/amupxm/pure-webserver/pkg/httpEngine"
)

// GetAll handler writes all products to output
func (e *engine) GetAll(c *httpEngine.ServerContext) {
	ee, err := e.ProductLogic.GetAllProducts()
	if err != nil {
		c.ErrorHandler(400, err)
		return
	}
	c.JSON(200, ee)
}

// GetOne handler writes one product by iid to output
func (e *engine) GetOne(c *httpEngine.ServerContext) {
	// check iid exists or not
	id, err := c.GetURLParam("iid")
	if err != nil {
		c.ErrorHandler(400, err)
		return
	}
	ee, err := e.ProductLogic.GetProductByID(id)
	if err != nil {
		c.ErrorHandler(400, errors.New(constants.NoData))
		return
	}
	c.JSON(200, ee)

}

// CreateProduct add one product to db
func (e *engine) CreateProduct(c *httpEngine.ServerContext) {
	// check iid exists or not
	id, err := c.GetURLParam("iid")
	if err != nil {
		c.ErrorHandler(400, err)
		return
	}
	var product = &domain.Product{}
	// check is valid json for product
	err = c.BindToJson(product)
	if err != nil {
		c.ErrorHandler(400, err)
		return

	}
	product.Iid = id
	res, err := e.ProductLogic.NewProduct(product)
	if err != nil {
		c.ErrorHandler(400, err)
		return

	}
	c.JSON(200, res)
}

// UpdateOne updates one product which exists in db
func (e *engine) UpdateOne(c *httpEngine.ServerContext) {
	// check iid exists or not
	id, err := c.GetURLParam("iid")
	if err != nil {
		c.ErrorHandler(400, err)
		return
	}

	var product = &domain.Product{}
	// check is valid json for product
	err = c.BindToJson(product)
	if err != nil {
		c.ErrorHandler(400, err)
		return

	}
	product.Iid = id
	res, err := e.ProductLogic.UpdateProduct(product)
	if err != nil {
		c.ErrorHandler(400, err)
		return

	}
	c.JSON(200, res)
}

// DeleteOne deletes one product from db
func (e *engine) DeleteOne(c *httpEngine.ServerContext) {
	// check iid exists or not
	id, err := c.GetURLParam("iid")
	if err != nil {
		c.ErrorHandler(400, err)
		return
	}
	var product = &domain.Product{}
	product.Iid = id
	err = e.ProductLogic.DeleteProduct(product)
	if err != nil {
		c.JSON(400,
			map[string]string{
				"error": "can't delete this product(invalid id)",
			},
		)
		return
	}
	c.JSON(200,
		map[string]string{
			"message": "product deletet successfully",
		},
	)
}
