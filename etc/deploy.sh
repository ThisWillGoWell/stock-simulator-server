DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $GOPATH/src/github.com/ThisWillGoWell/stock-simulator-server
go build
docker build