<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->
<p align="center">clix -- Generate markdown from go-flags</p>
<p align="center">
  <a href="https://github.com/lrstanley/clix/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/clix?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/clix?style=flat-square">
  </a>

  <a href="https://github.com/lrstanley/clix/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/workflow/status/lrstanley/clix/test/master?label=test&style=flat-square&event=push">
  </a>

  <a href="https://codecov.io/gh/lrstanley/clix">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/clix/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/clix">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/clix?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/clix">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/clix?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/clix/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/clix/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/clix/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/clix?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/clix/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->

## :link: Table of Contents

-   [TODO](#ballot_box_with_check-todo)
-   [Usage](#gear-usage)
    -   [Example output](#example-output)
        -   [Application Options](#application-options)
        -   [Example Group](#example-group)
        -   [Help Options](#help-options)
-   [Support &amp; Assistance](#raising_hand_man-support--assistance)
-   [Contributing](#handshake-contributing)
-   [License](#balance_scale-license)
<!-- template:end:toc -->

## :ballot_box_with_check: TODO

-   [ ] Generate commands/sub-commands/etc
-   [ ] Custom markdown formatting

## :gear: Usage

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->

```console
$ go get -u github.com/lrstanley/clix@latest
```

<!-- template:end:goget -->

## Generate Markdown

Generate markdown from a provided `flags.Parser`, that can be embedded in
a readme.

```go
package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
	clix "github.com/lrstanley/clix"
)

var (
	cli = &Flags{}
)

type Flags struct {
	VersionFlag bool `short:"v" long:"version" description:"display the version and exit"`
	Debug       bool `env:"DEBUG" short:"D" long:"debug" description:"enable debugging"`

	SubFlags struct {
		Username string `env:"USERNAME" short:"u" long:"username" default:"admin" description:"example username"`
		Password string `env:"PASSWORD" short:"p" long:"password" description:"example password"`
	} `group:"Example Group" namespace:"example" env-namespace:"EXAMPLE"`
}

func main() {
	var err error

	parser := flags.NewParser(cli, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)

	if _, err = parser.Parse(); err != nil {
		os.Exit(1)
	}

	clix.GenerateMarkdown(parser, os.Stdout)
	os.Exit(0)
}
```

### Example output

Raw:

```markdown
#### Application Options

| Environment vars | Flags           | Type | Description                  |
| ---------------- | --------------- | ---- | ---------------------------- |
| -                | `-v, --version` | bool | display the version and exit |
| `DEBUG`          | `-D, --debug`   | bool | enable debugging             |

#### Example Group

| Environment vars   | Flags                    | Type   | Description                           |
| ------------------ | ------------------------ | ------ | ------------------------------------- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password                      |

#### Help Options

| Environment vars | Flags        | Type | Description            |
| ---------------- | ------------ | ---- | ---------------------- |
| -                | `-h, --help` | -    | Show this help message |
```

Generated:

---

#### Application Options

| Environment vars | Flags           | Type | Description                  |
| ---------------- | --------------- | ---- | ---------------------------- |
| -                | `-v, --version` | bool | display the version and exit |
| `DEBUG`          | `-D, --debug`   | bool | enable debugging             |

#### Example Group

| Environment vars   | Flags                    | Type   | Description                           |
| ------------------ | ------------------------ | ------ | ------------------------------------- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password                      |

#### Help Options

| Environment vars | Flags        | Type | Description            |
| ---------------- | ------------ | ---- | ---------------------- |
| -                | `-h, --help` | -    | Show this help message |

---

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->

## :raising_hand_man: Support & Assistance

-   :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
    guidelines on ensuring everyone has the best experience interacting with
    the community.
-   :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
    guidelines for tips on how to ask the right questions.
-   :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/clix/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->

## :handshake: Contributing

-   :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
    on ensuring everyone has the best experience interacting with the
    community.
-   :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
    issues/a guide on submitting pull requests and helping out.
-   :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/clix/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->

## :balance_scale: License

```
MIT License

Copyright (c) 2021 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_

<!-- template:end:license -->
