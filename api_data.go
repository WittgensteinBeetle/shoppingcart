package main

type List struct {
	ID          int     `json:"id"`
	Food_id     int     `json:"food_id"`
	Type        string  `json:"type"`
	Quantity    int     `json:"quantity"`
	Coupon      float32 `json:"coupon,omitempty"`
	Customer_id int     `json:"customer_id"`
}

type Total struct {
	Name       string  `json:"name,omitempty"`
	SubTotal   float32 `json:"subTotal,omitempty"`
	Tax        string  `json:"tax,omitempty"`
	Coupon     float32 `json:"coupon,omitempty"`
	GrandTotal string  `json:"grandtotal,omitempty"`
}
