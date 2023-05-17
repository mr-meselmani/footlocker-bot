package main

import (
	"footlocker-bot/internal"
	"footlocker-bot/internal/logger"
	"footlocker-bot/internal/shared"
)

func main() {
	l := logger.NewLogger()
	footlocker := internal.NewFootlockerBot()

	l.EnableDebug()

	task := shared.Task{
		ProfileName:   "ASM",
		ProductURL:    "https://www.footlocker.com/product/~/38019001.html",
		Size:          "10.0",
		ProductID:     "",
		Quantity:      1,
		UseProxy:      false,
		Mode:          "",
		Aco:           false,
		Region:        "NY",
		RegionIsocode: "US-NY",
		Store:         "",
		Keywords:      "",
		Sku:           "37581101",
		Payment:       "",
		Profile: shared.Profile{
			ProfileName: "ASM Profile",
			FirstName:   "ASM",
			LastName:    "DEV",
			Age:         1,
			BirthDay:    1,
			BirthMonth:  1,
			BirthYear:   1,
			Gender:      "",
			Email:       "asm.dev29@gmail.com",
			Phone:       "2025961737",
			Address:     "83 Pendergast Street",
			Address2:    "",
			Zip:         "11225",
			City:        "BROOKLYN",
			Country:     "",
			CountryISO:  "US",
			CountryCode: "",
			State:       "",
			Cardname:    "Kai Avila",
			Cnb:         "4109703255583065",
			Month:       "01",
			Year:        "2027",
			Cvv:         "340",
			CardType:    "visa",
			Password:    "",
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
	GetHomeStatus, cid, _ := footlocker.GetHome(task)
	footlocker.Log.Info("GetHomeStatus: ", GetHomeStatus)
	footlocker.Log.Info("cid: ", cid)

	// GetProduct
	GetProductStatus, _ := footlocker.GetProduct(task, cid)
	footlocker.Log.Info("GetProductStatus: ", GetProductStatus)

	// TimeStamp
	TimeStampStatus, csrfToken, _ := footlocker.TimeStamp(task, cid)
	footlocker.Log.Info("TimeStampStatus: ", TimeStampStatus)
	footlocker.Log.Debug("csrfToken: ", csrfToken)

	// AddToCart
	AddToCartStatus, _ := footlocker.AddToCart(task, cid)
	footlocker.Log.Info("ATC Status: ", AddToCartStatus)

	// GetCheckoutpgae
	GetCheckoutPageStatus, _ := footlocker.GetCheckoutPage(task, cid)
	footlocker.Log.Info("GetCheckoutPageStatus: ", GetCheckoutPageStatus)

	// SubmitUserInfo
	SubmitUserInfoStatus, _ := footlocker.SubmitUserInfo(task, cid)
	footlocker.Log.Info("SubmitUserInfoStatus: ", SubmitUserInfoStatus)

	// LocationLookup
	LocationLookupStatus, _ := footlocker.LocationLookup(task, cid)
	footlocker.Log.Info("LocationLookupStatus: ", LocationLookupStatus)

	// AddAddress
	AddAddressStatus, _ := footlocker.AddAddress(task, cid)
	footlocker.Log.Info("AddAddressStatus: ", AddAddressStatus)

	// VerifyAddress
	VerifyAddressStatus, _ := footlocker.VerifyAddress(task, csrfToken, cid)
	footlocker.Log.Info("VerifyAddressStatus: ", VerifyAddressStatus)

	// SubmitVerifiedAddress
	SubmitVerifiedAddressStatus, _ := footlocker.SubmitVerifiedAddress(task, cid)
	footlocker.Log.Info("SubmitVerifiedAddressStatus: ", SubmitVerifiedAddressStatus)

	// GetAdyenPublicKey
	GetAdyenPublicKeyStatus, publicKey, dQuery, _ := footlocker.GetAdyenPublicKey(task, cid)
	footlocker.Log.Info("GetAdyenPublicKeyStatus: ", GetAdyenPublicKeyStatus)
	footlocker.Log.Info("publicKey: ", publicKey)

	// GetAdyenEncryption
	GetAdyenEncryptionStatus, encryptionKey, _ := footlocker.GetAdyenEncryptionKey(task, cid, publicKey, dQuery)
	footlocker.Log.Info("GetAdyenEncryptionStatus: ", GetAdyenEncryptionStatus)
	footlocker.Log.Info("encryptionKey: ", encryptionKey)

	// Adyen Encryption
	AdyenEncryptStatus, encryptedCardNumber, encryptedExpiryMonth, encryptedExpiryYear, encryptedSecurityCode, _ := footlocker.AdyenEncrypt(task, encryptionKey)
	footlocker.Log.Info("Adyen Encryption Status: ", AdyenEncryptStatus)
	footlocker.Log.Debug("encryptedCardNumber: ", encryptedCardNumber)
	footlocker.Log.Debug("encryptedExpiryMonth: ", encryptedExpiryMonth)
	footlocker.Log.Debug("encryptedExpiryYear: ", encryptedExpiryYear)
	footlocker.Log.Debug("encryptedSecurityCode: ", encryptedSecurityCode)

	if AdyenEncryptStatus != 200 {
		return
	}

	// PlaceOrder
	PlaceOrderStatus, _ := footlocker.PlaceOrder(task, encryptedCardNumber, encryptedExpiryMonth, encryptedExpiryYear, encryptedSecurityCode, cid)
	footlocker.Log.Info("PlaceOrderStatus: ", PlaceOrderStatus)

}
