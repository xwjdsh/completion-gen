run:
	docker run --rm -v "$(PWD)/dist:/workspace/dist" -it $(shell DOCKER_SCAN_SUGGEST=false docker build -q .) -c docker -w dist/docker.ts
