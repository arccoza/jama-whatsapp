package main

import (
	"fmt"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"time"
)

func main() {
	wac, err1 := whatsapp.NewConn(20 * time.Second)
	if err1 != nil {
		fmt.Println("err1")
		fmt.Println(err1)
		return
	}

	qrChan := make(chan string)
	go func() {
		fmt.Printf("qr code: %v\n", qrToURI(qrChan))
	}()

	wac.SetClientName("And1", "and1", "2,2121,6")

	sess, err2 := wac.Login(qrChan)
	if err2 != nil {
		fmt.Println("err2")
		fmt.Println(err2)
		return
	}
	fmt.Println(sess)
}
