# 6cord

## Todo

- [ ] Use TextView instead of List for Messages
	- [ ] Consider tv.Write() bytes
- [ ] Implement embed SIXEL images
- [ ] Implement inline emojis
- [ ] Implement auto-completion popups
	- Behavior: all keys except Enter and Esc belongs to the Input Field
	- Esc closes the popup, Enter puts the popup content into the box
	- When 0 results, hide dialog
	- Show dialog when: `@`, `#` and potentially `:` (: is pointless as I don't plan on adding emoji inputs any time soon)
- [ ] An actual channel browser

## Credits

XTerm from 
	- https://invisible-island.net/xterm/
	- https://gist.github.com/saitoha/7822989
