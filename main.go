package main

import (
	"footlocker-bot/internal"
	"footlocker-bot/internal/shared"
)

func main() {
	footlocker := internal.NewFootlockerBot()

	task := shared.Task{
		ProfileName: "ASM",
		ProductURL:  "https://www.footlocker.com/product/~/7895HMBC.html",
		Size:        "10.0",
		ProductID:   "270963",
		Quantity:    1,
		UseProxy:    false,
		Mode:        "",
		Aco:         false,
		Region:      "",
		Store:       "",
		Keywords:    "",
		Sku:         "7895HMBC",
		Payment:     "",
		Profile: shared.Profile{
			ProfileName: "ASM Profile",
			FirstName:   "mam",
			LastName:    "mam",
			Age:         22,
			BirthDay:    29,
			BirthMonth:  3,
			BirthYear:   2000,
			Gender:      "m",
			Email:       "tvtv8047@gmail.com",
			Phone:       "+44 020 8677 6161",
			Address:     "95 Thirsk Road",
			Address2:    "",
			Zip:         "SW11 5SU",
			City:        "Battersea",
			Country:     "United Kingdom",
			CountryISO:  "GB",
			CountryCode: "GBR",
			State:       "",
			Cardname:    "Kai Avila",
			Cnb:         "4109703255583065",
			Month:       "01",
			Year:        "2027",
			Cvv:         "340",
			CardType:    "visa",
			Password:    "Tv@Tv80",
		},
		Id: 0,
	}

	err := footlocker.GetFootlockerSettings(shared.Settings{})
	if err != nil {
		return
	}

	// product url
	footlocker.Log.Info(task.ProductURL)

	// GetHome
	GetHomeStatus, _ := footlocker.GetHome()
	footlocker.Log.Info(string(GetHomeStatus))

	// GetProduct
	GetProductStatus, _ := footlocker.GetProduct(task)
	footlocker.Log.Info(string(GetProductStatus))

	// AddToCart
	AddToCartStatus, _ := footlocker.AddToCart(task)
	footlocker.Log.Info(string(AddToCartStatus))

}
