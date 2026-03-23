package cli

import (
	"flag"
	"fmt"

	"github.com/onyz1/onyzify/internal/formatter"
	"github.com/onyz1/onyzify/internal/schema"
)

func usage(sch schema.CompiledSchema, formatter formatter.Formatter, fs *flag.FlagSet) {
	out := fs.Output()

	fmt.Fprintf(out, "Usage:\n  %s [options]\n\n", fs.Name())
	fmt.Fprintln(out, "Options:")

	for _, field := range sch {
		formatter(field, out)
	}
}
