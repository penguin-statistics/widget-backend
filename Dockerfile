# production service runner

FROM node:16-alpine AS webbuilderbase

RUN apk update && apk add git make gcc g++ yarn

FROM webbuilderbase AS docsbuilder

RUN git clone https://github.com/penguin-statistics/widget-docs /build/widget-docs

WORKDIR /build/widget-docs

RUN yarn install && yarn build

FROM webbuilderbase AS frontendbuilder

RUN git clone https://github.com/penguin-statistics/widget-frontend /build/widget-frontend

WORKDIR /build/widget-frontend

RUN yarn install && yarn build

FROM golang:1.18-alpine
WORKDIR /app

# builder
FROM base AS gobuilder
ENV GOOS linux
ENV GOARCH amd64

# modules: utilize build cache
COPY go.mod ./
COPY go.sum ./

# RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .

# build the binary
RUN go build -o widgetbackend .

# runner
FROM base AS runner
RUN apk add --no-cache libc6-compat

RUN apk add --no-cache tini
# Tini is now available at /sbin/tini

COPY --from=gobuilder /app/widgetbackend /app/widgetbackend
COPY --from=gobuilder /app/config.example.yml /app/config.yml

COPY --from=docsbuilder /build/widget-docs/dist /app/widget-docs
COPY --from=frontendbuilder /build/widget-frontend/dist /app/widget-frontend

ENV WIDGET_BACKEND_STATIC_DOCS_ROOT=/app/widget-docs
ENV WIDGET_BACKEND_STATIC_WIDGET_ROOT=/app/widget-frontend

ENTRYPOINT ["/sbin/tini", "--"]
CMD [ "/app/widgetbackend" ]
