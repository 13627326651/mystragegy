package util

/*import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
	"tinyquant/src/logger"
	"tinyquant/src/mod"

	"go.uber.org/zap"
)

var (
	client *http.Client
)

func init() {
	if client == nil {
		client = initHttpClient()
	}
}

func initHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Second,
			}).DialContext,
			MaxIdleConns:        30,               //最大空闲连接数
			MaxIdleConnsPerHost: 60,               //最大与服务器的连接数  默认是2
			IdleConnTimeout:     30 * time.Second, //空闲连接保持时间
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse("http://127.0.0.1:7890")
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
		},
	}
	return client
}

func HttpRequest(ctx context.Context, req *mod.ReqParam) (map[string]interface{}, error) {

	urlx := fmt.Sprintf("%s%s", BaseURL, req.URL)

	queryString := req.Query.Encode()
	if queryString != "" {
		urlx = fmt.Sprintf("%s?%s", urlx, queryString)
	}
	fmt.Println("full urlx : ", urlx)
	r, err := http.NewRequest(req.Method, urlx, nil)
	r = r.WithContext(ctx)
	if r.Header.Get("User-Agent") == "" {
		r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}
	if req.APIKEY != "" {
		r.Header.Add("X-MBX-APIKEY", req.APIKEY)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	if err != nil {
		logger.Logger.Error("http request failed ", zap.Error(err))
		return nil, err
	}
	res, err := client.Do(r)
	if err != nil {
		logger.Logger.Error("http Do failed ", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Logger.Error("io read failed ", zap.Error(err))
		return nil, err
	}
	fmt.Println(string(body))
	var msg map[string]interface{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		logger.Logger.Error("json unmarshal failed : ", zap.Error(err))
		return nil, err
	}
	return msg, nil

}*/
