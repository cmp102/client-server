all: target cpscripts target/send-file target/save-to-disk

target:
	mkdir target
	mkdir target/data

target/send-file: send-file-src/send-file.cpp
	g++ -o target/send-file send-file-src/send-file.cpp -lcurl

target/save-to-disk: save-to-disk-src/save-to-disk.cpp
	g++ -o target/save-to-disk save-to-disk-src/save-to-disk.cpp -lrealsense

cpscripts: scripts/*
	cp scripts/* target/

clean:
	rm -rf target send-file save-to-disk