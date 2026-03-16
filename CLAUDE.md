# CLAUDE.md — USM Compiler Framework

This file describes the codebase structure, development workflows, and conventions
for AI assistants working in this repository.

## Project Overview

**USM** is a universal, type-safe assembly language and compiler framework written
in Go. It provides:

- A typed assembly language (`.usm` source files)
- A multi-pass compilation pipeline (parsing → optimization → code generation)
- Extensible ISA and optimization infrastructure
- Current target: AArch64/ARM64 (Mach-O object files for macOS)

## Repository Structure

```
usm/
├── core/              # Result/error system, source views, shared types
├── lex/               # Lexer — tokenizes .usm source into typed tokens
├── parse/             # Parser — builds AST from token stream
├── gen/               # Code generation framework (generic, ISA-agnostic)
├── transform/         # Pipeline definitions: Target, Transformation, TargetData
├── opt/               # Optimization passes (DCE, liveness analysis)
├── usm/               # USM ISA definitions, managers, SSA construction
│   ├── isa/          # Instruction definitions (add, sub, j, jz, phi, …)
│   ├── managers/     # Context managers for USM compilation
│   └── ssa/          # SSA construction pass
├── aarch64/           # AArch64 backend
│   ├── isa/          # AArch64 instruction definitions
│   ├── codegen/      # Machine code emission
│   ├── managers/     # AArch64 context managers and register info
│   └── translation/  # USM → AArch64 translation + Mach-O output
├── examples/          # Example .usm files (fib, add, loops, dead-code, …)
├── justfile           # Task runner (build, test, fmt, cover)
├── .pre-commit-config.yaml
├── go.mod / go.sum
└── README.md
```

## Build & Development Commands

All common tasks are managed with [just](https://github.com/casey/just).

| Command       | Description                                   |
| ------------- | --------------------------------------------- |
| `just build`  | Compile to `usm.out` binary                   |
| `just test`   | Run all tests (uses `richgo` if available)    |
| `just cover`  | Run tests and generate `coverage.out`         |
| `just fmt`    | Format code (`go fmt`, `go mod tidy`, mdformat) |
| `just setup`  | Install dev tools (`richgo`, `mdformat`)      |
| `just cloc`   | Count lines of code (non-test files)          |

Before committing, run `just fmt` and `just test`.

## CLI Usage

```
usm <input.usm> [transformation...]
```

Available transformations (in order):

| Name                               | Effect                              |
| ---------------------------------- | ----------------------------------- |
| `static-single-assignment` / `ssa` | Convert to SSA form                 |
| `constant-propagation` / `cp`      | Propagate and fold constants        |
| `dead-code-elimination` / `dce`    | Remove unused instructions          |
| `aarch64` / `arm64`                | Translate USM → AArch64 assembly    |
| `macho` / `macho-obj`              | Emit Mach-O `.o` object file        |

Typical pipeline: `usm input.usm ssa cp dce aarch64 macho`

## Architecture: Key Patterns

### Hierarchical Context Managers

Code generation uses four nested context levels:

1. `GenerationContext` — global (target architecture, globals)
2. `FileGenerationContext` — per-file (type registry, global manager)
3. `FunctionGenerationContext` — per-function (registers, labels, instructions)
4. `InstructionGenerationContext` — per-instruction (argument access)

Each level exposes typed managers (e.g., `RegisterManager`, `LabelManager`,
`TypeManager`). New ISA backends must implement manager interfaces.

### Instruction Traits

Instructions declare capabilities via trait interfaces:

- `CriticalInstruction` — never eliminated by DCE (e.g., `ret`, `call`)
- `NonBranchingInstruction` — falls through to the next instruction
- `Uses(i int)` — returns the i-th used value (for liveness analysis)
- `Defines()` — returns the defined value (for SSA/DCE)
- `PropagateConstants(info) []ConstantDefinition` — for constant propagation:
  returns (register, immediate) pairs that are known-constant after this
  instruction executes; embed `opt.PropagatesNoConstants` as the default no-op

### Result/Error System (`core`)

All errors propagate as `core.ResultList`. Each `Result` contains one or more
`ResultDetail` entries with:

- Severity: `Error`, `Warning`, `Hint`
- Source location (file + byte offset)
- Human-readable message

Special constructors: `core.InternalErrorResult(...)`, `core.DebugResult(...)`.

### Pipeline (`transform`)

`Target` → `Transformation` → `Target` forms a linear pipeline. Each
`Transformation` receives `TargetData` carrying the compiled IR and returns
`TargetData` in the next target's format. Errors from any step abort the
pipeline.

## Testing Conventions

- Test files: `*_test.go`, co-located with the package under test.
- Framework: Go standard `testing` + `testify/assert` and `testify/require`.
- Run with `just test` (or `just cover` for coverage).
- Integration-style tests live in `aarch64/translation/` and produce machine
  code verified via test-specific expected byte sequences.
- Representative test names: `TestAddOne`, `TestSimpleFunctionGeneration`,
  `TestIfElseFunctionGeneration`.

When adding new instructions or passes, add a corresponding test in the same
package.

## Code Style & Conventions

### Go Conventions

- Standard Go naming: `CamelCase` for exported identifiers, `camelCase` for
  unexported.
- Constructor functions: `NewXxx(...)` pattern throughout (e.g.,
  `NewTokenizer()`, `NewGenerationContext()`).
- One primary type per file; helpers and methods may live in the same file.
- File names use underscores for multi-word: `token_view.go`, `type_field.go`.

### Package Responsibilities

- `core` — only shared primitives; no imports from other `usm/*` packages.
- `lex` — no AST types; only tokens and positions.
- `parse` — produces AST nodes; no codegen or ISA knowledge.
- `gen` — ISA-agnostic; no `usm/` or `aarch64/` imports.
- `opt` — operates on `gen` IR only; ISA-agnostic.
- ISA packages (`usm/`, `aarch64/`) import `gen` but not each other.

### Error Messages

- Include source location wherever possible.
- Provide `Hint` details suggesting how to fix the issue.
- Use `levenshtein`-based suggestions for unknown identifiers/types.

### Formatting

- Go: `gofmt` (enforced by pre-commit hook).
- Markdown: `mdformat` with GFM plugin, 80-char line wrap.
- Run `just fmt` before committing.

## Adding a New ISA Backend

1. Create a top-level package `<arch>/`.
2. Implement instruction definitions in `<arch>/isa/` satisfying the `gen`
   instruction interfaces (including relevant traits).
3. Implement context managers in `<arch>/managers/` satisfying `gen` manager
   interfaces.
4. Implement a translation pass in `<arch>/translation/` from the USM
   `TargetData` to your architecture's `TargetData`.
5. Register a `transform.Target` and `transform.Transformation` and wire them
   into the CLI in `usm.go`.

## Adding a New Optimization Pass

1. Create a file in `opt/`.
2. The pass receives a `*gen.FunctionInfo` (or similar IR).
3. Implement liveness/dataflow analysis if needed — see `opt/dead_code_elimination.go`
   and `opt/constant_propagation.go` as references.
4. Passes must run **after SSA construction** if they rely on SSA properties.
5. Register the pass as a `transform.Transformation`.

### Constant Propagation Pass (`opt/constant_propagation.go`)

Propagates known-constant registers to their use sites and folds constant
expressions. Uses a DFS over the CFG (`ControlFlowGraph.Dfs(0).Timeline`) with
a per-register reaching-constants stack (mirrors `opt/ssa.ReachingDefinitionsSet`):

- Only tracks registers with exactly **one reachable definition**. Unreachable
  definitions (dead code, isolated loops) do not count. Registers with multiple
  reachable definitions are propagated only when all those definitions are
  sequential redefinitions in the same basic block (later definitions shadow
  earlier ones on the DFS stack). Diamond joins and cross-block redefinitions
  are never propagated.
- After substituting arguments, calls `PropagateConstants` on the instruction to
  fold constant expressions (e.g. `add #2 #3 → #5`). Folded results propagate
  to downstream uses in the same DFS scope.
- Uses CFG DFS (not dominator-tree DFS) to correctly handle unreachable blocks:
  Lengauer-Tarjan assigns `PreOrder=0` to unvisited nodes, which can corrupt the
  dominator tree when unreachable blocks are present.

To support a new instruction in CP: implement `PropagateConstants` or embed
`opt.PropagatesNoConstants`. Binary arithmetic helpers (`foldBinaryConstants`) live
in `usm/isa/binary_calculation.go`.

## CI/CD

GitHub Actions (`.github/workflows/ci.yml`):

1. Setup Go 1.23
2. Install `just` + run `just setup`
3. Run `just cover`
4. Upload coverage to Codecov

All PRs must pass CI. Run `just cover` locally before pushing.

## Go Module

Module path: `alon.kr/x/usm`
Go version: 1.23.0

Key dependencies:

| Package                         | Purpose                              |
| ------------------------------- | ------------------------------------ |
| `alon.kr/x/aarch64codegen`      | AArch64 machine code encoding        |
| `alon.kr/x/aarch64-macho`       | Mach-O object file generation        |
| `alon.kr/x/faststringmap`       | Fast string lookup                   |
| `alon.kr/x/list`, `set`, `stack`, `graph`, `view` | Generic data structures |
| `spf13/cobra`                   | CLI framework                        |
| `fatih/color`                   | Colored terminal output              |
| `agnivade/levenshtein`          | Edit distance for suggestions        |
| `stretchr/testify`              | Test assertions                      |

## USM Language Quick Reference

```usm
; Function definition: func <return-type> @<name> <param-type> %<param> { ... }
func $64 @fib $64 %n {
    $64 %prev = $64 #0    ; typed register assignment
    $64 %cur  = $64 #1
.loop                      ; label
    %n = sub %n $64 #1
    jz %n .end             ; conditional branch
    $64 %next = add %cur %prev
    %prev = %cur
    %cur = %next
    j .loop                ; unconditional branch
.end
    ret %cur
}
```

Token sigils: `%` register, `$` type/size, `@` global/function, `#` immediate,
`.` label, `*` pointer, `^` repeat.
