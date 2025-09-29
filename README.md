<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/clix?icon=material-symbols%3Aterminal&icon.height=100&layout=left&icon.color=rgba%2839%2C+132%2C+85%2C+1%29&font=1.2)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Features](#sparkles-features)
  - [TODO](#ballot_box_with_check-todo)
  - [Usage](#gear-usage)
  - [Example help output](#example-help-output)
  - [Generate Markdown](#generate-markdown)
    - [Example output](#example-output)
      - [Application Options](#application-options)
      - [Example Group](#example-group)
      - [Logging Options](#logging-options)
<!-- template:end:toc -->

## :sparkles: Features

- go-flags wrapper, that handles parsing and decoding, with additional
  helpers. Uses Go 1.18's generics to allow embedding your own custom struct.
- Auto-generated logger, and auto-provided logger flags that:
  - allows users to switch between plain-text, colored/pretty, JSON output
    (and quiet for no output).
  - Allows configuring the level dynamically.
  - Uses the built-in debug flag to automatically change the logging level.
- Built-in `--debug` flag.
- Built-in `--version` flag that provides a lot of useful features:
  - Using Go 1.18's build metadata, the ability to use that as the version
    info, automatically using VCS information if available.
  - Printing dependencies and build flags.
  - Embedding useful links (support, repo, homepage, etc) in both version
    output, and help output.
  - Colored output!
- `--generate-markdown` flag (hidden) that allows generating markdown from
  the CLI's help information (see below!).
- Uses [godotenv](github.com/joho/godotenv) to auto-load environment variables
  from `.env` files, before parsing flags.
- Many flags to enable/disable functionality to suit your needs.

## :ballot_box_with_check: TODO

- [ ] Generate commands/sub-commands/etc ([example](https://github.com/jessevdk/go-flags/pull/364))
- [ ] Custom markdown formatting

## :gear: Usage

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<!-- template:end:goget -->

Example:

```go
package main

import (
 "github.com/apex/log"
 clix "github.com/lrstanley/clix"
)

var (
  cli = &clix.CLI[Flags]{
  Links: clix.GithubLinks("github.com/lrstanley/myproject", "master", "https://mysite.com"),
 }
 logger log.Interface
)

type Flags struct {
 EnableHTTP bool   `short:"e" long:"enable-http" description:"enable the http server"`
 File       string `env:"FILE" short:"f" long:"file" description:"some file that does something"`

 SubFlags struct {
  Username string `env:"USERNAME" short:"u" long:"username" default:"admin" description:"example username"`
  Password string `env:"PASSWORD" short:"p" long:"password" description:"example password"`
 } `group:"Example Group" namespace:"example" env-namespace:"EXAMPLE"`
}

func main() {
 // Initializes cli flags, and a pre-configured logger, based off user-provided
 // flags (e.g. --log.json, --log.level, etc). Also automatically handles
 // --version, --help, etc.
 cli.Parse()
 logger = cli.Logger

 logger.WithFields(log.Fields{
  "debug":     cli.Debug,
  "file_path": cli.Flags.File,
 }).Info("hello world")
}
```

## Example help output

Below shows an example of a non-tagged revision of `myproject`. When using
Git tags, the appropriate tag should be applied as the version, rather than
`(devel)`.

```console
$ ./myproject --version
github.com/lrstanley/myproject :: (devel)
  build commit :: 2a00b14d2ff16b79ecbb2afc54f480c2c1e28172
    build date :: unknown
    go version :: go1.18.1 linux/amd64

helpful links:
      homepage :: https://myproject
        github :: https://github.com/lrstanley/myproject
        issues :: https://github.com/lrstanley/myproject/issues/new/choose
       support :: https://github.com/lrstanley/myproject/blob/master/.github/SUPPORT.md
  contributing :: https://github.com/lrstanley/myproject/blob/master/.github/CONTRIBUTING.md
      security :: https://github.com/lrstanley/myproject/security/policy

build options:
     -compiler :: gc
   CGO_ENABLED :: 1
    CGO_CFLAGS ::
  CGO_CPPFLAGS ::
  CGO_CXXFLAGS ::
   CGO_LDFLAGS ::
        GOARCH :: amd64
          GOOS :: linux
       GOAMD64 :: v1
           vcs :: git
  vcs.revision :: 2a00b14d2ff16b79ecbb2afc54f480c2c1e28172
      vcs.time :: 2022-04-26T00:08:11Z
  vcs.modified :: true

dependencies:
  h1:IVj9dxSeAC0CRxjM+AYLIKbNdCzAjUnsUjAp/td7kYo= :: ariga.io/atlas :: v0.3.8-0.20220424181913-f64001131c0e
  h1:Jlkg6X37VI/k5U02yTBB3MjKGniiBAmUGfA1TC1+dtU= :: ariga.io/entcache :: v0.0.0-20211014200019-283c566a429b
  h1:XE5df6hfIlK/YeRxY6ynRxoMKCqn2mIOcOjId8JrQN8= :: entgo.io/ent :: v0.10.2-0.20220424193633-04e0dc936be9
  h1:YB2fHEn0UJagG8T1rrWknE3ZQzWM06O8AMAatNn7lmo= :: github.com/agext/levenshtein :: v1.2.3
  h1:FHtw/xuaM8AgmvDDTI9fiwoAL25Sq2cxojnZICUU8l0= :: github.com/apex/log :: v1.9.0
  [...]
```

You can also use `./myproject --version-json` for a more programmatic approach
to the above information.

## Generate Markdown

When using **clix**, you can generate markdown for your commands by passing the
flag `--generate-markdown`, which is a hidden flag. This will print the markdown
to stdout.

```console
./<my-script> --generate-markdown > USAGE.md
```

### Example output

Raw:

```markdown
#### Application Options

| Environment vars | Flags               | Type   | Description                          |
| ---------------- | ------------------- | ------ | ------------------------------------ |
| -                | `-e, --enable-http` | bool   | enable the http server               |
| `FILE`           | `-f, --file`        | string | some file that does something        |
| -                | `-v, --version`     | bool   | prints version information and exits |
| `DEBUG`          | `-D, --debug`       | bool   | enables debug mode                   |

#### Example Group

| Environment vars   | Flags                    | Type   | Description                           |
| ------------------ | ------------------------ | ------ | ------------------------------------- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password                      |

#### Logging Options

| Environment vars | Flags          | Type   | Description                                                                      |
| ---------------- | -------------- | ------ | -------------------------------------------------------------------------------- |
| `LOG_QUIET`      | `--log.quiet`  | bool   | disable logging to stdout (also: see levels)                                     |
| `LOG_LEVEL`      | `--log.level`  | string | logging level [**default: info**] [**choices: debug, info, warn, error, fatal**] |
| `LOG_JSON`       | `--log.json`   | bool   | output logs in JSON format                                                       |
| `LOG_PRETTY`     | `--log.pretty` | bool   | output logs in a pretty colored format (cannot be easily parsed)                 |
```

Generated:

---

#### Application Options

| Environment vars | Flags               | Type   | Description                          |
| ---------------- | ------------------- | ------ | ------------------------------------ |
| -                | `-e, --enable-http` | bool   | enable the http server               |
| `FILE`           | `-f, --file`        | string | some file that does something        |
| -                | `-v, --version`     | bool   | prints version information and exits |
| `DEBUG`          | `-D, --debug`       | bool   | enables debug mode                   |

#### Example Group

| Environment vars   | Flags                    | Type   | Description                           |
| ------------------ | ------------------------ | ------ | ------------------------------------- |
| `EXAMPLE_USERNAME` | `-u, --example.username` | string | example username [**default: admin**] |
| `EXAMPLE_PASSWORD` | `-p, --example.password` | string | example password                      |

#### Logging Options

| Environment vars | Flags          | Type   | Description                                                                      |
| ---------------- | -------------- | ------ | -------------------------------------------------------------------------------- |
| `LOG_QUIET`      | `--log.quiet`  | bool   | disable logging to stdout (also: see levels)                                     |
| `LOG_LEVEL`      | `--log.level`  | string | logging level [**default: info**] [**choices: debug, info, warn, error, fatal**] |
| `LOG_JSON`       | `--log.json`   | bool   | output logs in JSON format                                                       |
| `LOG_PRETTY`     | `--log.pretty` | bool   | output logs in a pretty colored format (cannot be easily parsed)                 |

---

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<!-- template:end:license -->
