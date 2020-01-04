build-apiserver:
	docker build -f apiserver/Dockerfile -t apiserver:latest .

run-apiserver:
	docker run -it -p 8080:8080 apiserver:latest
