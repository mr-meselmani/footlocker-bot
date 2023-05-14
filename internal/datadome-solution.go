package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"footlocker-bot/internal/datadome"
	"footlocker-bot/internal/logger"
	"footlocker-bot/internal/shared"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

type Datadome struct {
	Client tls_client.HttpClient
	Log    logger.Logger
}

type queryParam struct {
	Key   string
	Value string
}

func NewDatadome() Datadome {
	return Datadome{}
}

func (d Datadome) GetDatadomeSettings(settings shared.Settings) error {
	jarOptions := []tls_client.CookieJarOption{}
	jar := tls_client.NewCookieJar(jarOptions...)

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(60),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		tls_client.WithProxyUrl("http://127.0.0.1:8888"),
		// tls_client.WithCharlesProxy("", ""),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	if err != nil {
		d.Log.Error("Get settings error: ", err)
		return err
	}

	d.Log.EnableDebug()

	d.Client = client

	return nil
}

func (d Datadome) GenCh(cid string) (int, error) {
	// data := datadome.GenChPayload{
	// 	Ttst:     rand.Float64()*(40.0-8.0) + 8.0,
	// 	Ifov:     false,
	// 	Tagpu:    rand.Float64()*(40.0-8.0) + 8.0,
	// 	Glvd:     "Google Inc. (Intel Inc.)",
	// 	Glrd:     "ANGLE (Intel Inc., Intel(R) Iris(TM) Plus Graphics 645, OpenGL 4.1)",
	// 	Hc:       8,
	// 	BrOh:     1080,
	// 	BrOw:     1920,
	// 	Ua:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
	// 	Wbd:      false,
	// 	Wdif:     false,
	// 	Wdifrm:   false,
	// 	Npmtm:    false,
	// 	BrH:      969,
	// 	BrW:      1920,
	// 	Nddc:     1,
	// 	RsH:      1080,
	// 	RsW:      1920,
	// 	RsCd:     24,
	// 	Phe:      false,
	// 	Nm:       false,
	// 	Jsf:      false,
	// 	Lg:       "it-IT",
	// 	Pr:       1,
	// 	ArsH:     1080,
	// 	ArsW:     1920,
	// 	Tz:       -120,
	// 	StrSs:    true,
	// 	StrLs:    true,
	// 	StrIdb:   true,
	// 	StrOdb:   true,
	// 	Plgod:    false,
	// 	Plg:      5,
	// 	Plgne:    true,
	// 	Plgre:    true,
	// 	Plgof:    false,
	// 	Plggt:    false,
	// 	Pltod:    false,
	// 	Hcovdr:   false,
	// 	Hcovdr2:  false,
	// 	Plovdr:   false,
	// 	Plovdr2:  false,
	// 	Ftsovdr:  false,
	// 	Ftsovdr2: false,
	// 	Lb:       false,
	// 	Eva:      33,
	// 	Lo:       false,
	// 	TsMtp:    0,
	// 	TsTec:    false,
	// 	TsTsa:    false,
	// 	Vnd:      "Google Inc.",
	// 	Bid:      "NA",
	// 	Mmt:      "application/pdf,text/pdf",
	// 	Plu:      "PDF Viewer,Chrome PDF Viewer,Chromium PDF Viewer,Microsoft Edge PDF Viewer,WebKit built-in PDF",
	// 	Hdn:      true,
	// 	Awe:      false,
	// 	Geb:      false,
	// 	Dat:      false,
	// 	Med:      "defined",
	// 	Aco:      "probably",
	// 	Acots:    false,
	// 	Acmp:     "probably",
	// 	Acmpts:   true,
	// 	Acw:      "probably",
	// 	Acwts:    false,
	// 	Acma:     "maybe",
	// 	Acmats:   false,
	// 	Acaa:     "probably",
	// 	Acaats:   true,
	// 	Ac3:      "",
	// 	Ac3Ts:    true,
	// 	Acf:      "probably",
	// 	Acfts:    false,
	// 	Acmp4:    "maybe",
	// 	Acmp4Ts:  false,
	// 	Acmp3:    "probably",
	// 	Acmp3Ts:  false,
	// 	Acwm:     "maybe",
	// 	Acwmts:   false,
	// 	Ocpt:     false,
	// 	Vco:      "probably",
	// 	Vcots:    false,
	// 	Vch:      "probably",
	// 	Vchts:    true,
	// 	Vcw:      "probably",
	// 	Vcwts:    true,
	// 	Vc3:      "maybe",
	// 	Vc3Ts:    false,
	// 	Vcmp:     "",
	// 	Vcmpts:   false,
	// 	Vcq:      "",
	// 	Vcqts:    false,
	// 	Vc1:      "probably",
	// 	Vc1Ts:    true,
	// 	Dvm:      8,
	// 	Sqt:      false,
	// 	So:       "landscape-primary",
	// 	Wdw:      true,
	// 	Cokys:    "bG9hZFRpbWVzY3NpYXBwL=",
	// 	Ecpc:     false,
	// 	Lgs:      true,
	// 	Lgsod:    false,
	// 	Psn:      true,
	// 	Edp:      true,
	// 	Addt:     true,
	// 	Wsdc:     true,
	// 	Ccsr:     true,
	// 	Nuad:     true,
	// 	Bcda:     true,
	// 	Idn:      true,
	// 	Capi:     false,
	// 	Svde:     false,
	// 	Vpbq:     true,
	// 	Ucdv:     false,
	// 	Spwn:     false,
	// 	Emt:      false,
	// 	Bfr:      false,
	// 	Dbov:     false,
	// 	Cfpfe:    "ZnVuY3Rpb24oKXt2YXIgXzB4MTBhYzhhPV8weDVmMDJlMSxfMHgxMzQ4NmQ9ZG9jdW1lbnRbJ1x4NzFceDc1XHg2NVx4NzJceDc5XHg1M1x4NjVceDZjXHg2NVx4NjNceDc0XHg2Zlx4NzInXShfMHgxMGFjOGEoNjA4KSk7XzB4MTM0ODZkJiYhZnVuY3Rpb24gXzB4",
	// 	Stcfp:    "dHRwczovL2pzLmRhdGFkb21lLmNvL3RhZ3MuanM6MjoyMDI4NzYpCiAgICBhdCBfMHgzNGY4YjAuZGRfWiAoaHR0cHM6Ly9qcy5kYXRhZG9tZS5jby90YWdzLmpzOjI6MjE5MDI5KQogICAgYXQgaHR0cHM6Ly9qcy5kYXRhZG9tZS5jby90YWdzLmpzOjI6MTcxNjg1",
	// 	Prm:      true,
	// 	Tzp:      "Europe/Rome",
	// 	Cvs:      true,
	// 	Usb:      "defined",
	// 	Jset:     int(time.Now().Unix()),
	// }

	currentPath, _ := os.Getwd()
	fmt.Println("Getwd: ", currentPath)

	example, err := os.ReadFile(currentPath + "/internal/datadome_payload.json")
	if err != nil {
		panic(err.Error())
	}

	finger := &datadome.GenChPayload{}
	err = json.Unmarshal(example, &finger)

	if err != nil {
		d.Log.Error("unmarshaling payload error: ", err)
		return 0, err
	}

	jsetNum := time.Now().Unix()

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	ttst := round(random.Float64()*(40-8)+8, 15)

	tagGpu := round(ttst-(0.20*ttst), 14)

	tagpu := math.Round((tagGpu-(0.20*tagGpu))*1e14) / 1e14

	finger.Tagpu = tagpu
	finger.Ttst = ttst
	finger.Jset = int(jsetNum)

	json, _ := json.Marshal(finger)
	buffer := bytes.NewBuffer(json)

	fmt.Println("strData: ", buffer.String())

	req, err := http.NewRequest(http.MethodPost, "https://api-js.datadome.co/js/", nil)
	if err != nil {
		d.Log.Debug("Failed to issue GenCh request: ", err)
		return 0, err
	}

	// Add query parameters to the request body in a specific order
	params := []queryParam{
		{Key: "jsData", Value: buffer.String()},
		{Key: "eventCounters", Value: "[]"},
		{Key: "jsType", Value: "ch"},
		{Key: "cid", Value: cid},
		{Key: "ddk", Value: "A55FBF4311ED6F1BF9911EB71931D5"},
		{Key: "Referer", Value: "https%3A%2F%2Fwww.footlocker.com%2F"},
		{Key: "request", Value: "%2F"},
		{Key: "responsePage", Value: "origin"},
		{Key: "ddv", Value: "4.8.1"},
	}
	body := &strings.Builder{}
	for i, p := range params {
		if i > 0 {
			body.WriteByte('&')
		}
		body.WriteString(p.Key)
		body.WriteByte('=')
		body.WriteString(p.Value)
	}
	req.Body = io.NopCloser(strings.NewReader(body.String()))

	req.Header = http.Header{
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":       {"application/x-www-form-urlencoded"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"cross-site"},
		"Referer":            {"https://www.footlocker.com/"},
		"Referrer-Policy":    {"strict-origin-when-cross-origin"},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"content-type",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"Referer",
			"Referrer-Policy",
		},
	}

	res, err := d.Client.Do(req)
	if err != nil {
		d.Log.Debug("GenCh response error: ", err)
		return 0, err
	}

	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	d.Log.Debug("resBody: ", string(resBody))

	return res.StatusCode, nil
}

func round(f float64, n int) float64 {
	shift := math.Pow(10, float64(n))
	return math.Round(f*shift) / shift
}
