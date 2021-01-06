install:
	go build -o kubectl-skew cmd/kubectl-skew/main.go
	mv ./kubectl-skew $$GOPATH/bin/

test:
	go test -timeout 30s -count=1 ./...

