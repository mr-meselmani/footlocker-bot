package internal

import (
	"footlocker-bot/internal/logger"

	tls_client "github.com/bogdanfinn/tls-client"
)

var headers = map[string]string{
	"cache-control":      "no-cache",
	"sec-ch-ua":          "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"112\", \"Google Chrome\";v=\"105\"",
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": "\"Windows\"",
	"user-agent":         "Mozilla/5.0 (Linux; Android 9; K21) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	"accept":             "*/*",
	"sec-fetch-site":     "cross-site",
	"sec-fetch-mode":     "cors",
	"sec-fetch-dest":     "empty",
	"referer":            "https://www.footlocker.com/",
	"accept-encoding":    "gzip, deflate, br",
	"accept-language":    "en-US,en;q=0.9",
	"content-type":       "application/x-www-form-urlencoded;charset=UTF-8",
}

type Datadome struct {
	httpClient tls_client.HttpClient
	endpoint   string
	version    string
	headers    map[string]string
	Log        logger.Logger
}

func NewDatadome() Datadome {
	return Datadome{}
}

func (d *Datadome) GetDatadomeSettings() error {
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
		d.Log.Error("Get settings error", err)
		return err
	}

	d.httpClient = client
	d.version = "4.6.0"
	d.endpoint = "https://api-js.datadome.co/js/"
	d.headers = headers

	return nil
}

func (d *Datadome) GetCookieM() (int, error) {
	

	return 0, nil
}
