# DiscordEmotes

This is a bot that will download the discord emojis of the channels it's added to. It will have the emojis available to the server if a port is given. If no port is given, it will just download the emojis.

## Build Requirements

* Go 1.12 or newer
* GNU Make (Optional)

## Install

### Makefile

goimports must be installed if the `Makefile` is used

```bash
go get golang.org/x/tools/cmd/goimports # only run once
git clone http://github.com/joeyak/discordemotes
cd discordemotes
make
./discordemotes
```

### Go build

```bash
git clone http://github.com/joeyak/discordemotes
cd discordemotes
go build .
./discordemotes
```

## Usage

Get a bot token from discord

```bash
./discordemotes -t <token>
```
