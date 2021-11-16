package main

import (
	"os"
	"fmt"
	"context"
	"strings"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/session"
	// "github.com/gotd/td/telegram/auth/qrlogin"
)

var ID, _ = strconv.Atoi(os.Getenv("TELEGRAM_API_ID"))
var HASH = os.Getenv("TELEGRAM_API_HASH")

func initTelegram() {
	ctx := context.Background()
	d := tg.NewUpdateDispatcher()
	// storage := &TelegramIntegration{}
	storage := &session.FileStorage{Path: "./session.tg"}

	client := telegram.NewClient(ID, HASH, telegram.Options{
		// DC: 2,
		UpdateHandler: d,
		SessionStorage: storage,
	})

	if err := client.Run(ctx, func(ctx context.Context) error {
		Login(ctx, client)

		// Keep client alive until context done
		<-ctx.Done()

		return nil
	}); err != nil {
		panic(err)
	}
}

func Login(ctx context.Context, client *telegram.Client) error {
	flow := auth.NewFlow(
		TelegramUserAuth{},
		auth.SendCodeOptions{},
	)

	if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		return err
	}

	auth, err := client.Auth().Status(ctx)
	fmt.Println(auth, err)

	return nil
}


type noSignUp struct{}

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}


type TelegramUserAuth struct {
	noSignUp
	phone string
}

func (a TelegramUserAuth) Phone(_ context.Context) (string, error) {
	var phone string
	fmt.Println("Phone: ")
	fmt.Scanln(&phone)
	return phone, nil
}

func (a TelegramUserAuth) Password(_ context.Context) (string, error) {
	var password string
	fmt.Println("Password: ")
	fmt.Scanln(&password)
	return strings.TrimSpace(password), nil
}

func (a TelegramUserAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	var code string
	fmt.Println("Code: ")
	fmt.Scanln(&code)
	return strings.TrimSpace(code), nil
}


// func main() {
// 	ctx := context.Background()
// 	d := tg.NewUpdateDispatcher()

// 	loggedIn := make(chan struct{})

// 	d.OnLoginToken(func(ctx context.Context, e tg.Entities, update *tg.UpdateLoginToken) error {
// 		fmt.Println("OnLoginToken")
// 		fmt.Printf("%+v\n", e)
// 		fmt.Printf("%+v\n", update)
// 		loggedIn <- struct{}{}
// 		return nil
// 	})

// 	// loggedIn := qrlogin.OnLoginToken(d)

// 	// https://core.telegram.org/api/obtaining_api_id
// 	client := telegram.NewClient(ID, HASH, telegram.Options{
// 		// DC: 2,
// 		UpdateHandler: d,
// 	})
// 	if err := client.Run(ctx, func(ctx context.Context) error {
// 		qr := client.QR()
// 		fmt.Printf("%+v\n", qr)
// 		authz, err := qr.Auth(
// 			ctx,
// 			loggedIn,
// 			func(ctx context.Context, token qrlogin.Token) error {
// 				fmt.Println("token:", token.URL())
// 				fmt.Println(token.Expires().Unix())
// 				qrToTerminal(token.URL())
// 				return nil
// 			},
// 		)

// 		fmt.Printf("%+v\n", authz)
// 		fmt.Printf("%+v\n", err)

// 		if err != nil {
// 			return err
// 		}

// 		u, ok := authz.User.AsNotEmpty()
// 		if !ok {
// 			return fmt.Errorf("unexpected type %T", authz.User)
// 		}
// 		fmt.Println("ID:", u.ID, "Username:", u.Username, "Bot:", u.Bot)

// 		return nil
// 	}); err != nil {
// 		panic(err)
// 	}
// 	// Client is closed.
// }




// import (
// 	"bufio"
// 	"context"
// 	// "flag"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"strconv"

// 	"github.com/go-faster/errors"
// 	"go.uber.org/zap"
// 	"golang.org/x/crypto/ssh/terminal"

// 	"github.com/gotd/td/examples"
// 	"github.com/gotd/td/telegram"
// 	"github.com/gotd/td/telegram/auth"
// 	"github.com/gotd/td/tg"
// )

// var ID, _ = strconv.Atoi(os.Getenv("TELEGRAM_API_ID"))
// var HASH = os.Getenv("TELEGRAM_API_HASH")

// // noSignUp can be embedded to prevent signing up.
// type noSignUp struct{}

// func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
// 	return auth.UserInfo{}, errors.New("not implemented")
// }

// func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
// 	return &auth.SignUpRequired{TermsOfService: tos}
// }

// // termAuth implements authentication via terminal.
// type termAuth struct {
// 	noSignUp

// 	phone string
// }

// func (a termAuth) Phone(_ context.Context) (string, error) {
// 	return a.phone, nil
// }

// func (a termAuth) Password(_ context.Context) (string, error) {
// 	fmt.Print("Enter 2FA password: ")
// 	bytePwd, err := terminal.ReadPassword(0)
// 	if err != nil {
// 		return "", err
// 	}
// 	return strings.TrimSpace(string(bytePwd)), nil
// }

// func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
// 	fmt.Print("Enter code: ")
// 	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
// 	if err != nil {
// 		return "", err
// 	}
// 	return strings.TrimSpace(code), nil
// }

// func main() {
// 	phone := "6588900241"

// 	examples.Run(func(ctx context.Context, log *zap.Logger) error {
// 		// Setting up authentication flow helper based on terminal auth.
// 		flow := auth.NewFlow(
// 			termAuth{phone: phone},
// 			auth.SendCodeOptions{},
// 		)

// 		client := telegram.NewClient(ID, HASH, telegram.Options{
// 			Logger: log,
// 		})
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		return client.Run(ctx, func(ctx context.Context) error {
// 			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
// 				return err
// 			}

// 			log.Info("Success")

// 			return nil
// 		})
// 	})
// }

type TelegramConnector struct {
	integ *Integration
	conn *telegram.Client
	subscribers map[*Handler]Handler
}

func NewTelegramConnector(integ *Integration) *TelegramConnector {
	c := &TelegramConnector{
		integ: integ,
		subscribers: map[*Handler]Handler{},
	}

	return c
}

func (c *TelegramConnector) Start() error {
	conn, err := initTelegram(c.integ, &tgHandler{
		notify: func(pay Payload) {
			c.notify(pay)
		},
	})

	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *TelegramConnector) Publish(pay Payload) {
	for _, msg := range pay.Messages {
		// c.conn.Send(msg.toTelegram())
	}
}

func (c *TelegramConnector) Subscribe(fn Handler) {
	c.subscribers[&fn] = fn
}

func (c *TelegramConnector) Unsubscribe(fn Handler) {

}

func (c *TelegramConnector) notify(pay Payload) {
	for _, fn := range c.subscribers {
		fn(pay)
	}
}

func (c *TelegramConnector) Query(q string) []Payload {
	return nil
}

type tgHandler struct {
	tg.UpdateDispatcher
	conn *whatsapp.Conn
	integ *Integration
	notify func(pay Payload)
}
