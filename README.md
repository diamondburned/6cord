# 6cord

![](http://ix.io/1ANj.png)
![](http://ix.io/1ANk.png)

## Behaviors

- From input, hit arrow up to go to autocompletion. Arrow up again to go to the message box.
- In the message box
  - Arrow up/down and Page up/down will be used for scrolling
  - Any other key focuses back to input
- Tab to hide channels, focusing on input
- Tab again to show channels, focusing on the channel list

## Todo

- [ ] [Fix paste not working](https://github.com/rivo/tview/issues/133) (workaround: Ctrl + V)
- [ ] Commands
  - [ ] `/goto`
  - [ ] `/edit`
  - [ ] `s//` with regexp
  - [ ] `/exit`, `/shrug`
  - [ ] Autocompletion for those commands
- [ ] Fix onTyping events
- [x] Use TextView instead of List for Messages
	- [x] Consider tv.Write() bytes
	- [ ] Split messages into Primitives or find a way to edit them individually (cordless does this, too much effort)
- [x] Fetch nicknames and colors (16-bit hex to 256 cols somehow...)
	- [x] Async should be for later, when Split messages is done
	- [x] Add a user store
- [ ] Implement embed SIXEL images
- [ ] Implement inline emojis
- [x] Implement auto-completion popups
	- Behavior: all keys except Enter and Esc belongs to the Input Field
	- Esc closes the popup, Enter puts the popup content into the box
	- When 0 results, hide dialog
	- Show dialog when: `@`, `#` and potentially `:` (: is pointless as I don't plan on adding emoji inputs any time soon)
- [x] An actual channel browser

## Credits

XTerm from 
	- https://invisible-island.net/xterm/
	- https://gist.github.com/saitoha/7822989
