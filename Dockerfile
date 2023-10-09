FROM golang:1.21.1 AS build
LABEL maintainer='Edigar Herculano <edigarhdev@gmail.com>'

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 go build -o /app

#--------
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app /app

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT [ "/app/socialnets-api" ]