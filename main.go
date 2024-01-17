package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
)

var db *sql.DB

type Product struct {
	ID    int
	Name  string
	Price int
}

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	db = sdb

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	// INSERT PROD
	// err = createProduct(&Product{Name: "Go product", Price: 444})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Create Successfully!")

	// GET PROD
	// p, err := getProduct(1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Get Successfully!", p)

	// UPDATE PROD
	// product, err := updateProduct(3, &Product{Name: "Goduct", Price: 123})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Update Successfully!", product)

	// DELETE PROD
	// err = deleteProduct(3)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Delete Successfully!")

	// GET PRODS
	products, err := getProducts()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Get Successfully!", products)
}

func createProduct(product *Product) error {
	_, err := db.Exec(
		"INSERT INTO public.products(name, price) VALUES ($1, $2);",
		product.Name,
		product.Price,
	)

	return err
}

func getProduct(id int) (Product, error) {
	var p Product
	row := db.QueryRow(
		"SELECT id, name, price FROM public.products WHERE id=$1;",
		id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func getProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products;")
	if err != nil {
		return nil, err
	}

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func updateProduct(id int, product *Product) (Product, error) {
	// _, err := db.Exec(
	// 	"UPDATE public.products SET name=$1, price=$2 WHERE id=$3;",
	// 	product.Name,
	// 	product.Price,
	// 	id,
	// )
	var p Product
	row := db.QueryRow(
		"UPDATE public.products SET name=$1, price=$2 WHERE id=$3 RETURNING id, name, price;",
		product.Name,
		product.Price,
		id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, err
	}

	return p, err
}

func deleteProduct(id int) error {
	_, err := db.Exec(
		"DELETE FROM public.products WHERE id=$1;",
		id,
	)
	return err
}
