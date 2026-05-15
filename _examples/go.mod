module github.com/lrstanley/clix/v2/_examples

go 1.26.0

replace github.com/lrstanley/clix/v2 => ../

require (
	github.com/alecthomas/kong v1.15.0
	github.com/lrstanley/clix/v2 v2.0.0
	github.com/lrstanley/x/sync v0.0.0-20260514072400-85ffc850d28d
)

require github.com/lmittmann/tint v1.1.3 // indirect
