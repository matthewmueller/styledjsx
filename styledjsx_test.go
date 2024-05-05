package styledjsx_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/matthewmueller/diff"
	"github.com/matthewmueller/styledjsx"
)

var update = flag.Bool("update", false, "update golden files")

func equalFile(t *testing.T, path string) {
	t.Helper()
	t.Run(path, func(t *testing.T) {
		t.Helper()
		input, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		rewriter := styledjsx.New()
		rewriter.Minify = true
		actual, err := rewriter.Rewrite(path, string(input))
		if err != nil {
			t.Fatal(err)
		}
		expected, err := os.ReadFile(path + ".golden")
		if err != nil {
			if !os.IsNotExist(err) {
				t.Fatal(err)
			}
			expected = []byte{}
		}
		if *update {
			expected = []byte(actual)
			err := os.WriteFile(path+".golden", []byte(actual), 0644)
			if err != nil {
				t.Fatal(err)
			}
		}
		diff.TestString(t, actual, string(expected))
		// Just ensure it can be parsed
		result := esbuild.Transform(actual, esbuild.TransformOptions{
			Loader: esbuild.LoaderTSX,
		})
		if len(result.Errors) > 0 {
			t.Fatalf("esbuild: %v", result.Errors)
		}
	})
}

func TestData(t *testing.T) {
	files, err := filepath.Glob("testdata/*.jsx")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		equalFile(t, file)
	}
}

func TestVendor(t *testing.T) {
	input := `
export default function () {
	return (
		<main>
			<style scoped>{` + "`" + `
				main {
					appearance: none;
					user-select: none;
					position: sticky;
				}
			` + "`" + `}</style>
		</main>
	)
}
	`
	rewriter := styledjsx.New()
	rewriter.Minify = true
	rewriter.Engines = []esbuild.Engine{
		{Name: esbuild.EngineChrome, Version: "0"},
		{Name: esbuild.EngineSafari, Version: "0"},
	}
	actual, err := rewriter.Rewrite("input.jsx", string(input))
	if err != nil {
		t.Fatal(err)
	}
	expected := `
	import Style from "styled-jsx";

	export default function () {
		return (
			<main class="jsx-2sQbAZ">
				<Style scoped id="jsx-2sQbAZ">{` + "`" + `main.jsx-2sQbAZ{-webkit-appearance:none;appearance:none;-webkit-user-select:none;-khtml-user-select:none;user-select:none;position:-webkit-sticky;position:sticky}` + "`" + `}</Style>
			</main>
		)
	}
	`
	diff.TestContent(t, actual, expected)
}

func ExampleRewrite() {
	const jsx = `
export default () => (
  <main class="main">
    <h1>hello</h1>
    <style scoped>{` + "`" + `
      .main {
        background: blue;
      }
      h1 {
        color: red;
      }
    ` + "`" + `}</style>
  </main>
)
	`
	rewriter := styledjsx.New()
	rewriter.Minify = true
	example, _ := rewriter.Rewrite("example.jsx", jsx)

	os.Stdout.WriteString(example)
	// Output:
	// import Style from "styled-jsx";
	//
	// export default () => (
	//   <main class="jsx-8mUTT main">
	//     <h1 class="jsx-8mUTT">hello</h1>
	//     <Style scoped id="jsx-8mUTT">{`.main.jsx-8mUTT{background:#00f}h1.jsx-8mUTT{color:red}`}</Style>
	//   </main>
	// )
}

func TestIgnoreHead(t *testing.T) {
	input := `
export default function () {
	return (
		<main>
			<head>
				<title>hello</title>
			</head>
			<style scoped>{` + "`" + `
				main {
					background: blue;
				}
			` + "`" + `}</style>
		</main>
	)
}
	`
	rewriter := styledjsx.New()
	actual, err := rewriter.Rewrite("input.jsx", string(input))
	if err != nil {
		t.Fatal(err)
	}
	expected := `
	import Style from "styled-jsx";

	export default function () {
		return (
			<main class="jsx-1zKxzN">
				<head>
					<title>hello</title>
				</head>
				<Style scoped id="jsx-1zKxzN">{` + "`" + `main.jsx-1zKxzN { background: blue }` + "`" + `}</Style>
			</main>
		)
	}
	`
	diff.TestContent(t, actual, expected)
}
