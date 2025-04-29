package csv

import (
	"io"
	"strings"
)

// Parser konfiguriert das Verhalten des CSV-Parsers.
type Parser struct {
	separator               rune
	quoteChar               rune
	escapeChar              rune
	strictQuotes            bool
	ignoreLeadingWhitespace bool

	pending    string
	inField    bool
	moreTokens bool
}

// NewParser erstellt eine neue CSVParser-Instanz mit Standardwerten.
func NewParser() *Parser {
	return NewParserWithOptions(
		',',
		'"',
		'\\',
		false,
		true,
	)
}

// NewParserWithOptions erstellt eine neue CSVParser-Instanz mit benutzerdefinierten Optionen.
func NewParserWithOptions(separator rune, quoteChar rune, escapeChar rune, strictQuotes bool, ignoreLeadingWhitespace bool) *Parser {
	if separator == '\000' {
		panic("Das Trennzeichen muss definiert sein!")
	}
	if separator == quoteChar || separator == escapeChar {
		panic("Separator, QuoteChar und EscapeChar müssen unterschiedlich sein!")
	}
	if quoteChar == escapeChar && quoteChar != '\000' {
		panic("QuoteChar und EscapeChar müssen unterschiedlich sein!")
	}

	return &Parser{
		separator:               separator,
		quoteChar:               quoteChar,
		escapeChar:              escapeChar,
		strictQuotes:            strictQuotes,
		ignoreLeadingWhitespace: ignoreLeadingWhitespace,
	}
}

// IsPending gibt true zurück, wenn vom letzten Aufruf Reste vorhanden sind.
func (p *Parser) IsPending() bool {
	return p.pending != ""
}

// ParseLineMulti parst eine Zeile, möglicherweise über mehrere Zeilen hinweg.
func (p *Parser) ParseLineMulti(line string) ([]string, error) {
	return p.parseLine(line, true)
}

// ParseLine parst eine einzelne Zeile.
func (p *Parser) ParseLine(line string) ([]string, error) {
	return p.parseLine(line, false)
}

func (p *Parser) parseLine(line string, multi bool) ([]string, error) {
	p.moreTokens = false

	if !multi && p.pending != "" {
		p.pending = ""
	}

	if line == "" {
		if p.pending != "" {
			s := p.pending
			p.pending = ""
			return []string{s}, nil
		} else {
			return nil, nil
		}
	}

	tokensOnThisLine := make([]string, 0)
	sb := strings.Builder{}
	inQuotes := false

	if p.pending != "" {
		sb.WriteString(p.pending)
		p.pending = ""
		inQuotes = true
	}

	var currentRune rune
	for i, r := range line {
		currentRune = r
		if currentRune == p.escapeChar {
			if p.isNextRuneEscapable(line, inQuotes || p.inField, i) {
				sb.WriteRune(rune(line[i+1]))
				i++ // Wir haben das nächste Zeichen auch verarbeitet
			} else {
				sb.WriteRune(currentRune) // Escape-Zeichen ohne etwas zu escapen
			}
			continue
		} else if currentRune == p.quoteChar {
			if p.isNextRuneEscapedQuote(line, inQuotes || p.inField, i) {
				sb.WriteRune(rune(line[i+1]))
				i++ // Wir haben das nächste Zeichen auch verarbeitet
			} else {
				if !p.strictQuotes {
					if i > 0 && line[i-1] != byte(p.separator) && (i+1 < len(line) && line[i+1] != byte(p.separator)) {
						if p.ignoreLeadingWhitespace && sb.Len() > 0 && isAllWhitespace(sb.String()) {
							sb.Reset()
						} else {
							sb.WriteRune(currentRune)
						}
					}
				}
				inQuotes = !inQuotes
			}
			p.inField = !p.inField
		} else if currentRune == p.separator && !inQuotes {
			tokensOnThisLine = append(tokensOnThisLine, strings.TrimSpace(sb.String()))
			sb.Reset()
			p.inField = false
		} else {
			if !p.strictQuotes || inQuotes {
				sb.WriteRune(currentRune)
				p.inField = true
			}
		}
	}

	if inQuotes {
		if multi {
			sb.WriteString("\n")
			p.pending = sb.String()
			return nil, nil // Noch nicht fertig
		} else {
			return nil, io.ErrUnexpectedEOF // Fehler: Nicht abgeschlossenes Anführungszeichen am Ende der Zeile
		}
	}

	if sb.Len() > 0 || currentRune == p.separator {
		if sb.Len() == 0 && currentRune == p.separator {
			p.moreTokens = true
		}
		tokensOnThisLine = append(tokensOnThisLine, strings.TrimSpace(sb.String()))
	}

	return tokensOnThisLine, nil
}

func (p *Parser) isNextRuneEscapedQuote(line string, inQuotes bool, currentIndex int) bool {
	return inQuotes && len(line) > currentIndex+1 && rune(line[currentIndex+1]) == p.quoteChar
}

func (p *Parser) isNextRuneEscapable(line string, inQuotes bool, currentIndex int) bool {
	return inQuotes && len(line) > currentIndex+1 && (rune(line[currentIndex+1]) == p.quoteChar || rune(line[currentIndex+1]) == p.escapeChar)
}

func isAllWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func (p *Parser) HasMoreTokens() bool {
	return p.moreTokens
}
