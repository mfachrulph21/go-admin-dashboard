package models

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName  string         `gorm:"required; not null" valid:"required~full_name: field is required" json:"full_name"`
	Email     string         `gorm:"required; not null; unique" valid:"required~email: field is required, email" json:"email"`
	Password  string         `gorm:"required; not null" valid:"required~password: field is required, minstringlength(6)" json:"password"`
	Role      string         `gorm:"required; not null" valid:"role" json:"role"`
	Balance   int            `gorm:"required; not null" valid:"range(0|100000000)" json:"balance"`
	DeletedAt gorm.DeletedAt `gorm:"-"`
}

type Category struct {
	gorm.Model
	Type              string         `gorm:"required; not null" json:"type"`
	SoldProductAmount int            `json:"sold_product_amount"`
	DeletedAt         gorm.DeletedAt `gorm:"-"`
	Products          []Product      `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	gorm.Model
	Title      string         `gorm:"required; not null" valid:"required~title: field is required" json:"title"`
	Price      int            `gorm:"required; not null" valid:"required~price: field is required, range(0|50000000)" json:"price"`
	Stock      int            `gorm:"required; not null" valid:"required~stock: field is required, minimalStock" json:"stock"`
	CategoryID uint           `valid:"required~category_Id: field is required" json:"category_Id"`
	DeletedAt  gorm.DeletedAt `gorm:"-"`
	Category   Category       `json:"-"`
}

type TransactionHistory struct {
	gorm.Model
	ProductID  uint           `json:"product_id"`
	UserID     uint           `json:"user_id"`
	Quantity   int            `gorm:"required; not null" json:"quantity"`
	TotalPrice int            `gorm:"required; not null" json:"total_price"`
	DeletedAt  gorm.DeletedAt `gorm:"-"`

	Product Product `gorm:"foreignKey:ProductID"`
	User    User    `gorm:"foreignKey:UserID"`
}

type LoginUser struct {
	Email    string `valid:"required~email: field is required" json:"email"`
	Password string `valid:"required~password: field is required" json:"password"`
}

type TopupBody struct {
	Balance int `valid:"required~balance: field is required, range(0|100000000)"`
}

type EditCategoryBody struct {
	Type string `valid:"required~type: field is required" json:"type"`
}

type EditProductBody struct {
	gorm.Model
	Title      string `valid:"required~title: field is required" json:"title"`
	Price      int    `valid:"required~price: field is required, range(0|50000000)" json:"price"`
	Stock      int    `valid:"required~stock: field is required, minimalStock" json:"stock"`
	CategoryID uint   `valid:"required~category_Id: field is required" json:"category_Id"`
}

type PostTransactionBody struct {
	ProductId uint `valid:"required~product_id: field is required" json:"product_id"`
	Quantity  int  `valid:"required~quantity: field is required" json:"quantity"`
}

type DataToken struct {
	ID    uint
	Email string
	Role  string
	jwt.RegisteredClaims
}
