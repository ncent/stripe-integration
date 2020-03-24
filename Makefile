.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/stripe-webhook handlers/aws/apigateway/stripe-webhook/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/stripe-subscription handlers/aws/apigateway/stripe-subscription/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/stripe-checkout-session handlers/aws/apigateway/stripe-checkout-session/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/stripe-get-checkout-session handlers/aws/apigateway/stripe-get-checkout-session/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/stripe-cancel handlers/aws/apigateway/stripe-cancel/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create-subscription handlers/aws/eventbridge/create-subscription/main.go


test: build
	go test ./...

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
