build-base:
	docker build -f Dockerfile -t audit-base:latest .

build-apiserver:
	docker build -f apiserver/Dockerfile -t apiserver:latest .

build-tp:
	docker build -f tp/Dockerfile -t tp:latest .
