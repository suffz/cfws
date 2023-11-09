Example Usage:

```go
func main() {
	conn := (&WebsocketOptions{
		URL:        "wss://ws.bloxflip.com/socket.io/?EIO=3&transport=websocket",
		ServerName: "ws.bloxflip.com", PORT: "443",
		Origin:    "https://bloxflip.com",
		Host:      "ws.bloxflip.com",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203",
	}).Dial()
	for {
		_, msg, err := conn.Conn.ReadMessage()
		fmt.Println(string(msg), err)
	}
}
```
