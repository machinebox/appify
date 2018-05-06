test:
	go build -o testdata/app testdata/app.go
	go build -o appify
	go test
	rm appify
	rm testdata/app
