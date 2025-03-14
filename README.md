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

## Similar Projects

- [LLVM](https://github.com/llvm/llvm-project) - A full-featured compiler
  backend
- [QBE](https://c9x.me/compile/) - A lightweight compiler backend
- [MIR](https://github.com/vnmakarov/mir) - A lightweight JIT compiler
- [MLIR](https://mlir.llvm.org/) - A multi-level intermediate representation

[ci]: https://github.com/RealA10N/usm/actions/workflows/ci.yml/badge.svg
[codecov]: https://codecov.io/gh/RealA10N/usm/graph/badge.svg?token=ZXVrTG9OxC
[pre-commit.ci status]: https://results.pre-commit.ci/badge/github/RealA10N/usm/main.svg
