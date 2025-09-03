FROM alpine:latest

COPY main /main
COPY .env /.env

CMD ["/main"]