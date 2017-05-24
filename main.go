package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PiotrPrzybylak/koop-wooch/domain"
	"github.com/PiotrPrzybylak/koop-wooch/infracture/persistance/memory"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Delivery struct {
	Supplier string
	Category string
	Price    float64
	Unit     string
	Quantity float64
}

var deliverys = []Delivery{}

var shoppingCart = domain.ShoppingCart{Items: map[string]domain.CartItem{}}

type Supplier struct {
	Name        string
	DeliveryDay time.Weekday
}

var suppliers = []Supplier{}

type Category struct {
	Name string
}

var categories = []Category{}

var templates = template.Must(template.ParseFiles("templates/suppliers.html",
	"templates/supplier_form.html", "templates/categories.html",
	"templates/category_form.html", "templates/product_form.html",
	"templates/products.html", "templates/delivery_form.html",
	"templates/delivery.html", "templates/error.html",
	"templates/cart.html"))

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))

func main() {

	productService := domain.NewProductService(memory.NewProductRepository())

	r := mux.NewRouter()

	addExampleData(productService)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		write(w, "<h2>Welcome to Koop!</h2>")
		write(w, " <a href=\"/product_form\">Add product</a>")
		write(w, " <a href=\"/products\">Show products</a>")
		write(w, " <a href='/supplier_form'>Add supplier</a>")
		write(w, " <a href='/suppliers'>Show suppliers</a>")
		write(w, " <a href='/category_form'>Add category</a>")
		write(w, " <a href='/categories'>Show categories</a>")
		write(w, " <a href='/delivery_form'>Add delivery</a>")
		write(w, " <a href='/delivery'>Show delivery</a>")
		write(w, " <a href='/cart'>Show cart items</a>")

	})

	r.HandleFunc("/add_product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		name := r.URL.Query().Get("name")
		category := r.URL.Query().Get("category")
		price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
		unit := r.URL.Query().Get("unit")
		quantity, _ := strconv.ParseFloat(r.URL.Query().Get("quantity"), 64)

		p := domain.Product{Name: name, Category: category, Unit: unit, Quantity: quantity, Price: price}
		_, err := productService.Create(p)

		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			renderTemplate(w, "error", err)
			return
		}

		http.Redirect(w, r, "/", 303)

	})

	r.HandleFunc("/product_form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "product_form", categories)
	})

	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		products, err := productService.All()
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		renderTemplate(w, "products", products)
	})

	r.HandleFunc("/add_delivery", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		supplier := r.URL.Query().Get("supplier")
		category := r.URL.Query().Get("category")
		price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
		unit := r.URL.Query().Get("unit")
		quantity, _ := strconv.ParseFloat(r.URL.Query().Get("quantity"), 64)

		d := Delivery{Supplier: supplier, Category: category, Unit: unit, Quantity: quantity, Price: price}

		deliverys = append(deliverys, d)
		http.Redirect(w, r, "/", 303)
	})

	r.HandleFunc("/delivery_form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "delivery_form", deliverys)
	})

	r.HandleFunc("/delivery", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "delivery", deliverys)
	})

	r.HandleFunc("/Put_in", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		products, err := productService.All()
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		for _, product := range products {
			if name == product.Name {
				shoppingCart.Add(domain.CartItem{
					Product: product,
					Quantity: 1,
				})
			}
		}

		// TODO show total amount of shopping cart

		//http.Redirect(w, r, "/shopping_cart", 303)
	})

	r.HandleFunc("/suppliers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "suppliers", suppliers)
	})

	r.HandleFunc("/supplier_form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "supplier_form", nil)
	})

	r.HandleFunc("/add_supplier", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		day := MustParseWeekday(r.URL.Query().Get("delivery_day"))
		suppliers = append(suppliers, Supplier{name, day})
		http.Redirect(w, r, "/suppliers", 303)
	})

	r.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {

		type CategoriesAndProducts struct {
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "categories", categories)
	})
	r.HandleFunc("/category_form", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "category_form", nil)
	})

	r.HandleFunc("/add_category", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		categories = append(categories, Category{name})
		http.Redirect(w, r, "/categories", 303)
	})

	r.HandleFunc("/show_session", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		write(w, "<ul>")
		for key, value := range session.Values {
			write(w, fmt.Sprintf("<li>%v : %v</li>", key, value))
		}
		write(w, "</ul>")
	})

	r.HandleFunc("/add_session_param", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		name := r.URL.Query().Get("name")
		value := r.URL.Query().Get("value")

		session.Values[name] = value

		session.Save(r, w)

		http.Redirect(w, r, "/show_session", 303)

	})

	r.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		renderTemplate(w, "cart", shoppingCart.Items)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1234"
	}

	http.ListenAndServe("0.0.0.0:"+port, r)

}

func write(w http.ResponseWriter, text string) {
	w.Write([]byte(text))
}

func addExampleData(productService domain.ProductService) {
	productService.Create(domain.Product{"", "Carrot", "Vegetables", 123, "piece", 100})
	productService.Create(domain.Product{"", "Apple", "Fruits", 666, "kg", 200})

	suppliers = append(suppliers, Supplier{"Zdzisław Sztacheta", time.Monday})
	suppliers = append(suppliers, Supplier{"Tesco", time.Friday})

	categories = append(categories, Category{"Vegetables"})
	categories = append(categories, Category{"Fruits"})

	deliverys = append(deliverys, Delivery{"Zdzisław Sztacheta", "Carrot", 123, "kg", 100})
	deliverys = append(deliverys, Delivery{"Tesco", "Apple", 666, "kg", 200})
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
