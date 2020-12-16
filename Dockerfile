FROM golang:1.14.3 AS build
WORKDIR /src
COPY . .

RUN export GIT_COMMIT=$(git rev-list -1 HEAD)
RUN go build \
	-ldflags "-X main.GitCommit=$GIT_COMMIT" \
	-o /go/bin/sshspawner

#FROM scratch
#COPY --from=build /go/bin/sshspawner /go/bin/sshspawner
#COPY --from=build /bin/ls /bin/ls
EXPOSE 22
ENTRYPOINT ["/go/bin/sshspawner"]