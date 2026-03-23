package wizard

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/onyz1/onyzify/internal/formatter"
	inio "github.com/onyz1/onyzify/internal/io"
	"github.com/onyz1/onyzify/internal/schema"
)

func Run(sch schema.CompiledSchema, formatter formatter.Formatter, dst io.Writer, src io.Reader) (inio.Inputs, error) {
	inputs := inio.New(len(sch))

	for fieldName, field := range sch {
		input := &inio.Input{
			Name:  fieldName,
			Type:  field.Type,
			Value: field.Default,
		}
		inputs[fieldName] = input

		reader := bufio.NewReader(src)

		var isSet bool
		for {
			var promptBuilder strings.Builder
			formatter(field, &promptBuilder)

			fmt.Fprint(dst, promptBuilder.String())

			userInput, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("read user input: %w", err)
			}

			userInput = strings.TrimSpace(userInput)

			if userInput == "" {
				break
			}

			err = input.Set(userInput)
			if err != nil {
				fmt.Fprintf(dst, "Invalid input: %v. Please try again.\n", err)
				continue
			}

			isSet = true

			break
		}

		if err := field.CheckVal(&input.Value, isSet); err != nil {
			return nil, fmt.Errorf("field: %q: check value validity: %w", fieldName, err)
		}
	}

	return inputs, nil
}
