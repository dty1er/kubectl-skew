install:
	go build -o kubectl-ver cmd/kubectl-ver/main.go
	mv ./kubectl-ver $$GOPATH/bin/

test:
	go test -timeout 30s -count=1 ./...

