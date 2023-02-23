build: packweb
	go build .

packweb:
	zip -r -j web.zip web/*