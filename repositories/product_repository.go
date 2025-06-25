package repositories

import (
	"awesomeProject/datamodels"
	"gorm.io/gorm"
)

type IProduct interface {
	Insert(product *datamodels.Product) (int64, error)         // 插入产品，返回插入的ID
	Delete(productID int64) (bool, error)                      // 删除产品，返回是否成功
	Update(product *datamodels.Product) (int64, error)         // 更新产品，返回影响的行数
	SelectByKey(productKey int64) (*datamodels.Product, error) // 根据ID查询单个产品
	SelectAll() ([]*datamodels.Product, error)                 // 查询所有产品
}

type ProductManager struct {
	db *gorm.DB
}

// NewProductManager 结构体的初始化函数
func NewProductManager(db *gorm.DB) IProduct {
	return &ProductManager{
		db: db,
	}
}

// Insert 插入新产品
func (p *ProductManager) Insert(product *datamodels.Product) (int64, error) {
	result := p.db.Create(product)
	return product.ID, result.Error
}

// Delete 删除产品
func (p *ProductManager) Delete(productID int64) (bool, error) {
	result := p.db.Delete(&datamodels.Product{}, productID)
	return result.RowsAffected > 0, result.Error
}

// Update 更新产品信息
func (p *ProductManager) Update(product *datamodels.Product) (int64, error) {
	result := p.db.Save(product)
	return result.RowsAffected, result.Error
}

// SelectByKey 根据ID查询单个产品
func (p *ProductManager) SelectByKey(productKey int64) (*datamodels.Product, error) {
	var product datamodels.Product
	err := p.db.First(&product, productKey).Error
	return &product, err
}

// SelectAll 查询所有产品
func (p *ProductManager) SelectAll() ([]*datamodels.Product, error) {
	var products []*datamodels.Product
	err := p.db.Find(&products).Error
	return products, err
}
