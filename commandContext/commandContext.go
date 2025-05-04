package commandContext

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type CommandContext interface {
	ChannelID() string
	AuthorID() string
	Content() string

	Latency() time.Duration

	Options() []*discordgo.ApplicationCommandInteractionDataOption

	Send(content string) error
	Reply(content string) error

	SendEmbed(embed *discordgo.MessageEmbed) error
	ReplyEmbed(embed *discordgo.MessageEmbed) error
}

type MessageContext struct {
	Session       *discordgo.Session
	MessageCreate *discordgo.MessageCreate
	OptionsList   []*discordgo.ApplicationCommandInteractionDataOption
}

func (c *MessageContext) ChannelID() string {
	return c.MessageCreate.ChannelID
}

func (c *MessageContext) AuthorID() string {
	return c.MessageCreate.Author.ID
}

func (c *MessageContext) Content() string {
	return c.MessageCreate.Content
}

func (c *MessageContext) Latency() time.Duration {
	return c.Session.HeartbeatLatency()
}

func (c *MessageContext) Options() []*discordgo.ApplicationCommandInteractionDataOption {
	return c.OptionsList
}

func (c *MessageContext) Send(content string) error {
	_, err := c.Session.ChannelMessageSend(c.ChannelID(), content)
	return err
}

func (c *MessageContext) Reply(content string) error {
	_, err := c.Session.ChannelMessageSendReply(c.ChannelID(), content, &discordgo.MessageReference{
		MessageID:       c.MessageCreate.ID,
		ChannelID:       c.MessageCreate.ChannelID,
		GuildID:         c.MessageCreate.GuildID,
		FailIfNotExists: nil,
	})
	return err
}

func (c *MessageContext) SendEmbed(embed *discordgo.MessageEmbed) error {
	_, err := c.Session.ChannelMessageSendEmbed(c.ChannelID(), embed)
	return err
}

func (c *MessageContext) ReplyEmbed(embed *discordgo.MessageEmbed) error {
	_, err := c.Session.ChannelMessageSendEmbedReply(c.ChannelID(), embed, &discordgo.MessageReference{
		MessageID:       c.MessageCreate.ID,
		ChannelID:       c.MessageCreate.ChannelID,
		GuildID:         c.MessageCreate.GuildID,
		FailIfNotExists: nil,
	})
	return err
}

type InteractionContext struct {
	Session           *discordgo.Session
	InteractionCreate *discordgo.InteractionCreate
}

func (c *InteractionContext) ChannelID() string {
	return c.InteractionCreate.ChannelID
}

func (c *InteractionContext) AuthorID() string {
	return c.InteractionCreate.Member.User.ID
}

func (c *InteractionContext) Content() string {
	return c.InteractionCreate.ApplicationCommandData().Name
}

func (c *InteractionContext) Latency() time.Duration {
	return c.Session.HeartbeatLatency()
}

func (c *InteractionContext) Options() []*discordgo.ApplicationCommandInteractionDataOption {
	return c.InteractionCreate.ApplicationCommandData().Options
}

func (c *InteractionContext) Send(content string) error {
	return c.Session.InteractionRespond(c.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func (c *InteractionContext) Reply(content string) error {
	return c.Send(content)
}

func (c *InteractionContext) SendEmbed(embed *discordgo.MessageEmbed) error {
	return c.Session.InteractionRespond(c.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func (c *InteractionContext) ReplyEmbed(embed *discordgo.MessageEmbed) error {
	return c.SendEmbed(embed)
}
