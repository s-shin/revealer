
asset_files = \
	assets/normal.html assets/aggressive.html \
	assets/reveal/css/... assets/reveal/js/... assets/reveal/lib/... assets/reveal/plugin/... \
	assets/gemoji/db/emoji.json assets/gemoji/images/emoji/...

.PHONY: asset
asset:
	go-bindata -o assets.go $(asset_files)

.PHONY: asset_debug
asset_debug:
	go-bindata -debug -o assets.go $(asset_files)

.PHONY: emoji
emoji:
	go run script/preprocess-emoji.go assets/gemoji/
	go fmt emoji.go

.PHONY: build
build: clean emoji asset
	go build

.PHONY: clean
clean:
	rm assets.go revealer emoji.go

.PHONY: build_debug
build_debug:
	go build

.PHONY: debug
debug: build_debug
	./revealer --theme examples/custom-theme.css --separator-horizontal='^<!-- --- -->' --separator-vertical='^<!-- ___ -->' README.md

.PHONY: debug2
debug2: build_debug
	./revealer --mode aggressive --theme examples/custom-theme.css --separator-horizontal='^<!-- --- -->' --separator-vertical='^<!-- ___ -->' README.md
