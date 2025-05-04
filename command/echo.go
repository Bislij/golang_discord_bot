package command

import (
	"discordBot/commandContext"
	"github.com/bwmarrin/discordgo"
)

var Echo = Command{
	ApplicationCommand: discordgo.ApplicationCommand{
		Name:        "echo",
		Description: "Echos the user's message content.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "text",
				Description: "The message to be echoed.",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
	Execute: func(ctx commandContext.CommandContext) error {
		err := ctx.Reply(ctx.Options()[0].StringValue())
		return err
	},
}
