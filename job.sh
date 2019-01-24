cd job/workerserver/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../flowerworker main.go

cd ../jobserver/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ../../flowercron main.go