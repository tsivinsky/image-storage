out_dir = ./build
out_file = $(out_dir)/app

build:
	go build -o $(out_file) cmd/main.go

dev:
	go run cmd/main.go

run:
	$(out_file)

clean:
	rm -rf $(out_dir)
