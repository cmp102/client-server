all: target cpexamples target/server

target:
	mkdir target
	mkdir target/data

target/server: server.go
	go build -o target/server server.go

cpexamples: examples/*
	cp -r examples/* target/

clean:
	rm -rf target