package main

import (
	"footlocker-bot/internal"
	"footlocker-bot/internal/logger"
	"footlocker-bot/internal/shared"
	"time"
)

func main() {
	l := logger.NewLogger()
	footlocker := internal.NewFootlockerBot()
	// datadome := internal.NewDatadome()
	l.EnableDebug()

	task := shared.Task{
		ProfileName: "ASM",
		ProductURL:  "https://www.footlocker.com/product/~/7895HMBC.html",
		Size:        "10.0",
		ProductID:   "270963",
		Quantity:    1,
		UseProxy:    false,
		Mode:        "",
		Aco:         false,
		Region:      "NM",
		Store:       "",
		Keywords:    "",
		Sku:         "7895HMBC",
		Payment:     "",
		Profile: shared.Profile{
			ProfileName: "ASM Profile",
			FirstName:   "Peter",
			LastName:    "Valaxar",
			Age:         22,
			BirthDay:    29,
			BirthMonth:  3,
			BirthYear:   2000,
			Gender:      "m",
			Email:       "asm.dev29@gmail.com",
			Phone:       "2025961737",
			Address:     "S Santa Monica St",
			Address2:    "1721",
			Zip:         "88030",
			City:        "DEMING",
			Country:     "United Kingdom",
			CountryISO:  "US",
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
		l.Debug("Failed to GetFootlockerSettings: ", err)
		return
	}

	// product url
	footlocker.Log.Info("productUrl: ", task.ProductURL)

	// GetHome
	GetHomeStatus, _ := footlocker.GetHome(task)
	footlocker.Log.Info("GetHomeStatus: ", GetHomeStatus)

	time.Sleep(1 * time.Second)

	// GetProduct
	GetProductStatus, _ := footlocker.GetProduct(task)
	footlocker.Log.Info("GetProductStatus: ", GetProductStatus)

	time.Sleep(1 * time.Second)

	// TimeStamp
	TimeStampStatus, csrfToken, _ := footlocker.TimeStamp(task)
	footlocker.Log.Info("TimeStampStatus: ", TimeStampStatus)
	footlocker.Log.Debug("csrfToken: ", csrfToken)

	time.Sleep(2 * time.Second)

	// AddToCart
	AddToCartStatus, _ := footlocker.AddToCart(task)
	footlocker.Log.Info("ATC Status: ", AddToCartStatus)

	if AddToCartStatus != 200 {
		footlocker.Log.Error("ATC status not ok", nil)
		return
	}

	// GetCheckoutpgae
	GetCheckoutPageStatus, _ := footlocker.GetCheckoutPage(task)
	footlocker.Log.Info("GetCheckoutPageStatus: ", GetCheckoutPageStatus)

	// SubmitUserInfo
	SubmitUserInfoStatus, _ := footlocker.SubmitUserInfo(task)
	footlocker.Log.Info("SubmitUserInfoStatus: ", SubmitUserInfoStatus)

	// AddAddress
	AddAddressStatus, _ := footlocker.AddAddress(task)
	footlocker.Log.Info("AddAddressStatus: ", AddAddressStatus)

	// VerifyAddress
	VerifyAddressStatus, _ := footlocker.VerifyAddress(task, csrfToken)
	footlocker.Log.Info("VerifyAddressStatus: ", VerifyAddressStatus)

	// SubmitVerifiedAddress
	SubmitVerifiedAddressStatus, _ := footlocker.SubmitVerifiedAddress(task)
	footlocker.Log.Info("SubmitVerifiedAddressStatus: ", SubmitVerifiedAddressStatus)

	// GetAdyen
	GetAdyenStatus, publicKey, _ := footlocker.GetAdyen(task)
	footlocker.Log.Info("GetAdyenStatus: ", GetAdyenStatus)
	footlocker.Log.Info("publicKey: ", publicKey)

}
