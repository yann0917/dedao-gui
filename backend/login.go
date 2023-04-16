package backend

import (
	"errors"

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
	Status int            `json:"status"` // 1-登录成功,2-二维码过期
	Cookie string         `json:"cookie"` // cookies string
	User   *services.User `json:"user"`
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
	code, err := Instance.GetQrcode(token)
	if err != nil {
		return
	}
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
	if check != nil {
		if check.Data.Status == 1 {
			result.User, err = app.LoginByCookie(cookie)
			if err != nil {
				return
			}

			// fmt.Println("扫码成功")
		} else if check.Data.Status == 2 {
			err = errors.New("登录失败，二维码已过期")
			return
		}
		result.Status = check.Data.Status
	}
	return
}

func (a *App) Logout() (err error) {
	err = config.Instance.DeleteConfigFile()
	return
}

func (a *App) UserInfo() (user *services.User, err error) {
	user, err = Instance.User()
	return
}

func (a *App) EbookUserInfo() (user *services.EbookVIPInfo, err error) {
	user, err = Instance.EbookUserInfo()
	return
}

func (a *App) OdobUserInfo() (user *services.OdobVip, err error) {
	user, err = Instance.OdobUserInfo()
	return
}
