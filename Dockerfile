FROM golang:1.15.2 as build-env

ARG BRANCH=""

ENV BRANCH=${BRANCH}

WORKDIR /tmp/github.com/tarkov-database/website
COPY . .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/website && \
    mv -t /usr/share/tarkov-database/website frontendserver view static

FROM gcr.io/distroless/base-debian10

LABEL homepage="https://tarkov-database.com"
LABEL repository="https://github.com/tarkov-database/website"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

COPY --from=build-env /usr/share/tarkov-database/website /

EXPOSE 8080

CMD ["/frontendserver"]
