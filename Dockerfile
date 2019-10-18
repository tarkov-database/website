FROM golang:1.13

EXPOSE 8080

WORKDIR /tmp/github.com/tarkov-database/website
COPY . .

RUN make bin && \
    mkdir -p /usr/share/tarkov-database/website && \
    mv -t /usr/share/tarkov-database/website frontendserver view static && \
    rm -rf /tmp/github.com/tarkov-database/website

WORKDIR /usr/share/tarkov-database/website

CMD ["/usr/share/tarkov-database/website/frontendserver"]
