package logic

import (
	"fmt"

	"github.com/amupxm/pure-webserver/domain"
	"github.com/amupxm/pure-webserver/repository"
)

type (

	// ProductLogic is the business logic for products
	ProductLogic interface {
		NewProduct(product *domain.Product) (*domain.Product, error)
		GetProductByID(id string) (*[]domain.Product, error)
		GetAllProducts() (*[]domain.Product, error)
		DeleteProduct(product *domain.Product) error
		UpdateProduct(product *domain.Product) (*domain.Product, error)
	}
	productLogic struct {
		productRepository repository.ProductRepository
	}
)

func NewProductLogic(productRepository repository.ProductRepository) ProductLogic {
	return &productLogic{
		productRepository: productRepository,
	}
}

// NewProduct creates a new product
func (pl *productLogic) NewProduct(product *domain.Product) (*domain.Product, error) {
	// TODO : check for duplicated iid
	result, err := pl.productRepository.CreateProduct(product)
	if err != nil {
		fmt.Println(1, result)

		return nil, err
	}
	fmt.Println(1, result)

	return result, nil
}

// GetProductByID returns product with iid
func (pl *productLogic) GetProductByID(id string) (*[]domain.Product, error) {
	list, err := pl.productRepository.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (pl *productLogic) GetAllProducts() (*[]domain.Product, error) {
	var list *domain.Product
	result, err := pl.productRepository.GetAllProducts(list)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (pl *productLogic) UpdateProduct(product *domain.Product) (*domain.Product, error) {
	result, err := pl.productRepository.UpdateProduct(product)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (pl *productLogic) DeleteProduct(product *domain.Product) error {
	return pl.productRepository.DeleteProduct(product.Iid)
}
