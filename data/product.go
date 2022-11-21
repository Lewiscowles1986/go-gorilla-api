package data

import (
	"encoding/json"

	"database/sql"

	"github.com/satori/go.uuid"
)

type product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Product interface {
	GetID() string
	GetName() string
	GetPrice() float64
	ChangeName(newName string)
	SetPrice(newPrice float64)
}

func (p *product) GetID() string {
	if len(p.ID) != 36 {
		p.ID = uuid.Must(uuid.NewV4(), nil).String()
	}
	return p.ID
}

func (p *product) GetName() string {
	return p.Name
}

func (p *product) GetPrice() float64 {
	return p.Price
}

func (p *product) ChangeName(newName string) {
	p.Name = newName
}

func (p *product) SetPrice(newPrice float64) {
	p.Price = newPrice
}

func CreateProduct(name string, price float64) Product {
	p := &product{}
	p.ID = uuid.Must(uuid.NewV4(), nil).String()
	p.Name = name
	p.Price = price
	return p
}

func ParseProductDataJSON(data []byte) (Product, error) {
	product := &product{}
	if err := json.Unmarshal(data, product); err != nil {
		return nil, err
	}
	return product, nil
}

func ParseProductListData(rs *sql.Rows) ([]Product, error) {
	products := []Product{}

	for rs.Next() {
		p := &product{}
		if err := rs.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func ParseProductData(r *sql.Row) (Product, error) {
	p := &product{}
	err := r.Scan(&p.ID, &p.Name, &p.Price)
	return p, err
}
