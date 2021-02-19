.PHONY: build clean

BINARY="shovel"

.PHONY: build
build:
	GOOS=linux GOARCH="amd64" go build -o ${BINARY} ./src/shovel/main/main.go

.PHONY: install
install:
	@govendor sync -v

.PHONY: clean
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
