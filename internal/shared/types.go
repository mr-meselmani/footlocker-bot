package shared

import (
	http "github.com/bogdanfinn/fhttp"

	"net/url"
)

type Response struct {
	Error       error
	RotateProxy bool
	Success     bool
	Message     string
	Product     *Product
	Cookies     []*http.Cookie
	Extras      map[string]any
}

type Product struct {
	Name        string
	Size        string
	Url         string
	ImageUrl    string
	CheckoutUrl string
	Payment     string
	Store       string
}

type Task struct {
	ProfileName string
	ProductURL  string
	Size        string
	ProductID   string
	Quantity    int
	UseProxy    bool
	Mode        string
	Aco         bool
	Region      string
	Store       string
	Keywords    string
	Sku         string
	Payment     string
	Profile     Profile
	Id          int
}

type Profile struct {
	ProfileName string
	FirstName   string
	LastName    string
	Age         int
	BirthDay    int
	BirthMonth  int
	BirthYear   int
	Gender      string
	Email       string
	Phone       string
	Address     string
	Address2    string
	Zip         string
	City        string
	Country     string
	CountryISO  string
	CountryCode string
	State       string
	Cardname    string
	Cnb         string
	Month       string
	Year        string
	Cvv         string
	CardType    string
	Password    string
}

type Proxy struct {
	Url url.URL
}

type Settings struct {
	MonitorDelay    int
	RetryCount      int
	WebhookURL      string
	License         string
	ProxyHttps      bool
	CaptchaProvider string
	CaptchaKey      string
}