FROM golang

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go get github.com/botanio/sdk/go 
RUN go get github.com/Syfaro/telegram-bot-api

RUN go build -o RandomNumberGo_bot .

CMD ["/app/RandomNumberGo_bot"]