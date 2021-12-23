build:
	docker build -t wendellsun/completion-gen .

run:
	docker run --rm -v "$(PWD)/tmpls:/workspace/tmpls" -v "$(PWD)/dist:/workspace/dist" -it wendellsun/completion-gen -c docker -w dist/docker.ts
