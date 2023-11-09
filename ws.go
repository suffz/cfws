package cfws

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	tls "github.com/bogdanfinn/utls"
	gorilla "github.com/gorilla/websocket"
)

/*
func main() {
	conn := (&WebsocketOptions{
		URL:        "wss://ws.bloxflip.com/socket.io/?EIO=3&transport=websocket",
		ServerName: "ws.bloxflip.com", PORT: "443",
		Origin:    "https://bloxflip.com",
		Host:      "ws.bloxflip.com",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203",
		ReadSize:  MBs(5),
		WriteSize: MBs(5),
		// Extensions: "permessage-deflate; client_max_window_bits",
	}).Dial()
	for {
		_, msg, err := conn.Conn.ReadMessage()
		fmt.Println(string(msg), err)
	}
}
*/

type WebsocketOptions struct {
	URL                 string
	ServerName, PORT    string
	Origin              string
	Host                string
	Extensions          string
	UserAgent           string
	ReadSize, WriteSize int
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

	if conn, err := net.Dial("tcp", Info.ServerName+":"+strings.ReplaceAll(Info.PORT, ":", "")); err == nil {
		conn, resp, err := gorilla.NewClient(tls.UClient(conn, &tls.Config{
			ServerName: Info.ServerName,
		}, tls.HelloChrome_112, true, true).NetConn(), i, map[string][]string{
			"Origin":                   {Info.Origin},
			"Host":                     {Info.Host},
			"User-Agent":               {Info.UserAgent},
			"Sec-WebSocket-Extensions": {Info.Extensions},
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
