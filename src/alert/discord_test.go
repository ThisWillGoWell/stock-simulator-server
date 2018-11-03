package alert

import (
	"os"
	"testing"
	"time"
)

func TestDiscordWriter(t *testing.T) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		panic("no token set")
	}
	dw := Init(token, "504397270075179029")
	_, err := dw.Write([]byte("Hello"))
	if err != nil {
		//panic(err)
	}
	<-time.After(time.Second * 30)
}
