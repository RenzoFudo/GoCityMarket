package models

type User struct {
	UID string `json:"

"` // юзер id
	Name  string `json:"name" validate:"required"`         // юзер нейм
	Email string `json:"email" validate:"required, email"` // электронная почта
	Pass  string `json:"pass" validate:"required"`         // юзер пароль
}

type Product struct {
	PID         string `json:"P_id"`                             // id продукта
	PName       string `json:"Product_name" validate:"required"` // имя продукта
	Description string `json:"Description" validate:"required"`  // описание продукта
	Price       string `json:"price" validate:"required"`        // цена продукта
	Quantity    string `json:"quantity" validate:"required"`     // количество продукта
}

type Purchase struct {
	PurID       string `json:"pur_id"`       // уникальный идентификатор покупки
	PurUID      string `json:"pur_uid"`      // идентификатор пользователя, совершившего покупку
	PurPID      string `json:"pur_pid"`      // идентификатор купленного товара
	PurQuantity string `json:"pur_quantity"` // количество купленного товара
	Timestamp   string `json:"timestamp"`    // временная метка покупки
}
