package main

import (
	"encoding/base64"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/mdp/qrterminal/v3"
	"os"
)

func qrToURI(val string) string {
	png, _ := qrcode.Encode(val, qrcode.Medium, 256)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}

func qrToTerminal(val string) {
	qrterminal.GenerateHalfBlock(val, qrterminal.L, os.Stdout)
}
