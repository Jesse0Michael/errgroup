
test:
	go test -cover ./... 

test-coverage:
	go test -coverpkg ./... -coverprofile coverage.out ./... && go tool cover -html=coverage.out

