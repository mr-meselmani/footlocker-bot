package internal

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"footlocker-bot/internal/footlocker"
	"footlocker-bot/internal/logger"
	"footlocker-bot/internal/shared"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/CrimsonAIO/adyen"
	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/google/uuid"
)

type Footlocker struct {
	Client tls_client.HttpClient
	Log    logger.Logger
}

func NewFootlockerBot() Footlocker {
	return Footlocker{}
}

func (f *Footlocker) IsActive() bool {
	return true
}

func (f *Footlocker) GetFootlockerSettings(settings shared.Settings) error {
	jarOptions := []tls_client.CookieJarOption{}
	jar := tls_client.NewCookieJar(jarOptions...)

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(60),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		// tls_client.WithProxyUrl("http://127.0.0.1:8888"),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	if err != nil {
		f.Log.Error("Get settings error: ", err)
		return err
	}

	f.Log.EnableDebug()

	f.Client = client

	return nil
}

func (f *Footlocker) GetHome(task shared.Task) (int, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.footlocker.com/", nil)
	if err != nil {
		f.Log.Debug("Failed to issue GetHome request: ", err)
		return 0, err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"sec-ch-device-memory":        {"8"},
		"sec-ch-ua":                   {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-arch":              {`"x86"`},
		"sec-ch-ua-full-version-list": {`"Chromium";v="112.0.5615.138", "Google Chrome";v="112.0.5615.138", "Not:A-Brand";v="99.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-ch-ua-platform":          {`"Windows"`},
		"sec-fetch-dest":              {"document"},
		"sec-fetch-mode":              {"navigate"},
		"sec-fetch-site":              {"none"},
		"sec-fetch-user":              {"?1"},
		"upgrade-insecure-requests":   {"1"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"sec-ch-device-memory",
			"sec-ch-ua",
			"sec-ch-ua-arch",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"sec-fetch-user",
			"upgrade-insecure-requests",
			"user-agent",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("Failed to GetHome response: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) GetProduct(task shared.Task) (int, error) {
	req, err := http.NewRequest(http.MethodGet, task.ProductURL, nil)
	if err != nil {
		f.Log.Debug("Failed to issue GetProduct request: ", err)
		return 0, err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"cache-control":               {"max-age=0"},
		"sec-ch-device-memory":        {"8"},
		"sec-ch-ua":                   {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-arch":              {`"x86"`},
		"sec-ch-ua-platform":          {`"Windows"`},
		"sec-ch-ua-model":             {""},
		"sec-ch-ua-full-version-list": {`"Chromium";v="112.0.5615.138", "Google Chrome";v="112.0.5615.138", "Not:A-Brand";v="99.0.0.0"`},
		"upgrade-insecure-requests":   {"1"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"accept":                      {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-fetch-site":              {"same-origin"},
		"sec-fetch-mode":              {"navigate"},
		"sec-fetch-user":              {"?1"},
		"sec-fetch-dest":              {"document"},
		"referer":                     {"https://www.footlocker.com/"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		http.HeaderOrderKey: {
			"authority",
			"cache-control",
			"sec-ch-device-memory",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-arch",
			"sec-ch-ua-platform",
			"sec-ch-ua-model",
			"sec-ch-ua-full-version-list",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("GetProduct response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) TimeStamp(task shared.Task) (int, string, error) {
	// Get the current time
	currentTime := time.Now()

	// Convert it to the required timestamp format
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.footlocker.com/zgw/session?timestamp=%d", timestamp), nil)
	if err != nil {
		f.Log.Debug("Failed to issue TimeStamp request: ", err)
		return 0, "", err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"referer":                     {task.ProductURL},
		"sec-ch-device-memory":        {"8"},
		"sec-ch-ua":                   {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-arch":              {`"x86"`},
		"sec-ch-ua-full-version-list": {`"Chromium";v="112.0.5615.138", "Google Chrome";v="112.0.5615.138", "Not:A-Brand";v="99.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-ch-ua-platform":          {`"Windows"`},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"x-fl-request-id":             {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"referer",
			"sec-ch-device-memory",
			"sec-ch-ua",
			"sec-ch-ua-arch",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("TimeStamp response error: ", err)
		return 0, "", err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, "", nil
	}

	var Data footlocker.TimestampResponse
	err = json.Unmarshal(b, &Data)
	if err != nil {
		return 0, "", nil
	}

	csrfToken := Data.Data.CsrfToken

	return res.StatusCode, csrfToken, nil
}

func (f *Footlocker) AddToCart(task shared.Task) (int, error) {
	data := footlocker.AddToCartPayload{
		Size:            task.Size,
		Sku:             task.Sku,
		ProductQuantity: task.Quantity,
		FulfillmentMode: "SHIP",
		ResponseFormat:  "AllItems",
	}

	json, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/zgw/carts/co-cart-aggregation-service/site/fl/cart/cartItems/addCartItem", bytes.NewBuffer(json))
	if err != nil {
		f.Log.Debug("Failed to issue AddToCart request: ", err)
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"authority":       {"www.footlocker.com"},
		"accept":          {"application/json"},
		"accept-encoding": {"gzip, deflate, br"},
		"accept-language": {"en-US,en;q=0.9,ar;q=0.8"},
		"content-type":    {"application/json"},
		"referer":         {"https://www.footlocker.com/"},
		"cache-control":   {"no-cache"},
		"sec-ch-ua":       {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		// "sec-ch-ua-full-version-list": {`"Chromium";v="112.0.5615.138", "Google Chrome";v="112.0.5615.138", "Not:A-Brand";v="99.0.0.0"`},
		"sec-ch-ua-mobile": {"?0"},
		"sec-ch-ua-model":  {""},
		"sec-fetch-dest":   {"empty"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-site":   {"same-origin"},
		"user-agent":       {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"x-api-lang":       {"en-US"},
		"x-api-country":    {"US"},
		"x-fl-request-id":  {uuid.New().String()},
		"x-fl-size":        {task.Size},
		"x-fl-sku":         {task.Sku},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"referer",
			"cache-control",
			"sec-ch-ua",
			// "sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-api-country",
			"x-fl-request-id",
			"x-fl-size",
			"x-fl-sku",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("AddToCart response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) GetCheckoutPage(task shared.Task) (int, error) {
	// Get the current time
	currentTime := time.Now()

	// Convert it to the required timestamp format
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.footlocker.com/api/pages/en/checkout.page.mainpage.json?timestamp=%d", timestamp), nil)
	if err != nil {
		f.Log.Debug("Failed to issue GetCheckoutPage request: ", err)
		return 0, err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"referer":                     {"https://www.footlocker.com/checkout"},
		"sec-ch-ua":                   {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-full-version-list": {`"Google Chrome";v="113.0.5672.93", "Chromium";v="113.0.5672.93", "Not-A.Brand";v="24.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-ch-ua-platform":          {"Windows"},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-fl-request-id":             {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"referer",
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("GetCheckoutPage response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) SubmitUserInfo(task shared.Task) (int, error) {
	data := footlocker.SubmitUserInfoPayload{
		FirstName:    task.Profile.FirstName,
		LastName:     task.Profile.LastName,
		Email:        task.Profile.Email,
		Phone:        task.Profile.Phone,
		PhoneCountry: task.Profile.CountryISO,
	}

	json, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/zgw/carts/co-cart-aggregation-service/site/fl/cart/userInfo", bytes.NewBuffer(json))
	if err != nil {
		f.Log.Debug("Failed to issue SubmitUserInfo request: ", err)
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"content-type":                {"application/json"},
		"referer":                     {"https://www.footlocker.com/checkout"},
		"cache-control":               {"no-cache"},
		"sec-ch-ua":                   {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-full-version-list": {`"Google Chrome";v="113.0.5672.93", "Chromium";v="113.0.5672.93", "Not-A.Brand";v="24.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":                  {"en"},
		"x-fl-request-id":             {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"referer",
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("SubmitUserInfo response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) AddAddress(task shared.Task) (int, error) {
	data := footlocker.AddAddressPayload{
		CountryIsocode:     task.Profile.CountryISO,
		City:               task.Profile.City,
		PostalCode:         task.Profile.Zip,
		RegionIsocodeShort: task.Region,
		CheckoutType:       "EXPRESS",
		IsShipping:         true,
	}

	json, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodGet, "https://www.footlocker.com/zgw/carts/co-cart-aggregation-service/site/fl/cart/address", bytes.NewBuffer(json))
	if err != nil {
		f.Log.Debug("Failed to issue AddAddress request: ", err)
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"content-type":                {"application/json"},
		"referer":                     {"https://www.footlocker.com/checkout"},
		"cache-control":               {"no-cache"},
		"sec-ch-ua":                   {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-full-version-list": {`"Google Chrome";v="113.0.5672.93", "Chromium";v="113.0.5672.93", "Not-A.Brand";v="24.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":                  {"en"},
		"x-fl-request-id":             {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"referer",
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("AddAddress response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) VerifyAddress(task shared.Task, csrfToken string) (int, error) {
	data := footlocker.VerifyAddressPayload{
		Country: footlocker.VerifyAddressCountry{
			Isocode: task.Profile.CountryISO,
		},
		Region: footlocker.VerifyAddressRegion{
			IsocodeShort: task.Region,
		},
		Line1:      task.Profile.Address,
		Line2:      task.Profile.Address2,
		PostalCode: task.Profile.Zip,
		Town:       task.Profile.City,
	}

	json, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/api/users/addresses/verification", bytes.NewBuffer(json))
	if err != nil {
		f.Log.Debug("Failed to issue VerifyAddress request: ", err)
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"content-type":                {"application/json"},
		"referer":                     {"https://www.footlocker.com/checkout"},
		"cache-control":               {"no-cache"},
		"sec-ch-ua":                   {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-full-version-list": {`"Google Chrome";v="113.0.5672.93", "Chromium";v="113.0.5672.93", "Not-A.Brand";v="24.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":                  {"en"},
		"x-csrf-token":                {csrfToken},
		"x-fl-request-id":             {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"referer",
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-csrf-token",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("VerifyAddress response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) SubmitVerifiedAddress(task shared.Task) (int, error) {
	data := footlocker.VerifiedAddressPayload{
		CountryIsocode:     task.Profile.CountryISO,
		FirstName:          task.Profile.FirstName,
		LastName:           task.Profile.LastName,
		Line1:              task.Profile.Address,
		Line2:              task.Profile.Address2,
		PostalCode:         task.Profile.Zip,
		City:               task.Profile.City,
		RegionIsocodeShort: task.Profile.CountryISO,
		IsBilling:          true,
		IsShipping:         true,
		RegionIsocode:      task.Region,
		Phone:              task.Profile.Phone,
		AddressType:        " ",
		Residential:        false,
	}

	json, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/zgw/carts/co-cart-aggregation-service/site/fl/cart/address", bytes.NewBuffer(json))
	if err != nil {
		f.Log.Debug("Failed to issue SubmitVerifiedAddress request: ", err)
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"content-type":                {"application/json"},
		"referer":                     {"https://www.footlocker.com/checkout"},
		"cache-control":               {"no-cache"},
		"sec-ch-ua":                   {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-full-version-list": {`"Google Chrome";v="113.0.5672.93", "Chromium";v="113.0.5672.93", "Not-A.Brand";v="24.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":                  {"en"},
		"x-fl-request-id":             {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"referer",
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("SubmitVerifiedAddress response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) GetAdyen(task shared.Task) (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.footlocker.com/apigate/payment/origin-key", nil)
	if err != nil {
		f.Log.Debug("Failed to issue GetAdyen request: ", err)
		return 0, "", err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"referer":                     {"https://www.footlocker.com/checkout"},
		"sec-ch-ua":                   {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-full-version-list": {`"Google Chrome";v="113.0.5672.93", "Chromium";v="113.0.5672.93", "Not-A.Brand";v="24.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":                  {"en"},
		"x-fl-request-id":             {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"referer",
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("GetAdyen response error: ", err)
		return 0, "", err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, "", err
	}

	htmlString, _ := doc.Html()

	var publicKey string

	// public key
	publicKeyRegex := regexp.MustCompile(`"adyenResponse":{"https://www.footlocker.com":"([\s\S]*?)"`)
	publicKeyMatch := publicKeyRegex.FindStringSubmatch(htmlString)
	if len(publicKeyMatch) > 0 {
		fmt.Println("publicKey found")
		publicKey = publicKeyMatch[1]
	}

	f.Log.Debug("Publickey: ", publicKey)

	return res.StatusCode, publicKey, nil
}

func (f *Footlocker) PlaceOrder(task shared.Task) (int, error) {
	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/zgw/carts/co-cart-aggregation-service/site/fl/cart/placeOrder", nil)
	if err != nil {
		f.Log.Debug("Failed to issue PlaceOrder request: ", err)
		return 0, err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"content-type":                {"application/json"},
		"referer":                     {"https://www.footlocker.com/checkout"},
		"cache-control":               {"no-cache"},
		"sec-ch-ua":                   {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-full-version-list": {`"Google Chrome";v="113.0.5672.93", "Chromium";v="113.0.5672.93", "Not-A.Brand";v="24.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":                  {"en"},
		"x-fl-request-id":             {uuid.New().String()},
		"x-flgw-channel-info":         {"WEB"},
		"x-mobile-device":             {"false"},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"referer",
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-fl-request-id",
			"x-flgw-channel-info",
			"x-mobile-device",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Debug("PlaceOrder response error: ", err)
		return 0, err
	}

	return res.StatusCode, nil
}

// Helpers
func (f *Footlocker) AdyenEncrypt(task shared.Task, publicKey string) (int, string, string, string, string, error) {

	var plaintextKey = publicKey
	key := strings.Split(plaintextKey, "|")[1]
	fmt.Println("key: ", key)
	b, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}

	// create new encrypter
	enc, err := adyen.NewEncrypter("0_1_25", adyen.PubKeyFromBytes(b))
	if err != nil {
		panic(err)
	}

	encryptedCardNumber, _ := enc.EncryptField("encryptedCardNumber", task.Profile.Cnb)
	encryptedExpiryMonth, _ := enc.EncryptField("encryptedExpiryMonth", task.Profile.Month)
	encryptedExpiryYear, _ := enc.EncryptField("encryptedExpiryYear", task.Profile.Year)
	encryptedSecurityCode, _ := enc.EncryptField("encryptedSecurityCode", task.Profile.Cvv)

	return 200, encryptedCardNumber, encryptedExpiryMonth, encryptedExpiryYear, encryptedSecurityCode, nil
}
