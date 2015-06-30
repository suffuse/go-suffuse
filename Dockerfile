# I give up on supporting extended attributes in a container.
# https://github.com/docker/docker/issues/1070
# If we do ever want to, we need python-xattr for an xattr command.
FROM paulp/debian

ENV GOPATH /go
ENV PATH "$GOPATH/bin:$PATH"

ADD go /go
RUN go get -t -d -v github.com/paulp/suffuse/...
RUN go build github.com/paulp/suffuse/...
RUN go install github.com/paulp/suffuse/...

ENTRYPOINT [ "go", "test", "-v", "github.com/paulp/suffuse" ]
