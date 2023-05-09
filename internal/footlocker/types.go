package footlocker

type AddToCartPayload struct {
	Size            string `json:"size"`
	Sku             string `json:"sku"`
	ProductQuantity int    `json:"productQuantity"`
	FulfillmentMode string `json:"fulfillmentMode"`
	ResponseFormat  string `json:"responseFormat"`
}