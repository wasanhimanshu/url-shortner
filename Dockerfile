FROM golang:1.22.3 AS build-stage
  WORKDIR /app

  COPY go.mod go.sum ./
  RUN go mod download

  COPY *.go ./

  RUN CGO_ENABLED=0 GOOS=linux go build -o /urlshortner

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
  WORKDIR /app

  COPY --from=build-stage /urlshortner /urlshortner

  EXPOSE 8080

  ENTRYPOINT ["/urlshortner"]