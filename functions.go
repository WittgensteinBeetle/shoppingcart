package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func itemizedlist(w http.ResponseWriter, r *http.Request) {
	//set header to json payload
	w.Header().Set("Content-Type", "application/json")
	var fullList []List

	//obtain parameter from request
	params := mux.Vars(r)
	customer := params["cust_id"]

	//query the customer id shopping list
	result, err := db.Query("SELECT id, food_id, type, quantity,customer_id,coupon from groceries_list WHERE customer_id=?", customer)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var list List
		err := result.Scan(&list.ID, &list.Food_id, &list.Type, &list.Quantity, &list.Customer_id, &list.Coupon)
		if err != nil {
			panic(err.Error())
		}
		fullList = append(fullList, list)
	}
	json.NewEncoder(w).Encode(fullList)
}

func itemizedtaxtotal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var fullList []Total
	var query string

	//obtain parameter from request
	params := mux.Vars(r)
	customer := params["cust_id"]

	//else only report metrics without tax information

	//read the body of the HTTP request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//store the body results in json string map
	//payload needs to be unmarshaled
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	//variable for new quantity sent
	isTax := keyVal["isTaxable"]

	if isTax == "true" {
		//is_taxable specified
		fmt.Println("Found true taxable value ")
		query = `select products.name,price*quantity as 'Sub Total', 
		products.price * quantity * if(is_taxable=1,0.0812,0) as 'Tax',
		products.price * quantity * if(coupon > 0,coupon,0) as 'Coupon',
		price*quantity + products.price * quantity * if(is_taxable=1,0.0825,0) - products.price * quantity *  if(coupon > 0,coupon,0)  as 'GT' 
		from groceries_list join products on groceries_list.food_id=products.food_id 
		where customer_id=?`

		result, err := db.Query(query, customer)
		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var total Total
			err := result.Scan(&total.Name, &total.SubTotal, &total.Tax, &total.Coupon, &total.GrandTotal)
			if err != nil {
				panic(err.Error())
			}
			fullList = append(fullList, total)
		}

	} else {
		query = `select products.name,price*quantity as 'Sub Total', 
		products.price * quantity * if(coupon > 0,coupon,0) as 'Coupon', 
		price*quantity  - products.price * quantity *  if(coupon > 0,coupon,0)  as 'GT' 
		from groceries_list join products on groceries_list.food_id=products.food_id 
		where customer_id=?`
		//how does the app hold/interact with non-taxable items?  These could be removed by adding products.is_taxable=1
		//in the filter conditions though this depends on the data access method

		result, err := db.Query(query, customer)
		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var total Total
			err := result.Scan(&total.Name, &total.SubTotal, &total.Coupon, &total.GrandTotal)
			if err != nil {
				panic(err.Error())
			}
			fullList = append(fullList, total)
		}

	}

	json.NewEncoder(w).Encode(fullList)

}

func fulltaxtotal(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var fullList []Total
	var query string

	//obtain parameter from request
	params := mux.Vars(r)
	customer := params["cust_id"]

	//else only report metrics without tax information

	//read the body of the HTTP request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//store the body results in json string map
	//payload needs to be unmarshaled
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	//variable for new quantity sent
	isTax := keyVal["isTaxable"]

	if isTax == "true" {
		//is_taxable specified
		fmt.Println("Found true taxable value ")
		query = `select sum(GT) as 'Grand Total' 
		from (select products.name,price*quantity as 'Sub Total', 
		products.price * quantity * if(is_taxable=1,0.0825,0) as 'Tax',
		price*quantity + products.price * quantity * if(is_taxable=1,0.0825,0) as 'GT' 
		      from groceries_list join products on groceries_list.food_id=products.food_id 
		      where customer_id=?) as grand_total`
		result, err := db.Query(query, customer)
		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var total Total
			err := result.Scan(&total.GrandTotal)
			if err != nil {
				panic(err.Error())
			}
			fullList = append(fullList, total)
		}

	} else {
		query = `select sum(GT) as 'Grand Total' 
		from (select products.name,price*quantity as 'Sub Total', 
		products.price * quantity * if(is_taxable=1,0.0825,0) as 'Tax',
		price*quantity - products.price * quantity * if(is_taxable=1,0.0812,0) as 'GT' 
		      from groceries_list join products on groceries_list.food_id=products.food_id 
		      where customer_id=?) as grand_total`
		result, err := db.Query(query, customer)
		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var total Total
			err := result.Scan(&total.GrandTotal)
			if err != nil {
				panic(err.Error())
			}
			fullList = append(fullList, total)
		}

	}

	json.NewEncoder(w).Encode(fullList)

}

func couponAdd(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//obtain parameter from request
	params := mux.Vars(r)
	customer := params["cust_id"]

	//read the body of the HTTP request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//store the body results in json string map
	//payload needs to be unmarshaled
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	//variable for new quantity sent
	couponValue := keyVal["couponValue"]
	var coup_id int
	var cart_coup int
	var checkcart int
	var coup_val float32

	//initialize logic values
	checkcart = 0

	//query used for coupon search in cart
	q_coup_in_cart := `select if (c.food_id is NULL,0,c.food_id),
	if(c.value is NULL, 0, trim(c.value)) 
	from groceries_list gl left join coupons c on (gl.food_id=c.food_id and c.coup_id=?) 
	where gl.customer_id=?`

	//query used for whether coupon exists and value
	q_coupon_count := `select if(d.coup_id is NULL, 0,d.coup_id)  
	from coupons c left join coupons d on (c.coup_id=d.coup_id and c.coup_id=?)`

	//Obtain coupon information
	rows, err := db.Query(q_coupon_count, couponValue)
	if err != nil {
		panic(err.Error())
	}

	//check to see whether coupon is valid
	for rows.Next() {
		rows.Scan(&coup_id)
		if err != nil {
			panic(err.Error())
		}
		if coup_id > 0 {
			fmt.Fprintf(w, "\nCoupon Found")
			checkcart = 1
			break

		}

	}

	if coup_id == 0 {
		fmt.Fprintf(w, "\nCoupon expired")
	}

	//see whether coupon item is in cart
	lrows, err := db.Query(q_coup_in_cart, couponValue, customer)
	if err != nil {
		panic(err.Error())
	}

	for lrows.Next() && checkcart == 1 {
		lrows.Scan(&cart_coup, &coup_val)
		if err != nil {
			panic(err.Error())
		}
		if cart_coup > 0 {
			fmt.Fprintf(w, "\nFound in cart")
			fmt.Println(cart_coup, customer, coup_val)

			//input sanitization needed
			//update customer shopping cart items with coupon
			_, err = db.Query("update groceries_list set coupon=? where food_id=? and customer_id=?", coup_val, cart_coup, customer)
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(w, "\nCoupon applied")

			break

		}

	}

	if cart_coup == 0 {
		fmt.Fprintf(w, "\nCoupon is not found in cart")
	}

}

func itemAdd(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//obtain parameter from request
	params := mux.Vars(r)
	customer := params["cust_id"]

	//read the body of the HTTP request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//store the body results in json string map
	//payload needs to be unmarshaled
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	//obtain data from json payload
	foodId := keyVal["food_id"]
	quantity := keyVal["amount"]
	itemCat := keyVal["type"]

	//variables used in function
	var food_id int

	//obtain food_id from payload to see whether it is already in cart
	q_food_id_check := `
		select food_id from groceries_list 
               where food_id=? and customer_id=?`

	result, err := db.Query(q_food_id_check, foodId, customer)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		result.Scan(&food_id)
		if err != nil {
			panic(err.Error())
		}
		if food_id > 0 {
			fmt.Fprintf(w, "\nItem already in cart")
			break
		}
	}
	if food_id == 0 {
		fmt.Fprintf(w, "\nItem can be inserted in cart")
		qitemInsert := `insert into groceries_list(type,quantity,customer_id,food_id) values (?,?,?,?)`
		db.Query(qitemInsert, itemCat, quantity, customer, foodId)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintf(w, "\nItem added to cart")

	}
	//used for debugging purposes to inspect payload
	fmt.Println(customer, foodId, food_id, quantity, itemCat)

	fmt.Fprintf(w, "\nEnd")

}

func itemRemove(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//obtain parameter from request
	params := mux.Vars(r)
	customer := params["cust_id"]

	//read the body of the HTTP request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//store the body results in json string map
	//payload needs to be unmarshaled
	// Note: the payload here is an int value, will need to verify this
	// in other payloads and format a api standard with verification
	keyVal := make(map[string]int)
	json.Unmarshal(body, &keyVal)

	//variables used in function
	var food_id int
	var payloadfoodID int
	var removeFood int

	//obtain data from json payload
	payloadfoodID = keyVal["food_id"]

	//obtain food_id from payload to see whether it is already in cart
	q_food_id_check := `select food_id from groceries_list where food_id=? and customer_id=?`

	result, err := db.Query(q_food_id_check, payloadfoodID, customer)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		result.Scan(&food_id)
		if err != nil {
			panic(err.Error())
		}
		if food_id > 0 {
			fmt.Fprintf(w, "\nItem in cart")
			removeFood = 1
			break
		}
	}
	if food_id > 0 && removeFood == 1 {
		fmt.Fprintf(w, "\nItem can be removed from cart")
		qitemRemove := "delete from groceries_list where food_id=? and customer_id=?"
		db.Query(qitemRemove, food_id, customer)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintf(w, "\nItem removed from cart")
		removeFood = 0

	}

	//used for debugging purposes to inspect payload variables
	fmt.Println(customer, payloadfoodID, food_id)

	fmt.Fprintf(w, "\nItem is not in cart")

}
