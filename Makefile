run: build
	./zk-view

build:
	go build -o zk-view main.go
clean:
	rm zk-view
