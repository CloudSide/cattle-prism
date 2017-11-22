package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
    "cattle-prism/utils/wsutil"
	"github.com/gorilla/websocket"
	"cattle-prism/models"
	"fmt"
	"regexp"
	"time"
	"errors"
	"encoding/json"
	"net/http"
	"net/http/httputil"
    "net/url"
	"log"
	// "net/http"
	"strings"
)

// Operations about Users
type AppController struct {
	beego.Controller
	Token 			string
	UserInfo		models.TokenDataItemUserIdentity
}

var RancherEndpointHost string
var UserInfoCache, _ = cache.NewCache("memory", `{"interval":60}`)

func init() {
	if RancherEndpointHost = beego.AppConfig.String("RancherEndpointHost"); RancherEndpointHost == "" {
		RancherEndpointHost = "127.0.0.1:8080"
	}
}

func (this *AppController) ServeError(status int, err error, message string) {
	logs.Error(err)
	this.Data["json"] = &models.CattleError{
		Type:		"error",
		BaseType:	"error",
		Status:		status,
		Code:		message,
		Message:	message,
		Detail:		err.Error(),
	}
	this.Ctx.Output.SetStatus(status)
	this.ServeJSON()
	this.StopRun()
}

func (this *AppController) GetUserInfo() {
	this.Token = this.Ctx.Input.Cookie("token")
	if this.Token != "" {
		var tokenData models.TokenResource
		cacheKey := fmt.Sprintf("userinfo_%s", this.Token)
		cacheData := UserInfoCache.Get(cacheKey)
		if cacheData == nil {
			tokenRequest := httplib.Get(`http://` + RancherEndpointHost + `/v2-beta/token`)
			for headerName, _ := range this.Ctx.Request.Header {
				if headerName != "User-Agent" {
					tokenRequest.Header(headerName, this.Ctx.Input.Header(headerName))
				}
			}
			body, err := tokenRequest.Bytes()
			if err != nil {
				this.ServeError(500, err, "Internal Server Error")
			}
			err = json.Unmarshal(body, &tokenData)
			if err != nil {
				this.ServeError(500, err, "Internal Server Error")
			}
			if len(tokenData.Data) > 0 {
				this.UserInfo = tokenData.Data[0].UserIdentity
				this.UserInfo.Cached = false
				if this.UserInfo.Id != "" {
					if err := UserInfoCache.Put(cacheKey, body, 60 * time.Second); err != nil {
						this.ServeError(500, err, "Internal Server Error")
					}
				}
				/*
				this.Data["json"] = this.UserInfo
				this.ServeJSON()
				this.StopRun()
				*/
			}
		} else {
			if cacheBytes, ok := cacheData.([]byte); ok {
				err := json.Unmarshal(cacheBytes, &tokenData)
				if err != nil {
					UserInfoCache.Delete(cacheKey)
					this.ServeError(500, err, "Internal Server Error")
				}
				if len(tokenData.Data) > 0 {
					this.UserInfo = tokenData.Data[0].UserIdentity
					this.UserInfo.Cached = true
					/*
					this.Data["json"] = this.UserInfo
					this.ServeJSON()
					this.StopRun()
					*/
				}
			} else {
				UserInfoCache.Delete(cacheKey)
			    this.ServeError(500, errors.New("Decoding Cache Data Error"), "Internal Server Error")
			}
		}
	}
}


func (this *AppController) Prepare() {
	this.GetUserInfo()

	if wsutil.IsWebSocketRequest(this.Ctx.Request) && this.Ctx.Input.IsGet() {
		re := regexp.MustCompile(`^/v2-beta/projects/[a-z0-9]+/subscribe$`)
		if matched := re.MatchString(this.Ctx.Input.URL()); matched {
			// fmt.Println(this.Ctx.Input.URL())
			this.Subscribe()
			return
		}
	}

	this.Proxy()
}


func (this *AppController) Proxy() {
	if this.UserInfo.CompanyId != "" {
		this.Ctx.Request.Header.Set("X-Company-Id", this.UserInfo.CompanyId)
		this.Ctx.ResponseWriter.Header().Set("X-Company-Id", this.UserInfo.CompanyId)
	}

	/*
	if this.UserInfo.Cached {
		this.Ctx.ResponseWriter.Header().Set("X-Userinfo-Cached", "True")
	} else {
		this.Ctx.ResponseWriter.Header().Set("X-Userinfo-Cached", "False")
	}
	*/

    // this.Ctx.Request.SetBasicAuth("E37CE8E5038A794B25FC", "wQCfoBTrT8PBFUobT4oQJRHEDaBkuv3wwZi4fVCb")
    if wsutil.IsWebSocketRequest(this.Ctx.Request) {
        remoteWs := &url.URL{
            Scheme: "ws://",
            Host: RancherEndpointHost,
        }
        proxyWs := wsutil.NewSingleHostReverseProxy(remoteWs)
        proxyWs.ServeHTTP(this.Ctx.ResponseWriter, this.Ctx.Request)
    } else {
        remoteHttp, err := url.Parse(`http://` + RancherEndpointHost)
        if err != nil {
            // panic(err)
			logs.Error(err)
        }
        proxyHttp := httputil.NewSingleHostReverseProxy(remoteHttp)
        proxyHttp.ServeHTTP(this.Ctx.ResponseWriter, this.Ctx.Request)
    }
}

func (this *AppController) Subscribe() {
	// fmt.Println("Subscribe")
	// remoteWs := &url.URL{
	// 	Scheme: "ws://",
	// 	Host: RancherEndpointHost,
	// }
	// proxyWs := wsutil.NewSingleHostReverseProxy(remoteWs)
	// proxyWs.ServeHTTP(this.Ctx.ResponseWriter, this.Ctx.Request)

	c, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 4096, 4096)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	subscribeMessage := make(chan []byte)

	filterUrl := `ws://` +  RancherEndpointHost + this.Ctx.Input.URI()
	wsHeader := this.Ctx.Request.Header
	wsHeader.Del("Sec-Websocket-Version")
	wsHeader.Del("Connection")
	wsHeader.Del("Sec-Websocket-Key")
	wsHeader.Del("Sec-Websocket-Extensions")
	wsHeader.Del("Upgrade")
	wsClient, _, err := websocket.DefaultDialer.Dial(filterUrl, wsHeader)
	if err != nil {
		log.Println("read:", err)
		return
	}

	defer c.Close()
	defer wsClient.Close()

	go func() {
		for {
			_, message, err := wsClient.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			// 过滤
			if true {
				var subscribe models.SubscribeResource
				if err := json.Unmarshal(message, &subscribe); err == nil {
					// if subscribe.Name == "ping" || subscribe.ResourceType == "containerEvent" || subscribe.ResourceType == "host" {
					if subscribe.Name == "ping" {
						subscribeMessage <- message
						log.Printf("recv ping --------:\n %s", message)
					} else if subscribe.Data.Resource.CompanyId == this.UserInfo.CompanyId {
						stringMsg := string(message)
						stringMsg = strings.Replace(stringMsg, RancherEndpointHost, fmt.Sprintf("%s:%d", this.Ctx.Input.Host(), this.Ctx.Input.Port()), -1)
						subscribeMessage <- []byte(stringMsg)
						log.Printf("recv change --------:\n %s", message)
					} else {
						log.Printf("interception change --------:\n %s", message)
					}
				} else {
					log.Println("read json:", err)
				}
			} else {
				subscribeMessage <- message
				log.Printf("recv --------:\n %s", message)
			}
		}

	}()

	for {
		select {
		case msg := <- subscribeMessage:
			err = c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
		// mt, message, err := c.ReadMessage()
		// if err != nil {
		// 	log.Println("read:", err)
		// 	break
		// }
		// log.Printf("recv: %s", message)
		// time.Sleep(time.Second)
	}
}

func (this *AppController) Finish() {

}
