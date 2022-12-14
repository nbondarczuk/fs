FROM docker.github.azc.ext.hp.com/krypton/utils/krypton-go-builder as builder

ADD . /go/src/fs
WORKDIR /go/src/fs

# build the source
RUN make tidy && make build-binaries

# use a minimal alpine image
FROM docker.github.azc.ext.hp.com/krypton/utils/krypton-go-base

# set working directory
WORKDIR /go/bin

ENV GIN_MODE=release
COPY --from=builder /go/src/fs/bin .
COPY --from=builder /go/src/fs/config/config.yaml .

USER 1001
EXPOSE 8989/tcp

# run the binary
CMD ["./fs"]
