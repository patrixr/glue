package runtime

import "fmt"

func AssertValidSymbolName(name string) {
	if err := ValidSymbolName(name); err != nil {
		panic(err)
	}
}

func ValidSymbolName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("Field name cannot be empty")
	}

	if name[0] >= '0' && name[0] <= '9' {
		return fmt.Errorf("Field name cannot start with a digit")
	}

	for _, char := range name {
		if !(char == '_' || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return fmt.Errorf("Field name \"%s\" contains invalid characters (%c)\n", name, char)
		}
	}

	return nil
}
