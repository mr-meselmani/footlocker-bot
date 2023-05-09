package internal

import (
	"bytes"
	"encoding/json"
	"footlocker-bot/internal/footlocker"
	"footlocker-bot/internal/logger"
	"footlocker-bot/internal/shared"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
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
		tls_client.WithProxyUrl("http://127.0.0.1:8888"),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	if err != nil {
		f.Log.Error("Get settings error")
		return err
	}

	f.Client = client

	return nil
}

func (f *Footlocker) GetHome() (int, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.footlocker.com/", nil)
	if err != nil {
		f.Log.Debug("Failed to issue GetHome request")
		return 0, err
	}

	req.Header = http.Header{
		"host":            {"www.footlocker.com"},
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-encoding": {"gzip, deflate, br"},
		"accept-language": {"en-US,en;q=0.9,ar;q=0.8"},
		// "if-none-match":               {`W/"84318-4KDNT3sIxHfWkNLb4Ru/VkDh15w"`},
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
			"host",
			"accept",
			"accept-encoding",
			"accept-language",
			// "if-none-match",
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
		f.Log.Debug("Failed to GetHome response")
		return 0, err
	}

	return res.StatusCode, nil
}

func (f *Footlocker) GetProduct(task shared.Task) (int, error) {
	req, err := http.NewRequest(http.MethodGet, task.ProductURL, nil)
	if err != nil {
		f.Log.Debug("Failed to issue GetProduct request")
		return 0, err
	}

	req.Header = http.Header{
		"host":                        {"www.footlocker.com"},
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
			"host",
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
		f.Log.Debug("GetProduct response error")
		return 0, err
	}

	return res.StatusCode, nil
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
		f.Log.Debug("Failed to issue AddToCart request")
		return 0, err
	}

	req.ContentLength = int64(len(json))

	req.Header = http.Header{
		"host":                        {"www.footlocker.com"},
		"sec-ch-ua":                   {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"x-fl-size":                   {"14.0"},
		"x-api-lang":                  {"en"},
		"sec-ch-ua-mobile":            {"?0"},
		"user-agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
		"sec-ch-ua-arch":              {`"x86"`},
		"sec-ch-device-memory":        {"8"},
		"content-type":                {"application/json"},
		"x-fl-sku":                    {task.Sku},
		"accept":                      {"application/json"},
		"sec-ch-ua-full-version-list": {`"Chromium";v="112.0.5615.138", "Google Chrome";v="112.0.5615.138", "Not:A-Brand";v="99.0.0.0"`},
		"sec-ch-ua-model":             {""},
		"x-fl-request-id":             {"87ae7830-ee61-11ed-958a-cf681a6f9726"},
		"sec-ch-ua-platform":          {`"Windows"`},
		"origin":                      {"https://www.footlocker.com"},
		"sec-fetch-site":              {"same-origin"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-dest":              {"empty"},
		"referer":                     {task.ProductURL},
		"accept-encoding":             {"gzip, deflate, br"},
		"accept-language":             {"en-US,en;q=0.9,ar;q=0.8"},
		http.HeaderOrderKey: {
			"host",
			"content-length",
			"sec-ch-ua",
			"x-fl-size",
			"x-api-lang",
			"sec-ch-ua-mobile",
			"user-agent",
			"sec-ch-ua-arch",
			"sec-ch-device-memory",
			"content-type",
			"x-fl-sku",
			"accept",
			"sec-ch-ua-full-version-list",
			"sec-ch-ua-model",
			"x-fl-request-id",
			"sec-ch-ua-platform",
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
		f.Log.Debug("AddToCart response error")
		return 0, err
	}

	return res.StatusCode, nil
}
