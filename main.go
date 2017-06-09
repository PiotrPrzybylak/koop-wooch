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

var shoppingCart = domain.ShoppingCart{Items: map[string]domain.CartItem{}}

var templates = template.Must(template.ParseFiles("templates/suppliers.html",
	"templates/supplier_form.html", "templates/categories.html",
	"templates/category_form.html", "templates/product_form.html",
	"templates/products.html", "templates/delivery_form.html",
	"templates/delivery.html", "templates/error.html",
	"templates/cart.html", "templates/home.html", "templates/footer.html"))

var store = sessions.NewCookieStore([]byte("something-very-very-secret"))

func main() {

	productService := domain.NewProductService(memory.NewProductRepository())
	deliveryService := domain.NewDeliverytService(memory.NewDeliveryRepository())

	supplierService := domain.NewSupplierService(memory.NewSupplierRepository())

	categoryService := domain.NewCategoryService(memory.NewCategoryRepository())

	r := mux.NewRouter()

	addExampleData(productService, supplierService, deliveryService, categoryService)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "home", nil)
	})

	r.HandleFunc("/add_product", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		category := r.URL.Query().Get("category")
		price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
		unit := r.URL.Query().Get("unit")
		quantity, _ := strconv.ParseFloat(r.URL.Query().Get("quantity"), 64)

		p := domain.Product{Name: name, Category: category, Unit: unit, Quantity: quantity, Price: price}
		_, err := productService.Create(p)

		if err != nil {
			renderTemplate(w, "error", err)
			return
		}

		http.Redirect(w, r, "/", 303)

	})

	r.HandleFunc("/product_form", func(w http.ResponseWriter, r *http.Request) {
		categories, err := categoryService.All()
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		renderTemplate(w, "product_form", categories)
	})

	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		products, err := productService.All()
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		renderTemplate(w, "products", products)
	})

	r.HandleFunc("/add_delivery", func(w http.ResponseWriter, r *http.Request) {
		supplier := r.URL.Query().Get("supplier")
		category := r.URL.Query().Get("category")
		price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
		unit := r.URL.Query().Get("unit")
		quantity, _ := strconv.ParseFloat(r.URL.Query().Get("quantity"), 64)

		d := domain.Delivery{Supplier: supplier, Category: category, Unit: unit, Quantity: quantity, Price: price}
		_, err := deliveryService.Create(d)

		if err != nil {
			renderTemplate(w, "error", err)
			return
		}

		http.Redirect(w, r, "/delivery", 303)
	})

	r.HandleFunc("/delivery_form", func(w http.ResponseWriter, r *http.Request) {

		suppliers, err := supplierService.ListAll()
		categories, err :=categoryService.All()

		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		type data struct{
			Suppliers []domain.Supplier
			Categories []domain.Category
		}
		renderTemplate(w, "delivery_form", data{suppliers, categories})
	})

	r.HandleFunc("/delivery", func(w http.ResponseWriter, r *http.Request) {
		deliverys, err := deliveryService.All()
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
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
					Product:  product,
					Quantity: 1,
				})
			}
		}
		// TODO show total amount of shopping cart

		http.Redirect(w, r, "/cart", 303)
	})

	r.HandleFunc("/suppliers", func(w http.ResponseWriter, r *http.Request) {
		suppliers, err := supplierService.ListAll()
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		renderTemplate(w, "suppliers", suppliers)
	})

	r.HandleFunc("/supplier_form", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "supplier_form", nil)
	})

	r.HandleFunc("/add_supplier", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		day := MustParseWeekday(r.URL.Query().Get("delivery_day"))
		_, err := supplierService.Create(domain.Supplier{"", name, day})
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		http.Redirect(w, r, "/suppliers", 303)
	})

	r.HandleFunc("/delete_supplier", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		err := supplierService.Delete(id)
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		http.Redirect(w, r, "/suppliers", 303)
	})

	r.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {

		type CategoriesAndProducts struct {
		}

		categories, err := categoryService.All()
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
		renderTemplate(w, "categories", categories)
	})
	r.HandleFunc("/category_form", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "category_form", nil)
	})

	r.HandleFunc("/add_category", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		c := domain.Category{Name: name}
		_, err := categoryService.Create(c)
		if err != nil {
			renderTemplate(w, "error", err)
			return
		}
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
		sum :=shoppingCart.Sum(map[string]domain.CartItem{})

		type cartData struct {
			Items  map[string]domain.CartItem
			Sum float64
		}

		renderTemplate(w, "cart",cartData{shoppingCart.Items, sum} )
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

func addExampleData(productService domain.ProductService, supplierService domain.SupplierService, deliveryService domain.DeliveryService, categoryService domain.CategoryService) {
	productService.Create(domain.Product{"", "Carrot", "Vegetables", 123, "piece", 100})
	productService.Create(domain.Product{"", "Apple", "Fruits", 666, "kg", 200})

	supplierService.Create(domain.Supplier{"", "Zdzisław Sztacheta", time.Monday})
	supplierService.Create(domain.Supplier{"", "Tesco", time.Friday})

	categoryService.Create(domain.Category{Name: "Vegetables"})
	categoryService.Create(domain.Category{Name: "Fruits"})

	deliveryService.Create(domain.Delivery{"Zdzisław Sztacheta", "Carrot", 123, "kg", 100})
	deliveryService.Create(domain.Delivery{"Tesco", "Apple", 666, "kg", 200})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
