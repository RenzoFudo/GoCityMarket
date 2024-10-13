package server

import (
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type storage interface {
	GetUserProfile(models.User) error
	RegisterUser(models.User) (string, error)
	LoginUser(models.User) (string, error)
	GetProductByID(models.Product) error
	AddProduct(models.Product) error
	GetProducts() ([]models.Product, error)
	UpdateProduct(models.Product) error
	DeleteProduct(string) error
	GetAllProducts() ([]models.Product, error)
	GetUserPurchases(models.Purchase) error
	MakePurchase(models.Purchase) error
	GetProductPurchases(models.Purchase) error
}

type Server struct {
	host   string
	router *gin.Engine
}

func New(host string) *Server {
	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register") //(username, password, email): регистрация нового пользователя
		userGroup.POST("/login")    //(username, password): авторизация пользователя.
		userGroup.GET("/")          //(userID): получение информации о профиле пользователя
	}
	productGroup := r.Group("/product")
	{
		productGroup.POST("/addProd")     //(name, description, price, quantity): добавление нового товара
		productGroup.POST("/delete/:PID") // (productID): удаление товара.
		productGroup.POST("/updateProd")  // (productID, name, description, price, quantity): редактирование информации о товаре.
		productGroup.GET("/getPID")       //(productID): получение информации о конкретном товаре
		productGroup.GET("/allProd")      //(name, description, price, quantity): добавление нового товара.
	}
	purchaseGroup := r.Group("/purchase")
	{
		purchaseGroup.GET("/getPurID")   //(userID, productID, quantity): совершение покупки товара.
		purchaseGroup.GET("/getUserPur") //(userID): получение списка всех покупок пользователя.
		purchaseGroup.GET("/getProdPur") //(productID): получение списка всех покупок по конкретному товару
	}
	return &Server{
		Host:   host,
		router: r,
	}
}
func (s *Server) Run() {
	s.router.Run(s.host)
	if err := s.router.Run(); err != nil {
		return err
	}
	return nil

}
