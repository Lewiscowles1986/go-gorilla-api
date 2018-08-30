package repositories

import (
	"database/sql"

	"../data"
)

func GetProduct(db *sql.DB, id string) (data.Product, error) {
	return data.ParseProductData(
		db.QueryRow("SELECT id, name, price FROM products WHERE id=?", id))
}

func UpdateProduct(db *sql.DB, id string, p data.Product) error {
	_, err :=
		db.Exec("UPDATE products SET name=?, price=? WHERE id=?",
			p.GetName(), p.GetPrice(), id)

	return err
}

func DeleteProduct(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM products WHERE id=?", id)

	return err
}

func CreateProduct(db *sql.DB, p data.Product) error {
	_, err := db.Exec(
		"INSERT INTO products(id, name, price) VALUES(?, ?, ?)",
		p.GetID(), p.GetName(), p.GetPrice())

	if err != nil {
		return err
	}

	return nil
}

func GetProducts(db *sql.DB, page uint64, count uint8) ([]data.Product, error) {
	if page > 0 {
		page--
	}
	pageOffset := page * uint64(count)
	rows, err := db.Query(
		"SELECT id, name, price FROM products LIMIT ? OFFSET ?",
		count, pageOffset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return data.ParseProductListData(rows)
}

func GetProductCount(db *sql.DB) uint64 {
	i := uint64(0)
	r := db.QueryRow("SELECT COUNT(id) FROM products")
	err := r.Scan(&i)
	if err != nil {
		i = uint64(0)
	}
	return i
}
