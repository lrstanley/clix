<p align="center">goflags-markdown -- generate markdown from a [go-flags](https://github.com/jessevdk/go-flags) parser</p>
<p align="center">
  <a href="https://pkg.go.dev/github.com/lrstanley/goflags-markdown"><img src="https://pkg.go.dev/badge/github.com/lrstanley/goflags-markdown" alt="pkg.go.dev"></a>
  <a href="https://github.com/lrstanley/goflags-markdown/actions"><img src="https://github.com/lrstanley/goflags-markdown/workflows/test/badge.svg" alt="test status"></a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/goflags-markdown"><img src="https://goreportcard.com/badge/github.com/lrstanley/goflags-markdown" alt="goreportcard"></a>
  <a href="https://gocover.io/github.com/lrstanley/goflags-markdown"><img src="http://gocover.io/_badge/github.com/lrstanley/goflags-markdown" alt="gocover"></a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/Community-Chat%20with%20us-green.svg" alt="Community Chat"></a>
</p>

## TODO

- [ ] Commands/sub-commands/etc
- [ ] Custom formatting

## Usage

```go
package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
	goflagsmarkdown "github.com/lrstanley/goflags-markdown"
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

	goflagsmarkdown.Generate(parser, os.Stdout)
	os.Exit(0)
}
```

### Example output

Raw:

```markdown
#### Application Options
| Environment vars | Flags | Description |
| --- | --- | --- |
|  | `-v, --version` | display the version and exit |
| `DEBUG` | `-D, --debug` | enable debugging |

#### Example Group
| Environment vars | Flags | Description |
| --- | --- | --- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | example password |

#### Help Options
| Environment vars | Flags | Description |
| --- | --- | --- |
|  | `-h, --help` | Show this help message |
```

Generated:

------------

#### Application Options
| Environment vars | Flags | Description |
| --- | --- | --- |
|  | `-v, --version` | display the version and exit |
| `DEBUG` | `-D, --debug` | enable debugging |

#### Example Group
| Environment vars | Flags | Description |
| --- | --- | --- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | example password |

#### Help Options
| Environment vars | Flags | Description |
| --- | --- | --- |
|  | `-h, --help` | Show this help message |

------------

## Contributing

Please review the [CONTRIBUTING](CONTRIBUTING.md) doc for submitting issues/a guide
on submitting pull requests and helping out.


## License

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
