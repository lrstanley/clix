module github.com/lrstanley/clix/v2/examples

go 1.25.1

replace github.com/lrstanley/clix/v2 => ../

require (
	github.com/alecthomas/kong v1.12.1
	github.com/lrstanley/clix/v2 v2.0.0-alpha.1
	github.com/lrstanley/x/scheduler v0.0.0-20251011064715-0e7e5a75ccaa
)

require github.com/lmittmann/tint v1.1.2 // indirect
