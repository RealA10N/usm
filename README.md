# USM - The Universal Assembly Language

[![CI]](https://github.com/RealA10N/usm/actions/workflows/ci.yml)
[![codecov]](https://codecov.io/gh/RealA10N/usm)
[![pre-commit.ci status]](https://results.pre-commit.ci/latest/github/RealA10N/usm/main)

USM is a universal assembly language designed to bridge the gap between
high-level programming languages and machine code. It provides a type-safe,
platform-independent assembly language that can be compiled to any target
architecture.

```usm
func $64 @fib $64 %n {
    $64 %prev = $64 #0
    $64 %cur = $64 #1
.loop
    %n = sub %n $64 #1
    jz %n .end
    %prev = %cur
    %cur = add %cur %prev
    j .loop
.end
    ret %c
}
```

## Key Features

- **Type Safety**: Unlike traditional assembly languages, USM enforces type
  checking at compile time
- **Platform Independence**: Write once, compile to any target architecture
- **Modern Syntax**: Clean, readable syntax that maintains the power of assembly
- **Unlimited Registers**: No artificial limits on register count or naming

## Use Cases

USM transcends the role of a mere assembly language, functioning as a
comprehensive framework for systems programming and compiler design:

### Optimization and Transformation of Code

USM's framework enables implementing any instruction set in its syntax,
providing a powerful foundation for compiler development. When an ISA is
implemented in USM, it automatically gains access to a rich set of optimizations
including dead code elimination, liveness analysis, and SSA transformations.
Developers can define custom optimization passes or create transformations
between different instruction sets, such as compiling a virtual ISA to
hardware-specific code like x86_64.

### Enhanced Assembly Programming

When writing low-level code for existing architectures like Aarch64 or x86_64
using the USM syntax, USM offers significant advantages:

- Strong type checking prevents common assembly errors
- Modern developer tooling including formatters and linters
- Static analysis capabilities not available in traditional assemblers

The power of USM lies in its flexibility - serving equally well as a robust
assembly language and as a framework for creating compiler infrastructure
components.

## Similar Projects

- [LLVM](https://github.com/llvm/llvm-project) - A full-featured compiler
  backend
- [QBE](https://c9x.me/compile/) - A lightweight compiler backend
- [MIR](https://github.com/vnmakarov/mir) - A lightweight JIT compiler
- [MLIR](https://mlir.llvm.org/) - A multi-level intermediate representation

[ci]: https://github.com/RealA10N/usm/actions/workflows/ci.yml/badge.svg
[codecov]: https://codecov.io/gh/RealA10N/usm/graph/badge.svg?token=ZXVrTG9OxC
[pre-commit.ci status]: https://results.pre-commit.ci/badge/github/RealA10N/usm/main.svg
