export UID := `id -u`
export GID := `id -g`

COMPOSE := 'docker compose -f docker/compose.yml'
COMPOSE-RUN := COMPOSE + ' run --rm'
GO-RUN := COMPOSE-RUN + ' go'
SVELTEKIT-RUN := COMPOSE-RUN + ' sveltekit'
DB-RUN := COMPOSE-RUN + ' -T --no-deps db'

default:
	just --list

# start all containers
start *args:
	{{COMPOSE}} up -d {{args}}
	@echo URL: http://localhost:5173

# stop all containers
stop:
	{{COMPOSE}} stop

# remove all containers and unversioned changes
clean:
	{{COMPOSE}} rm -vsf
	git clean -fdqx -e .idea

# run just clean, then rebuild containers with pulling anew & install dependencies
rebuild: clean
	{{COMPOSE}} build --pull
	just install

# rebuild containers without cleaning up and pulling anew
soft-rebuild:
	{{COMPOSE}} build

# install svelteApp dependencies
install:
	{{SVELTEKIT-RUN}} pnpm install

# seed DB with testdata
init-db:
    cat migrations/*.sql | {{DB-RUN}} mysql -uroot -psoccer -hdb zendo

# exec into a running container
exec-ti *args:
	{{COMPOSE}} exec -ti {{args}}

# direct access to pnpm command inside a svelteApp container
pnpm *args:
    {{SVELTEKIT-RUN}} pnpm {{args}}

# short for: docker compose -f docker/compose.yml
compose *args:
	{{COMPOSE}} {{args}}

# direct access to mysql command inside a db container
mysql *args:
    {{DB-RUN}} mysql {{args}}

# run eslint on sveltekit app  code
eslint-sveltekit *args:
    {{SVELTEKIT-RUN}} pnpm eslint {{args}}
