# 6cord

## Todo

- [x] Use TextView instead of List for Messages
	- [x] Consider tv.Write() bytes
	- [ ] Split messages into Primitives or find a way to edit them individually
- [ ] Fetch nicknames and colors (16-bit hex to 256 cols somehow...)
	- Async should be for later, when Split messages is done
- [ ] Implement embed SIXEL images
- [ ] Implement inline emojis
- [ ] Implement auto-completion popups
	- Behavior: all keys except Enter and Esc belongs to the Input Field
	- Esc closes the popup, Enter puts the popup content into the box
	- When 0 results, hide dialog
	- Show dialog when: `@`, `#` and potentially `:` (: is pointless as I don't plan on adding emoji inputs any time soon)
- [x] An actual channel browser

## Credits

XTerm from 
	- https://invisible-island.net/xterm/
	- https://gist.github.com/saitoha/7822989
