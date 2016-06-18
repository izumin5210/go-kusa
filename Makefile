all: kusa
	docker build \
		--tag izumin5210/kusa \
		--build-arg run_at="00\t22\t*\t*\t*" \
		.

kusa:
	GOOS=linux GOARCH=386 go build kusa.go

clean:
	docker rmi izumin5210/kusa
	[ -e kusa ] && rm kusa
