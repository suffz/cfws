Made with the help of https://github.com/suffz/http2

Example Usage:

```go
func main() {
	conn := (&cfws.WebsocketOptions{
		URL:        "wss://ws.bloxflip.com/socket.io/?EIO=3&transport=websocket",
		ServerName: "ws.bloxflip.com", PORT: "443",
		Origin:     "https://bloxflip.com",
		Host:       "ws.bloxflip.com",
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		ReadSize:   cfws.MBs(5),
		WriteSize:  cfws.MBs(5),
		Extensions: "client_max_window_bits",
		Headers: []cfws.Headers{
			{Name: "Pragma", Value: "no-cache"},
			{Name: "Cache-Control", Value: "no-cache"},
			{Name: "Accept-Language", Value: "en-US,en;q=0.9"},
			{Name: "Accept-Encoding", Value: "gzip, deflate, br, zstd"},
		},
		Cert: cfws.ReturnCertBasedOnBytes([]byte(`
		-- GlobalSign Root R2, valid until Dec 15, 2021
		-----BEGIN CERTIFICATE-----
		MIIDujCCAqKgAwIBAgILBAAAAAABD4Ym5g0wDQYJKoZIhvcNAQEFBQAwTDEgMB4G
		A1UECxMXR2xvYmFsU2lnbiBSb290IENBIC0gUjIxEzARBgNVBAoTCkdsb2JhbFNp
		Z24xEzARBgNVBAMTCkdsb2JhbFNpZ24wHhcNMDYxMjE1MDgwMDAwWhcNMjExMjE1
		MDgwMDAwWjBMMSAwHgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMjETMBEG
		A1UEChMKR2xvYmFsU2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjCCASIwDQYJKoZI
		hvcNAQEBBQADggEPADCCAQoCggEBAKbPJA6+Lm8omUVCxKs+IVSbC9N/hHD6ErPL
		v4dfxn+G07IwXNb9rfF73OX4YJYJkhD10FPe+3t+c4isUoh7SqbKSaZeqKeMWhG8
		eoLrvozps6yWJQeXSpkqBy+0Hne/ig+1AnwblrjFuTosvNYSuetZfeLQBoZfXklq
		tTleiDTsvHgMCJiEbKjNS7SgfQx5TfC4LcshytVsW33hoCmEofnTlEnLJGKRILzd
		C9XZzPnqJworc5HGnRusyMvo4KD0L5CLTfuwNhv2GXqF4G3yYROIXJ/gkwpRl4pa
		zq+r1feqCapgvdzZX99yqWATXgAByUr6P6TqBwMhAo6CygPCm48CAwEAAaOBnDCB
		mTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUm+IH
		V2ccHsBqBt5ZtJot39wZhi4wNgYDVR0fBC8wLTAroCmgJ4YlaHR0cDovL2NybC5n
		bG9iYWxzaWduLm5ldC9yb290LXIyLmNybDAfBgNVHSMEGDAWgBSb4gdXZxwewGoG
		3lm0mi3f3BmGLjANBgkqhkiG9w0BAQUFAAOCAQEAmYFThxxol4aR7OBKuEQLq4Gs
		J0/WwbgcQ3izDJr86iw8bmEbTUsp9Z8FHSbBuOmDAGJFtqkIk7mpM0sYmsL4h4hO
		291xNBrBVNpGP+DTKqttVCL1OmLNIG+6KYnX3ZHu01yiPqFbQfXf5WRDLenVOavS
		ot+3i9DAgBkcRcAtjOj4LaR0VknFBbVPFd5uRHg5h6h+u/N5GJG79G+dwfCMNYxd
		AfvDbbnvRG15RjF+Cv6pgsH/76tuIMRQyV+dTZsXjAzlAcmgQWpzU/qlULRuJQ/7
		TBj0/VLZjmmx6BEP3ojY+x1J96relc8geMJgEtslQIxq/H5COEBkEveegeGTLg==
		-----END CERTIFICATE-----`)),
		InitMessages: []string{
			`40/chat,`,
			`40/cloud-games,`,
			fmt.Sprintf(`42/chat,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/blackjack,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/cups,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/jackpot,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/roulette,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/rouletteV2,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/crash,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/wallet,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/marketplace,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/case-battles,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/mod-queue,["auth","%v"]`, jwt),
			fmt.Sprintf(`42/feed,["auth","%v"]`, jwt),
		},
	}).Dial()
	for {
		_, msg, err := conn.Conn.ReadMessage()
		fmt.Println(string(msg), err)
	}
}
```
