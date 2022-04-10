<!-- template:begin:header -->
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- template:end:toc -->

## :ballot_box_with_check: TODO

- [ ] Commands/sub-commands/etc
- [ ] Custom formatting

## :gear: Usage

<!-- template:begin:goget -->
<!-- template:end:goget -->

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
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| - | `-v, --version` | bool | display the version and exit |
| `DEBUG` | `-D, --debug` | bool | enable debugging |

#### Example Group
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password |

#### Help Options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| - | `-h, --help` | - | Show this help message |
```

Generated:

------------

#### Application Options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| - | `-v, --version` | bool | display the version and exit |
| `DEBUG` | `-D, --debug` | bool | enable debugging |

#### Example Group
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password |

#### Help Options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| - | `-h, --help` | - | Show this help message |

------------

<!-- template:begin:support -->
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- template:end:license -->
