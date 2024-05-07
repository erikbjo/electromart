package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

/*
SearchProducts retrieves rows from the Product table in the database based on the query string provided.
The query string is used to search for products by name, description, brand name, and category name.
Returns a slice of Product structs if successful, or an error if not.
*/
func SearchProducts(query string) ([]structs.Product, error) {
	lowerQuery := strings.ToLower(query) // Convert query to lowercase
	rows, err := Client.Query(
		`SELECT * FROM Product
      WHERE LOWER(Name) LIKE ? OR Description LIKE ? OR BrandName IN (SELECT ID FROM Brand WHERE LOWER(Name) LIKE ?)
        AND CategoryName IN (SELECT ID FROM Category WHERE LOWER(Name) LIKE ?);
  `, "%"+lowerQuery+"%", "%"+lowerQuery+"%", "%"+lowerQuery+"%", "%"+lowerQuery+"%",
	)
	if err != nil {
		log.Println("Error when querying for products: ", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err2 := rows.Close()
		if err2 != nil {
			log.Println("Error when closing rows: ", err2)
		}
	}(rows)

	products, err3 := rowsToProductSlice(rows)
	if err3 != nil {
		log.Println("Error when converting rows to slice: ", err3)
		return nil, err3
	}
	return products, nil
}

/*
GetAllProducts retrieves all rows from the Product table in the database and returns them as a slice of Product structs.
*/
func GetAllProducts() ([]structs.Product, error) {
	rows, err := Client.Query(`SELECT * FROM Product`)
	if err != nil {
		log.Println("Error when selecting all products: ", err)
		return nil, err
	}

	foundProducts, err2 := rowsToProductSlice(rows)
	if err2 != nil {
		log.Println("Error when converting rows to slice: ", err2)
		return nil, err2
	}
	return foundProducts, nil
}

/*
GetProductByID retrieves a single row from the Product table in the database based on the ID and returns it as a Product struct.
*/
func GetProductByID(id string) (*structs.Product, error) {
	exists, err := productExists(id)
	if err != nil {
		log.Println("Error when checking if product exists: ", err)
		return nil, err
	}
	if exists {
		var product structs.Product
		err2 := Client.QueryRow("SELECT * FROM Product WHERE ID = ?", id).Scan(
			&product.ID,
			&product.Name,
			&product.BrandName,
			&product.CategoryName,
			&product.Description,
			&product.QtyInStock,
			&product.Price,
		)
		if err2 != nil {
			log.Println("Error when selecting product by ID: ", err2)
			return nil, err2
		}
		return &product, nil
	} else {
		log.Println("Product does not exist")
		return nil, fmt.Errorf("product ID does not match any products in DB")
	}
}

/*
AddProduct adds a single row to the Product table in the database. Returns the ID if successful, or an error if not.
*/
func AddProduct(product structs.Product) (string, error) {
	// Generate and retrieve new UUID
	id, err := GenerateUUID(Client)
	if err != nil {
		log.Println("Error generating UUID: ", err)
		return "", err
	}

	// Insert product
	_, err2 := Client.Exec(
		`INSERT INTO Product (ID, Name, BrandName, CategoryName, Description, QtyInStock, 
Price) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id,
		product.Name,
		product.BrandName,
		product.CategoryName,
		product.Description,
		product.QtyInStock,
		product.Price,
	)
	if err2 != nil {
		log.Println("Error inserting product: ", err2)
		return "", err2
	}

	return id, nil // Return the UUID
}

/*
UpdateProduct updates a single row in the Product table in the database based on the ID in the provided Product struct.
Returns nil if successful, or an error if not.
*/
func UpdateProduct(product structs.Product) error {
	// Check if product exists
	exists, err := productExists(product.ID)
	if err != nil {
		log.Println("Error when checking if product exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec(
			`UPDATE Product SET Name = ?, BrandName = ?, CategoryName = ?, Description = ?, QtyInStock = ?, 
Price = ? WHERE ID = ?`,
			product.Name,
			product.BrandName,
			product.CategoryName,
			product.Description,
			product.QtyInStock,
			product.Price,
			product.ID,
		)
		if err2 != nil {
			log.Println("Error when updating product: ", err2)
			return err2
		}
	} else {
		log.Println("Product does not exist")
		return fmt.Errorf("product ID does not match any products in DB")
	}
	return nil
}

/*
DeleteProductByID deletes a single row from the Product table in the database based on the ID.
Returns nil if successful, or an error if not.
*/
func DeleteProductByID(id string) error {
	exists, err := productExists(id)
	if err != nil {
		log.Println("Error when checking if product exists: ", err)
		return err
	}
	if exists {
		_, err2 := Client.Exec("DELETE FROM Product WHERE ID = ?", id)
		if err2 != nil {
			log.Println("Error when deleting product: ", err2)
			return err2
		}
	} else {
		log.Println("Product does not exist")
		return fmt.Errorf("product ID does not match any products in DB")
	}
	return nil
}

/*
rowsToProductSlice converts the rows from a SQL query to a slice of Product structs.
*/
func rowsToProductSlice(rows *sql.Rows) ([]structs.Product, error) {
	var productSlice []structs.Product
	for rows.Next() {
		var product structs.Product
		err2 := rows.Scan(
			&product.ID,
			&product.Name,
			&product.BrandName,
			&product.CategoryName,
			&product.Description,
			&product.QtyInStock,
			&product.Price,
		)
		if err2 != nil {
			return nil, err2
		}
		productSlice = append(productSlice, product)
	}
	return productSlice, nil
}

/*
ProductExists checks if a product with the provided ID exists in the database.
*/
func productExists(id string) (bool, error) {
	var exists bool
	err := Client.QueryRow(`SELECT EXISTS(SELECT * FROM Product WHERE ID = ?)`, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
