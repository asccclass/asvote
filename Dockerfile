FROM golang

WORKDIR /app
COPY ./app /app
ENV PORT=80
ENTRYPOINT ["/app/app"]
