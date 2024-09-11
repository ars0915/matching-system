gen:
	go generate ./...

docker-build:
	docker build . --build-arg APP_NAME=matching-system -f docker/Dockerfile -t matching-system

docker-run:
	docker run --name matching-system -d -p 8080:8080 matching-system

tests:
	go test -v  ./...