package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milkey coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
		DeletedOn:   "",
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and string coffee without milk",
		Price:       2.45,
		SKU:         "fj34",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
		DeletedOn:   "",
	},
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNexID()
	productList = append(productList, p)
}

func getNexID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func UpdateProduct(id int, prod *Product) error {
	// find if product is existed with that ID
	err, position, _ := findProductByID(id)
	if err != nil {
		return err
	}
	// perform the update
	prod.ID = id
	productList[position] = prod
	return nil

}

var ProdNotFoundError = fmt.Errorf("product not found")

func findProductByID(id int) (err error, position int, product *Product) {
	for indx, product := range productList {
		if product.ID == id {
			return nil, indx, product
		}
	}

	return ProdNotFoundError, -1, product
}
