FROM golang:1.21.1 AS build
LABEL maintainer='Edigar Herculano <edigarhdev@gmail.com>'

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 go build -o /app /app/cmd/api

RUN addgroup nonroot &&  \
    adduser --ingroup nonroot --uid 19998 --shell /bin/false nonroot &&  \
    cat /etc/passwd | grep nonroot > /etc/passwd_nonroot

#--------
#FROM gcr.io/distroless/base-debian10
FROM scratch AS prod

WORKDIR /app

COPY --from=build /etc/passwd_nonroot /etc/passwd
COPY --from=build /app /app

EXPOSE 8000

USER nonroot

ENTRYPOINT [ "./api" ]