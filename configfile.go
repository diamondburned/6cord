package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stevenroose/gonfig"
	fast "github.com/valyala/fasttemplate"
	"gitlab.com/diamondburned/6cord/md"
)

var cfg Config

// Properties ..
type Properties struct {
	CompactMode                bool   `id:"compact-mode"     default:"true"  desc:"Compact Mode"`
	TrueColor                  bool   `id:"true-color"       default:"true"  desc:"Enable True color mode instead of 256 color mode"`
	DefaultNameColor           string `id:"default-name-color" default:"#CCCCCC" desc:"Sets the default name color, format: #XXXXXX"`
	MentionColor               string `id:"mention-color"      default:"#0D4A91" desc:"The default mention background color"`
	MentionSelfColor           string `id:"mention-self-color" default:"#17AC86" desc:"The default mention background color, when the target is you"`
	ShowChannelsOnStartup      bool   `id:"show-channels"    default:"true"  desc:"Show the left channel bar on startup."`
	ChatPadding                int    `id:"chat-padding"     default:"2"     desc:"Determine the default indentation of messages from the left side."`
	SidebarRatio               int    `id:"sidebar-ratio"    default:"3"     desc:"The sidebar width in ratio of 1:N, whereas N is the ratio for the message box. The higher the number is, the narrower the sidebar is."`
	SidebarIndent              int    `id:"sidebar-indent"   default:"2"     desc:"Width in spaces each indentation level on the sidebar adds."`
	HideBlocked                bool   `id:"hide-blocked"     default:"true"  desc:"Ignore all blocked users."`
	TriggerTyping              bool   `id:"trigger-typing"   default:"true"  desc:"Send a TypingStart event periodically to the Discord server, default behavior of clients."`
	ForegroundColor            int    `id:"foreground-color" default:"15"    desc:"Default foreground color, 0-255, 0 is black, 15 is white."`
	BackgroundColor            int    `id:"background-color" default:"-1"    desc:"Acceptable values: tcell.Color*, -1, 0-255 (terminal colors)."`
	AuthorFormat               string `id:"author-format"    default:"[#{color}::b]{name}" desc:"The formatting of message authors"`
	CommandPrefix              string `id:"command-prefix"   default:"[${GUILD}${CHANNEL}] " desc:"The prefix of the input box"`
	DefaultStatus              string `id:"default-status"   default:"Send a message or input a command" desc:"The message in the status bar."`
	SyntaxHighlightColorscheme string `id:"syntax-highlight-colorscheme" default:"emacs" desc:"The color scheme for syntax highlighting, refer to https://xyproto.github.io/splash/docs/all.html."`
	ShowEmojiURLs              bool   `id:"show-emoji-urls"  default:"true"  desc:"Converts emojis into clickable URLs."`
	ObfuscateWords             bool   `id:"obfuscate-words"  default:"false" desc:"Insert a zero-width space to obfuscate word-logging telemetry."`
	ChatMaxWidth               int    `id:"chat-max-width"   default:"0"     desc:"The maximum width of the chat box, if smaller, will be centered."`
	ImageFetchTimeout          int    `id:"image-fetch-timeout" default:"1"  desc:"The timeout to fetch images, in seconds."`
	ImageWidth                 int    `id:"image-width"      default:"400"   desc:"The maximum width for an image."`
	ImageHeight                int    `id:"image-height"     default:"400"   desc:"The maximum height for an image."`
	ShortenURL                 bool   `id:"shorten-url"      default:"true"  desc:"Opens a webserver to redirect URLs"`
	RPCServer                  bool   `id:"rpc-server"       default:"true"  desc:"Start a Rich Presence server for applications to use. Experimental. Source: https://gitlab.com/diamondburned/drpc-server"`
}

type Config struct {
	Username string `id:"username" short:"u" default:"" desc:"Used when token is empty, avoid if 2FA"`
	Password string `id:"password" short:"p" default:"" desc:"Used when token is empty"`
	Token    string `id:"token" short:"t" default:"" desc:"Authentication Token, recommended way of using"`

	Login bool `id:"login" short:"l" default:"false" desc:"Force pop up a login prompt"`

	Prop Properties `id:"properties"`

	Debug bool `id:"debug" short:"d" default:"false" desc:"Enables debug mode"`

	Config string `short:"c"`
}

var (
	authorRawFormat  string
	authorPrefix     string
	messageRawFormat string

	// color, name, time
	authorTmpl *fast.Template

	// ID, content
	messageTmpl *fast.Template

	chatPadding string

	defaultNameColor int
)

func loadCfg() error {
	// Get the XDG paths
	var xdg = os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		if h, err := os.UserHomeDir(); err == nil {
			xdg = filepath.Join(h, ".config")
		}
	}

	err := gonfig.Load(&cfg, gonfig.Conf{
		ConfigFileVariable:  "config",
		FileDefaultFilename: filepath.Join(xdg, "6cord", "6cord.toml"),
		FileDecoder:         gonfig.DecoderTOML,
		EnvPrefix:           "sixcord_",
	})

	if err != nil {
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

	if cfg.Prop.TrueColor {
		term := os.Getenv("TERM")
		if strings.Contains(term, "-") {
			os.Setenv("TERM", strings.Split(term, "-")[0]+"-truecolor")
		}
	}

	hex := cfg.Prop.DefaultNameColor
	if len(hex) < 6 {
		return errors.New("Invalid format for name color, refer to help")
	}

	if hex[0] == '#' {
		hex = hex[1:]
	}

	hex64, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return err
	}

	defaultNameColor = int(hex64)
	messageRawFormat = `["{ID}"][-]{content}[-::-]["ENDMESSAGE"]`

	if cfg.Prop.CompactMode {
		messageRawFormat = " " + messageRawFormat
		authorPrefix = "\n[\"author\"]"
		authorRawFormat = authorPrefix + cfg.Prop.AuthorFormat + `[-:-:-][""]`
	} else {
		messageRawFormat = "\n" + messageRawFormat
		authorPrefix = "\n\n[\"author\"]"
		authorRawFormat = authorPrefix + cfg.Prop.AuthorFormat + ` [-:-:-][::d]{time}[::-][""]`
	}

	authorTmpl = fast.New(authorRawFormat, "{", "}")
	messageTmpl = fast.New(messageRawFormat, "{", "}")

	chatPadding = strings.Repeat(" ", cfg.Prop.ChatPadding)

	showChannels = cfg.Prop.ShowChannelsOnStartup
	md.HighlightStyle = cfg.Prop.SyntaxHighlightColorscheme

	return nil
}
