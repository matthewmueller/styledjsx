# styledjsx

[![Go Reference](https://pkg.go.dev/badge/github.com/matthewmueller/styledjsx.svg)](https://pkg.go.dev/github.com/matthewmueller/styledjsx)

A native implementation of [styled-jsx](https://github.com/vercel/styled-jsx) for Go. Works with both JSX and TSX files.

It's built on:

- [jsx](https://github.com/matthewmueller/jsx): A JSX parser
- [css](https://github.com/matthewmueller/css): A CSS parser

## Features

- Very fast (native Go!)
- Supports CSS minification via ESBuild
- Supports vendor prefixing via ESBuild
- Supports `class` and `className` variants
- Supports `<style scoped>` and `<style jsx>` variants

## Usage

Given the following JSX:

```css
export default () => (
  <main class="main">
    <h1>hello</h1>
    <style scoped>{`
      .main {
        background: blue;
      }
      h1 {
        color: red;
      }
    `}</style>
  </main>
)
```

Rewrite that JSX with the following:

```go
out, _ = styledjsx.Rewrite("input.jsx", jsx)
```

Resulting in:

```jsx
import Style from "styled-jsx"

export default () => (
  <main class="jsx-8mUTT main">
    <h1 class="jsx-8mUTT">hello</h1>
    <Style scoped>{`
      .main.jsx-8mUTT{background:#00f}
      h1.jsx-8mUTT{color:red}
    `}</Style>
  </main>
)
```

**Note:** you'll need to bring your own `<Style/>` component. The JS library `styled-jsx/style` might work, but it's currently untested.

### Planned

- [ ] Scope `@keyframe` names
- [ ] `:global()` support
- [ ] VSCode support for `<style scoped>` (fork [this repo](https://github.com/iChenLei/vscode-styled-jsx))
- [ ] Example client library instead of `styled-jsx`

### Unplanned

[styled-jsx](https://github.com/vercel/styled-jsx) has a plethora of ways to include CSS. I tend to stick with the very basics. Here are some features that are not planned and what the alternative is:

- **Dynamic styles**: Use toggling class names or CSS variables to provide this.
- `<style global>`: Use `<style>` instead.
- **Embedding constants**: This would be nice, but it's not worth the effort.
- **Plugins**: Write an ESBuild plugin instead.

## Why?

Styled JSX is still the most natural way to write CSS in JSX.

You'll probably use this library alongside a bundler like [ESBuild](https://github.com/evanw/esbuild).

## Install

```sh
go get github.com/matthewmueller/styledjsx
```

## Contributors

- Matt Mueller ([@mattmueller](https://twitter.com/mattmueller))

## License

MIT
