package main

import (
	"os"
	"path/filepath"

	"github.com/stevenroose/gonfig"
)

var cfg Config

// Properties ..
type Properties struct {
	ShowChannelsOnStartup      bool   `id:"show-channels"    default:"true"  desc:"Show the left channel bar on startup."`
	ChatPadding                int    `id:"chat-padding"     default:"2"     desc:"Determine the default indentation of messages from the left side."`
	SidebarRatio               int    `id:"sidebar-ratio"    default:"3"     desc:"The sidebar width in ratio of 1:N, whereas N is the ratio for the message box. The higher the number is, the narrower the sidebar is."`
	SidebarIndent              int    `id:"sidebar-indent"   default:"2"     desc:"Width in spaces each indentation level on the sidebar adds."`
	HideBlocked                bool   `id:"hide-blocked"     default:"true"  desc:"Ignore all blocked users."`
	TriggerTyping              bool   `id:"trigger-typing"   default:"true"  desc:"Send a TypingStart event periodically to the Discord server, default behavior of clients."`
	ForegroundColor            int    `id:"foreground-color" default:"15"    desc:"Default foreground color, 0-255, 0 is black, 15 is white."`
	BackgroundColor            int    `id:"background-color" default:"-1"    desc:"Acceptable values: tcell.Color*, -1, 0-255 (terminal colors)."`
	CommandPrefix              string `id:"command-prefix"   default:"[${GUILD}${CHANNEL}] " desc:"The prefix of the input box"`
	DefaultStatus              string `id:"default-status"   default:"Send a message or input a command" desc:"The message in the status bar."`
	SyntaxHighlightColorscheme string `id:"syntax-highlight-colorscheme" default:"emacs" desc:"The color scheme for syntax highlighting, refer to https://xyproto.github.io/splash/docs/all.html."`
	ShowEmojiURLs              bool   `id:"show-emoji-urls"  default:"true"  desc:"Converts emojis into clickable URLs."`
	ObfuscateWords             bool   `id:"obfuscate-words"  default:"false" desc:"Insert a zero-width space to obfuscate word-logging telemetry."`
	ImageFetchTimeout          int    `id:"image-fetch-timeout" default:"1"  desc:"The timeout to fetch images, in seconds."`
	ImageWidth                 int    `id:"image-width"      default:"400"   desc:"The maximum width for an image."`
	ImageHeight                int    `id:"image-height"     default:"400"   desc:"The maximum height for an image."`
	ShortenURL                 bool   `id:"shorten-url"      default:"true"  desc:"Opens a webserver to redirect URLs"`
}

type Config struct {
	Username string `id:"username" short:"u" default:"" desc:"Used when token is empty, avoid if 2FA"`
	Password string `id:"password" short:"p" default:"" desc:"Used when token is empty"`
	Token    string `id:"token" short:"t" default:"" desc:"Authentication Token, recommended way of using"`

	Prop Properties `id:"properties"`

	Debug bool `id:"debug" short:"d" default:"false" desc:"Enables debug mode"`

	Config string `short:"c"`
}

func loadCfg() error {
	// Get the XDG paths
	var xdg = os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		if h, err := os.UserHomeDir(); err == nil {
			xdg = filepath.Join(h, ".config")
		}
	}

	if err := gonfig.Load(&cfg, gonfig.Conf{
		ConfigFileVariable:  "config",
		FileDefaultFilename: filepath.Join(xdg, "6cord", "6cord.toml"),
		FileDecoder:         gonfig.DecoderTOML,
		EnvPrefix:           "sixcord_",
	}); err != nil {
		return err
	}

	if cfg.Config != "" {
		return gonfig.Load(&cfg, gonfig.Conf{
			ConfigFileVariable:  "config",
			FileDefaultFilename: cfg.Config,
			FileDecoder:         gonfig.DecoderTOML,
			EnvPrefix:           "sixcord_",
		})
	}

	return nil
}
