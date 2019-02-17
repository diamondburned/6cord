package main

import (
	"fmt"
	"strings"

	"github.com/RumbleFrog/discordgo"
	"github.com/rivo/tview"
)

func fmtMessage(m *discordgo.Message) string {
	ct, emojiMap := ParseAll(m)

	if m.EditedTimestamp != "" {
		ct += " [::d](edited)[::-]"

		// Prevent cases where the message is empty
		// " (edited)"
		ct = strings.TrimPrefix(ct, " ")
	}

	var (
		c []string
		l = strings.Split(ct, "\n")
	)

	if m.Content == "¯\\_(ツ)_/¯" {
		ct = "¯\\_(ツ)_/¯"
	}

	if ct != "" {
		for i := 0; i < len(l); i++ {
			c = append(c, "\t"+l[i])
		}
	}

	for _, arr := range emojiMap {
		m.Attachments = append(
			m.Attachments,
			&discordgo.MessageAttachment{
				Filename: arr[0],
				URL:      arr[1],
			},
		)
	}

	for _, e := range m.Embeds {
		var embed = []string{""}

		if e.URL != "" {
			m.Attachments = append(
				m.Attachments,
				&discordgo.MessageAttachment{
					Filename: "EmbedURL",
					URL:      e.URL,
				},
			)
		}

		if e.Author != nil {
			embed = append(
				embed,
				"[::u]"+e.Author.Name+"[::-]",
			)

			if e.Author.IconURL != "" {
				m.Attachments = append(
					m.Attachments,
					&discordgo.MessageAttachment{
						Filename: "AuthorIcon",
						URL:      e.Author.IconURL,
					},
				)
			}

			if e.Author.URL != "" {
				m.Attachments = append(
					m.Attachments,
					&discordgo.MessageAttachment{
						Filename: "AuthorURL",
						URL:      e.Author.URL,
					},
				)
			}
		}

		if e.Title != "" {
			embed = append(
				embed,

				/*
					Sure, there's a bug here, but
					it'll rarely happen anyway lul

					if you don't know what it is,
					if L1 > 45 chars, it will line
					break and L2 will have 50 chars,
					which that looks inconsistent
				*/
				splitEmbedLine(e.Title, "[::b]", "[#0096cf]")...,
			)
		}

		if e.Description != "" {
			var desc, emojis = parseEmojis(e.Description)

			embed = append(
				embed,
				splitEmbedLine(desc)...,
			)

			for _, arr := range emojis {
				m.Attachments = append(
					m.Attachments,
					&discordgo.MessageAttachment{
						Filename: arr[0],
						URL:      arr[1],
					},
				)
			}
		}

		if len(e.Fields) > 0 {
			embed = append(embed, "")

			for _, f := range e.Fields {
				embed = append(embed,
					splitEmbedLine(f.Name, " [::b]")...)
				embed = append(embed,
					splitEmbedLine(f.Value, " [::d]")...)
				embed = append(embed, "")
			}
		}

		var footer []string
		if e.Footer != nil {
			footer = append(
				footer,
				"[::d]"+tview.Escape(e.Footer.Text)+"[::-]",
			)

			if e.Footer.IconURL != "" {
				m.Attachments = append(
					m.Attachments,
					&discordgo.MessageAttachment{
						Filename: "FooterIcon",
						URL:      e.Footer.IconURL,
					},
				)
			}
		}

		if e.Timestamp != "" {
			footer = append(
				footer,
				"[::d]"+e.Timestamp+"[::-]",
			)
		}

		if len(footer) > 0 {
			embed = append(
				embed,
				strings.Join(footer, " - "),
			)
		}

		if e.Thumbnail != nil {
			m.Attachments = append(
				m.Attachments,
				&discordgo.MessageAttachment{
					Filename: "Thumbnail",
					URL:      e.Thumbnail.URL,
				},
			)
		}

		if e.Image != nil {
			m.Attachments = append(
				m.Attachments,
				&discordgo.MessageAttachment{
					Filename: "Image",
					URL:      e.Image.URL,
				},
			)
		}

		if e.Video != nil {
			m.Attachments = append(
				m.Attachments,
				&discordgo.MessageAttachment{
					Filename: "Video",
					URL:      e.Video.URL,
				},
			)
		}

		c = append(
			c, strings.Join(
				embed, fmt.Sprintf("\n [#%06X]┃[-::] ", e.Color),
			), "",
		)
	}

	for _, a := range m.Attachments {
		c = append(
			c,
			"\t[::d]"+tview.Escape(fmt.Sprintf("[%s]: %s", a.Filename, a.URL))+"[::-]",
		)
	}

	return strings.Join(c, "\n")
}
