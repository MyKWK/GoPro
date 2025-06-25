package services

import (
	"awesomeProject/datamodels"
	"awesomeProject/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProducts() ([]*datamodels.Product, error)
	DeleteProduct(int64) bool
	InsertProduct(*datamodels.Product) (int64, error)
	UpdateProduct(*datamodels.Product) error
}

type ProductService struct {
	ProductRepo repositories.IProduct
}

// 结构体的初始化函数
func NewProductService(productRepo repositories.IProduct) *ProductService {
	return &ProductService{
		ProductRepo: productRepo,
	}
}

func (p *ProductService) GetProductByID(id int64) (*datamodels.Product, error) {
	return p.ProductRepo.SelectByKey(id)
}

func (p *ProductService) GetAllProducts() ([]*datamodels.Product, error) {
	return p.ProductRepo.SelectAll()
}

func (p *ProductService) DeleteProduct(id int64) bool {
	ok, _ := p.ProductRepo.Delete(id)
	return ok
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.ProductRepo.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	_, err := p.ProductRepo.Update(product)
	return err
}
