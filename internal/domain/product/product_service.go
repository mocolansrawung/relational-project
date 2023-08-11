package product

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type ProductService interface {
	Create(requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error)
	Update(id uuid.UUID, requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error)
	SoftDelete(id uuid.UUID, userID uuid.UUID) (product Product, err error)
	ResolveByID(id uuid.UUID) (product Product, err error)
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

func (s *ProductServiceImpl) Create(requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error) {
	product, err = product.NewFromRequestFormat(requestFormat, userID)
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
func (s *ProductServiceImpl) Update(id uuid.UUID, requestFormat ProductRequestFormat, userID uuid.UUID) (product Product, err error) {
	product, err = s.ProductRepository.ResolveByID(id)
	if err != nil {
		return
	}

	err = product.Update(requestFormat, userID)
	if err != nil {
		return
	}

	err = s.ProductRepository.Update(product)
	return
}
func (s *ProductServiceImpl) SoftDelete(id uuid.UUID, userID uuid.UUID) (product Product, err error) {
	product, err = s.ProductRepository.ResolveByID(id)
	if err != nil {
		return
	}

	images, err := s.ProductRepository.ResolveImagesByProductIDs([]uuid.UUID{product.ID})
	if err != nil {
		return product, err
	}
	variants, err := s.ProductRepository.ResolveVariantsByProductIDs([]uuid.UUID{product.ID})
	if err != nil {
		return product, err
	}

	product.AttachImages(images)
	product.AttachVariants(variants)

	err = product.SoftDelete(userID)
	if err != nil {
		return
	}

	err = s.ProductRepository.Update(product)

	return
}
func (s *ProductServiceImpl) ResolveByID(id uuid.UUID) (product Product, err error) {
	product, err = s.ProductRepository.ResolveByID(id)

	if product.IsDeleted() {
		return product, failure.NotFound("foo")
	}

	images, err := s.ProductRepository.ResolveImagesByProductIDs([]uuid.UUID{product.ID})
	if err != nil {
		return product, err
	}
	variants, err := s.ProductRepository.ResolveVariantsByProductIDs([]uuid.UUID{product.ID})
	if err != nil {
		return product, err
	}

	product.AttachImages(images)
	product.AttachVariants(variants)

	return
}
