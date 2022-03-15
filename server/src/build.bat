SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
SET GODEBUG=madvdontneed=1
go clean
go build -o src -x  main.go