FROM golang:1.12
COPY ./*.go /app/
COPY ./go.* /app/
ENV GOPROXY=http://goproxy.io
EXPOSE 8888
WORKDIR /app

ENTRYPOINT ["go","run","main.go"]
