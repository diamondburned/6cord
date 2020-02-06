<h1 align="center">6cord</h1>
<p align="center">
	<a href="https://liberapay.com/diamondburned/donate">
		<img align="center" alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg">
	</a>
	<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/latest.png" />
	<img src="https://u.cubeupload.com/diamondburned/MCz9fP.png" />
</p>

## Installation

### Method 1. (recommended, precompiled binaries)

**[Linux](https://gitlab.com/diamondburned/6cord/builds/artifacts/master/file/6cord?job=linux)**
**[Linux (no dbus)](https://gitlab.com/diamondburned/6cord/builds/artifacts/master/file/6cord_nk?job=linux)**
**[Linux (arm64)](https://gitlab.com/diamondburned/6cord/builds/artifacts/master/file/6cord_arm64?job=linux_arm64)**
**[Windows](https://gitlab.com/diamondburned/6cord/builds/artifacts/master/file/6cord.exe?job=windows)**

Only do this if the CI passed (a green tick in the commit bar)

### Method 2. (building from source)

```sh
git clone https://gitlab.com/diamondburned/6cord
cd 6cord && go build
./6cord

# Optional
mkdir -p ~/bin/
mv ./6cord ~/bin/
echo PATH="$HOME/bin:$PATH" >> ~/.bashrc && . ~/.bashrc # or any shellrc
```

### Method 3. (package manager)

```sh
# Arch Linux, using your favourite AUR helper:
yay install 6cord
# Alternatively you can install '6cord-git'
# which is the latest development version.

# FreeBSD:
pkg install 6cord
```

## Getting the token

This is possible from both the web client and the Electron client.

1. Hit <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>I</kbd>
2. Switch to the `Network` tab
3. Find Discord API requests. This is usually called `messages`, `ack`, `typing`, etc
4. Search for the `Authorization` header. This is the token.

## Running 6cord with the token

`./6cord -t "TOKEN_HERE"`

- If you have Gnome Keyring (usually the case on most DEs), the token would automatically be stored securely. This could be tested by running `./6cord` without any arguments.
	- To reset the token, override it with a new one using `-t`
- It is also possible to move the `6cord.toml` file from the root of this Git repository to `~/.config/6cord/`, then run without any arguments.

## Additional things

### Quirks

- The <kbd>~</kbd> key could be used to both preview images and select a message ID
- `/mentions` is useless at the moment. This is planned to change in the future.
- There is currently no global emoji support. This is also planned to change, along with emoji previews.

### Additional keybinds

- Refer to the Quick Start section displayed when starting 6cord
- <kbd>Tab</kbd> to show/hide the server list
- Input field history is cycled with <kbd>Alt</kbd> + <kbd>Up</kbd>/<kbd>Down</kbd>
- <kbd>PgUp</kbd> and <kbd>PgDn</kbd> can be used to jump between servers in the list
- There are some Vim binds available ie ^n and ^p to move between fuzzy listed items

### `command-prefix`

- The following variables are available: `CHANNEL`, `GUILD`, `USERNAME` and `DISCRIM`
- This follows `tview`'s rich text format:
	- Coloring text with `[#424242]`
	- Bold text with `[::b]`
	- Both can be done with `[#424242::b]`
	- Reset with `[-]`, `[::-]` or `[-::-]`
- You need to manually escape square brackets by adding an opening (`[`) bracket before a closing (`]`) bracket
	- Example: `[${guild}]` to `[${guild}[]`

### Color support

6cord runs in 256 color mode most of the time. To force true color, run:

```sh
TERM=xterm-truecolor ./6cord`
```

(`xterm-truecolor` is known to break a lot of applications including `htop`, only use it with `6cord`)

To limit 6cord to strictly 16 colors, run:

```sh
TERM=xterm-basic ./6cord
```

To run 6cord in monochrome mode:

```sh
TERM=xterm ./6cord
```

### Supported Image backends

Currently, Xorg is the only supported image backend. SIXEL support proved itself to be challenging with how `tcell` and `tview` call redraws. There is no Kitty terminal implementation in Golang that is available as a library yet (`termui` has a PR with Kitty support). There are things in my priority list right now. That said, PRs are welcomed.

## Screenshots

<p align="center">
<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/guildview.png" />
<img src="https://gitlab.com/diamondburned/6cord/raw/master/_screenshots/img.png" />
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

