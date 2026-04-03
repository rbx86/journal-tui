# Journal TUI

A terminal-based journal written in Go with tview.

```bash
git clone https://github.com/rbx86/journal-tui.git && cd journal-tui
go build -o ~/.local/bin/journal
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc #if you're using zsh, otherwise change to bashrc, etc.
source ~/.zshrc
```

Type `journal` to start using it :). Entries saved in `~/.journaltui/entries/`.

Requires Go and is only compatible on Linux.