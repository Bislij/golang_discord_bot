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

const Prefix = "g;"

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if err := handleInteractionCommand(s, i); err != nil {
			panic(err)
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}

	if m.Message.Content[0:len(Prefix)] == Prefix {
		if err := handleMessageCommand(s, m); err != nil {
			panic(err)
		}
	}
}

func handleInteractionCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return Commands[i.ApplicationCommandData().Name].Execute(&commandContext.InteractionContext{
		Session:           s,
		InteractionCreate: i,
	})
}

func handleMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	rgx := regexp.MustCompile(`"[^"]+"|\S+`)
	args := rgx.FindAllString(m.Message.Content, -1)
	for i, match := range args {
		args[i] = strings.Trim(match, `"`)
	}

	cmd := strings.TrimPrefix(args[0], Prefix)
	value, ok := Commands[cmd]

	if !ok {
		log.Printf("Unrecognized command: %s", cmd)
		return nil
	}

	args = args[1:]

	if len(args) != len(value.ApplicationCommand.Options) {
		fmt.Printf("Non-fatal: incorrect number of options for command '%s'\n", cmd)
		return nil
	}

	var options []*discordgo.ApplicationCommandInteractionDataOption
	for i, opt := range value.ApplicationCommand.Options {
		switch opt.Type {
		case discordgo.ApplicationCommandOptionString:
			options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
				Name:  opt.Name,
				Type:  opt.Type,
				Value: args[i],
			})
		case discordgo.ApplicationCommandOptionInteger:
			val, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				fmt.Printf("Non-fatal: could not parse value at index %d into integer: %e\n", i, err)
				return nil
			}
			options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
				Name:  opt.Name,
				Type:  opt.Type,
				Value: val,
			})
		case discordgo.ApplicationCommandOptionBoolean:
			val, err := strconv.ParseBool(args[i])
			if err != nil {
				fmt.Printf("Non-fatal: could not parse value at index %d into boolean: %e\n", i, err)
				return nil
			}
			options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
				Name:  opt.Name,
				Type:  opt.Type,
				Value: val,
			})
		case discordgo.ApplicationCommandOptionRole:
			guild, err := s.Guild(m.GuildID)
			if err != nil {
				fmt.Printf("Non-fatal: could not retrieve guild: %v\n", err)
				return nil
			}
			var val *discordgo.Role
			for _, role := range guild.Roles {
				if role.Name == args[i] {
					val = role
					break
				}
			}
			if val == nil {
				fmt.Printf("Non-fatal: could not find role named '%s'\n", args[i])
				return nil
			}
			options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
				Name:  opt.Name,
				Type:  opt.Type,
				Value: val,
			})
		default:
			fmt.Printf("Non-fatal: unimplemented option type: %d\n", opt.Type)
			return nil
		}
	}

	return value.Execute(&commandContext.MessageContext{
		Session:       s,
		MessageCreate: m,
		OptionsList:   options,
	})
}
