package repository

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/amupxm/pure-webserver/constants"
	"github.com/amupxm/pure-webserver/domain"
	"github.com/amupxm/pure-webserver/pkg/database"
)

// import (
// 	"github.com/amupxm/pure-webserver/domain"
// 	"github.com/amupxm/pure-webserver/pkg/database"
// )

type (
	// ProductRepository is the interface for product repository
	ProductRepository interface {
		// CreateProduct writes a new product to the database
		CreateProduct(product *domain.Product) (*domain.Product, error)
		// GetProductByID gets a product by id
		GetProductByID(id string) ([]domain.Product, error)
		// GetAllProducts gets all products
		GetAllProducts(product *domain.Product) (*[]domain.Product, error)
		// UpdateProduct updates a product
		UpdateProduct(product *domain.Product) (*domain.Product, error)
		// DeleteProduct deletes a product
		DeleteProduct(iid string) error
	}
	productRepository struct {
		db database.Database
	}
)

// NewProductRepository creates a new product repository
func NewProductRepository(db database.Database) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// GetProductByID gets a product by id
func (pl *productRepository) GetProductByID(id string) ([]domain.Product, error) {
	var product *domain.Product
	var result []domain.Product
	products, err := pl.db.GetFromCollection(product)
	if err != nil {
		return result, err
	}
	if products == nil || len(products) == 0 {

		return nil, errors.New(constants.NoData)

	}
	productFilteredList := products.Where("iid", id)
	log.Print(productFilteredList)
	//marshal and unmarshal to get the struct
	s, err := json.Marshal(productFilteredList)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(s, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// CreateProduct creates a new product
func (pl *productRepository) CreateProduct(product *domain.Product) (*domain.Product, error) {
	err := pl.db.WriteToCollection(product)
	return product, err
}

// GetAllProducts gets all products
func (pl *productRepository) GetAllProducts(product *domain.Product) (*[]domain.Product, error) {
	var result *[]domain.Product
	res, err := pl.db.GetFromCollection(product)
	if err != nil {
		return result, err
	}
	s, err := json.Marshal(res)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(s, &result)
	if err != nil {
		return result, err
	}
	return result, err
}

// UpdateProduct updates a product
func (pl *productRepository) UpdateProduct(product *domain.Product) (*domain.Product, error) {
	var allProducts *[]domain.Product
	res, err := pl.db.GetFromCollection(product)
	if err != nil {
		return product, err
	}
	s, err := json.Marshal(res)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(s, &allProducts)
	if err != nil {
		return product, err
	}
	var tempProducts []domain.Product
	// find product by same iid
	for _, v := range *allProducts {
		if v.Iid == product.Iid {
			// update product
			v.Name = product.Name
			v.Brand = product.Brand
			v.Company = product.Company
			v.UpdatedAt = product.UpdatedAt
			product = &v
		}
		tempProducts = append(tempProducts, v)
	}
	err = pl.db.UpdateCollection(&tempProducts)
	if err != nil {
		return product, err
	}
	return product, nil
}

// DeleteProduct deletes a product
func (pl *productRepository) DeleteProduct(iid string) error {
	var allProducts *[]domain.Product
	var product *domain.Product
	res, err := pl.db.GetFromCollection(product)
	if err != nil {
		return err
	}
	s, err := json.Marshal(res)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(s, &allProducts)
	if err != nil {
		return err
	}
	var tempProducts []domain.Product
	// find product by same iid
	for _, v := range *allProducts {
		if v.Iid != iid {
			tempProducts = append(tempProducts, v)
		}
	}
	err = pl.db.UpdateCollection(&tempProducts)
	return err
}
