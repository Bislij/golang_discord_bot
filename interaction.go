package main

import (
	"discordBot/commandContext"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const Prefix string = "g;"

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		err := handleInteractionCommand(s, i)
		if err != nil {
			panic(err)
			return
		}
	}

	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}

	if m.Message.Content[0:len(Prefix)] == Prefix {
		err := handleMessageCommand(s, m)
		if err != nil {
			panic(err)
			return
		}
	}
}

func handleInteractionCommand(s *discordgo.Session, i *discordgo.InteractionCreate) (err error) {
	err = Commands[i.ApplicationCommandData().Name].Execute(&commandContext.InteractionContext{Session: s, InteractionCreate: i})
	return
}

func handleMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {
	rgx := regexp.MustCompile(`"[^"]+"|\S+`)
	args := rgx.FindAllString(m.Message.Content, -1)

	for i, match := range args {
		args[i] = strings.Trim(match, `"`)
	}

	cmd := strings.TrimPrefix(args[0], Prefix)
	value, ok := Commands[cmd]

	if ok {
		if len(args) == 0 || len(args) != len(value.ApplicationCommand.Options)+1 {
			fmt.Printf("Non-fatal: incorrect amount of options in command %s\n", cmd)
			return
		}
		args = args[1 : len(value.ApplicationCommand.Options)+1]
		var options []*discordgo.ApplicationCommandInteractionDataOption
		for index, el := range value.ApplicationCommand.Options {
			switch el.Type {
			case discordgo.ApplicationCommandOptionSubCommand:
				panic("unimplemented")
			case discordgo.ApplicationCommandOptionSubCommandGroup:
				panic("unimplemented")
			case discordgo.ApplicationCommandOptionString:
				options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
					Name:    value.ApplicationCommand.Options[index].Name,
					Type:    el.Type,
					Value:   args[index],
					Options: options,
					Focused: false,
				})
			case discordgo.ApplicationCommandOptionInteger:
				parseInt, e := strconv.ParseInt(args[index], 10, 64)
				if e != nil {
					fmt.Printf("Non-fatal: could not parse value at index %d into integer: %e\n", index, e)
					return
				} else {
					options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
						Name:    value.ApplicationCommand.Options[index].Name,
						Type:    el.Type,
						Value:   parseInt,
						Options: options,
						Focused: false,
					})
				}
			case discordgo.ApplicationCommandOptionBoolean:
				parseBool, e := strconv.ParseBool(args[index])
				if e != nil {
					fmt.Printf("Non-fatal: could not parse value at index %d into boolean: %e\n", index, e)
					return
				} else {
					options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
						Name:    value.ApplicationCommand.Options[index].Name,
						Type:    el.Type,
						Value:   parseBool,
						Options: options,
						Focused: false,
					})
				}
			case discordgo.ApplicationCommandOptionUser:
				panic("unimplemented")
			case discordgo.ApplicationCommandOptionChannel:
				panic("unimplemented")
			case discordgo.ApplicationCommandOptionRole:
				panic("unimplemented")
			case discordgo.ApplicationCommandOptionMentionable:
				panic("unimplemented")
			case discordgo.ApplicationCommandOptionNumber:
				panic("unimplemented")
			case discordgo.ApplicationCommandOptionAttachment:
				panic("unimplemented")
			}
		}
		err = value.Execute(&commandContext.MessageContext{Session: s, MessageCreate: m, OptionsList: options})
	} else {
		log.Printf("Unrecognized command: %s", cmd)
	}

	return
}
