
.ONESHELL:

CONFIG := ./examples/config_sqlite3.toml
GO := go
GOFLAGS :=-race -v -tags=jsoniter 
MAIN := cmd/simplerest.go
BENCH := ./tests/benchmark.js
DBVERSION := latest
TOKEN := 467aa100-7883-4cbd-8152-b3478a0c3d0d
BINS := $(wildcard cmd/*.go)
BINDIR := bin
VERBOSITY := -vvvvvvvv

define rundb
	docker run -d \
		--env-file .env-file \
		--hostname database \
		--name database \
		--publish ${2}:${2} \
		${1}:${DBVERSION}
endef

build: $(foreach f,$(BINS),$(BINDIR)/$(basename $(notdir $(f))))
$(BINDIR)/%: cmd/%.go
	${GO} build ${GOFLAGS} -o $@ $<

clean:
	rm -rdf $(BINDIR)

run:
	${GO} run ${GOFLAGS} ${MAIN} --config ${CONFIG}

benchmark:
	k6 run ${BENCH}

examples/certificate/self-signed.csr:
	mkdir -p examples/certificate/
	openssl req -new -subj "/C=US/ST=California/CN=localhost" \
    -newkey rsa:2048 -nodes -keyout "examples/certificate/self-signed.key" -out "examples/certificate/self-signed.csr"
	openssl x509 -req -days 365 -in "examples/certificate/self-signed.csr" -signkey "examples/certificate/self-signed.key" \
		-out "examples/certificate/self-signed.crt" -extfile "examples/certificate/self-signed.ext"

database-postgres:
	-$(call rundb,postgres,5432)

database-mysql:
	-$(call rundb,mariadb,3306)

database-oracle:
	-$(call rundb,gvenzl/oracle-xe,1521)

database-sqlite3:
	@echo "nothing to do for sqlite3"

stopdb:
	-docker stop database
	-docker rm database

test:
	http ${VERBOSITY} GET  :8888/todos
	http ${VERBOSITY} POST :8888/todos name=qui description=nux
	http ${VERBOSITY} PUT :8888/todos/1 name=bar description=foo
	http ${VERBOSITY} GET :8888/todos/1
	http ${VERBOSITY} DELETE :8888/todos/5
	http ${VERBOSITY} GET :8888/todos/5
	http ${VERBOSITY} GET :8888/error/missing_table
	http ${VERBOSITY} GET :8888/error/missing_param
	http ${VERBOSITY} GET  :8888/todos "Accept:application/x-yaml"
	http ${VERBOSITY} GET  :8888/todos "Accept:text/yaml"
	http ${VERBOSITY} GET  :8888/todos "Accept:application/xml"
	http ${VERBOSITY} GET  :8888/todos "Accept:text/xml"
	http ${VERBOSITY} GET  :8888/todos "Accept:application/toml"
	http ${VERBOSITY} GET  :8888/todos "Accept:text/csv"
	http ${VERBOSITY} GET  :8888/todos "Accept:text/html"
