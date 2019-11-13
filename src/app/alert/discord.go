package alert

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

type DiscordWriter struct {
	session   *discordgo.Session
	ChannelId string
}

func Init(token, channel string) io.Writer {
	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		panic(err)
	}
	ap := &discordgo.Application{}
	ap.Name = "Alerts"
	dg.ApplicationCreate(ap)
	return &DiscordWriter{
		session:   dg,
		ChannelId: channel,
	}
}

func (discordWriter *DiscordWriter) Write(p []byte) (int, error) {
	_, err := discordWriter.session.ChannelMessageSend(discordWriter.ChannelId, string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
