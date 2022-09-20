

.PHONY: build
build:
	go build -o bin/admission_controller .

.PHONY: certs
certs:
	./dev/scripts/generate-certs.sh ./dev/certs

.PHONY: run
run: build
	./bin/admission_controller
