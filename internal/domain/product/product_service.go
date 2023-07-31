package product

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
)

type ProductService interface {
	Create(requestFormat ProductRequestFormat) (product Product, err error)
	ResolveByID()
}

type ProductServiceImpl struct {
	ProductRepository ProductRepository
	Config            *configs.Config
}

func ProvideProductServiceImpl(productRepository ProductRepository, config *configs.Config) *ProductServiceImpl {
	s := new(ProductServiceImpl)
	s.ProductRepository = productRepository
	s.Config = config

	return s
}

func (s *ProductServiceImpl) Create(requestFormat ProductRequestFormat) (product Product, err error) {
	product, err = product.NewFromRequestFormat(requestFormat)
	if err != nil {
		return
	}

	if err != nil {
		return product, failure.BadRequest(err)
	}

	err = s.ProductRepository.Create(product)
	if err != nil {
		return
	}

	return
}
