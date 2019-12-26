build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create-timer create-timer/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/turn-on-aircon turn-on-aircon/main.go

.PHONY: clean
clean:
	rm -rf ./bin 

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
