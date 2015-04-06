.PHONY: build clean extension

build:
	mkdir -p ./build
	go build -o ./build/gemsearch-alfred

workflow: build
	mkdir -p ./tmp
	rm -rf ./tmp/workflow
	cp -Rp ./workflow ./tmp/workflow
	cp -fp ./build/gemsearch-alfred ./tmp/workflow/gemsearch-alfred
	perl -ne 's/README_INSERTION/`cat README.md`/e;print' ./workflow/info.plist > ./tmp/workflow/info.plist
	cd ./tmp/workflow && zip -r Gemsearch.alfredworkflow *
	mv ./tmp/workflow/Gemsearch.alfredworkflow ./build/Gemsearch.alfredworkflow

clean:
	rm -rf ./tmp
	rm -rf ./build
	go clean
