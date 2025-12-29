MAIN_DIRECTORY = ./cmd/gator
OUTPUT = gator

.PHONY: clean

setup: cmd/gator/main.go
	go build -o $(OUTPUT) $(MAIN_DIRECTORY)
	cd $(MAIN_DIRECTORY) && go install .

clean: gator
	rm $(OUTPUT)
