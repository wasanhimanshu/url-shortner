build:
	@go build -o bin/url-shortner .

run: build
	@./bin/url-shortner