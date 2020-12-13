package sfv

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

func Marshal(v interface{}) (string, error) {
	var w strings.Builder

	// todo: other types
	if v, ok := v.(Item); ok {
		if err := marshalItem(&w, v); err != nil {
			return "", err
		}
	}

	if v, ok := v.([]Member); ok {
		if err := marshalList(&w, v); err != nil {
			return "", err
		}
	}

	return w.String(), nil
}

func marshalItem(w *strings.Builder, v Item) error {
	if err := marshalBareItem(w, v.Value); err != nil {
		return err
	}

	if err := marshalParams(w, v.Params); err != nil {
		return err
	}

	return nil
}

func marshalList(w *strings.Builder, v []Member) error {
	for i, m := range v {
		if m.IsItem {
			if err := marshalItem(w, m.Item); err != nil {
				return err
			}
		} else {
			if err := marshalInnerList(w, m.InnerList); err != nil {
				return err
			}
		}

		if i != len(v)-1 {
			fmt.Fprint(w, ", ")
		}
	}

	return nil
}

func marshalInnerList(w *strings.Builder, v InnerList) error {
	fmt.Fprint(w, "(")

	for i, m := range v.Items {
		if err := marshalItem(w, m); err != nil {
			return err
		}

		if i != len(v.Items)-1 {
			fmt.Fprintf(w, " ")
		}
	}

	fmt.Fprint(w, ")")

	if err := marshalParams(w, v.Params); err != nil {
		return err
	}

	return nil
}

func marshalBareItem(w *strings.Builder, v interface{}) error {
	switch v := v.(type) {
	case float64:
		return marshalDecimal(w, v)
	case int64:
		return marshalInteger(w, v)
	case string:
		return marshalString(w, v)
	case Token:
		return marshalToken(w, v)
	case []byte:
		return marshalByteSequence(w, v)
	case bool:
		return marshalBoolean(w, v)
	default:
		return fmt.Errorf("unsupported bare item type: %v", v)
	}
}

func marshalDecimal(w *strings.Builder, v float64) error {
	// todo: check precision
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.ContainsRune(s, '.') {
		s += ".0"
	}

	fmt.Fprint(w, s)
	return nil
}

func marshalInteger(w *strings.Builder, v int64) error {
	// todo: check range
	fmt.Fprintf(w, "%d", v)
	return nil
}

func marshalString(w *strings.Builder, v string) error {
	// todo: check all chars ascii
	fmt.Fprint(w, "\"")
	for _, c := range v {
		if c == '\\' || c == '"' {
			fmt.Fprintf(w, "\\%s", string(c))
		} else {
			fmt.Fprintf(w, "%s", string(c))
		}
	}
	fmt.Fprint(w, "\"")
	return nil
}

func marshalToken(w *strings.Builder, v Token) error {
	// todo: check chars ok for token
	fmt.Fprintf(w, "%s", string(v))
	return nil
}

func marshalByteSequence(w *strings.Builder, v []byte) error {
	fmt.Fprintf(w, ":%s:", base64.StdEncoding.EncodeToString(v))
	return nil
}

func marshalBoolean(w *strings.Builder, v bool) error {
	n := 0
	if v {
		n = 1
	}

	fmt.Fprintf(w, "?%d", n)
	return nil
}

func marshalParams(w *strings.Builder, v Params) error {
	for _, k := range v.Keys {
		fmt.Fprintf(w, ";")
		if err := marshalKey(w, k); err != nil {
			return err
		}

		if v.Map[k] != true {
			fmt.Fprint(w, "=")
			if err := marshalBareItem(w, v.Map[k]); err != nil {
				return err
			}
		}
	}

	return nil
}

func marshalKey(w *strings.Builder, v string) error {
	// todo: check chars ok for key
	fmt.Fprintf(w, "%s", v)
	return nil
}