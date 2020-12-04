FROM node:latest as prebuild-env

WORKDIR /tmp/github.com/tarkov-database/website
COPY . .

RUN npm ci

RUN npx tsc

FROM golang:1.15.6 as build-env

WORKDIR /tmp/github.com/tarkov-database/website
COPY --from=prebuild-env /tmp/github.com/tarkov-database/website .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/website  && \
    mv -t /usr/share/tarkov-database/website frontendserver view

RUN make statics && \
    mkdir -p /usr/share/tarkov-database/website/static && \
    mv -t /usr/share/tarkov-database/website/static static/public

FROM gcr.io/distroless/base-debian10

LABEL homepage="https://tarkov-database.com"
LABEL repository="https://github.com/tarkov-database/website"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

COPY --from=build-env /usr/share/tarkov-database/website /

EXPOSE 8080

CMD ["/frontendserver"]
