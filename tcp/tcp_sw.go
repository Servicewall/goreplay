package tcp

//curl -X "POST" http://api.servicewall.net/ganon/v1/collector/access_log -H 'Content-Type: application/json; charset=utf-8' -d $'{"source": "test"}'
//sudo ./goreplay --input-raw :8000 --output-http "api.servicewall.net" --output-http-sw-api "http://api.servicewall.net/ganon/v1/collector/access_log" --output-http-sw-source "tUrH27_iHhL"
//sudo ./goreplay --input-raw :8000 --output-kafka-host '172.31.44.1:9092,172.31.44.2:9092,172.31.44.3:9092'  -output-kafka-sw-source "tUrH27_iHhL"  -input-raw-realip-header sw-ip

func getSYNFp(pckt *Packet) string {

	return "test"
}

//SW: type from Ganon
//type collectedData struct {
//	RequestParam map[string][]string `json:"request_param"`
//	RequestBody  string              `json:"request_body"`
//	Header       map[string][]string `json:"header"`
//	Cookie       map[string]string   `json:"cookie"`
//	Method       string              `json:"method"`
//	Uri          string              `json:"uri"`
//	Timestamp    int64               `json:"timestamp"`
//	TraceID      string              `json:"trace_id"`
//	ClientIP     string              `json:"client_ip"`
//}
//
//type basicInfo struct {
//	//req
//	Timestamp  int64             `json:"timestamp"`
//	Ip         string            `json:"ip"`
//	Ua         string            `json:"ua"`
//	Url        string            `json:"url"`
//	StatusCode int               `json:"statusCode"`
//	Method     string            `json:"method"`
//	Headers    map[string]string `json:"headers"`
//	Cookies    map[string]string `json:"cookies"`
//	Referer    string            `json:"referer"`
//	Host       string            `json:"host"`
//}
//
//type requestBody struct {
//	//req
//	Source    string    `json:"source"`
//	LogType   string    `json:"logType"`
//	RawData   string    `json:"rawData"`
//	UserId    string    `json:"userId"`
//	SessionId string    `json:"sessionId"`
//	BasicInfo basicInfo `json:"basicInfo"`
//}
