
SRCS=gclife.go

myapp: $(SRCS)
	docker run --rm -v "$$(pwd)":/usr/src/myapp \
			-w /usr/src/myapp golang:1.3.1 go build -v

run: myapp
	@./myapp $(WIDTH) $(HEIGHT)

build-image:
	docker build -t dochan-life .

run-image: build-image
	docker run --rm -it gochan-life $(WIDTH) $(HEIGHT)

clean:
	rm -f myapp
