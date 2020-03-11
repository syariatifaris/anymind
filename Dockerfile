FROM golang:latest

COPY ./ /
WORKDIR /
EXPOSE 9091
RUN go build -o anymind main.go

CMD ["./anymind"]