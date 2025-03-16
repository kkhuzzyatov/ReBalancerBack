FROM golang:1.23
WORKDIR /app
COPY . .
RUN go mod tidy 
WORKDIR /app/cmd
RUN go build -o calc
CMD ["./calc"]
EXPOSE 8080