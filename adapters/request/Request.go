package request

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type FileItem struct {
	Key      string //image_content
	FileName string //test.jpg
	Content  []byte //[]byte
}

type sRequest struct {
	url               string
	method            string
	params            map[string]string
	headers           map[string]string
	body              io.Reader
	suppressParseData bool
}

type SRequest interface {
	send() (SResponse, error)
	SetMethod(string) *sRequest
	SetParams(map[string]string) *sRequest
	SetHeaders(map[string]string) *sRequest
	SetBody(map[string]interface{}) *sRequest
	SetMultipartForm(map[string]string, FileItem) *sRequest
	SetBasicAuth(string, string) *sRequest
	Get() (SResponse, error)
	Post() (SResponse, error)
	Put() (SResponse, error)
	Patch() (SResponse, error)
	Delete() (SResponse, error)
}

func Make(url string) SRequest {
	return &sRequest{
		url:     url,
		method:  "GET",
		headers: map[string]string{},
		body:    &bytes.Buffer{},
	}
}
func (sR *sRequest) send() (SResponse, error) {
	var rs SResponse
	req, err := http.NewRequest(sR.method, sR.url, sR.body)
	if err != nil {
		return rs, err
	}
	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   600 * time.Second,
			KeepAlive: 300 * time.Second,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		TLSHandshakeTimeout: 600 * time.Second,
	}
	for key, value := range sR.headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{Transport: t}
	response, errDo := client.Do(req)
	if errDo != nil {
		return rs, errDo
	}
	defer response.Body.Close()
	data, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		return rs, errRead
	}
	if !sR.suppressParseData {
		err1 := json.Unmarshal(data, &rs)
		return rs, err1
	} else {
		rs.Data = string(data)
		rs.Status = response.StatusCode
	}
	return rs, nil
}

func (sR *sRequest) Get() (SResponse, error) {
	sR.method = "GET"
	queryStr := getUrlQueryString(sR.params)
	if queryStr != "" {
		if strings.Contains(sR.url, "?") {
			sR.url += "&" + queryStr
		} else {
			sR.url += "?" + queryStr
		}
	}
	return sR.send()
}

func getUrlQueryString(query map[string]string) string {
	str := ""
	for k, v := range query {
		str += k + "=" + v + "&"
	}
	str = strings.TrimSuffix(str, "&")
	return str
}

func (sR *sRequest) Post() (SResponse, error) {
	sR.method = "POST"
	return sR.send()
}
func (sR *sRequest) Put() (SResponse, error) {
	sR.method = "PUT"
	return sR.send()
}
func (sR *sRequest) Patch() (SResponse, error) {
	sR.method = "PATCH"
	return sR.send()
}
func (sR *sRequest) Delete() (SResponse, error) {
	sR.method = "DELETE"
	return sR.send()
}

func (sR *sRequest) SetMethod(m string) *sRequest {
	sR.method = m
	return sR
}
func (sR *sRequest) SetParams(p map[string]string) *sRequest {
	sR.params = p
	return sR
}
func (sR *sRequest) SetHeaders(h map[string]string) *sRequest {
	for k, v := range h {
		sR.headers[k] = v
	}
	return sR
}
func (sR *sRequest) SetBody(b map[string]interface{}) *sRequest {
	var body []byte
	body, _ = json.Marshal(b)
	sR.body = bytes.NewBuffer(body)
	return sR
}
func (sR *sRequest) SetBodyUrlEncoded(b map[string]interface{}) *sRequest {
	data := url.Values{}
	for k, v := range b {
		data.Set(k, fmt.Sprintf("%v", v))
	}
	sR.body = strings.NewReader(data.Encode())
	return sR
}
func (sR *sRequest) SetSuppressParseData(b bool) *sRequest {
	sR.suppressParseData = b
	return sR
}
func (sR *sRequest) SetMultipartForm(m map[string]string, fileItem FileItem) *sRequest {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range m {
		writer.WriteField(k, v)
	}
	fileWriter, err := writer.CreateFormFile(fileItem.Key, fileItem.FileName)
	if err != nil {
		fmt.Println(err)
	}
	fileWriter.Write(fileItem.Content)
	writer.Close()
	sR.headers["Content-Type"] = writer.FormDataContentType()
	sR.body = body
	return sR
}
func (sR *sRequest) SetBasicAuth(userName, password string) *sRequest {
	sR.headers["Authorization"] = "Basic " + basicAuth(userName, password)
	return sR
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
