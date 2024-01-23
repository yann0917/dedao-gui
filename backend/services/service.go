package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
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
	GuardDeviceID string `json:"_guard_device_id" mapstructure:"_guard_device_id"`
	SID           string `json:"_sid" mapstructure:"_sid"`
	AcwTc         string `json:"acw_tc" mapstructure:"acw_tc"`
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
	client.SetDebug(true)
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

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.New("404 NotFound")
	}
	if resp.StatusCode() == http.StatusBadRequest {
		return nil, errors.New("400 BadRequest")
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, errors.New("401 Unauthorized")
	}
	if resp.StatusCode() == 496 {
		return nil, errors.New("496 NoCertificate")
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
	// fmt.Printf("result.C:=%#v", result.C)
	if !result.isSuccess() {
		// 未登录或者登录凭证无效
		err = errors.New("服务异常，请稍后重试。errMsg:" + result.H.E)
		return err
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
	for _, v := range list {
		item := strings.Split(v, "=")
		if len(item) > 1 {
			if item[1] != "" {
				cookieM[strings.TrimSpace(item[0])] = item[1]
			}
		}
	}
	err = mapstructure.Decode(cookieM, v)
	return
}
