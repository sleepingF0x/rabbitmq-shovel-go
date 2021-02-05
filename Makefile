.PHONY: build clean

BINARY="shovel"

.PHONY: build
build:
	GOOS=linux GOARCH="amd64" go build -o ${BINARY}-linux-amd64 ./src/shovel/

.PHONY: install
install:
	@govendor sync -v

.PHONY: clean
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
