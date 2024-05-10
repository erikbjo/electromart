# ElectroMart ⚡️🛒

## Description 📜
ElectroMart is a web application that allows users to buy electronic products.
The application is built as a part of the course IDATG2204 Datamodellering og databasesystemer (2024 VÅR) at NTNU.

## Local Development 🛠
1. Install GO from https://golang.org/dl/

2. Install MySQL:
    - Linux: https://www.geeksforgeeks.org/how-to-install-mysql-on-linux/
    - MacOS: https://www.geeksforgeeks.org/how-to-install-mysql-on-macos/
    - Windows: https://www.geeksforgeeks.org/how-to-install-mysql-in-windows/

3. Login to MySQL:
    - Use the following command:
    ```bash
    mysql -u root -p
    ```
    - Enter your password
    - Create a database with the following command:
    ```sql
    CREATE DATABASE databasename;
    ```
    - Create a user with the following command:
    ```sql
    CREATE USER 'username'@'localhost' IDENTIFIED BY 'password';
    ```
    - Grant privileges to the user with the following command:
    ```sql
    GRANT ALL PRIVILEGES ON databasename
    TO 'username'@'localhost';
    ```
    - Create the tables by running the 2204Database.sql script in the database folder

4. Clone the repository:
    - Navigate to the directory where you want to clone the repository
    - Use the following command:
    ```bash
    git clone https://github.com/EskilAl/database_2024
    ```
   
5. Navigate to the project directory:
    ```bash
    cd database_2024
    ```
   
6. Build the project:
    ```bash
    go build
    ```
   
7. Run the project:
    ```bash
    ./database_2024
    ```

## Usage 🤓

### Website 🌐 
```plaintext
In the current version the user is able to do the following:
 - Create a user account
 - Update user information (e.g. name, address, email)
 - View all products
 - View products by category
 - Search for products
 - Add items to cart
 ```

---
### API Endpoints 🛠️

The service provides the following endpoints:

```plaintext
/api/v1/products/
/api/v1/categories/
/api/v1/brands/
```
Note: where the user needs to provide an ID, it can be found using the `View all` endpoints or the `Search` endpoints.

#### Products
The initial endpoint focuses on the management of products stored in the database.

##### Add new product to database

Manages the registration of new product and its details. This includes brand, category, quantity in stock, price name and description. Before adding the user has to make sure that both the category and brand that is referenced already exists in the database.

###### Request

```http
Method: POST
Path: /api/v1/products/
Content-Type: application/json
```

Body example:

```json lines
{
    "name": "Supatest",
    "brandName": "Supabrand",
    "categoryName": "Supacategory",
    "description": "Supatest Supatest Supadescription",
    "qtyInStock": 10,
    "price": 999
}
```

###### Response

The response to the POST request on the endpoint stores the product on the server and returns the associated ID. 

* Content type: application/json

Body (exemplary code for registered product):

```json
{
    "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6"
}
```

##### View a specific product

Enables retrieval of a specific product stored in the database.
###### Request

The following shows a request for an individual product identified by its ID.

```text
Method: GET
Path: /api/v1/products/{id}
```

* `id` is the ID associated with the specific product.

Example request:

```text
/api/v1/product/7acef38c-0d18-11ef-96c4-fa163ecc81b6
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json
{
    "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6",
    "name": "Supatest",
    "brandName": "Apple",
    "categoryName": "Smartphones",
    "description": "Supatest Supatest Supadescription",
    "qtyInStock": 10,
    "price": 999
}
```
##### Search for products

Lets the user search for products using a query. The query is used to find a match in either product name, brand name, category name or description. The search is not case-sensitive.

###### Request

A `GET` request to the endpoint should return all products in the search result.

```text
Method: GET
Path: /api/v1/products/search/{query}
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json lines
[
    {
        "id": "ca868913-0d16-11ef-96c4-fa163ecc81b6",
        "name": "iPhone 13",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Latest iPhone model",
        "qtyInStock": 100,
        "price": 9999.99
    },
    {
        "id": "ca86c85d-0d16-11ef-96c4-fa163ecc81b6",
        "name": "MacBook Pro",
        "brandName": "Apple",
        "categoryName": "Laptops",
        "description": "Latest MacBook model",
        "qtyInStock": 100,
        "price": 1299.99
    },
    {
        "id": "ca86d4fc-0d16-11ef-96c4-fa163ecc81b6",
        "name": "AirPods Pro",
        "brandName": "Apple",
        "categoryName": "Audio",
        "description": "Wireless earbuds with noise cancellation",
        "qtyInStock": 50,
        "price": 2499.99
    }
]
```

The response above is the result from the request url `/api/v1/products/search/apple` 

##### View all registered products

Enables retrieval of all products.

###### Request

A `GET` request to the endpoint should return all products including their IDs.

```text
Method: GET
Path: /api/v1/products/
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json lines
[
    {
        "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6",
        "name": "Supatest",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Supatest Supatest Supadescription",
        "qtyInStock": 10,
        "price": 999
    },
    {
        "id": "ca868913-0d16-11ef-96c4-fa163ecc81b6",
        "name": "iPhone 13",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Latest iPhone model",
        "qtyInStock": 100,
        "price": 9999.99
    }
]
```

The response should return a collection of return all stored products.

##### View all registered products by a Category

Enables retrieval of all products filtered by a category.

###### Request

A `GET` request to the endpoint should return all products including their IDs.

```text
Method: GET
Path: /api/v1/products/category/{category}
```

###### Response

* Content type: `application/json`

Example request: `/api/v1/products/category/smartphones`

Body (exemplary code):

```json lines
[
    {
        "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6",
        "name": "Supatest",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Supatest Supatest Supadescription",
        "qtyInStock": 10,
        "price": 999
    },
    {
        "id": "ca868913-0d16-11ef-96c4-fa163ecc81b6",
        "name": "iPhone 13",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Latest iPhone model",
        "qtyInStock": 100,
        "price": 9999.99
    }
]
```

The response should return a collection of return all products filtered by the provided category.

##### View all registered products by a Brand

Enables retrieval of all products filtered by a brand.

###### Request

A `GET` request to the endpoint should return all products including their IDs.

```text
Method: GET
Path: /api/v1/products/brand/{brand}
```

###### Response

* Content type: `application/json`

Example request: `/api/v1/products/brand/apple`

Body (exemplary code):

```json lines
[
    {
        "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6",
        "name": "Supatest",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Supatest Supatest Supadescription",
        "qtyInStock": 10,
        "price": 999
    },
    {
        "id": "ca868913-0d16-11ef-96c4-fa163ecc81b6",
        "name": "iPhone 13",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Latest iPhone model",
        "qtyInStock": 100,
        "price": 9999.99
    }
]
```

The response should return a collection of return all products filtered by the provided brand.

##### Update a specific product
Enables the replacing of specific product.

###### Request

The following shows a request for an updated product identified by its ID.

```
Method: PUT
Path: /api/v1/products/{id}
```

* `id` is the ID associated with the specific product and must match in the provided url and body.

Example request: ```/api/v1/products/ca86bc9b-0d16-11ef-96c4-fa163ecc81b6```

Body (exemplary code):

```json lines
{
        "id": "ca86bc9b-0d16-11ef-96c4-fa163ecc81b6",
        "name": "Galaxy Book Pro",
        "brandName": "Samsung",
        "categoryName": "Laptops",
        "description": "Latest Samsung laptop model",
        "qtyInStock": 100,
        "price": 1199.99
}
```

###### Response

This is the response to the change request.

* Status code: `204 No Content` if the update was successful.
* Body: empty

##### Delete a specific product

Enables the deletion of a specific product.

###### Request

The following shows a request for deletion of an individual product identified by its ID. This update should lead
to a deletion of the product from the server.

```text
Method: DELETE
Path: /api/v1/products/{id}
```

* `id` is the ID associated with the specific product.

Example request:

```text
/api/v1/products/ca86bc9b-0d16-11ef-96c4-fa163ecc81b6
```

###### Response

This is the response to the delete request.

* Status code: `204 No Content` if deletion is successful
* Body: empty

---

#### Categories

This endpoint focuses on the management of categories stored in the database.

##### Add new category to database

Manages the registration of new category and its details. This includes name and description.

###### Request

```http
Method: POST
Path: /api/v1/categories/
Content-Type: application/json
```

Body example:

```json lines
{
    "name": "Supacategory",
    "description": "Supa supa supa"
}
```

###### Response

The response to the POST request on the endpoint stores the category on the server and returns the associated ID. 

* Status Code: `204 No Content`
* Body: empty

##### View a specific category

Enables retrieval of a specific category stored in the database.

###### Request

The following shows a request for an individual category identified by its ID.

```text
Method: GET
Path: /api/v1/categories/{name}
```

* `name` is the name associated with the specific category. This is not case-sensitive.

Example request:

```text
/api/v1/category/smartphones
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json
{
    "name": "Smartphones",
    "description": "Mobile devices with advanced computing capabilities"
}
```

##### View all registered categories

Enables retrieval of all categories.

###### Request

A `GET` request to the endpoint should return all categories.

```text
Method: GET
Path: /api/v1/categories/
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json lines
[
    {
        "name": "Accessories",
        "description": "Additional items that complement electronic devices"
    },
    {
        "name": "Appliances",
        "description": "Electrical devices used for performing household tasks"
    },
    {
        "name": "Audio",
        "description": "Electronic devices for reproducing sound"
    },
]
```

The response should return a collection of all stored categories.

##### Update a specific category
Enables the replacing of specific category.

###### Request

The following shows a request for an updated category.

```
Method: PUT
Path: /api/v1/categories/{name}
```

* `name` is the name associated with the specific category and must match in the provided url and body.

Example request: ```/api/v1/categories/smartphones```

Body (exemplary code):

```json lines
{
    "name": "Smartphones",
    "description": "Mobile devices with advanced computing capabilities"
}
```

###### Response

This is the response to the change request.

* Status code: `204 No Content` if the update was successful.
* Body: empty

##### Delete a specific category

Enabling the deletion of a specific category.

###### Request

The following shows a request for deletion of an individual category identified by its name.

```text
Method: DELETE
Path: /api/v1/categories/{name}
```

* `name` is the name associated with the specific category.

Example request:

```text
/api/v1/categories/smartphones
```

###### Response

This is the response to the delete request.

* Status code: `204 No Content`.
* Body: empty

---

#### Brands

This endpoint focuses on the management of brands stored in the database.

##### Add new brand to database

Manages the registration of new brand and its details. This includes name and description.

###### Request

```http
Method: POST
Path: /api/v1/brands/
Content-Type: application/json
```

Body example:

```json lines
{
    "name": "Supabrand",
    "description": "Hes insane :D"
}
```

###### Response

The response to the POST request on the endpoint stores the brand on the server. 

* Status Code: `204 No Content`
* Body: empty

##### View a specific brand

Enables retrieval of a specific brand stored in the database.

###### Request

The following shows a request for an individual brand identified by its name.

```text
Method: GET
Path: /api/v1/brands/{name}
```

* `name` is the name associated with the specific brand. This is not case-sensitive.

Example request:

```text
/api/v1/brands/apple
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json
{
    "name": "Apple",
    "description": "American multinational technology company"
}
```

##### View all registered brands

Enables retrieval of all brands.

###### Request

A `GET` request to the endpoint should return all brands.

```text
Method: GET
Path: /api/v1/brands/
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json lines
[
    {
        "name": "Acer",
        "description": "Taiwanese multinational electronics company"
    },
    {
        "name": "Apple",
        "description": "American multinational technology company"
    },
    {
        "name": "Asus",
        "description": "Taiwanese multinational computer hardware and consumer electronics company"
    }
]
```

The response should return a collection of all stored brands.

##### Update a specific brand
Enables the replacing of specific brand.

###### Request

The following shows a request for an updated brand.

```
Method: PUT
Path: /api/v1/brands/{name}
```

* `name` is the name associated with the specific brand and must match in the provided url and body.

Example request: ```/api/v1/brands/Apple```

Body (exemplary code):

```json lines
{
        "name": "Apple",
        "description": "This is apple company"
}
```

###### Response

This is the response to the change request.

* Status code: `204 No Content` if the update was successful.
* Body: empty

##### Delete a specific brand

Enabling the deletion of a specific brand.

###### Request

The following shows a request for deletion of an individual brand identified by its name. This update should lead
to a deletion of the brand from the database.

```text
Method: DELETE
Path: /api/v1/brands/{name}
```

* `name` is the name associated with the specific brand.

Example request:

```text
/api/v1/brands/apple
```

###### Response

This is the response to the delete request.

* Status code: `204 No Content` if deletion is successful
* Body: empty

---


## Contact 📧

- [Eskil Alstad](mailto:eskil.alstad@ntnu.no)
- [Erik Bjørnsen](mailto:erbj@stud.ntnu.no)
- [Simon Husås Houmb](mailto:simon.h.houmb@ntnu.no)
- [Maja Melby](mailto:maja.melby@ntnu.no)
- [Jon André Solberg](mailto:jon.a.h.solberg@ntnu.no)
