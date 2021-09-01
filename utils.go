package main

import (
	"encoding/base64"
	qrcode "github.com/skip2/go-qrcode"
)

func qrToURI(qrChan chan string) string {
	tmp := <-qrChan
	png, _ := qrcode.Encode(tmp, qrcode.Medium, 256)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}
