.PHONY: up down build logs restart status test clean

COMPOSE := docker compose

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

clean:
	$(COMPOSE) down --volumes --remove-orphans