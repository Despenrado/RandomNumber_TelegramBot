# RandomNumberGo_bot

This is a simple telegram bot, that generates diferent quantity of random numbers and calculates average of it.

## HELP

```
structure of queries: /command [parametr] [parametr]
example: /setmin 10 - this command sets minimum border to 10
commands:
help - shows this message
settemplate - shows list of templates
setmin [number]- sets minimum border
setmax [number] - sets maximum border
setquantity [number] - sets number of random numberssetminmaxqua [min] [max] [quantity] - sets minimum, maximum and quantity
setwords [word1;word2;word3] - sets words for random choice
status - shows your current tamplate of rundom
random
roll - genertes random nuber/numbers via current template
```

## Compile

- arm:

```
env GOOS=linux GOARCH=arm64 go build -o ./deploy/arm
```

- x86:

```
env GOOS=linux GOARCH=amd64 go build -o ./deploy/arm
```

## Run

- Docker:

```
docker run
```

- linux:

```
./deploy/amd64/RandomNumberGo_bot
```

- Without compile

```
go run RandomNumberGo_bot.go
```

[](easter_egg:/serverip,/nextcloud)
