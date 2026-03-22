# onyzify

Stop writing CLI boilerplate.

`onyzify` turns a simple YAML schema into a fully functional CLI with built-in validation and an interactive wizard.

No flag wiring. No validation code. No duplicated logic.


`onyzify` lets you define your CLI inputs once using a simple YAML schema - and automatically get:

* CLI flags
* Interactive wizard prompts
* Type-safe parsing
* Validation
* Structured output (YAML)

No more manually wiring flags, parsing values, or validating inputs.

---

## Why this exists

Building CLIs in Go usually means:

* Writing repetitive flag parsing code
* Manually validating inputs
* Maintaining separate logic for interactive prompts (if any)

`onyzify` removes all of that.

Define a schema → get a fully working CLI + wizard instantly.

---

## Installation

### As a CLI tool (binary)

```bash
go install github.com/onyz1/onyzify/cmd/onyzify@latest
```

---

### As a library

```bash
go get github.com/onyz1/onyzify
```

---

## Usage

There are **two main ways** to use `onyzify`:

---

### 1. Ship it as a builder (recommended for tools)

Embed a schema in your project and compile a binary that:

* Guides users via CLI or wizard
* Generates valid config/output

Think of it as:

> a smarter, typed, user-friendly replacement for Makefiles or custom scripts

---

### 2. Use it as a library

Use the engine directly in your Go code if you want:

* Custom wizard formatting
* Full control over execution
* No YAML if you build schema programmatically

---

## Example Schema

```yaml
students:
  type: "[]string"
  required: true
  description: "A list of student names."

ages:
  type: "[]int"
  required: true
  description: "A list of student ages."

isComputerScienceMajor:
  type: "[]bool"
  required: true
  description: "A list indicating whether each student is a computer science major."

teacherName:
  type: "string"
  required: false
  default: "Dr. Smith"
  description: "The name of the teacher."
```

---

## Example (minimal usage)

```go
opts := &onyzify.Options{
    Args:   os.Args[1:],
    Wizard: true,
    WizardOptions: onyzify.WizardOptions{
        Dst: os.Stdout,
        Src: os.Stdin,
    },
}

opts, _ = opts.WithSchemaFile("schema.yaml")

engine, _ := onyzify.New(opts)
result, _ := engine.Run(context.Background())

yamlData, _ := result.YAML()
fmt.Println(string(yamlData))
```

---

## What you get automatically

From the schema above, `onyzify` generates:

* CLI flags:

  ```bash
  -students
  -ages
  -isComputerScienceMajor
  -teacherName
  ```

* OR interactive wizard:

  ```
  Field: students (comma separated)
  Type: []string
  Required: true
  Enter value:
  ```

* Validation:

  * Required fields enforced
  * Type-safe parsing
  * Enum validation (if defined)

---

## Contributing

Contributions are **very welcome**.

If you:

* find bugs
* want to improve DX/UX
* want to add new types/features
* or just want to clean things up

Open a PR or issue.

The project is intentionally designed to be extensible - so feel free to build on top of it.

---

## TODO (Roadmap)

Planned improvements and good starting points for contributors:

* Add more types such as timestamps and other things YAML supports.
* Add support for different configuration schemas such as JSON, TOML etc.
* Add support for custom validation rules.
* Add support for object values e.g.:

  ```yaml
  user:
    type: object
    fields:
      name:
        type: string
      age:
        type: int
  ```
* Add support for only verifying a configuration file with the schema rules (currently only building is supported).
* Currently `TypeList`s are supported but their `Default` and `Enum` values are not fully implemented.

---

## Nerd Section (Architecture & Design)

If you want to understand or contribute, this matters.

### Core idea

Everything revolves around:

```
Schema → Compile → Inputs → Validation → Output
```

---

### Key components

#### 1. Schema layer (`internal/schema`)

* Parses YAML into structured fields
* Compiles into strongly-typed `CompiledField`
* Handles:

  * type parsing
  * default values
  * enum validation

---

#### 2. Type system (`internal/types`)

* Custom enum-based type system
* Supports primitives + list variants
* Explicit mapping (`string → Type → runtime value`)

---

#### 3. Value system (`internal/value`)

* Centralized typed value container
* Handles:

  * parsing from string
  * equality checks
  * zero checks
  * conversion to `any`

This is what makes validation and CLI parsing consistent.

---

#### 4. CLI layer (`internal/cli`)

* Dynamically builds flags from schema
* Uses `flag.Value` interface for custom parsing
* Handles:

  * parsing
  * required checks
  * enum validation

---

#### 5. Wizard layer (`internal/wizard`)

* Interactive alternative to CLI flags
* Uses a pluggable formatter
* Same validation logic as CLI

---

#### 6. Engine

The `Engine` ties everything together:

* Loads schema
* Chooses CLI or wizard mode
* Produces a unified result

---

### Design philosophy

* **Schema-first**: everything is derived from a single source of truth
* **Strong typing**: avoid `interface{}` chaos
* **Separation of concerns**: each package does one thing
* **Extensibility**: easy to add new types, formats, or behaviors
* **Consistency**: CLI and wizard share the same validation pipeline

---

## Why not Cobra / Viper?

Tools like Cobra CLI and Viper are great - but they solve a *different problem*.

They give you **building blocks**.
`onyzify` gives you a **complete system**.

### With Cobra/Viper you:

* Manually define every flag
* Manually parse values
* Manually validate inputs
* Manually build interactive flows (if needed)
* Maintain CLI + config logic separately

### With onyzify you:

* Define a schema once
* Get CLI flags automatically
* Get an interactive wizard automatically
* Get validation automatically
* Keep everything in sync by design

---

### The real difference

| Feature                | Cobra/Viper  | onyzify        |
| ---------------------- | ------------ | -------------- |
| CLI flags              | Manual       | Auto-generated |
| Validation             | Manual       | Built-in       |
| Interactive wizard     | Not included | Built-in       |
| Schema-driven          | ❌            | ✅              |
| Single source of truth | ❌            | ✅              |

---

## Real-world use case

### Shipping a configurable app to non-technical users

Let’s say you built:

* a backend service
* a CLI tool
* a local automation tool

You want:

* Developers → full control via CLI flags
* Non-technical users → simple guided setup

Normally, you’d need:

* CLI flags (Cobra)
* Config parsing (Viper)
* A custom interactive setup script (manual work)

That’s messy and duplicated.

---

### With onyzify

You ship a binary with a schema:

```yaml
port:
  type: int
  required: true

env:
  type: string
  enum: ["dev", "prod"]
  default: "dev"
```

Now users can:

#### Option 1: CLI (power users)

```bash
./app -port=8080 -env=prod
```

#### Option 2: Wizard (non-technical users)

```
Field: port
Type: int
Required: true
Enter value:
```

Same schema. Same validation. No extra code.

---

### Why this matters

You’re solving a real pain:

> “How do I make my app both configurable *and* easy to use?”

Most tools force you to pick one.

`onyzify` gives you both.

---

## Why I built this

I build apps.

And when you build apps, two things matter:

* **Configuration & flexibility**
* **Ease of use**

The problem is - these usually conflict.

If you make something configurable:

* You expose flags
* You expect users to understand CLI usage
* It becomes unfriendly for non-technical users

If you make it easy:

* You hide configuration
* You lose flexibility
* Power users get frustrated

---

I wanted both.

I wanted to ship a program where:

* A non-technical user can run it and just follow prompts
* A developer can fully configure it via CLI flags
* Both paths use the same logic, same validation, same structure

And I didn’t want to:

* Write flag parsing every time
* Duplicate validation logic
* Maintain separate CLI + interactive flows

So I built `onyzify`.

---

At its core, it’s simple:

> Define your inputs once -> get everything else for free.

---

## License

Apache License 2.0
