package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	STRING = "STRING"
	BOOL   = "BOOL"

	COLON     = ":"
	SEMICOLON = ";"
)

// Estructura de un Token
type Token struct {
	Type    TokenType
	Literal string
}