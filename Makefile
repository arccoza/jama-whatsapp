run:
	set -a && . ./.env.local && set +a && go run main.go firebase.go bridge.go message.go chat.go payload.go contact.go jama.go whatsapp.go utils.go
