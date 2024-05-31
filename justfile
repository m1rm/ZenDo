export UID := `id -u`
export GID := `id -g`

COMPOSE := 'docker compose -f docker/compose.yml'
COMPOSE-RUN := COMPOSE + ' run --rm'
GO-RUN := COMPOSE-RUN + ' go'
SVELTEKIT-RUN := COMPOSE-RUN + ' sveltekit'
DB-RUN := COMPOSE-RUN + ' --no-deps db'

default:
	just --list

# lifecycle commands
start:
	{{COMPOSE}} up -d
	{{DB-RUN}} mysqladmin -uroot -pChangeMe -hdb --wait=10 ping
	@echo URL: http://localhost:8080

stop:
	{{COMPOSE}} stop

clean:
	{{COMPOSE}} rm -vsf
	git clean -fdqx -e .idea


rebuild: clean
	{{COMPOSE}} build --pull
	just install

install:
	{{SVELTEKIT-RUN}} pnpm install

pnpm *args:
    {{SVELTEKIT-RUN}} pnpm {{args}}

	

