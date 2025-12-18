package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringArray almacena un text[] de Postgres.
// Implementa sql.Scanner y driver.Valuer para que GORM lo pueda persistir sin convertirlo a "record".
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	// nil slice => NULL
	if a == nil {
		return nil, nil
	}
	// empty slice => {}
	if len(a) == 0 {
		return "{}", nil
	}

	parts := make([]string, 0, len(a))
	for _, item := range a {
		escaped := strings.ReplaceAll(item, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		parts = append(parts, `"`+escaped+`"`)
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *StringArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}

	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("StringArray: tipo no soportado: %T", src)
	}

	values, err := parsePostgresTextArray(raw)
	if err != nil {
		return err
	}
	*a = values
	return nil
}

func parsePostgresTextArray(raw string) ([]string, error) {
	if raw == "" {
		return []string{}, nil
	}
	if raw == "{}" {
		return []string{}, nil
	}
	if len(raw) < 2 || raw[0] != '{' || raw[len(raw)-1] != '}' {
		return nil, fmt.Errorf("StringArray: formato inválido: %q", raw)
	}

	s := raw[1 : len(raw)-1]
	if s == "" {
		return []string{}, nil
	}

	var out []string
	var buf strings.Builder
	inQuotes := false
	escaped := false

	flush := func() {
		val := buf.String()
		buf.Reset()
		// NULL sin comillas representa NULL real en Postgres; como no podemos representar NULL en []string,
		// lo convertimos a string vacía.
		if !inQuotes && val == "NULL" {
			out = append(out, "")
			return
		}
		out = append(out, val)
	}

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if escaped {
			buf.WriteByte(ch)
			escaped = false
			continue
		}

		if inQuotes {
			switch ch {
			case '\\':
				escaped = true
			case '"':
				inQuotes = false
			default:
				buf.WriteByte(ch)
			}
			continue
		}

		switch ch {
		case '"':
			inQuotes = true
		case ',':
			flush()
		default:
			buf.WriteByte(ch)
		}
	}
	flush()

	return out, nil
}

