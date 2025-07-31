.PHONY: app

update:
	@go get "github.com/{{ORG_NAME}}/{{ORG_LOWER}}-sdk/v5@latest"
	@go get {{CHIFRA}}@latest
	@go mod tidy

lint:
	@yarn lint

test:
	@export $(grep -v '^#' ../.env | xargs) >/dev/null && go test ./...

clean:
	@rm -fR node_modules
	@rm -fR build/bin

