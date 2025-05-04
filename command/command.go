package command

import (
	"discordBot/commandContext"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	ApplicationCommand discordgo.ApplicationCommand
	Execute            func(commandContext.CommandContext) error
}
