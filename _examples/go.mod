module github.com/lrstanley/clix/v2/_examples

go 1.26.0

replace github.com/lrstanley/clix/v2 => ../

require (
	github.com/alecthomas/kong v1.14.0
	github.com/lrstanley/clix/v2 v2.0.0
	github.com/lrstanley/x/sync v0.0.0-20260430041455-3dc196cd51f1
)

require github.com/lmittmann/tint v1.1.3 // indirect
