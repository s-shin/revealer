
asset_files = \
  assets/index.html \
  assets/reveal/css/... assets/reveal/js/... assets/reveal/lib/... assets/reveal/plugin/... \
	assets/gemoji/db/emoji.json assets/gemoji/images/emoji/...

.PHONY: asset
asset: assets.go
	go-bindata -o assets.go $(asset_files)

.PHONY: asset_debug
asset_debug: assets.go
	go-bindata -debug -o assets.go $(asset_files)

.PHONY: build
build: assets.go
	go build

.PHONY: clean
clean:
	rm assets.go

.PHONY: build_debug
build_debug: asset_debug
	go build

.PHONY: debug
debug: build_debug
	./revealer --theme examples/custom-theme.css --separator-horizontal='^<!-- --- -->' --separator-vertical='^<!-- ___ -->' README.md
