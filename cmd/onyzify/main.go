package main

import (
	"context"
	"os"

	"github.com/onyz1/infonyz"
	"github.com/onyz1/onyzify"
)

func main() {
	logger := infonyz.New(&infonyz.Config{
		Backend: infonyz.Charm,
		Level:   infonyz.DebugLevel,
	}, os.Stderr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = infonyz.WithLogger(ctx, logger)

	opts := &onyzify.Options{
		Logger: logger,
		Args:   os.Args[1:],
		// Wizard: true,
		WizardOptions: onyzify.WizardOptions{
			Dst: os.Stdout,
			Src: os.Stdin,
		},
	}

	var err error
	opts, err = opts.WithSchemaFile("../../examples/basic.yaml")
	if err != nil {
		logger.Error("load schema file", infonyz.F("error", err))
		return
	}

	engine, err := onyzify.New(opts)
	if err != nil {
		logger.Error("create engine", infonyz.F("error", err))
		return
	}

	result, err := engine.Run(ctx)
	if err != nil {
		logger.Error("run engine", infonyz.F("error", err))
		return
	}

	yamlData, err := result.YAML()
	if err != nil {
		logger.Error("convert to YAML", infonyz.F("error", err))
		return
	}

	logger.Info("YAML output:\n" + string(yamlData))
}
