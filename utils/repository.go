package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"vocatrueid/entity"
	"vocatrueid/helpers"

	"github.com/go-rod/rod"
	"github.com/google/uuid"
)

type Repository interface {
	GetUsernameByCode(userId, zoneId, code string) entity.CheckNicknameResponse
}

type repository struct {
	browser *rod.Browser
}

func NewRepository(browser *rod.Browser) *repository {
	return &repository{browser}
}

const (
	NOT_VERIFY              string = "NOVERIFY"
	MOBILE_LEGENDS          string = "MOBILE_LEGENDS"
	FREEFIRE                string = "FREEFIRE"
	HIGGS_DOMINO            string = "HIGGS"
	RAGNAROK_X              string = "RAGNAROK_X"
	GENSHIN_IMPACT          string = "GENSHIN_IMPACT"
	CLOUD_SONG              string = "VNG_CLOUD_SONG"
	LORDS_MOBILE            string = "LORDS_MOBILE"
	RAGNAROK_ETERNAL        string = "GRAVITY_RAGNAROK_M"
	DRAGON_RAJA             string = "ZULONG_DRAGON_RAJA"
	LIFE_AFTER              string = "NETEASE_LIFEAFTER"
	HYPER_FRONT             string = "HYPER_FRONT"
	LAPLACE_M               string = "ZLONGAME"
	MIRRAGE_PERFECT_SKYLINE string = "PERFECT_SKYLINE"
	SAUSAGE_MAN             string = "SAUSAGE_MAN"
	VALORANT                string = "VALORANT"
	CALL_OF_DUTY_MOBILE     string = "CALL_OF_DUTY"
	APEX_LEGENDS_MOBILE     string = "APEX_LEGENDS"
	SUPER_SUS               string = "SUPER_SUS"
	ARENA_OF_VALOR          string = "AOV"
	TOWER_OF_FANTASY        string = "TOWER_OF_FANTASY"
	PLN                     string = "PLN"
	PUBG_ID                 string = "PUBG_ID"
)

var CODASHOP_TARGET_URL = "https://order-sg.codashop.com/initPayment.action"
var SMILEONE_TARGET_URL = "https://www.smile.one/merchant/mobilelegends/checkrole"
var SEPULSA_TARGET_URL = "https://api.sepulsa.com/api/v1/carts/add/"
var MIDASBUY_TARGET_URL = "https://www.midasbuy.com/unipin/id/buy/pubgm"

func (r repository) GetUsernameByCode(userId, zoneId, code string) entity.CheckNicknameResponse {
	var username string
	// if code == MOBILE_LEGENDS {
	// 	username = Mobapay(userId, zoneId)
	// } else
	if code == PLN {
		username = Sepulsa(userId)
	} else if code == PUBG_ID {
		username = Midasbuy(userId, zoneId, r.browser)
	} else {
		username = Codashop(userId, zoneId, code)
	}
	if username != "" {
		if username == "Required Zone ID" {
			return entity.CheckNicknameResponse{
				IsSuccess: false,
				Username:  username,
			}
		} else {
			return entity.CheckNicknameResponse{
				IsSuccess: true,
				Username:  username,
				UserId:    userId,
				ZoneId:    zoneId,
			}
		}
	} else {
		return entity.CheckNicknameResponse{
			IsSuccess: false,
			Username:  "",
		}
	}
}

func Mobapay(userId, zoneId string) string {
	baseUrl := "https://api.mobapay.com/api/app_shop?app_id=100000&shop_id=1001&user_id=" + userId + "&server_id=" + zoneId + "&country=ID&language=id"

	req, err := http.NewRequest("GET", baseUrl, strings.NewReader(""))

	if err != nil {
		println("(Mobapay)Log => error: " + err.Error())
	}

	req.Header.Add("Authority", "api.mobapay.com")
	req.Header.Add("Origin", "https://www.mobapay.com")
	req.Header.Add("Referer", "https://www.mobapay.com/")

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	var data entity.MobapayResponse
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		println("(Mobapay)Log => error: " + err.Error())
		return ""
	}

	return data.Data.User_info.User_name
}

func Sepulsa(userId string) string {
	payload := []byte(`{"url":"http://api.sepulsa.com/api/v1/oscar/products/13/","quantity":"1","options":[{"option":"https://api.sepulsa.com/api/v1/oscar/options/1/","value":` + userId + `},{"option":"https://api.sepulsa.com/api/v1/oscar/options/3/","value":"08123454562"}]}`)
	req, err := http.NewRequest("POST", SEPULSA_TARGET_URL, bytes.NewBuffer(payload))

	if err != nil {
		println("(Sepulsa)Log => error: " + err.Error())
		return ""
	}

	req.Header.Set("X-Chital-Api-Key", "qQFAFT8d.6Yt44sZWZdkd1P4jFwAv4E5UyEp9QYNw")
	req.Header.Set("X-Chital-Order-Source", "web")
	req.Header.Set("X-Chital-Requester", "https://www.sepulsa.com")
	req.Header.Set("Authority", "api.sepulsa.com")
	req.Header.Set("Origin", "https://www.sepulsa.com")
	req.Header.Set("Referer", "https://www.sepulsa.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	proxyURL := helpers.GetEnv("PROXY_URL")
	var client http.Client
	if proxyURL != "" {
		proxyURL, err := url.Parse(proxyURL)
		if err != nil {
			println("(Sepulsa)Log => error: " + err.Error())
			return ""
		}
		transport := http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{},
		}
		client = http.Client{
			Transport: &transport,
			Timeout:   10 * time.Second,
		}

	} else {
		client = http.Client{
			Timeout: 10 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}

	var data entity.SepulsaResponse
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		println("(Sepulsa)Log => error: " + err.Error())
		return ""
	}
	if len(data.Data.Lines) == 0 {
		return ""
	}
	var value struct {
		NamaPelanggan string `json:"Nama Pelanggan"`
	}
	err = json.Unmarshal([]byte(data.Data.Lines[0].Attributes[2].Value), &value)
	if err != nil {
		println("(Sepulsa)Log => error: " + err.Error())
		return ""
	}
	return value.NamaPelanggan
}

func Codashop(userId, zoneId, code string) string {
	payload, message := GetPayload(userId, zoneId, code)
	if message != "" {
		return message
	}
	req, err := http.NewRequest("POST", CODASHOP_TARGET_URL, bytes.NewBuffer(payload))

	if err != nil {
		println("(Codashop)Log => error: " + err.Error())
		return ""
	}

	req.Header.Set("Proxy", "false")

	proxyURL := helpers.GetEnv("PROXY_URL")
	var client http.Client
	if proxyURL != "" {
		proxyURL, err := url.Parse(proxyURL)
		if err != nil {
			println("(Codashop)Log => error: " + err.Error())
			return ""
		}
		transport := http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{},
		}
		client = http.Client{
			Transport: &transport,
			Timeout:   10 * time.Second,
		}

	} else {
		client = http.Client{
			Timeout: 10 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		println("(Codashop)Log => error: " + err.Error())
		return ""
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		println("(Codashop)Log => error: " + err.Error())
		return ""
	}
	var data entity.CodashopResponse
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		println("(Codashop)Log => error: " + err.Error())
		return ""
	}
	return data.ConfirmationFields.Username
}

func GetPayload(userId, zoneId, code string) ([]byte, string) {
	id := uuid.New()
	var payload []byte
	var message string
	switch code {
	case MOBILE_LEGENDS:
		if zoneId == "" {
			message = "Required Zone ID"
			break
		}
		payload = []byte(`{"voucherPricePoint.id":"4151","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":` + zoneId + `,"voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case FREEFIRE:
		payload = []byte(`{"voucherPricePoint.id":"270281","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"voucherTypeId":"17","gvtId":"33","shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case HIGGS_DOMINO:
		payload = []byte(`{"voucherPricePoint.id":"27577","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case GENSHIN_IMPACT:
		if zoneId == "" {
			message = "Required Zone ID"
			break
		}
		payload = []byte(`{"voucherPricePoint.id":"116054","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":` + zoneId + `,"voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case RAGNAROK_X:
		if zoneId == "" {
			message = "Required Zone ID"
			break
		}
		payload = []byte(`{"voucherPricePoint.id":"195773","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":` + zoneId + `,"voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case CLOUD_SONG:
		payload = []byte(`{"voucherPricePoint.id":"210657","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case LORDS_MOBILE:
		payload = []byte(`{"voucherPricePoint.id":"49955","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"1051","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case DRAGON_RAJA:
		payload = []byte(`{"voucherPricePoint.id":"75566","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case LIFE_AFTER:
		if zoneId == "" {
			message = "Required Zone ID"
			break
		}
		payload = []byte(`{"voucherPricePoint.id":"45706","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":` + zoneId + `,"voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case HYPER_FRONT:
		if zoneId == "" {
			message = "Required Zone ID"
			break
		}
		payload = []byte(`{"voucherPricePoint.id":"267156","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":` + zoneId + `,"voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case LAPLACE_M:
		payload = []byte(`{"voucherPricePoint.id":"25471","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case MIRRAGE_PERFECT_SKYLINE:
		if zoneId == "" {
			message = "Required Zone ID"
			break
		}
		payload = []byte(`{"voucherPricePoint.id":"255831","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":` + zoneId + `,"voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case SAUSAGE_MAN:
		payload = []byte(`{"voucherPricePoint.id":"256513","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"global-release","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case CALL_OF_DUTY_MOBILE:
		payload = []byte(`{"voucherPricePoint.id":"270249","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case APEX_LEGENDS_MOBILE:
		payload = []byte(`{"voucherPricePoint.id":"355073","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"voucherTypeId":"203","gvtId":"257","shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case SUPER_SUS:
		payload = []byte(`{"voucherPricePoint.id":"266077","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case VALORANT:
		payload = []byte(`{"voucherPricePoint.id":"115691","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"voucherTypeId":"109","gvtId":"139","shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	case ARENA_OF_VALOR:
		payload = []byte(`{"voucherPricePoint.id":"270294","voucherPricePoint.price":"0","voucherPricePoint.variablePrice":"0","user.userId":` + userId + `,"user.zoneId":"","voucherTypeName":` + code + `,"voucherTypeId":"16","gvtId":"32","shopLang":"id_ID","checkoutId":` + id.String() + `}`)
		break
	}

	return payload, message
}

func Midasbuy(userId, zoneId string, browser *rod.Browser) string {
	page := browser.MustPage(MIDASBUY_TARGET_URL).MustWaitLoad()
	isPopup, popupEl, _ := page.Has("#uc_landing_pop > div > div")
	if isPopup {
		popupEl.MustClick()
	}
	page.MustElement("#app > div.content > div.x-main > div.section.tab-nav-box.sub-game-account.sub-id.g-clr > div > div > div > div.input-box.have-select-input-box > input").MustInput(userId)
	page.MustElement("#app > div.content > div.x-main > div.section.tab-nav-box.sub-game-account.sub-id.g-clr > div > div > div > div.btn").MustClick()
	partition := page.MustElements("div.li")
	for i := 0; i < len(partition); i++ {
		v := partition[i]
		if v.MustText() == zoneId {
			v.MustClick()
			break
		}
	}
	time.Sleep(time.Second * 2)
	isSuccess := page.MustHas("div.name")
	if isSuccess {
		nickname := page.MustElement("div.name").MustText()
		page.MustElement("li[data-id='os_midaspay_id_unipin_wallet']").MustClick()
		page.MustElement("li[cr='amount_select.250']").MustClick()
		time.Sleep(1 * time.Second)
		waitOpen := page.MustWaitOpen()
		page.MustElement("#buy-payBtn").MustClick()
		birthday, _, _ := page.Has("#birthday-pop[style='display: block;']")
		if birthday {
			page.MustElement("#birthdayDateInput").MustClick()
			time.Sleep(100 * time.Millisecond)
			page.MustElement(`div.time-picker-box > div.bd > ul > li:nth-child(2)`).MustClick()
			page.MustElement("#birthday-pop > div > div.btn-wrap > div").MustClick()
		}
		agreement, _, _ := page.Has("#pop-box[style='display: block;']")
		if agreement {
			cb := page.MustElements("div.check-box")
			for _, v := range cb {
				v.MustClick()
			}
			wait := page.MustWaitOpen()
			page.MustElement("#pop-box > div > div > div.bottom-fixed.none-border > div > div").MustClick()
			newPage := wait()
			time.Sleep(1 * time.Second)
			isErrPopup, errEl, _ := page.Has("div.pop-mode.pay-fail-pop.show > div.mess")
			if isErrPopup {
				msg := errEl.MustText()
				if msg != "Pembayaran gagal." {
					page.Close()
					return ""
				} else {
					newPage.Close()
					page.Close()
					return nickname
				}
			} else {
				newPage.Close()
				page.Close()
				return nickname
			}
		} else {
			newPage := waitOpen()
			time.Sleep(1 * time.Second)
			isErrPopup, errEl, _ := page.Has("div.pop-mode.pay-fail-pop.show > div.mess")
			if isErrPopup {
				msg := errEl.MustText()
				if msg != "Pembayaran gagal." {
					page.Close()
					return ""
				} else {
					newPage.Close()
					page.Close()
					return nickname
				}
			} else {
				newPage.Close()
				page.Close()
				return nickname
			}

		}
	} else {
		page.Close()
		return ""
	}
}
