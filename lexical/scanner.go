package lexical

import (
	"fmt"
)


type Scanner struct {
	input    string
	position int
	line     int
	column   int 
	keywords map[string]TokenType
}

func NewScanner(input string) *Scanner {
	scanner := &Scanner{
		input:  input,
		line:   1,
		column: 1,
		keywords: map[string]TokenType{
			"int":    INT,
			"float":  FLOAT,
			"if":     IF,
			"else":   ELSE,
			"while":  WHILE,
			"for":    FOR,
			"return": RETURN,
			"print":  PRINT,
		},
	}
	return scanner
}

func (s *Scanner) currentChar() byte {
	if s.position >= len(s.input) {
		return 0
	}
	return s.input[s.position]
}

func (s *Scanner) peekChar() byte {
	if s.position+1 >= len(s.input) {
		return 0
	}
	return s.input[s.position+1]
}

func (s *Scanner) advance() {
	if s.position < len(s.input) && s.input[s.position] == '\n' {
		s.line++
		s.column = 1
	} else {
		s.column++
	}
	s.position++
}

func (s *Scanner) skipWhitespace() {
	for s.currentChar() != 0 && isWhitespace(s.currentChar()) {
		s.advance()
	}
}

func (s *Scanner) skipLineComment() {
	for s.currentChar() != 0 && s.currentChar() != '\n' {
		s.advance()
	}
}

func (s *Scanner) skipBlockComment() error {
	s.advance()
	s.advance()
	
	for s.currentChar() != 0 {
		if s.currentChar() == '*' && s.peekChar() == '/' {
			s.advance()
			s.advance()
			return nil
		}
		s.advance()
	}
	
	return fmt.Errorf("unterminated block comment")
}

func (s *Scanner) readIdentifier() string {
	start := s.position
	
	for s.currentChar() != 0 && (isAlpha(s.currentChar()) || isDigit(s.currentChar()) || s.currentChar() == '_') {
		s.advance()
	}
	
	return s.input[start:s.position]
}

func (s *Scanner) readNumber() string {
	start := s.position
	
	if s.currentChar() == '.' && isDigit(s.peekChar()) {
		s.advance()
	}
	for s.currentChar() != 0 && isDigit(s.currentChar()) {
		s.advance()
	}
	
	if s.currentChar() == '.' && isDigit(s.peekChar()) {
		s.advance() 
		
		for s.currentChar() != 0 && isDigit(s.currentChar()) {
			s.advance()
		}
	}
	
	return s.input[start:s.position]
}

func (s *Scanner) NextToken() Token {
	for s.currentChar() != 0 {
		s.skipWhitespace()
		
		if s.currentChar() == 0 {
			break
		}
		
		line := s.line
		column := s.column
		position := s.position
		
		char := s.currentChar()
		
		if char == '#' {
			s.skipLineComment()
			continue
		}

		if char == '/' {
			if s.peekChar() == '/' {
				s.skipLineComment()
				continue
			} else if s.peekChar() == '*' {
				if err := s.skipBlockComment(); err != nil {
					return Token{ERROR, err.Error(), line, column, position}
				}
				continue
			}
		}
			
		if isAlpha(char) || char == '_' {
			value := s.readIdentifier()
			tokenType := IDENTIFIER
			
			if keyword, exists := s.keywords[value]; exists {
				tokenType = keyword
				value = "" 
			}
			
			return Token{tokenType, value, line, column, position}
		}
		
		if isDigit(char) || (char == '.' && isDigit(s.peekChar())) {
			value := s.readNumber()
			return Token{NUMBER, value, line, column, position}
		}
		
		if char == '=' && s.peekChar() == '=' {
			s.advance()
			s.advance()
			return Token{EQUAL, "", line, column, position}
		}

		if char == '+' && s.peekChar() == '+' {
			s.advance()
			s.advance()
			return Token{INCREMENT, "", line, column, position}
		}

		if char == '-' && s.peekChar() == '-' {
			s.advance()
			s.advance()
			return Token{DECREMENT, "", line, column, position}
		}
		
		if char == '!' && s.peekChar() == '=' {
			s.advance()
			s.advance()
			return Token{NOT_EQUAL, "", line, column, position}
		}
		
		if char == '<' && s.peekChar() == '=' {
			s.advance()
			s.advance()
			return Token{LESS_EQUAL, "", line, column, position}
		}
		
		if char == '>' && s.peekChar() == '=' {
			s.advance()
			s.advance()
			return Token{GTE, "", line, column, position}
		}
		
		if char == '&' && s.peekChar() == '&' {
			s.advance()
			s.advance()
			return Token{AND, "", line, column, position}
		}
		
		if char == '|' && s.peekChar() == '|' {
			s.advance()
			s.advance()
			return Token{OR, "", line, column, position}
		}
		
		s.advance()
		
		switch char {
		case '=':
			return Token{ASSIGN, "", line, column, position}
		case '+':
			return Token{PLUS, "", line, column, position}
		case '-':
			return Token{MINUS, "", line, column, position}
		case '*':
			return Token{MULTIPLY, "", line, column, position}
		case '/':
			return Token{DIVIDE, "", line, column, position}
		case '%':
			return Token{MODULO, "", line, column, position}
		case '<':
			return Token{LESS, "", line, column, position}
		case '>':
			return Token{GREATER, "", line, column, position}
		case '!':
			return Token{NOT, "", line, column, position}
		case ';':
			return Token{SEMICOLON, "", line, column, position}
		case ',':
			return Token{COMMA, "", line, column, position}
		case '(':
			return Token{LPAREN, "", line, column, position}
		case ')':
			return Token{RPAREN, "", line, column, position}
		case '{':
			return Token{LBRACE, "", line, column, position}
		case '}':
			return Token{RBRACE, "", line, column, position}
		case '[':
			return Token{LBRACKET, "", line, column, position}
		case ']':
			return Token{RBRACKET, "", line, column, position}
		default:
			return Token{ERROR, fmt.Sprintf("unexpected character: '%c'", char), line, column, position}
		}
	}
	
	return Token{EOF, "", s.line, s.column, s.position}
}
		
func (s *Scanner) ScanAll() []Token {
	var tokens []Token
	
	for {
		token := s.NextToken()
		if token.Type == EOF {
			break
		}
		if token.Type == ERROR {
			fmt.Printf("Lexical error: %s\n", token.Value)
			continue
		}
		tokens = append(tokens, token)
	}
	
	return tokens
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

