package command

import (
	"discordBot/commandContext"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var Ping = Command{
	ApplicationCommand: discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with Pong!",
	},
	Execute: func(ctx commandContext.CommandContext) error {
		err := ctx.ReplyEmbed(&discordgo.MessageEmbed{
			Title:       "Latency Time",
			Description: strconv.FormatInt(ctx.Latency().Milliseconds(), 10) + "ms",
			Color:       0,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "pong you bumbling pillock",
			},
		})
		return err
	},
}
