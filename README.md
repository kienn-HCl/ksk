# ksk

[日本語](README.ja.md)

A terminal-based web search TUI app. Browse search results with vim-like keybindings and open them in your browser.

![screenshot](screenshot/search.png)

## Install

```
go install github.com/frort/ksk@latest
```

Or build from source:

```
git clone https://github.com/frort/ksk.git
cd ksk
go build -o ksk .
```

## Usage

```
ksk                          # launch and type a query
ksk "search terms"           # launch with an initial query
ksk -e brave "search terms"  # search with Brave
ksk -r jp "search terms"     # search with region set to Japan
ksk -e brave -r de "query"   # Brave + Germany
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-e` | `duckduckgo` | Search engine (`duckduckgo` / `ddg`, `brave` / `b`) |
| `-r` | _(none)_ | Region / country code (`jp`, `us`, `de`, `fr`, `kr`, etc.) |

### Supported engines

| Engine | Alias | Notes |
|--------|-------|-------|
| DuckDuckGo | `ddg` | Default. HTML scraping via `html.duckduckgo.com` |
| Brave Search | `b` | HTML scraping via `search.brave.com` |

## Keybindings

### Results mode

| Key | Action |
|-----|--------|
| `j` / `↓` | Next result |
| `k` / `↑` | Previous result |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `l` / `→` | Next page |
| `h` / `←` | Previous page |
| `Enter` / `o` | Open in browser |
| `/` | Search |
| `q` / `Ctrl+C` | Quit |

### Input mode

| Key | Action |
|-----|--------|
| `Enter` | Execute search |
| `Escape` | Back to results |
| `Ctrl+C` | Quit |

## See Also

- [ddgr](https://github.com/jarun/ddgr) - DuckDuckGo from the terminal
- [googler](https://github.com/jarun/googler) - Google from the terminal
- [searxngr](https://github.com/scross01/searxngr) - SearXNG from the command line
