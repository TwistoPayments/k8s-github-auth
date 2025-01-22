IMAGE_NAME := public.ecr.aws/twisto/k8s-github-auth:latest

.PHONY: build
build:
	GOOS=linux GARCH=amd64 go build -o build/main main.go

.PHONY: clean
clean:
	rm -rf build

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)
