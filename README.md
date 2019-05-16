<h1 align="center">6cord</h1>

<p align="center">
	<a href="https://liberapay.com/diamondburned/donate">
		<img alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg">
	</a>
	<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/guildview.png" />
	<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/img.png" />
	<img src="https://u.cubeupload.com/diamondburned/MCz9fP.png" />
</p>

## Installation

### 1. [From CI (only when the tick-mark is green)](https://gitlab.com/diamondburned/6cord/builds/artifacts/master/file/6cord?job=compile)

### 2. `go get -u gitlab.com/diamondburned/6cord`

## Behaviors (possibly outdated)

- From input, hit arrow up to go to autocompletion. Arrow up again to go to the message box.
- In the message box
- Arrow up/down and Page up/down will be used for scrolling
- Any other key focuses back to input
- Tab to hide channels, focusing on input
- Tab again to show channels, focusing on the channel list
- To clear the keyring, feed `6cord` a new token with `-t`
- Plans:
	- Disable focus on the message view and use Alt+Arrows instead (considering)

## Stuff

- Sample config is in `6cord.toml`, use with `-c`
- Command history is cycled through with <kbd>Alt</kbd> + <kbd>Up</kbd>/<kbd>Down</kbd>/<kbd>j</kbd>/<kbd>k</kbd>
- To get the following colors, use the variable
- Monochrome: `TERM=xtermm`
- Terminal colors: `TERM=xterm-basic`
- Templating for `command-prefix`
	- This can be templated, for example: "`[${GUILD}@${CHANNEL}] `"
	- Avaiable variables are `CHANNEL`, `GUILD`, `USERNAME` and `DISCRIM`
- To install the config, move  `6cord.toml` to `$HOME/.config/6cord/`

## Todo

- [x] [Fix paste not working](https://github.com/rivo/tview/issues/133) (workaround: Ctrl + V)
	- [x] Better paste library with image support (Linux only)
- [x] Syntax highlighting, better markdown parsing
- [x] Message Delete and Edit
- [x] Full reaction support
- [x] Command history (refer to Plans on above section)
- [ ] A separate user view
- [ ] Guild member _list_
	- Should be combined with all guild infos imo
	- Can also contain pinned messages, though I'm not sure
	- A method to call this, preferably by
	- A keybind when on the guild tree
	- Commands: `/pins` , `/members`, etc
- [x] Typing events
	- [x] The client sends the typing event
	- [x] The client receives and indicates typing events
- [x] Commands
	- [x] `/goto`
	- [x] `/edit`
	- [x] `s//` with regexp
	- [x] `/exit`, `/shrug`
	- [x] Autocompletion for those commands
		- (refer to the screenshot)
- [x] Use TextView instead of List for Messages
	- [x] Consider tv.Write() bytes
	- [x] Proper inline message edit renders
	- ~~Split messages into Primitives or find a way to edit them individually (cordless does this, too much effort)~~
- [x] Fetch nicknames and colors (16-bit hex to 256 cols somehow...)
	- [x] Async should be for later, when Split messages is done
	- [x] Add a user store
- [ ] Implement embed SIXEL images
	- [ ] Port library to [termui](https://github.com/gizak/termui)
	- [ ] Work on [issue #213](https://github.com/gizak/termui/issues/213)
- [x] Implement inline emojis
- [x] Implement auto-completion popups
	- Behavior: all keys except Enter and Esc belongs to the Input Field
	- Esc closes the popup, Enter puts the popup content into the box
	- When 0 results, hide dialog
	- Show dialog when: `@`, `#` and potentially `:` (`:` is pointless as I don't plan on adding emoji inputs any time soon)
	- Auto-completed items:
		- Mentions `@`
		- Stock emojis `:`
		- Commands `/`
		- Channels `#`
		- Messages `~`
- [x] An actual channel browser
- [x] Message acknowledgements (read/unread)
	- Isn't fully working yet, channel overrides are still janky
- [x] Message mentions
	- Partially working (only counts future mentions)
	- Past mentions using the endpoint (`/mentions`)
- [x] Scrolling up gets more messages
- [ ] Port current user stores into only Discord state caches
- [ ] Voice support (partially atm)
	- [x] Show who's in, muted, deafened and ignored
	- [ ] [Actual microphone handling](https://github.com/gordonklaus/portaudio/blob/master/examples/record.go)
	- [ ] [Auto volume](https://dsp.stackexchange.com/questions/46147/how-to-get-the-volume-level-from-pcm-audio-data)
		- Basically, I need to time so that an array of PCM int16s will contain data for 400ms
		- Then, I'll need to either root-mean-square it or calculate decibels 
		- Finally, I will compare the calculated value to the one in `config.go`
		- If it's louder, send it over to the buffer
- ~~Keyboard event handling~~
- [x] Fix `discordgo` spasming out when a goroutine panics
	- A solution could be `./6cord 2> /dev/null`
- [x] Confirm Windows compatibility
	- `/upload` fuzzy match doesn't work, wontfix

## Screenshots

<p align="center">
<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/clean.png" />
<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/mentions.png" />
<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/commands.png" />
<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/highlight.png" />
<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/reactions.png" />
</p>

## Credits

- XTerm from 
	- https://invisible-island.net/xterm/
	- https://gist.github.com/saitoha/7822989
- Fishy ([RumbleFrog](https://github.com/rumblefrog)) for his
	- [discordgo fork](https://github.com/rumblefrog/discordgo)
	- [Channel sort lib ~~that he stole from my shittercord~~](https://gist.github.com/rumblefrog/c9ebd9fb84a8955495d4fb7983345530)
- Some people on unixporn and nix nest (ym555, tdeo, ...)
- [cordless](https://github.com/Bios-Marcel/cordless) [(author)](https://github.com/Bios-Marcel) for some of the functions

