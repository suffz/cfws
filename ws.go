package cfws

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	tls "github.com/bogdanfinn/utls"
	gorilla "github.com/gorilla/websocket"
)

type WebsocketOptions struct {
	URL                 string
	ServerName, PORT    string
	Origin              string
	Host                string
	Extensions          string
	UserAgent           string
	CF_Clearance        string
	ReadSize, WriteSize int
	KeepAlive           bool
}

// If you use the MBs() func just know it scales up to Megabits, 5 = 5MB.
func MBs(i int) int {
	return i * 1024 * 1024
}

type WebsocketConnection struct {
	Conn *gorilla.Conn
	Resp *http.Response
	Err  error
}

func (Info *WebsocketOptions) Dial() WebsocketConnection {

	i, err := url.Parse(Info.URL)
	if err != nil {
		return WebsocketConnection{
			Err: err,
		}
	}

	var conn net.Conn
	if Info.KeepAlive {
		conn, err = (&net.Dialer{KeepAlive: time.Hour * 999999}).Dial("tcp", Info.ServerName+":"+strings.ReplaceAll(Info.PORT, ":", ""))
	} else {
		conn, err = net.Dial("tcp", Info.ServerName+":"+strings.ReplaceAll(Info.PORT, ":", ""))
	}

	if err == nil {
		conn, resp, err := gorilla.NewClient(tls.UClient(conn, &tls.Config{
			ServerName: Info.ServerName,
		}, tls.HelloChrome_120, true, true).NetConn(), i, map[string][]string{
			"Origin":                   {Info.Origin},
			"Host":                     {Info.Host},
			"User-Agent":               {Info.UserAgent},
			"Sec-WebSocket-Extensions": {Info.Extensions},
			"Cookie":                   {"cf_clearance=" + Info.CF_Clearance},
		}, Info.ReadSize, Info.WriteSize) // 5mb of allocated storage.
		return WebsocketConnection{
			Conn: conn,
			Resp: resp,
			Err:  err,
		}
	} else {
		return WebsocketConnection{
			Err: err,
		}
	}
}
