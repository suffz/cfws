package cfws

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
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
	Proxy               string
}

// If you use the MBs() func just know it scales up to Megabits, 5 = 5MB.
func MBs(i int) int {
	return i * 1024 * 1024
}

type WebsocketConnection struct {
	Conn      *gorilla.Conn
	Resp      *http.Response
	Err       error
	ProxyInfo ProxyData
}

type ProxyData struct {
	ip, port, user, pass string
}

func (Info *WebsocketOptions) Dial() WebsocketConnection {

	i, err := url.Parse(Info.URL)
	if err != nil {
		return WebsocketConnection{
			Err: err,
		}
	}

	var conn net.Conn

	switch true {
	case Info.Proxy != "":
		if strings.Contains(Info.Proxy, "http") {
			return WebsocketConnection{Err: errors.New("Proxy: invalid format | use > ip:port:user:pass OR ip:port")}
		}
		conn, err, _, _ = Info.Connect()
	case Info.KeepAlive:
		conn, err = (&net.Dialer{KeepAlive: time.Hour * 999999}).Dial("tcp", Info.ServerName+":"+strings.ReplaceAll(Info.PORT, ":", ""))
	default:
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
			Conn:      conn,
			Resp:      resp,
			Err:       err,
			ProxyInfo: GetProxyStrings(Info.Proxy),
		}
	} else {
		return WebsocketConnection{
			Err: err,
		}
	}
}

func (Info *WebsocketOptions) Connect() (net.Conn, error, bool, string) {

	proxy := Info.Proxy
	ip := strings.Split(proxy, ":")
	if conn, err := net.Dial("tcp", ip[0]+":"+ip[1]); err == nil {
		if len(ip) > 2 {
			conn.Write([]byte(fmt.Sprintf("CONNECT %v:%v HTTP/1.1\r\nHost: %v:%v\r\nProxy-Authorization: Basic %v\r\nProxy-Connection: keep-alive\r\nUser-Agent: %v\r\n\r\n", Info.ServerName, Info.PORT, Info.ServerName, Info.PORT, base64.RawStdEncoding.EncodeToString([]byte(ip[2]+":"+ip[3])), Info.UserAgent)))
		} else {
			conn.Write([]byte(fmt.Sprintf("CONNECT %v:%v HTTP/1.1\r\nHost: %v:%v\r\nProxy-Connection: keep-alive\r\nUser-Agent: %v\r\n\r\n", Info.ServerName, Info.PORT, Info.ServerName, Info.PORT, Info.UserAgent)))
		}
		var junk = make([]byte, 20)
		conn.Read(junk)
		switch Status := string(junk); Status[9:12] {
		case "200":
			return conn, nil, true, ip[0]
		case "407":
			return nil, errors.New(fmt.Sprintf("[%v] Proxy <%v> Failed to authorize: Username/Password invalid.", Status[9:12], ip[0])), false, ip[0]
		default:
			return nil, errors.New("Unknown status code " + Status + " Returned.. body length " + strconv.Itoa(len(junk)) + " data: " + string(junk)), false, ip[0]
		}
	} else {
		return nil, err, false, ip[0]
	}
}

func GetProxyStrings(proxy string) ProxyData {
	var ip, port, user, pass string
	switch data := strings.Split(proxy, ":"); len(data) {
	case 2:
		ip = data[0]
		port = data[1]
	case 4:
		ip = data[0]
		port = data[1]
		user = data[2]
		pass = data[3]
	}
	return ProxyData{ip: ip, port: port, user: user, pass: pass}
}
