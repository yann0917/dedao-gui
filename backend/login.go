package backend

import (
	"errors"
	"fmt"

	"github.com/yann0917/dedao/app"
	"github.com/yann0917/dedao/config"
	"github.com/yann0917/dedao/services"
)

type QrCodeResp struct {
	Token        string `json:"token"`
	QrCode       string `json:"qrCode"`
	QrCodeString string `json:"qrCodeString"`
}

type LoginResult struct {
	Status int    `json:"status"` // 1-登录成功,2-二维码过期
	Cookie string `json:"cookie"` // cookies string
}

var Instance *services.Service

func init() {
	Instance = config.Instance.ActiveUserService()
}
func (a *App) GetQrcode() (qrCode QrCodeResp, err error) {
	token, err := Instance.LoginAccessToken()
	if err != nil {
		return
	}
	// fmt.Printf("token:%#v\n", token)
	code, err := Instance.GetQrcode(token)
	if err != nil {
		return
	}
	fmt.Printf("code:%#v\n", code)
	qrCode.Token = token
	if code != nil {
		qrCode.QrCode = code.Data.QrCode
		qrCode.QrCodeString = code.Data.QrCodeString
	}
	return
}

func (a *App) CheckLogin(token, qrCodeString string) (result LoginResult, err error) {
	check, cookie, err := Instance.CheckLogin(token, qrCodeString)
	if err != nil {
		return
	}
	result.Cookie = cookie
	fmt.Println(cookie)
	if check != nil {
		if check.Data.Status == 1 {
			err = app.LoginByCookie(cookie)
			if err != nil {
				return
			}
			fmt.Println("扫码成功")
		} else if check.Data.Status == 2 {
			err = errors.New("登录失败，二维码已过期")
			return
		}
		result.Status = check.Data.Status
	}
	return
}

func (a *App) UserInfo() (user *services.User, err error) {
	user, err = Instance.User()
	if err != nil {
		return
	}
	return
}

func (a *App) EbookUserInfo() (user *services.EbookVIPInfo, err error) {
	user, err = Instance.EbookUserInfo()
	if err != nil {
		return
	}
	fmt.Println(user)
	return
}

func (a *App) OdobUserInfo() (user *services.OdobVip, err error) {
	user, err = Instance.OdobUserInfo()
	if err != nil {
		return
	}
	fmt.Println(user)
	return
}
