package server

import (
	"github.com/RenzoFudo/GoCityMarket/cmd/internal/domain/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Storage interface {
	RegisterUser(models.User) error
	LoginUser(models.User) (string, error)
	GetProducts() ([]models.Product, error)
	GetProductByID(pId string) (models.Product, error)
	AddProduct(models.Product) error
	DeleteProduct(string) error

	UpdateProduct(pId string) (models.Product, error)
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
		userGroup.POST("/login", s.loginHandler)       //(username, password): авторизация пользователя.
		userGroup.GET("/", s.ProfileHandler)           //(userID): получение информации о профиле пользователя
	}
	productGroup := r.Group("/product")
	{
		productGroup.POST("/addProd", s.AddProductHandler)        //(name, description, price, quantity): добавление нового товара
		productGroup.POST("/delete/:PID", s.DeleteProductHandler) // (productID): удаление товара.
		productGroup.POST("/updateProd", s.UpdateProductHandler)  // (productID, name, description, price, quantity): редактирование информации о товаре.
		productGroup.GET("/getPID", s.ProductHandler)             //(productID): получение информации о конкретном товаре
		productGroup.GET("/allProd", s.ProductsHandler)           //(name, description, price, quantity): добавление нового товара.
	}
	purchaseGroup := r.Group("/purchase")
	{
		purchaseGroup.GET("/getPurID", s.MakePurchaseHandler)       //(userID, productID, quantity): совершение покупки товара.
		purchaseGroup.GET("/getUserPur", s.UserPurchasesHandler)    //(userID): получение списка всех покупок пользователя.
		purchaseGroup.GET("/getProdPur", s.ProductPurchasesHandler) //(productID): получение списка всех покупок по конкретному товару
	}
	if err := r.Run(s.host); err != nil {
		return err
	}
	return nil
}

func (s *Server) RegisterHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (s *Server) loginHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.storage.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) ProfileHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	profile, err := s.storage.GetUserProfile(models.User{Token: token})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (s *Server) ProductsHandler(c *gin.Context) {
	products, err := s.storage.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (s *Server) ProductHandler(c *gin.Context) {
	id := c.Query("id")

	product, err := s.storage.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (s *Server) AddProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.AddProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

func (s *Server) DeleteProductHandler(c *gin.Context) {
	id := c.Param("id")

	err := s.storage.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (s *Server) UpdateProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := s.storage.UpdateProduct(product.PID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func (s *Server) UserPurchasesHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	purchases, err := s.storage.GetUserPurchases(models.Purchase{Token: token})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"purchases": purchases})
}

func (s *Server) MakePurchaseHandler(c *gin.Context) {
	var purchase models.Purchase
	if err := c.ShouldBindBodyWithJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.MakePurchase(purchase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Purchase created"})
}

func (s *Server) ProductPurchasesHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing product id"})
		return
	}

}

/*
func (s *Server) RegisterHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.storage.RegisterUser(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()()})
		return
	}
	ctx.JSON(http.StatusOK, "user is registered")
}

func (s *Server) loginHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := s.storage.LoginUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidLoginData) {
			ctx.String(http.StatusUnauthorized, err.Error()())
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()()})
		return
	}
	ctx.JSON(http.StatusOK, "login completed")
}
*/
/*
func (ms *MemStorage) GetProductByID(product models.Product) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MemStorage) DeleteProduct(s string) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MemStorage) UpdateProduct(product models.Product) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MemStorage) GetUserProfile(user models.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (ms *MemStorage) GetAllProducts() ([]models.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (ms *MemStorage) GetUserPurchases(purchase models.Purchase) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MemStorage) MakePurchase(purchase models.Purchase) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MemStorage) GetProductPurchases(purchase models.Purchase) error {
	//TODO implement me
	panic("implement me")
}
*/
