FROM golang:1.14 AS build

WORKDIR /src/
COPY go.mod go.sum /src/
RUN go mod download

COPY internal /src/internal
COPY cmd /src/cmd
COPY .git /src/.git

RUN export GIT_COMMIT=$(git rev-list -1 HEAD)

RUN CGO_ENABLED=0 go build -o /bin/sshspawner -ldflags "-X main.GitCommit=$(git rev-list -1 HEAD)" /src/cmd/spawner/main.go

FROM scratch
EXPOSE 22
COPY --from=build /bin/sshspawner /bin/sshspawner
ENTRYPOINT ["/bin/sshspawner"]
