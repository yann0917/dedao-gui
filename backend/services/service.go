package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/yann0917/dedao-gui/backend/utils"
)

var (
	dedaoCommURL = &url.URL{
		Scheme: "https",
		Host:   "dedao.cn",
	}
	baseURL   = "https://www.dedao.cn"
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
)

// Response dedao success response
type Response struct {
	H respH `json:"h"`
	C respC `json:"c"`
}

type respH struct {
	C   int    `json:"c"`
	E   string `json:"e"`
	S   int    `json:"s"`
	T   int    `json:"t"`
	Apm string `json:"apm"`
}

// respC response content
type respC []byte

func (r *respC) UnmarshalJSON(data []byte) error {
	*r = data

	return nil
}

func (r respC) String() string {
	return string(r)
}

// Service dedao service
type Service struct {
	client *resty.Client
}

// CookieOptions dedao cookie options
type CookieOptions struct {
	GAT           string `json:"gat"`
	ISID          string `json:"isid"`
	Iget          string `json:"iget"`
	Token         string `json:"token"`
	CsrfToken     string `json:"csrfToken"`
	GuardDeviceID string `json:"_guard_device_id"`
	SID           string `json:"_sid"`
	AcwTc         string `json:"acw_tc"`
	AliyungfTc    string `json:"aliyungf_tc"`
	CookieStr     string `json:"cookieStr"`
}

// NewService new service
func NewService(co *CookieOptions) *Service {
	var cookies []*http.Cookie
	if co.GAT != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "GAT",
			Value:  co.GAT,
			Domain: "." + dedaoCommURL.Host,
		})
	}

	if co.ISID != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "ISID",
			Value:  co.ISID,
			Domain: "." + dedaoCommURL.Host,
		})
	}

	if co.GuardDeviceID != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "_guard_device_id",
			Value:  co.GuardDeviceID,
			Domain: "www." + dedaoCommURL.Host,
		})
	}

	if co.SID != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "_sid",
			Value:  co.SID,
			Domain: "www." + dedaoCommURL.Host,
		})
	}

	if co.AcwTc != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "acw_tc",
			Value:  co.AcwTc,
			Domain: "www." + dedaoCommURL.Host,
		})
	}

	if co.Iget != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "iget",
			Value:  co.Iget,
			Domain: "www." + dedaoCommURL.Host,
		})
	}

	if co.Token != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "token",
			Value:  co.Token,
			Domain: "www." + dedaoCommURL.Host,
		})
	}

	if co.CsrfToken != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "csrfToken",
			Value:  co.CsrfToken,
			Domain: "www." + dedaoCommURL.Host,
		})
	}

	if co.AliyungfTc != "" {
		cookies = append(cookies, &http.Cookie{
			Name:   "aliyungf_tc",
			Value:  co.AliyungfTc,
			Domain: "www." + dedaoCommURL.Host,
		})
	}

	client := resty.New()
	client.SetDebug(false)
	client.SetBaseURL(baseURL).
		SetCookies(cookies).
		SetHeaderVerbatim("User-Agent", UserAgent).
		SetHeaderVerbatim("Xi-DT", "web")

	if co.CsrfToken != "" {
		client.SetHeaderVerbatim("Xi-Csrf-Token", co.CsrfToken)
	}
	return &Service{client: client}
}

func (r *Response) isSuccess() bool {
	return r.H.C == 0
}

func handleHTTPResponse(resp *resty.Response, err error) (io.ReadCloser, error) {
	if err != nil {
		return nil, err
	}

	status := resp.StatusCode()

	// 403/429 常以 JSON body 出现（如 h.c=0 且 c=null），必须在解析内容之前拦截，
	// 否则会被当作成功解析出空页面
	if status == http.StatusForbidden || status == http.StatusTooManyRequests {
		return nil, &HTTPStatusError{Code: status, URL: resp.Request.URL}
	}

	// HTML 错误页（如 <h2>403 Forbidden</h2>）
	contentType := resp.Header().Get("Content-Type")
	if strings.Contains(contentType, "text/html") && status != http.StatusOK {
		return nil, &HTTPStatusError{Code: status, URL: resp.Request.URL, FromHTMLPage: true}
	}

	// Permanent errors that shouldn't be retried
	switch status {
	case http.StatusNotFound, http.StatusBadRequest, http.StatusUnauthorized, 496:
		return nil, &HTTPStatusError{Code: status, URL: resp.Request.URL}
	}

	data := resp.Body()
	reader := bytes.NewReader(data)
	result := io.NopCloser(reader)
	return result, nil
}

func handleJSONParse(reader io.Reader, v interface{}) error {
	result := new(Response)

	err := utils.UnmarshalReader(reader, &result)
	if err != nil {
		fmt.Printf("err1: %s \n", err.Error())
		return err
	}
	if !result.isSuccess() {
		// 业务错误（如 user no legal 的 code 4000），带 Code/Msg 供上层分类处理
		return &BusinessError{Code: result.H.C, Msg: result.H.E}
	}
	err = utils.UnmarshalJSON(result.C, v)
	if err != nil {
		fmt.Printf("err2: %s", err.Error())
		return err
	}

	return nil
}

// ParseCookies parse cookie string to cookie options
func ParseCookies(cookie string, v interface{}) (err error) {
	if cookie == "" {
		return errors.New("cookie is empty")
	}
	list := strings.Split(cookie, ";")
	cookieM := make(map[string]string, len(list))
	for _, item := range list {
		parts := strings.Split(item, "=")
		if len(parts) > 1 {
			if parts[1] != "" {
				cookieM[strings.TrimSpace(parts[0])] = parts[1]
			}
		}
	}

	// 创建大小写不敏感的 map（为了兼容 mapstructure 的行为）
	cookieMInsensitive := make(map[string]string)
	for k, v := range cookieM {
		cookieMInsensitive[strings.ToLower(k)] = v
	}

	// 使用反射将 map 的值赋给结构体
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a pointer to struct")
	}

	elem := value.Elem()
	structType := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := structType.Field(i)

		// 获取 json tag
		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		// 处理逗号后面的选项
		jsonName := strings.Split(tag, ",")[0]
		if jsonName == "" {
			jsonName = fieldType.Name
		}

		// 查找 map 中的值（大小写不敏感）
		if mapValue, ok := cookieMInsensitive[strings.ToLower(jsonName)]; ok {
			if field.Kind() == reflect.String {
				field.SetString(mapValue)
			}
		}
	}

	return nil
}
