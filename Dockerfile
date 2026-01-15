FROM node:20-bookworm AS web-build
WORKDIR /src/web

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build


FROM golang:1.25-bookworm AS go-build
WORKDIR /src

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
 && rm -rf /var/lib/apt/lists/*

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src .
COPY integrations ./integrations

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app

WORKDIR /src/integrations/digi4school
RUN chmod +x build.sh && ./build.sh


FROM debian:bookworm-slim


RUN apt-get update && apt-get install -y --no-install-recommends \
    librsvg2-bin \
    ghostscript \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*


WORKDIR /app

COPY --from=web-build /src/web/dist /app/dist

COPY --from=go-build /src/app /app/app

COPY --from=go-build /src/integrations/d4s /app/d4s
COPY --from=go-build /src/integrations/d4s /app/integrations/d4s
RUN ls -lah /app/
RUN mkdir -p /app/data

ENTRYPOINT ["/app/app"]
