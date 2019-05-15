all: target cpscripts target/send-file target/save-to-disk

target:
	mkdir target
	mkdir target/data

target/server: server.go
	go build -o target/server server.go

cpexamples: examples/*
	cp examples/* target/

clean:
	rm -rf target