build:
	docker build --rm -f "Dockerfile" -t gocker-project:latest .
update:
	docker run -it -v $PWD:/app gocker-project main update