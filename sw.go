package main

//curl -X "POST" http://api.servicewall.net/ganon/v1/collector/access_log -H 'Content-Type: application/json; charset=utf-8' -d $'{"source": "test"}'
//sudo ./goreplay --input-raw :8000 --output-http "api.servicewall.net" --output-http-sw-api "http://api.servicewall.net/ganon/v1/collector/access_log" --output-http-sw-source "tUrH27_iHhL"
//sudo ./goreplay --input-raw :8000 --output-kafka-host '172.31.44.1:9092,172.31.44.2:9092,172.31.44.3:9092'  -output-kafka-sw-source "tUrH27_iHhL"  -input-raw-realip-header sw-ip

import (
	"encoding/json"
	"github.com/buger/goreplay/byteutils"
	"github.com/buger/goreplay/proto"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseCookie(cookie string) map[string]string {
	rawCookies := cookie

	header := http.Header{}
	header.Add("Cookie", rawCookies)
	request := http.Request{Header: header}

	//fmt.Println(request.Cookies()) // [cookie1=value1 cookie2=value2]

	cookieMap := make(map[string]string)
	cookies := request.Cookies()
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}

	return cookieMap
}

func parseIp(headerMap map[string][]string) string {
	ip := ""
	if len(headerMap["x-real-ip"]) > 0 {
		ip = headerMap["x-real-ip"][0]
	} else if len(headerMap["x-forwarded-for"]) > 0 {
		ips := strings.Split(headerMap["x-forwarded-for"][0], ",")
		if len(ips) != 0 && ips[0] != "" {
			ip = ips[0]
		}
	} else if len(headerMap["sw-ip"]) > 0 {
		ip = headerMap["sw-ip"][0]
	}
	return ip
}
func buildRequestBody(source string, mimeHeader map[string][]string, meta [][]byte, req []byte) []byte {
	// build up header map.
	headerMap := make(map[string][]string)
	for k, v := range mimeHeader {
		// convert header key name to lower case.
		headerMap[strings.ToLower(k)] = v
	}

	ts, _ := strconv.ParseInt(byteutils.SliceToString(meta[2]), 10, 64)
	ip := parseIp(headerMap)
	//println(headerMap["cookie"][0])
	//, err = os.Stdout.Write(msg.Data)
	cookieMap := make(map[string]string)
	if len(headerMap["cookie"]) > 0 {
		cookieMap = parseCookie(headerMap["cookie"][0])
	}

	header := make(map[string]string)
	for k, v := range mimeHeader {
		header[k] = strings.Join(v, ", ")
	}

	ua := ""
	if len(headerMap["user-agent"]) > 0 {
		ua = headerMap["user-agent"][0]
	}
	referer := ""
	if len(headerMap["referer"]) > 0 {
		referer = headerMap["referer"][0]
	}
	host := ""
	if len(headerMap["host"]) > 0 {
		host = headerMap["host"][0]
	}
	body := requestBody{
		Source:  source,
		LogType: "access_log",
		BasicInfo: basicInfo{
			Timestamp: ts / int64(time.Millisecond),
			Ip:        ip,
			Ua:        ua,
			Url:       byteutils.SliceToString(proto.Path(req)),
			Method:    byteutils.SliceToString(proto.Method(req)),
			Headers:   header,
			Cookies:   cookieMap,
			Referer:   referer,
			Host:      host,
		},
		RawData: string(req),
	}
	data, _ := json.Marshal(body)
	//SW: DEBUG
	os.Stdout.Write(data)

	return data
}

func buildSwMessage(source string, mimeHeader map[string][]string, meta [][]byte, req []byte) collectedData {

	headerMap := make(map[string][]string)
	for k, v := range mimeHeader {
		// convert header key name to lower case.
		headerMap[strings.ToLower(k)] = v
	}

	ts, _ := strconv.ParseInt(byteutils.SliceToString(meta[2]), 10, 64)
	ip := parseIp(headerMap)
	//println(headerMap["cookie"][0])
	//, err = os.Stdout.Write(msg.Data)
	cookieMap := make(map[string]string)
	if len(headerMap["cookie"]) > 0 {
		cookieMap = parseCookie(headerMap["cookie"][0])
	}

	return collectedData{
		RequestParam: nil,
		RequestBody:  string(buildRequestBody(source, mimeHeader, meta, req)),
		Header:       headerMap,
		Cookie:       cookieMap,
		Method:       byteutils.SliceToString(proto.Method(req)),
		Uri:          byteutils.SliceToString(proto.Path(req)),
		Timestamp:    ts / int64(time.Millisecond),
		TraceID:      byteutils.SliceToString(meta[1]),
		ClientIP:     ip,
		//			ReqType:    byteutils.SliceToString(meta[0]),//
	}
}

//SW: type from Ganon
type collectedData struct {
	RequestParam map[string][]string `json:"request_param"`
	RequestBody  string              `json:"request_body"`
	Header       map[string][]string `json:"header"`
	Cookie       map[string]string   `json:"cookie"`
	Method       string              `json:"method"`
	Uri          string              `json:"uri"`
	Timestamp    int64               `json:"timestamp"`
	TraceID      string              `json:"trace_id"`
	ClientIP     string              `json:"client_ip"`
}

type basicInfo struct {
	//req
	Timestamp  int64             `json:"timestamp"`
	Ip         string            `json:"ip"`
	Ua         string            `json:"ua"`
	Url        string            `json:"url"`
	StatusCode int               `json:"statusCode"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Cookies    map[string]string `json:"cookies"`
	Referer    string            `json:"referer"`
	Host       string            `json:"host"`
}

type requestBody struct {
	//req
	Source    string    `json:"source"`
	LogType   string    `json:"logType"`
	RawData   string    `json:"rawData"`
	UserId    string    `json:"userId"`
	SessionId string    `json:"sessionId"`
	BasicInfo basicInfo `json:"basicInfo"`
}
