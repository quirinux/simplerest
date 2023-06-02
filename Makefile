
.ONESHELL:

CONFIG := ./examples/config_sqlite3.toml
GO := go
MAIN := main.go
BENCH := ./tests/benchmark.js
DBVERSION := latest
TOKEN := 467aa100-7883-4cbd-8152-b3478a0c3d0d

define rundb
	docker run -d \
		--env-file .env-file \
		--hostname database \
		--name database \
		--publish ${2}:${2} \
		${1}:${DBVERSION}
endef

run:
	${GO} run -tags=jsoniter ${MAIN} --config ${CONFIG}

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
	http GET :8888/todos
	http POST :8888/todos name=qui description=nux
	http PUT :8888/todos/1 name=bar description=foo
	http GET :8888/todos/1
	http DELETE :8888/todos/5
	http GET :8888/todos/5
	http GET :8888/error/missing_table
	http GET :8888/error/missing_param
