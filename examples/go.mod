module github.com/lrstanley/clix/v2/examples

go 1.25.1

replace github.com/lrstanley/clix/v2 => ../

require (
	github.com/alecthomas/kong v1.12.1
	github.com/lrstanley/clix/v2 v2.0.0-alpha.1
	github.com/lrstanley/x/scheduler v0.0.0-20250929040648-16084399ed30
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lmittmann/tint v1.1.2 // indirect
	github.com/lrstanley/x/logging/handlers v0.0.0-20250929040648-16084399ed30 // indirect
	golang.org/x/sync v0.17.0 // indirect
)
