FROM golang:1.8.3-alpine
RUN apk add --no-cache git gcc musl-dev

WORKDIR /go/src/github.com/elBroom/meteo
ADD . .

ENV PATH_CONFIG=/go/src/github.com/elBroom/meteo/config/
RUN go build -a -o app_ .
EXPOSE 3030
EXPOSE 5432
CMD ["./app_"]

# docker build -t elbroom/meteo .
# docker push elbroom/meteo
# docker run --rm -it -p 127.0.0.1:3000:3000 --network host elbroom/meteo