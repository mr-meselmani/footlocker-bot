package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"footlocker-bot/footlocker"
	"footlocker-bot/logger"
	"footlocker-bot/shared"
	"io"
	"net/url"
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
	DD     Datadome
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

	f.DD = NewDatadome()

	f.Client = client

	newProxy := f.Flip()
	f.Client.SetProxy(newProxy)

	return nil
}

// Requests
func (f *Footlocker) GetHome(task shared.Task) (int, string, error) {

	req, err := http.NewRequest(http.MethodGet, "https://www.footlocker.com/", nil)
	if err != nil {
		f.Log.Error("Failed to issue GetHome request: ", err)
		return 0, "", err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"sec-ch-ua":                   {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-full-version-list": {`"Chromium";v="112.0.5615.138", "Google Chrome";v="112.0.5615.138", "Not:A-Brand";v="99.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
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
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
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
		f.Log.Error("Failed to GetHome response: ", err)
		return 0, "", err
	}

	if res.StatusCode != 200 {
		newProxy := f.Flip()
		f.Client.SetProxy(newProxy)
		f.Log.Warning("GetHomeStatus: ", res.StatusCode)
		return f.GetHome(task)
	}

	defer res.Body.Close()

	cid := strings.Split(strings.Split(res.Header.Get("set-cookie"), ";")[0], "=")[1]

	return res.StatusCode, cid, nil
}

func (f *Footlocker) GetProduct(task shared.Task, cid string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, task.ProductURL, nil)
	if err != nil {
		f.Log.Error("Failed to issue GetProduct request: ", err)
		return 0, err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"cache-control":               {"max-age=0"},
		"sec-ch-ua":                   {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-mobile":            {"?0"},
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
			"sec-ch-ua",
			"sec-ch-ua-mobile",
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
		f.Log.Error("GetProduct response error: ", err)
		return 0, err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		return f.GetProduct(task, cid)
	}

	defer res.Body.Close()

	return res.StatusCode, nil
}

func (f *Footlocker) TimeStamp(task shared.Task, cid string) (int, string, error) {
	// Get the current time
	currentTime := time.Now()

	// Convert it to the required timestamp format
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.footlocker.com/zgw/session?timestamp=%d", timestamp), nil)
	if err != nil {
		f.Log.Error("Failed to issue TimeStamp request: ", err)
		return 0, "", err
	}

	req.Header = http.Header{
		"authority":                   {"www.footlocker.com"},
		"accept":                      {"application/json"},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		"referer":                     {task.ProductURL},
		"sec-ch-ua":                   {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-full-version-list": {`"Chromium";v="112.0.5615.138", "Google Chrome";v="112.0.5615.138", "Not:A-Brand";v="99.0.0.0"`},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {""},
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
			"sec-ch-ua",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-mobile",
			"sec-ch-ua-model",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Error("TimeStamp response error: ", err)
		return 0, "", err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {

			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("TimeStamp: ", res.StatusCode)

		return f.TimeStamp(task, cid)
	}

	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	var Data footlocker.TimestampResponse
	err = json.Unmarshal(b, &Data)
	if err != nil {
		return 0, "", nil
	}

	csrfToken := Data.Data.CsrfToken

	return res.StatusCode, csrfToken, nil
}

func (f *Footlocker) AddToCart(task shared.Task, cid string) (int, error) {
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
		f.Log.Error("Failed to issue AddToCart request: ", err)
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"authority":          {"www.footlocker.com"},
		"accept":             {"application/json"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-US,en;q=0.9"},
		"content-type":       {"application/json"},
		"origin":             {"https://www.footlocker.com"},
		"referer":            {task.ProductURL},
		"sec-ch-ua":          {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":         {"en"},
		"x-fl-request-id":    {uuid.New().String()},
		"x-fl-size":          {"10.0"},
		"x-fl-sku":           {task.Sku},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"origin",
			"referer",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-api-lang",
			"x-fl-request-id",
			"x-fl-size",
			"x-fl-sku",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Error("AddToCart response error: ", err)
		return 0, err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("AddToCartStatus: ", res.StatusCode)

		return f.AddToCart(task, cid)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	f.Log.Debug("ATC response: ", string(body))

	return res.StatusCode, nil
}

func (f *Footlocker) GetCheckoutPage(task shared.Task, cid string) (int, error) {
	// Get the current time
	currentTime := time.Now()

	// Convert it to the required timestamp format
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.footlocker.com/api/pages/en/checkout.page.mainpage.json?timestamp=%d", timestamp), nil)
	if err != nil {
		f.Log.Error("Failed to issue GetCheckoutPage request: ", err)
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
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
			"x-fl-request-id",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Error("GetCheckoutPage response error: ", err)
		return 0, err
	}

	defer res.Body.Close()

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("GetCheckoutPage: ", res.StatusCode)

		return f.GetCheckoutPage(task, cid)
	}

	return res.StatusCode, nil
}

func (f *Footlocker) SubmitUserInfo(task shared.Task, cid string) (int, error) {
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
		f.Log.Error("Failed to issue SubmitUserInfo request: ", err)
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
		f.Log.Error("SubmitUserInfo response error: ", err)
		return 0, err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("AddToCartStatus: ", res.StatusCode)

		return f.SubmitUserInfo(task, cid)
	}

	defer res.Body.Close()

	return res.StatusCode, nil
}

func (f *Footlocker) LocationLookup(task shared.Task, cid string) (int, error) {
	data := footlocker.LocationLokkupPayload{
		ZipCode: task.Region + " " + task.Profile.Zip,
	}

	j, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/api/satori/location-lookup/", bytes.NewBuffer(j))
	if err != nil {
		f.Log.Error("Failed to issue LocationLokkup request: ", err)
		return 0, err
	}

	req.Header = http.Header{
		"authority":        {"www.footlocker.com"},
		"sec-ch-ua":        {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"x-api-lang":       {"en"},
		"sec-ch-ua-mobile": {"?0"},
		"user-agent":       {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"content-type":     {"application/json"},
		"accept":           {"application/json"},
		"x-fl-request-id":  {uuid.New().String()},
		"origin":           {"https://www.footlocker.com"},
		"sec-fetch-site":   {"same-origin"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-dest":   {"empty"},
		"referer":          {"https://www.footlocker.com/checkout"},
		"accept-encoding":  {"gzip, deflate, br"},
		"accept-language":  {"en-US,en;q=0.9"},
		http.HeaderOrderKey: {
			"authority",
			"content-length",
			"sec-ch-ua",
			"x-api-lang",
			"sec-ch-ua-mobile",
			"user-agent",
			"content-type",
			"accept",
			"x-fl-request-id",
			"origin",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Error("Failed to issue LocationLookup request: ", err)
		return 0, err
	}

	// body, _ := io.ReadAll(res.Body)

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("LocationLookup: ", res.StatusCode)

		return f.LocationLookup(task, cid)
	}

	defer res.Body.Close()

	return res.StatusCode, nil
}

func (f *Footlocker) AddAddress(task shared.Task, cid string) (int, error) {
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
		f.Log.Error("Failed to issue AddAddress request: ", err)
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"authority":        {"www.footlocker.com"},
		"accept":           {"application/json"},
		"accept-encoding":  {"gzip, deflate, br"},
		"accept-language":  {"en-US,en;q=0.9"},
		"content-type":     {"application/json"},
		"origin":           {"https://www.footlocker.com"},
		"referer":          {"https://www.footlocker.com/checkout"},
		"sec-ch-ua":        {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-mobile": {"?0"},
		"sec-fetch-dest":   {"empty"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-site":   {"same-origin"},
		"user-agent":       {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":       {"en"},
		"x-fl-request-id":  {uuid.New().String()},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"origin",
			"referer",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
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
		f.Log.Error("AddAddress response error: ", err)
		return 0, err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("AddAddress: ", res.StatusCode)

		return f.AddAddress(task, cid)
	}

	defer res.Body.Close()

	return res.StatusCode, nil
}

func (f *Footlocker) VerifyAddress(task shared.Task, csrfToken string, cid string) (int, error) {
	data := footlocker.VerifyAddressPayload{
		Country: footlocker.VerifyAddressCountry{
			Isocode: task.Profile.CountryISO,
		},
		Region: footlocker.VerifyAddressRegion{
			IsocodeShort: task.Region,
		},
		Line1:      task.Profile.Address,
		Line2:      "",
		PostalCode: task.Profile.Zip,
		Town:       task.Profile.City,
	}

	json, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/api/users/addresses/verification", bytes.NewBuffer(json))
	if err != nil {
		f.Log.Error("Failed to issue VerifyAddress request: ", err)
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
		f.Log.Error("VerifyAddress response error: ", err)
		return 0, err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("VerifyAddress: ", res.StatusCode)

		return f.VerifyAddress(task, csrfToken, cid)
	}

	defer res.Body.Close()

	return res.StatusCode, nil
}

func (f *Footlocker) SubmitVerifiedAddress(task shared.Task, cid string) (int, error) {
	data := footlocker.VerifiedAddressPayload{
		CountryIsocode:     task.Profile.CountryISO,
		FirstName:          task.Profile.FirstName,
		LastName:           task.Profile.LastName,
		Line1:              task.Profile.Address,
		Line2:              "",
		PostalCode:         task.Profile.Zip,
		City:               task.Profile.City,
		RegionIsocodeShort: task.Region,
		IsBilling:          true,
		IsShipping:         true,
		RegionIsocode:      task.RegionIsocode,
		Phone:              task.Profile.Phone,
		AddressType:        " ",
		Residential:        false,
	}

	json, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/zgw/carts/co-cart-aggregation-service/site/fl/cart/address", bytes.NewBuffer(json))
	if err != nil {
		f.Log.Error("Failed to issue SubmitVerifiedAddress request: ", err)
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
		f.Log.Error("SubmitVerifiedAddress response error: ", err)
		return 0, err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("SubmitVerifiedAddress: ", res.StatusCode)

		return f.SubmitVerifiedAddress(task, cid)
	}

	defer res.Body.Close()

	return res.StatusCode, nil
}

func (f *Footlocker) GetAdyenPublicKey(task shared.Task, cid string) (int, string, string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.footlocker.com/apigate/payment/origin-key", nil)
	if err != nil {
		f.Log.Error("Failed to issue GetAdyen request: ", err)
		return 0, "", "", err
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
		f.Log.Error("GetAdyen response error: ", err)
		return 0, "", "", err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("GetAdyen: ", res.StatusCode)

		return f.GetAdyenPublicKey(task, cid)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var resData footlocker.GetAdyenRes
	err = json.Unmarshal(body, &resData)

	if err != nil {
		f.Log.Error("GetAdyen unmarshal data error: ", err)
		return 0, "", "", err
	}

	dQuery := strings.Split(resData.OKey, ".")[3]

	return res.StatusCode, resData.OKey, dQuery, nil
}

func (f *Footlocker) GetAdyenEncryptionKey(task shared.Task, cid string, publicKey string, dQuery string) (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://checkoutshopper-live.adyen.com/checkoutshopper/securedfields/%s/3.3.1/securedFields.html?type=card&d=%s%3D", publicKey, dQuery), nil)
	if err != nil {
		f.Log.Error("Failed to issue GetAdyenEncryptionKey request: ", err)
		return 0, "", err
	}

	req.Header = http.Header{
		"Host":                      {"checkoutshopper-live.adyen.com"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Encoding":           {"gzip, deflate, br"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Connection":                {"keep-alive"},
		"Referer":                   {"https://www.footlocker.com/"},
		"Sec-Fetch-Dest":            {"iframe"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"cross-site"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"sec-ch-ua":                 {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		http.HeaderOrderKey: {
			"Host",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Connection",
			"Referer",
			"Sec-Fetch-Dest",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Site",
			"Sec-Fetch-User",
			"Upgrade-Insecure-Requests",
			"User-Agent",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
		},
	}

	res, err := f.Client.Do(req)
	if err != nil {
		f.Log.Error("GetAdyenEncryptionKey response error: ", err)
		return 0, "", err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("GetAdyenEncryptionKey: ", res.StatusCode)

		return f.GetAdyenEncryptionKey(task, cid, publicKey, dQuery)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		f.Log.Error("Failed to read GetAdyenEncryptionKey doc: ", err)
		return 0, "", err

	}

	htmlString, _ := doc.Html()

	var encryptionKey string

	// public key
	encryptionKeyRegex := regexp.MustCompile(`var key = "([\s\S]*?)";`)
	encryptionKeyMatch := encryptionKeyRegex.FindStringSubmatch(htmlString)
	if len(encryptionKeyMatch) > 0 {
		f.Log.Debug("encryptionKey found", "")
		encryptionKey = encryptionKeyMatch[1]
	}

	return res.StatusCode, encryptionKey, nil
}

func (f *Footlocker) PlaceOrder(task shared.Task, encryptedCardNumber, encryptedExpiryMonth, encryptedExpiryYear, encryptedSecurityCode, cid string) (int, error) {
	payload := footlocker.PlaceOrderPayload{
		Payment: footlocker.PlaceOrderPayment{
			CcPaymentInfo: footlocker.PlaceOrderCcPaymentInfo{
				EncryptedCardNumber:   encryptedCardNumber,
				EncryptedExpiryMonth:  encryptedExpiryMonth,
				EncryptedExpiryYear:   encryptedExpiryYear,
				EncryptedSecurityCode: encryptedSecurityCode,
				SavePayment:           false,
			},
			BrowserInfo: footlocker.PlaceOrderBrowserInfo{
				ScreenWidth:    1536,
				ScreenHeight:   864,
				ColorDepth:     24,
				UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
				TimeZoneOffset: -180,
				Language:       "en-US",
				JavaEnabled:    false,
			},
		},
		IsNoChargeOrder: false,
		CheckoutType:    "NORMAL",
		OptIn:           false,
		DeviceID:        "",
	}

	j, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "https://www.footlocker.com/zgw/carts/co-cart-aggregation-service/site/fl/cart/placeOrder", bytes.NewBuffer(j))
	if err != nil {
		f.Log.Error("Failed to issue PlaceOrder request: ", err)
		return 0, err
	}

	req.Header = http.Header{
		"authority":           {"www.footlocker.com"},
		"accept":              {"application/json"},
		"accept-encoding":     {"gzip, deflate, br"},
		"accept-language":     {"en-US,en;q=0.9"},
		"content-type":        {"application/json"},
		"origin":              {"https://www.footlocker.com"},
		"referer":             {"https://www.footlocker.com/checkout"},
		"sec-ch-ua":           {`"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`},
		"sec-ch-ua-mobile":    {"?0"},
		"sec-fetch-dest":      {"empty"},
		"sec-fetch-mode":      {"cors"},
		"sec-fetch-site":      {"same-origin"},
		"user-agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":          {"en"},
		"x-fl-request-id":     {uuid.New().String()},
		"x-flgw-channel-info": {"WEB"},
		"x-mobile-device":     {"false"},
		http.HeaderOrderKey: {
			"authority",
			"accept",
			"accept-encoding",
			"accept-language",
			"content-length",
			"content-type",
			"origin",
			"referer",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
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
		f.Log.Error("PlaceOrder response error: ", err)
		return 0, err
	}

	domain := req.URL

	if res.StatusCode != 200 {

		ddCookie := f.GenerateCookies(cid, domain)

		// Modify a cookie value
		cookies := f.Client.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = ddCookie
				break
			}
		}

		// Print the updated cookies
		for _, cookie := range cookies {
			fmt.Printf("Name: %s, Value: %s\n", cookie.Name, cookie.Value)
		}

		f.Client.SetCookies(req.URL, cookies)

		f.Log.Warning("PlaceOrder: ", res.StatusCode)

		return f.PlaceOrder(task, encryptedCardNumber, encryptedExpiryMonth, encryptedExpiryYear, encryptedSecurityCode, cid)
	}

	defer res.Body.Close()

	return res.StatusCode, nil
}

// Helpers
func (f *Footlocker) AdyenEncrypt(task shared.Task, encryptionKey string) (int, string, string, string, string, error) {

	var plaintextKey = encryptionKey
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

func (f *Footlocker) Rotate() string {
	newProxy := ProxyRotator(true)

	if newProxy != " " {
		return newProxy
	}

	f.Log.Error("No proxies found: ", "check proxies.txt")

	return ""
}

func (f *Footlocker) Flip() string {
	p := f.Rotate()
	err := f.Client.SetProxy(p)

	if err != nil {
		f.Log.Error("Proxy failed: ", err)

	}

	return f.Client.GetProxy()
}

func (f *Footlocker) GenerateCookies(cid string, domain *url.URL) string {
	newProxy := f.Flip()

	f.Log.Debug("ddProxy: ", newProxy)

	// Generate ddCookie
	status, ddCookie, _ := f.DD.GenCh(cid, domain.String(), newProxy)
	if status != 200 && ddCookie == "" {
		f.Log.Warning("ddCookie status: ", status)
		f.Log.Error("ddCookie not generated: ", ddCookie)

		time.Sleep(2 * time.Second)

		return f.GenerateCookies(cid, domain)
	}

	f.Log.Info("ddCookie found: ", ddCookie)
	f.Client.SetProxy(newProxy)
	return ddCookie

}
