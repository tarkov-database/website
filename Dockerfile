FROM golang:1.13.8

LABEL homepage="https://tarkov-database.com"
LABEL repository="https://github.com/tarkov-database/website"
LABEL maintainer="Markus Wiegand <mail@morphy2k.dev>"

ARG BRANCH=""

ENV BRANCH=${BRANCH}

EXPOSE 8080

WORKDIR /tmp/github.com/tarkov-database/website
COPY . .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/website && \
    mv -t /usr/share/tarkov-database/website frontendserver view static && \
    rm -rf /tmp/github.com/tarkov-database/website

WORKDIR /usr/share/tarkov-database/website

CMD ["/usr/share/tarkov-database/website/frontendserver"]
