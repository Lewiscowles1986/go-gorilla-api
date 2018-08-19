package repositories

import (
    "database/sql"

    "../data"
)

func GetProduct(db *sql.DB, id string) (data.Product, error) {
    return data.ParseProductData(
        db.QueryRow("SELECT id, name, price FROM products WHERE id=$1", id))
}

func UpdateProduct(db *sql.DB, id string, p data.Product) error {
    _, err :=
        db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
            p.GetName(), p.GetPrice(), id)

    return err
}

func DeleteProduct(db *sql.DB, id string) error {
    _, err := db.Exec("DELETE FROM products WHERE id=$1", id)

    return err
}

func CreateProduct(db *sql.DB, p data.Product) error {
    _, err := db.Exec(
        "INSERT INTO products(id, name, price) VALUES($1, $2, $3)",
        p.GetID(), p.GetName(), p.GetPrice())

    if err != nil {
        return err
    }

    return nil
}

func GetProducts(db *sql.DB, start uint64, count uint8) ([]data.Product, error) {
    rows, err := db.Query(
        "SELECT id, name, price FROM products LIMIT $1 OFFSET $2",
        count, start)

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
