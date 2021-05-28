package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	host string
	port string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "Host on which to run")
	flag.StringVar(&port, "port", "8080", "Port on which to run")
}

type Categories struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Catid  string `json: "catid"`
	Name   string `json: "name"`
	Image  string `json: "image"`
	Descn  string
	BanImg string
}

type Banners struct {
	Banners []Banner `json:"banners"`
}

type Banner struct {
	Favcategory string `json:"favcategory"`
	Bannername  string `json:"bannername"`
	Descn       string `json:"descn"`
	Image       string `json:"image"`
}

type Product struct {
	ProductId   string `json:"productId"`
	CategoryId  string `json:"categoryId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Products struct {
	Products []Product `json:"products"`
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func forbidden(w http.ResponseWriter, r *http.Request) {
	// see http://golang.org/pkg/net/http/#pkg-constants
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("403 HTTP status code returned!"))
}

func find_banners() Banners {
	url := "http://banners:8080/banners"

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data Banners
	json.Unmarshal(body, &data.Banners)
	fmt.Printf("Results: %v\n", data)

	return data
}

func find_products(id string) Products {
	url := "http://products:8080/products"

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data Products
	json.Unmarshal(body, &data.Products)
	fmt.Printf("Results: %v\n", data)

	var result Products

	for i := 0; i < len(data.Products); i++ {
		if string(data.Products[i].CategoryId) == id {
			//data.Products = append(data.Products[:i], data.Products[i+1:]...)
			p := Product{}
			p.ProductId = data.Products[i].ProductId
			p.CategoryId = data.Products[i].CategoryId
			p.Name = data.Products[i].Name
			p.Description = data.Products[i].Description

			fmt.Printf("Results: %v\n", p)
			result.Products = append(result.Products, p)
		}
	}

	fmt.Printf("Results: %v\n", data)

	return result
}

func find() Categories {

	categories := Categories{}
	banners := Banners{}
	data, err := ioutil.ReadFile("data/categories.json")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened categories.json")

	err = json.Unmarshal(data, &categories)
	if err != nil {
		log.Fatal(err)
	}

	banners = find_banners()

	for i := range categories.Categories {
		for x := range banners.Banners {
			if string(categories.Categories[i].Catid) == banners.Banners[x].Favcategory {
				categories.Categories[i].Descn = banners.Banners[x].Descn
				categories.Categories[i].BanImg = banners.Banners[x].Image
			}
		}
	}

	return categories
}

func one(id string) Category {

	categories := Categories{}
	category := Category{}
	data, err := ioutil.ReadFile("data/categories.json")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Opened categories.json")

	err = json.Unmarshal(data, &categories)
	if err != nil {
		log.Fatal(err)
	}

	for i := range categories.Categories {
		//fmt.Println("|",categories.Categories[i].Catid , "|==|" , id,"|")
		if string(categories.Categories[i].Catid) == id {
			//fmt.Println(categories.Categories[i])
			category = categories.Categories[i]
		}
	}

	return category
}

func findAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getCategory")

	categories := find()

	output, err := json.Marshal(categories.Categories)
	//output, err := json.MarshalIndent(categories.Categories, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
	fmt.Println(string(output))

	fmt.Println("Endpoint Hit: getCategory")
}

func findOne(w http.ResponseWriter, r *http.Request) {
	//id := strings.TrimPrefix(r.URL.Path, "/categories/")

	v := strings.Split(string(r.URL.Path), "/")

	id := v[2]

	if len(v) > 3 {
		handler := v[3]
		fmt.Println(handler)
		products := find_products(id)

		output, err := json.Marshal(products.Products)
		//output, err := json.MarshalIndent(categories.Categories, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
		fmt.Println(string(output))

	} else {
		categories := one(id)

		output, err := json.Marshal(categories)
		//output, err := json.MarshalIndent(categories.Categories, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(output)
		fmt.Println(string(output))
	}
}

func handleRequests() {

	http.HandleFunc("/favicon.ico", doNothing)
	http.HandleFunc("/categories", findAll)
	http.HandleFunc("/categories/", findOne)

	address := ":" + port

	log.Println("Starting server on address", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	handleRequests()
}
