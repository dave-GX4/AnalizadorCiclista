package parser

import (
	"compilerciclista/src/lexer"
	"compilerciclista/src/token"
	"fmt"
	"strconv"
)

type ParticipantData map[string]interface{}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram es el punto de entrada principal
func (p *Parser) ParseProgram() (ParticipantData, []string) {
	data := make(ParticipantData)

	for p.curToken.Type != token.EOF {
		key, value := p.parseStatement()
		if key != "" {
			data[key] = value
		}
		p.nextToken()
	}

	return data, p.errors
}

// Parsea una declaración completa, ej: nombre: "Juan";
func (p *Parser) parseStatement() (string, interface{}) {
	// Debe empezar con un identificador (la clave)
	if p.curToken.Type != token.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("Error de sintaxis: se esperaba un identificador, se obtuvo %s", p.curToken.Literal))
		return "", nil
	}
	key := p.curToken.Literal

	// Después debe venir un ':'
	if !p.expectPeek(token.COLON) {
		return "", nil
	}

	p.nextToken() // Avanzamos al valor

	// Parseamos el valor
	var value interface{}
	switch p.curToken.Type {
	case token.STRING:
		value = p.curToken.Literal
	case token.BOOL:
		boolValue, err := strconv.ParseBool(p.curToken.Literal)
		if err != nil {
			p.errors = append(p.errors, fmt.Sprintf("Error de sintaxis: no se pudo convertir '%s' a booleano", p.curToken.Literal))
			return "", nil
		}
		value = boolValue
	default:
		p.errors = append(p.errors, fmt.Sprintf("Error de sintaxis: se encontró un tipo de valor no válido %s", p.curToken.Type))
		return "", nil
	}

	// La declaración debe terminar con ';'
	if !p.expectPeek(token.SEMICOLON) {
		return "", nil
	}

	return key, value
}

// expectPeek revisa el tipo del siguiente token. Si es correcto, avanza.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Error de sintaxis: se esperaba el token %s, se obtuvo %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}