## Shopping Cart Project 

Goal:
Project for building a small shopping cart API endpoint for json payloads 

## Features 

##### Calculate itemized list of grand and subtotal with optional inclusion of tax on items in shopping cart 

##### Remove/add items to cart 

##### Using established coupon codes apply coupons to existing cart 


## Go file structure and database schema  

```
---shoppingcart
|--api.go
|--functions.go 
|--main.go
```

Assumptions 
This project assumes the endpoint will be using a secure network. 
The products list will be static and also established in the database before hand 
along with the available coupon codes and associated values 

Customers will be authenticated in a seperate system and given a unique customer-id value 

Backend: MySQL 8.0.33 

Schema: 

```
* coupons 

CREATE TABLE `coupons` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type` varchar(200) NOT NULL,
  `value` float NOT NULL,
  `coup_id` bigint NOT NULL,
  `food_id` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_coupd_id` (`coup_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


* groceries_list

 CREATE TABLE `groceries_list` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type` varchar(100) NOT NULL,
  `quantity` bigint NOT NULL,
  `customer_id` bigint NOT NULL,
  `food_id` int NOT NULL,
  `coupon` float NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_id` (`customer_id`,`food_id`),
  KEY `idx_food_id` (`food_id`),
  KEY `idx_customer_id` (`customer_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

* products 

 CREATE TABLE `products` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `price` decimal(10,2) DEFAULT '0.00',
  `is_taxable` tinyint(1) DEFAULT NULL,
  `name` varchar(200) DEFAULT NULL,
  `food_id` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


```

## Test Cases 

The customer being verified and the list established from a GUI is a large missing piece of this project to say the least. 
I would like to build more upon this perhaps with a small ncurses text interface as the project builds  

Customer ID: 42 

#### Itemized list of shopping cart


Test endpoint: http://localhost:8001/cart_items/42
Method: GET 

* Expected results sample: 

```
[
    {
        "id": 2,
        "food_id": 4,
        "type": "protien",
        "quantity": 2,
        "coupon": 0.2,
        "customer_id": 42
    },
    {
        "id": 1,
        "food_id": 8,
        "type": "condiment",
        "quantity": 3,
        "coupon": 0.1,
        "customer_id": 42
    },
    {
        "id": 7,
        "food_id": 10,
        "type": "snack",
        "quantity": 4,
        "customer_id": 42
    }
]
```
#### Itemized tax total 

Test endpoint: http://localhost:8001/cart_items_tax/42
Method: PUT 

The goal of this was to provide an itemized list with the total,subtotal, and grandtotal of each item depending on whether 
the isTaxable json payload is set to true else false.  I would like to build upon this, the value should be a true boolean 
rather than a string, though this is not required. When a coupon is applied to the item it will include the total value deducted  

* Expected json payload

```
{
"isTaxable": "true"
}

```
* Expected Results
```
[
    {
        "name": "Boneless Chicken Breasts",
        "subTotal": 24,
        "tax": "0.000000",
        "coupon": 4.8,
        "grandtotal": "19.199999928474426"
    },
    {
        "name": "Stir Fry Sauce",
        "subTotal": 24,
        "tax": "1.948800",
        "coupon": 2.4,
        "grandtotal": "23.579999964237214"
    },
    {
        "name": "Cookies",
        "subTotal": 8,
        "tax": "0.000000",
        "grandtotal": "8"
    }
]
```

#### Entire cart grandtotal values 

Test endpoint: http://localhost:8001/cart_items_full/42
Method: PUT 

The goal of this endpoint was to provide a non itemized grandtotal of the entire cart with an optional inclusion of tax. 
The use case would be aking to "tax free" shopping weekends. 

* Expected json payload

```
{
"isTaxable": "true"
}
```
* Expected results

```
[
    {
        "grandtotal": "57.980000"
    }
]
```

#### Send existing coupon code to be applied to cart 

Test endpoint: http://localhost:8001/cart_items_coupon/42
Method: PUT 

The use case here is rather simple at the moment.  It will take the 'couponValue' from the json payload the search the 
MySQL backend for a specific type and value.  At the moment this is limited to type "percept" with a given value 
to apply to a unique food identification number.  Though I would like to build upon this in the future.  

When the coupon is not valid is an item is not in the customers cart it will return a message to the client 

* Expected json payload

```
{
"couponValue": "78777"
}

```
* Expected results

```
Item found in cart
Coupon applied 
```
  
This can be tested with the full itemized list endpoints listed above.  Another method for localized testing is the monitoring of the MySQL
database tables and potentially the general query log depending on traffic levels and bench marktesting conditions 

#### Add items to cart 

Test endpoint: http://localhost:8001/cart_items_add/42
Method: POST 

This will check to see whether the provided food identification number, type and quantity is valid.  Should the item not exist in the 
customers cart it will be added.  There is a unique key constraint in effect for food_id,customer_id that ensures each item has a single entry 

*Expected json payload 

```
{
    "food_id" : "10",
    "amount" : "4" ,
    "type"  : "snack"
}

```
* Expected results
```
Item can be inserted in cart
Item added to cart
```
This can be tested with the full itemized list endpoints listed above.  Another method for localized testing is the monitoring of the MySQL
database tables and potentially the general query log depending on traffic levels and bench marktesting conditions 


#### Delete items from cart 

Test endpoint: http://localhost:8001/cart_items_remove/42
Method: POST

This will check to see whether the item is an existing valid food_id value, then whether it exists within the cart.  
Should both these criterias be met it will then remove this from the customers shopping cart.  

* Expected json payload
```
{
"food_id" : 10
}
```

*Expected results: 

```
Item in cart
Item can be removed from cart
Item removed from cart
```
This can be tested with the full itemized list endpoints listed above.  Another method for localized testing is the monitoring of the MySQL
database tables and potentially the general query log depending on traffic levels and bench marktesting conditions 




