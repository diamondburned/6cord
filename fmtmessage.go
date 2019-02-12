package main

import (
	"fmt"
	"strings"

	"github.com/RumbleFrog/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/eidolon/wordwrap"
	"github.com/rivo/tview"
	"gitlab.com/diamondburned/6cord/md"
)

func fmtMessage(m *discordgo.Message) string {
	spew.Dump()
	var (
		ct = md.Parse(
			m.ContentWithMentionsReplaced(),
		)

		edited string
		c      []string
		l      = strings.Split(ct, "\n")
	)

	if m.Content == "¯\\_(ツ)_/¯" {
		ct = "¯\\_(ツ)_/¯"
	}

	if ct != "" {
		for i := 0; i < len(l); i++ {
			c = append(c, "\t"+l[i])
		}
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
			embed = append(
				embed,
				splitEmbedLine(e.Description)...,
			)
		}

		if len(e.Fields) > 0 {
			embed = append(embed, "")

			for _, f := range e.Fields {
				embed = append(embed, splitEmbedLine(f.Name, " [::b]")...)
				embed = append(embed, splitEmbedLine(f.Value, " [::d]")...)
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
		c = append(c, "\t"+tview.Escape(
			fmt.Sprintf("[%s]: %s", a.Filename, a.URL),
		))
	}

	if m.EditedTimestamp != "" {
		edited = " [::d](edited)[::-]"
	}

	return strings.Join(c, "\n") + edited
}

// WordWrapper makes a global wrapper for embed use
var WordWrapper = wordwrap.Wrapper(EmbedColLimit, false)

// 2nd arg ::-
// 3rd arg -::
func splitEmbedLine(e string, customMarkup ...string) (spl []string) {
	lines := strings.Split(e, "\n")

	// Todo: clean this up ETA never

	var (
		cm = ""
		ce = ""
	)

	if len(customMarkup) > 0 {
		cm = customMarkup[0]
		ce = "[::-]"
	}

	if len(customMarkup) > 1 {
		cm += customMarkup[1]
		ce += "[-::]"
	}

	for _, l := range lines {
		splwrap := strings.Split(md.Parse(l), "\n")

		for _, s := range splwrap {
			spl = append(spl, cm+s+ce)
		}
	}

	return
}
