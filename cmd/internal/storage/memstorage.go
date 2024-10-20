package storage

import (
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/domain/models"
	"github.com/google/uuid"
	"log"
)

type MemStorage struct {
	usersMap    map[string]models.User
	productsMap map[string]models.Product
}

func New() *MemStorage {
	uMap := make(map[string]models.User)
	pMap := make(map[string]models.Product)
	return &MemStorage{
		usersMap:    uMap,
		productsMap: pMap,
	}
}
func (ms *MemStorage) RegisterUser(user models.User) error {
	uid := uuid.New().String()
	ms.usersMap[uid] = user
	return nil
}
func (ms *MemStorage) LoginUser(user models.User) (string, error) {
	for uid, value := range ms.usersMap {
		if value.Email == user.Email {
			if value.Pass != user.Pass {
				return "", ErrInvalidLoginData
			}
			return uid, nil
		}
	}
	return "", ErrUserNotFound
}
func (ms *MemStorage) GetProducts() ([]models.Product, error) {
	var products []models.Product
	for pid, value := range ms.productsMap {
		product := value
		product.PID = pid
		products = append(products, product)
	}
	if len(products) == 0 {
		return nil, ErrProductsListEmpty
	}
	return products, nil
}
func (ms *MemStorage) GetProductByID(pId string) (models.Product, error) {
	log.Printf("PID: %s\n", pId)
	for _, val := range ms.productsMap {
		log.Println(val.Description, val.PID)
	}
	product, ok := ms.productsMap[pId]
	if !ok {
		return models.Product{}, ErrProductNotFound
	}
	return product, nil
}
func (ms *MemStorage) AddProduct(product models.Product) error {
	pId := uuid.New().String()
	ms.productsMap[pId] = product
	return nil
}
func (ms *MemStorage) Deleteproduct(pId string) error {
	_, ok := ms.productsMap[pId]
	if !ok {
		return ErrProductNotFound
	}
	delete(ms.productsMap, pId)
	return nil
}
