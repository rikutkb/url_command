test-html:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
