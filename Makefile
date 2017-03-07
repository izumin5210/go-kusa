all: kusa
	docker build \
		--tag izumin5210/kusa \
		.

kusa:
	GOOS=linux GOARCH=386 go build kusa.go

clean:
	docker rmi izumin5210/kusa
	[ -e kusa ] && rm kusa
