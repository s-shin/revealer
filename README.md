
# revealer - instant presentation server for markdown by [reveal.js](https://github.com/hakimel/reveal.js).

<!-- --- -->

## Install

```sh
go get -u github.com/s-shin/revealer
```

<!-- --- -->

## Usage

```sh
# Start server
revealer slide.md
# then open localhost:3000 in browser.

# Custom theme and separators.
./revealer \
  --theme examples/custom-theme.css \
  --separator-horizontal='^<!-- --- -->' --separator-vertical='^<!-- ___ -->' \
  README.md

# View full help text.
revealer --help
```
