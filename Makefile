
LEVEL=NOTICE

all: build

build:
	mkdir -p bin
	go build -o bin/gofreedom gofreedom.go 
	
clean:
	rm -rf bin
