package handlers

import (
	"fmt"
	"go-micro-service/data"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetAllProducts(rw, r)
		return
	}

	// for not supported methods
	fmt.Println("--------", r.Method)
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)

		return
	}

	if r.Method == http.MethodPut {
		fmt.Println("inside PUT ", r.URL.Path)
		var id int
		regx := regexp.MustCompile(`/([0-9]+)`)
		matches := regx.FindAllStringSubmatch(r.URL.Path, -1)
		fmt.Println(matches)
		if len(matches) != 1 {
			http.Error(rw, "Invalid URL", http.StatusBadRequest)
			return
		}

		if len(matches[0]) != 2 {
			http.Error(rw, "Invalid Parameters in URL ", http.StatusBadRequest)
			return
		}

		idString := matches[0][1]

		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Unable to extract ID for URI", http.StatusBadRequest)
			return
		}

		p.l.Println("Got it", id)

		p.UpdateProducts(rw, r)

	}

	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	p.l.Println(err)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling POST products")

	prod := &data.Product{}
	fmt.Println(r.Body)
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("PROD: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling PUT products")

	vars := mux.Vars(r)
	fmt.Println("ID : ", vars)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to extract ID for URI", http.StatusBadRequest)
		return
	}

	prod := &data.Product{}
	fmt.Println(r.Body)
	err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ProdNotFoundError {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}

}
