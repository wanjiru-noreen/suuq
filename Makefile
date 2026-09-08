.PHONY: up down build logs restart status test format format-check clean

COMPOSE := docker compose
PRETTIER := npx --yes prettier@3.5.3
GO_FILES := $(shell find backend -type f -name '*.go' -print)
FRONTEND_FILES := $(shell find frontend -type f \( -name '*.html' -o -name '*.css' -o -name '*.js' \) -print)

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

build:
	$(COMPOSE) build

logs:
	$(COMPOSE) logs -f

restart:
	$(COMPOSE) restart

status:
	$(COMPOSE) ps

test:
	cd backend && go test ./...

format:
	gofmt -w $(GO_FILES)
	$(PRETTIER) --write $(FRONTEND_FILES)

format-check:
	@test -z "$$(gofmt -l $(GO_FILES))" || (echo "Go files need formatting. Run 'make format'." && exit 1)
	$(PRETTIER) --check $(FRONTEND_FILES)

clean:
	$(COMPOSE) down --volumes --remove-orphans
