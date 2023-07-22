package main

import (
	"encoding/json"
	"fmt"
	"footlocker-bot/datadome"
	"footlocker-bot/logger"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Datadome struct {
	Client http.Client
	Log    logger.Logger
}

func NewDatadome() Datadome {
	return Datadome{}
}

func (d Datadome) GenCh(cid string, domain string, proxy string) (int, string, error) {
	filename := "./internal/datadome_payload.json"
	data, err := d.readJSONFile(filename)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return 0, "", err
	}

	jsetNum := time.Now().Unix()

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	ttst := d.round(random.Float64()*(40-8)+8, 15)
	tagGpu := d.round(ttst-(0.20*ttst), 14)
	tagpu := math.Round((tagGpu-(0.20*tagGpu))*1e14) / 1e14

	data.Tagpu = tagpu
	data.Ttst = ttst
	data.Jset = int(jsetNum)

	j, _ := json.Marshal(data)

	urlSearchParams := map[string]string{
		"jsData":        string(j),
		"eventCounters": "[]",
		"jsType":        "ch",
		"cid":           cid,
		"ddk":           "A55FBF4311ED6F1BF9911EB71931D5",
		"Referer":       domain,
		"request":       "%2",
		"responsePage":  "origin",
		"ddv":           "4.8.1",
	}

	// Define the desired order of keys
	urlSearchParamskeyOrder := []string{
		"jsData",
		"eventCounters",
		"jsType",
		"cid",
		"ddk",
		"Referer",
		"request",
		"responsePage",
		"ddv",
	}

	// Create a URL values instance
	values := url.Values{}

	// Add the key-value pairs in the desired order
	for _, key := range urlSearchParamskeyOrder {
		if value, ok := urlSearchParams[key]; ok {
			values.Add(key, value)
		}
	}

	// Construct the query string
	queryString := values.Encode()

	// Define the Charles proxy URL
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}

	d.Log.Debug("current dd proxy: ", proxyURL)

	// Create an http.Transport with the proxy configuration
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// Create an http.Client with the transport
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest(http.MethodPost, "https://api-js.datadome.co/js/", strings.NewReader(queryString))
	if err != nil {
		d.Log.Debug("Failed to issue GenCh request: ", err)
		return 0, "", err
	}

	req.Header = http.Header{
		"accept":           {"*/*"},
		"accept-language":  {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":     {"application/x-www-form-urlencoded"},
		"sec-ch-ua":        {`"Chromium";v="112", "Google Chrome";v="112", "Not:A-Brand";v="99"`},
		"sec-ch-ua-mobile": {"?0"},
		"sec-fetch-dest":   {"empty"},
		"sec-fetch-mode":   {"cors"},
		"sec-fetch-site":   {"cross-site"},
		"Referer":          {"https://www.footlocker.com/"},
		"Referrer-Policy":  {"strict-origin-when-cross-origin"},
	}

	res, err := client.Do(req)
	if err != nil {
		d.Log.Debug("GenCh response error: ", err)
		return 0, "", err
	}

	if res.StatusCode != 200 {
		d.Log.Debug("StatusCode: ", res.StatusCode)
		return 0, "", nil
	}

	defer res.Body.Close()

	var ddResponse datadome.DDresponse
	resBody, _ := io.ReadAll(res.Body)

	json.Unmarshal(resBody, &ddResponse)

	ddCookie := strings.Split(strings.Split(ddResponse.Cookie, ";")[0], "=")[1]

	return res.StatusCode, ddCookie, nil
}

// Helpers
func (d Datadome) readJSONFile(filename string) (*datadome.GenChPayload, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	data := &datadome.GenChPayload{}
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d Datadome) round(f float64, n int) float64 {
	shift := math.Pow(10, float64(n))
	return math.Round(f*shift) / shift
}
