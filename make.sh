
export GOROOT=/usr/local/go
export GOBIN=$GOROOT/bin
export GOPATH=~/XRZP/golang/server:$GOROOT:~/XRZP/golang/server/gopath:~/go/
export PATH=$PATH:$GOROOT/bin:$GOBIN
export GOOS=linux GOARCH=amd64

echo $(go version)
go clean

cd bin
go build ../server/src
echo build src ok !
go build ../server/src/cache
echo build cache ok !
go build ../server/src/robot
echo build src robot !

