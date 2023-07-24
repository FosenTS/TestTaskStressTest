.PHONY: build
build:
	docker build -t testtask .

.PHONY: start
start:
	docker run -it --rm testtask