
# revealer - instant presentation server for markdown by [reveal.js](https://github.com/hakimel/reveal.js).

<!-- --- -->

## Install

```sh
go get -u github.com/s-shin/revealer
```

<!-- --- -->

## Usage

```sh
# Essentially, start server like this
revealer slide.md
# then open localhost:3000 in browser.

# Run with some options.
revealer \
  --mode=aggressive \
  --theme=examples/custom-theme.css \
  --separator-horizontal='^<!-- --- -->' --separator-vertical='^<!-- ___ -->' --separator-note='^Notes:' \
  README.md

# View other options and help text.
revealer --help
```

<!-- ___ -->

### Custom Theme Example

```css
@import "./black.css";

body {
  background-color: #000;
}

.reveal code {
  font-family: Ricty, Consolas, monospace;
  border-radius: 0.1em;
}

.reveal :not(pre) code {
  background-color: rgba(255, 255, 255, 0.2);
  padding: 0.05em 0.2em;
  font-size: 0.9em;
}

.reveal pre code {
  padding: 0.4em 0.6em;
}
```

<!-- --- -->

## Modes

Two modes are available:

* Normal
* Aggressive

<!-- ___ -->

### Normal

* Execute revealer with `--mode=normal` (or without `--mode` option because of default).
* Features
  * Converts markdown by plugin of Reveal.js.
  * Works just like [revealgo](https://github.com/yusukebe/revealgo).

<!-- ___ -->

### Aggressive

* Execute revealer with `--mode=aggressive`.
* Features
  * Converts markdown by processor in revealer.
  * Supports emoji aliases :smile: :octocat: :+1:
  * Ignores slide separators in code blocks.

```sh
<!-- --- -->
<!-- ___ -->
```

Note:
* `^<!-- --- -->` and `^<!-- ___ -->` are assumed as horizontal/vertical separators in this document.

<!-- --- -->

## For Developers

```sh
# For submodules
git submodule update --init

# For auto generated files
make asset
make emoji

# Build revealer
make build

# There are more tasks in Makefile, please see it.
```

<!-- --- -->

## Dependencies

* [reveal.js](https://github.com/hakimel/reveal.js)
* [gemoji](https://github.com/github/gemoji)
* [github.com/urfave/cli](https://github.com/urfave/cli)
* [github.com/jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata)
* [github.com/russross/blackfriday](https://github.com/russross/blackfriday)
* [github.com/labstack/echo](https://github.com/labstack/echo)

<!-- --- -->

## See Also

* [revealgo](https://github.com/yusukebe/revealgo)
  * Revealer respects the concept of this.

<!-- --- -->

## Author

Shintaro Seki https://github.com/s-shin/

<!-- --- -->

## License

The MIT License
