package server

import (
	"errors"
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/domain/models"
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Storage interface {
	RegisterUser(models.User) error
	LoginUser(models.User) (string, error)
	GetProducts() ([]models.Product, error)
	GetProductByID(models.Product) error
	AddProduct(models.Product) error
	DeleteProduct(string) error

	UpdateProduct(models.Product) error
	GetUserProfile(models.User) (string, error)
	GetAllProducts() ([]models.Product, error)
	GetUserPurchases(models.Purchase) error
	MakePurchase(models.Purchase) error
	GetProductPurchases(models.Purchase) error
}

type Server struct {
	host    string
	storage Storage
}

func New(host string, storage Storage) *Server {
	return &Server{
		host:    host,
		storage: storage,
	}
}

func (s *Server) Run() error {
	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", s.RegisterHandler) //(username, password, email): регистрация нового пользователя
		userGroup.POST("/login")                       //(username, password): авторизация пользователя.
		userGroup.GET("/")                             //(userID): получение информации о профиле пользователя
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
	if err := r.Run(s.host); err != nil {
		return err
	}
	return nil
}

func (s *Server) RegisterHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.storage.RegisterUser(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "user is registered")
}

func (s *Server) loginHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := s.storage.LoginUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidLoginData) {
			ctx.String(http.StatusUnauthorized, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "login completed")
}
