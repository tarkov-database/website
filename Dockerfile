FROM node:latest as prebuild-env

WORKDIR /tmp/github.com/tarkov-database/website
COPY . .

RUN npm ci

RUN npx tsc

FROM golang:1.18.3 as build-env

WORKDIR /tmp/github.com/tarkov-database/website
COPY --from=prebuild-env /tmp/github.com/tarkov-database/website .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/website && \
    mv -t /usr/share/tarkov-database/website frontendserver

RUN mkdir -p /usr/share/tarkov-database/website/view && \
    mv -t /usr/share/tarkov-database/website/view view/templates

RUN make statics && \
    mkdir -p /usr/share/tarkov-database/website/static && \
    mv -t /usr/share/tarkov-database/website/static static/dist

FROM gcr.io/distroless/base-debian11

LABEL homepage="https://tarkov-database.com"
LABEL repository="https://github.com/tarkov-database/website"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

LABEL org.opencontainers.image.source="https://github.com/tarkov-database/website"

COPY --from=build-env /usr/share/tarkov-database/website /

EXPOSE 8080

CMD ["/frontendserver"]
