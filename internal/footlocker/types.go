package footlocker

import "time"

type TimestampResponse struct {
	Data    TimestampData `json:"data"`
	Success bool          `json:"success"`
	Errors  []interface{} `json:"errors"`
}

type TimestampData struct {
	CsrfToken string        `json:"csrfToken"`
	User      TimestampUser `json:"user"`
}

type TimestampUser struct {
	FirstName        string    `json:"firstName"`
	ServerUTC        time.Time `json:"serverUTC"`
	OptIn            bool      `json:"optIn"`
	MilitaryVerified bool      `json:"militaryVerified"`
	LoyaltyStatus    bool      `json:"loyaltyStatus"`
	SsoComplete      bool      `json:"ssoComplete"`
	VipUser          bool      `json:"vipUser"`
	Recognized       bool      `json:"recognized"`
	Vip              bool      `json:"vip"`
	Loyalty          bool      `json:"loyalty"`
	Authenticated    bool      `json:"authenticated"`
}

type AddToCartPayload struct {
	Size            string `json:"size"`
	Sku             string `json:"sku"`
	ProductQuantity int    `json:"productQuantity"`
	FulfillmentMode string `json:"fulfillmentMode"`
	ResponseFormat  string `json:"responseFormat"`
}

type SubmitUserInfoPayload struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PhoneCountry string `json:"phoneCountry"`
}

type AddAddressPayload struct {
	CountryIsocode     string `json:"countryIsocode"`
	City               string `json:"city"`
	PostalCode         string `json:"postalCode"`
	RegionIsocodeShort string `json:"regionIsocodeShort"`
	CheckoutType       string `json:"checkoutType"`
	IsShipping         bool   `json:"isShipping"`
}

type VerifyAddressPayload struct {
	Country    VerifyAddressCountry `json:"country"`
	Region     VerifyAddressRegion  `json:"region"`
	Line1      string               `json:"line1"`
	Line2      string               `json:"line2"`
	PostalCode string               `json:"postalCode"`
	Town       string               `json:"town"`
}

type VerifyAddressCountry struct {
	Isocode string `json:"isocode"`
}

type VerifyAddressRegion struct {
	IsocodeShort string `json:"isocodeShort"`
}

type VerifiedAddressPayload struct {
	CountryIsocode     string `json:"countryIsocode"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Line1              string `json:"line1"`
	Line2              string `json:"line2"`
	PostalCode         string `json:"postalCode"`
	City               string `json:"city"`
	RegionIsocodeShort string `json:"regionIsocodeShort"`
	IsBilling          bool   `json:"isBilling"`
	IsShipping         bool   `json:"isShipping"`
	RegionIsocode      string `json:"regionIsocode"`
	Phone              string `json:"phone"`
	AddressType        string `json:"addressType"`
	Residential        bool   `json:"residential"`
}

type PlaceOrderPayload struct {
	Payment struct {
		CcPaymentInfo struct {
			EncryptedCardNumber   string `json:"encryptedCardNumber"`
			EncryptedExpiryMonth  string `json:"encryptedExpiryMonth"`
			EncryptedExpiryYear   string `json:"encryptedExpiryYear"`
			EncryptedSecurityCode string `json:"encryptedSecurityCode"`
			SavePayment           bool   `json:"savePayment"`
		} `json:"ccPaymentInfo"`
		BrowserInfo struct {
			ScreenWidth    int    `json:"screenWidth"`
			ScreenHeight   int    `json:"screenHeight"`
			ColorDepth     int    `json:"colorDepth"`
			UserAgent      string `json:"userAgent"`
			TimeZoneOffset int    `json:"timeZoneOffset"`
			Language       string `json:"language"`
			JavaEnabled    bool   `json:"javaEnabled"`
		} `json:"browserInfo"`
	} `json:"payment"`
	IsNoChargeOrder bool   `json:"isNoChargeOrder"`
	CheckoutType    string `json:"checkoutType"`
	OptIn           bool   `json:"optIn"`
	DeviceID        string `json:"deviceId"`
}
