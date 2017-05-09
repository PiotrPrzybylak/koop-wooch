package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Product struct {
	Name     string
	Category string
	Price    float64
}

var products = []Product{}

type Supplier struct {
	Name        string
	DeliveryDay time.Weekday
}

var suppliers = []Supplier{}

type Category struct {
	Name string
}

var categories = []Category{}

var templates = template.Must(template.ParseFiles("templates/suppliers.html", "templates/supplier_form.html", "templates/categories.html", "templates/category_form.html"))

func main() {

	_ = mux.NewRouter()

	addExampleData()

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		write(w, "<h2>Welcome to Koop!</h2>")
		write(w, " <a href=\"/product_form\">Add product</a>")
		write(w, " <a href=\"/products\">Show products</a>")
		write(w, " <a href='/supplier_form'>Add supplier</a>")
		write(w, " <a href='/suppliers'>Show suppliers</a>")
		write(w, " <a href='/category_form'>Add category</a>")
		write(w, " <a href='/categories'>Show categories</a>")

	})

	http.HandleFunc("/add_product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		name := r.URL.Query().Get("name")
		category := r.URL.Query().Get("category")
		price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)

		p := Product{Name: name, Category: category, Price: price}

		products = append(products, p)
		http.Redirect(w, r, "/", 303)

	})

	http.HandleFunc("/product_form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := "<form action='add_product'>PRODUCT<input name='name'>CATEGORY" +
			"<select name='category'>"
		for _, category := range categories {
			html += "<option>" + category.Name + "</option>"
		}
		html += "</select>" +
			"PRICE<input name='price'><input type='submit' value='Add'></form>"
		write(w, html)
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		write(w, "<table>")
		for _, product := range products {

			write(w, fmt.Sprintf("<tr><td>%v</td><td>%v</td><td>%v</td><td><form action='Put_in'><input type='hidden' name='name' value='%v'><input type='submit' value='Put'></form></td></tr>", product.Name, product.Category, product.Price, product.Name))

		}
		write(w, "</table>")

	})
	http.HandleFunc("/Put_in", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/suppliers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "suppliers", suppliers)
	})

	http.HandleFunc("/supplier_form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "supplier_form", nil)
	})

	http.HandleFunc("/add_supplier", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		day := MustParseWeekday(r.URL.Query().Get("delivery_day"))
		suppliers = append(suppliers, Supplier{name, day})
		http.Redirect(w, r, "/suppliers", 303)
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "categories", categories)
	})
	http.HandleFunc("/category_form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "category_form", nil)
	})

	http.HandleFunc("/add_category", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		categories = append(categories, Category{name})
		http.Redirect(w, r, "/categories", 303)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1234"
	}

	http.ListenAndServe("0.0.0.0:"+port, nil)

}

func write(w http.ResponseWriter, text string) {
	w.Write([]byte(text))
}

func addExampleData() {
	products = append(products, Product{"Carrot", "Vegetables", 123})
	products = append(products, Product{"Apple", "Fruits", 666})

	suppliers = append(suppliers, Supplier{"Zdzisław Sztacheta", time.Monday})
	suppliers = append(suppliers, Supplier{"Tesco", time.Friday})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MustParseWeekday(weekday string) time.Weekday {
	switch weekday {
	case "Monday":
		return time.Monday
	case "Tuesday":
		return time.Tuesday
	case "Wednesday":
		return time.Wednesday
	case "Thursday":
		return time.Thursday
	case "Friday":
		return time.Friday
	case "Saturday":
		return time.Saturday
	case "Sunday":
		return time.Sunday
	default:
		panic(fmt.Sprintf("Wrong weekday: %v", weekday))
	}
}
