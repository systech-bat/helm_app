FROM golang:1.20 AS build

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o app ./cmd/weather-bot

FROM alpine:latest    
                   
WORKDIR /app

COPY --from=build app .    
COPY --from=build /app/default-template.txt .    

CMD ["./app"]
