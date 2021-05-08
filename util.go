package HBDM_API

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"time"
)

func ComputeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func Map2UrlQueryBySort(mapParams map[string]string) string {
	var keys []string
	for key := range mapParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strParams string
	for _, key := range keys {
		strParams += key + "=" + mapParams[key] + "&"
	}
	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}
	return strParams
}
func CreateSign(mapParams map[string]string, strMethod, strHostUrl, strRequestPath, strSecretKey string) string {
	// 参数处理, 按API要求, 参数名应按ASCII码进行排序(使用UTF-8编码, 其进行URI编码, 16进制字符必须大写)
	mapCloned := make(map[string]string)
	for key, value := range mapParams {
		mapCloned[key] = url.QueryEscape(value)
	}

	strParams := Map2UrlQueryBySort(mapCloned)

	strPayload := strMethod + "\n" + strHostUrl + "\n" + strRequestPath + "\n" + strParams
	return ComputeHmac256(strPayload, strSecretKey)
}

func MapValueEncodeURI(mapValue map[string]string) map[string]string {
	for key, value := range mapValue {
		valueEncodeURI := url.QueryEscape(value)
		mapValue[key] = valueEncodeURI
	}

	return mapValue
}
func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	for key, value := range mapParams {
		strParams += (key + "=" + value + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}
func ApiKeyReady(strMethod, strRequestPath string, ACCESS_KEY, SECRET_KEY string) string {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams2Sign := make(map[string]string)
	mapParams2Sign["AccessKeyId"] = ACCESS_KEY
	mapParams2Sign["SignatureMethod"] = "HmacSHA256"
	mapParams2Sign["SignatureVersion"] = "2"
	mapParams2Sign["Timestamp"] = timestamp

	hostName := "api.hbdm.com"
	mapParams2Sign["Signature"] = CreateSign(mapParams2Sign, strMethod, hostName, strRequestPath, SECRET_KEY)
	strUrl := "https://api.hbdm.com" + strRequestPath + "?" + Map2UrlQuery(MapValueEncodeURI(mapParams2Sign))
	return strUrl
}

//--------------------------------------------------------------------------------------------------------------------
// 参数序列化
func ParamReady(url string, param interface{}, method string) *http.Request {
	jsonParams := ""
	if nil != param {
		bytesParams, _ := json.Marshal(param)
		jsonParams = string(bytesParams)
	}

	var req *http.Request
	if method == "POST" {
		req, _ = http.NewRequest(method, url, strings.NewReader(jsonParams))
	} else if method == "GET" {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Add("Content-Type", "application/json")
	return req
}

// Request请求封装
func RequestDo(client *http.Client, req *http.Request, param interface{}) (interface{}, error) {
	resp, err := client.Do(req)
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		return param, fmt.Errorf(" %v() Request error %v", pc, err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, param)
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		return param, fmt.Errorf(" %v() Unmarshal error %v", pc, err)
	}
	return param, nil
}
