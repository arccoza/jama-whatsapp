package main

import (
	"encoding/base64"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/mdp/qrterminal/v3"
	"os"
)

func qrToURI(qrChan chan string) string {
	tmp := <-qrChan
	png, _ := qrcode.Encode(tmp, qrcode.Medium, 256)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}

func qrToTerminal(qrChan chan string) {
	val := <-qrChan
	qrterminal.GenerateHalfBlock(val, qrterminal.L, os.Stdout)
}
