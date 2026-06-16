.PHONY: build
build:
		go build -o bin/viv ./cmd

.PHONY: clean
clean:
		rm bin/*

.PHONY: doc
doc: build
	bin/viv doc
