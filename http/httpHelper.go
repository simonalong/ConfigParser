package http

import (
	"encoding/json"
	"github.com/lunny/log"
	"github.com/simonalong/gole/util"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var httpClient = createHTTPClient()

const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 100
	IdleConnTimeout     int = 90
)

type NetError struct {
	ErrMsg string
}

func (error *NetError) Error() string {
	return error.ErrMsg
}

type StandardResponse struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},

		Timeout: 20 * time.Second,
	}
	return client
}

func SetHttpClient(httpClientOuter *http.Client) {
	httpClient = httpClientOuter
}

// ------------------ get ------------------

func GetSimple(url string) ([]byte, error) {
	return Get(url, nil, nil)
}

func GetSimpleOfStandard(url string) ([]byte, error) {
	return GetOfStandard(url, nil, nil)
}

func Get(url string, header http.Header, parameterMap map[string]string) ([]byte, error) {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return call(httpRequest, url)
}

func GetOfStandard(url string, header http.Header, parameterMap map[string]string) ([]byte, error) {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return callToStandard(httpRequest, url)
}

// ------------------ head ------------------

func HeadSimple(url string) error {
	return Head(url, nil, nil)
}

func Head(url string, header http.Header, parameterMap map[string]string) error {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return callIgnoreReturn(httpRequest, url)
}

// ------------------ post ------------------

func PostSimple(url string, body interface{}) ([]byte, error) {
	return Post(url, nil, nil, body)
}

func PostSimpleOfStandard(url string, body interface{}) ([]byte, error) {
	return PostOfStandard(url, nil, nil, body)
}

func Post(url string, header http.Header, parameterMap map[string]string, body interface{}) ([]byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("POST", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	return call(httpRequest, url)
}

func PostOfStandard(url string, header http.Header, parameterMap map[string]string, body interface{}) ([]byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("POST", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	return callToStandard(httpRequest, url)
}

// ------------------ put ------------------

func PutSimple(url string, body interface{}) ([]byte, error) {
	return Put(url, nil, nil, body)
}

func PutSimpleOfStandard(url string, body interface{}) ([]byte, error) {
	return PutOfStandard(url, nil, nil, body)
}

func Put(url string, header http.Header, parameterMap map[string]string, body interface{}) ([]byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PUT", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	return call(httpRequest, url)
}

func PutOfStandard(url string, header http.Header, parameterMap map[string]string, body interface{}) ([]byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PUT", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	return callToStandard(httpRequest, url)
}

// ------------------ delete ------------------

func DeleteSimple(url string) ([]byte, error) {
	return Get(url, nil, nil)
}

func DeleteSimpleOfStandard(url string) ([]byte, error) {
	return GetOfStandard(url, nil, nil)
}

func Delete(url string, header http.Header, parameterMap map[string]string) ([]byte, error) {
	httpRequest, err := http.NewRequest("DELETE", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return call(httpRequest, url)
}

func DeleteOfStandard(url string, header http.Header, parameterMap map[string]string) ([]byte, error) {
	httpRequest, err := http.NewRequest("DELETE", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return callToStandard(httpRequest, url)
}

// ------------------ patch ------------------

func PatchSimple(url string, body interface{}) ([]byte, error) {
	return Post(url, nil, nil, body)
}

func PatchSimpleOfStandard(url string, body interface{}) ([]byte, error) {
	return PostOfStandard(url, nil, nil, body)
}

func Patch(url string, header http.Header, parameterMap map[string]string, body interface{}) ([]byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PATCH", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	return call(httpRequest, url)
}

func PatchOfStandard(url string, header http.Header, parameterMap map[string]string, body interface{}) ([]byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PATCH", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Errorf("NewRequest err, %v", err.Error())
		return nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	return callToStandard(httpRequest, url)
}

func call(httpRequest *http.Request, url string) ([]byte, error) {
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil && httpResponse == nil {
		log.Printf("Error sending request to API endpoint. %+v", err)
		return nil, &NetError{ErrMsg: "Error sending request, url: " + url + ", err" + err.Error()}
	} else {
		if httpResponse == nil {
			log.Printf("httpResponse is nil")
			return nil, nil
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Infof("Body close err. %+v", err.Error())
			}
		}(httpResponse.Body)

		code := httpResponse.StatusCode
		if code != http.StatusOK {
			body, _ := ioutil.ReadAll(httpResponse.Body)
			return nil, &NetError{ErrMsg: "remote error, url: " + url + ", code " + strconv.Itoa(code) + ", message: " + string(body)}
		}

		// We have seen inconsistencies even when we get 200 OK response
		body, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			log.Printf("Couldn't parse response body %+v", err)
			return nil, &NetError{ErrMsg: "Couldn't parse response body, err: " + err.Error()}
		}

		return body, nil
	}
}

// ------------------ trace ------------------
// ------------------ options ------------------
// 暂时先不处理

func callIgnoreReturn(httpRequest *http.Request, url string) error {
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil && httpResponse == nil {
		log.Printf("Error sending request to API endpoint. %+v", err)
		return &NetError{ErrMsg: "Error sending request, url: " + url + ", err" + err.Error()}
	} else {
		if httpResponse == nil {
			log.Printf("httpResponse is nil")
			return nil
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Infof("Body close err. %+v", err.Error())
			}
		}(httpResponse.Body)

		code := httpResponse.StatusCode
		if code != http.StatusOK {
			body, _ := ioutil.ReadAll(httpResponse.Body)
			return &NetError{ErrMsg: "remote error, url: " + url + ", code " + strconv.Itoa(code) + ", message: " + string(body)}
		}

		return nil
	}
}

func callToStandard(httpRequest *http.Request, url string) ([]byte, error) {
	return parseStandard(call(httpRequest, url))
}

func parseStandard(responseResult []byte, errs error) ([]byte, error) {
	if errs != nil {
		return nil, errs
	}
	var standRsp StandardResponse
	err := json.Unmarshal(responseResult, &standRsp)
	if err != nil {
		return nil, err
	}

	if standRsp.Code == nil {
		return nil, &NetError{ErrMsg: "code is nil"}
	}

	// 判断业务的失败信息
	codeKind := reflect.ValueOf(standRsp.Code).Kind()
	if codeKind == reflect.String {
		if standRsp.Code != "success" {
			return nil, &NetError{ErrMsg: "remote err, bizCode=" + standRsp.Code.(string) + ", message=" + standRsp.Message}
		}
	} else if codeKind == reflect.Int || codeKind == reflect.Int8 || codeKind == reflect.Int16 || codeKind == reflect.Int32 || codeKind == reflect.Int64 || codeKind == reflect.Float32 || codeKind == reflect.Float64 {
		code := util.ToInt(standRsp.Code)
		if code != 0 && code != 200 {
			return nil, &NetError{ErrMsg: "remote err, bizCode=" + standRsp.Code.(string) + ", message=" + standRsp.Message}
		}
	}

	data, err := json.Marshal(standRsp.Data)
	if nil != err {
		return nil, err
	}

	return data, nil
}

func urlWithParameter(url string, parameterMap map[string]string) string {
	if parameterMap == nil || len(parameterMap) == 0 {
		return url
	}

	url += "?"

	var parameters []string
	for key, value := range parameterMap {
		parameters = append(parameters, key+"="+value)
	}

	return url + strings.Join(parameters, "&")
}
