module github.com/lrstanley/clix/v2/_examples

go 1.27.0

replace github.com/lrstanley/clix/v2 => ../

require (
	github.com/alecthomas/kong v1.16.1
	github.com/lrstanley/clix/v2 v2.0.5
	github.com/lrstanley/x/sync v0.0.0-20261027055443-d021fec05428
)

require github.com/lmittmann/tint v1.2.1 // indirect
