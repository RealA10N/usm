# USM - The Universal Assembly Language<a name="usm---the-universal-assembly-language"></a>

[![CI]](https://github.com/RealA10N/usm/actions/workflows/ci.yml)
[![codecov]](https://codecov.io/gh/RealA10N/usm)
[![pre-commit.ci status]](https://results.pre-commit.ci/latest/github/RealA10N/usm/main)

One Universal assembly language to rule them all.

<!-- mdformat-toc start --slug=github --maxlevel=6 --minlevel=2 -->

- [Registers](#registers)
- [Types](#types)
  - [Standard Integer Types](#standard-integer-types)
  - [Custom Type Definitions](#custom-type-definitions)
  - [Function Pointer Type Definition](#function-pointer-type-definition)
- [Functions](#functions)
- [Instructions](#instructions)
  - [Expressions](#expressions)
- [Immediate Values](#immediate-values)
  - [Integer Immediate Value](#integer-immediate-value)
  - [Character Immediate Value](#character-immediate-value)
  - [Pointer Immediate Value](#pointer-immediate-value)
- [Globals](#globals)
  - [Global Initialization](#global-initialization)

<!-- mdformat-toc end -->

## Registers<a name="registers"></a>

A register is a location that can store values of a certain type. Registers are
defined and bounded to the context of a single function. The first assignment of
the register to a value (possibly, as a function parameter) defines the register
type, and the type of the register cannot be changed afterwards. A register type
can be any valid type, and the size of the register (in bits) is unbounded.

Unlike in other, machine specific, assembly languages, the number of available
registers are not bounded by USM, and their names can be any sequence of non
whitespace[^1] unicode characters, prefixed by `%`.

Registers are not necessarily stored in memory, and thus can't be directly
dereferenced.

## Types<a name="types"></a>

Each value in USM has a distinct type. A type name is prefixed with `$`.

### Standard Integer Types<a name="standard-integer-types"></a>

For any strictly positive integer `n`, there exists a builtin standard integer
type named `$<n>` where `<n>` is the decimal representation of `n`. The `$<n>`
type is a `n` bit integer. USM does not distinguish between the signed and
unsigned values.

```usm
%boolean = $1
%integer = $32
```

### Custom Type Definitions<a name="custom-type-definitions"></a>

A custom type declaration begins with the top level token `type`. Then, follows
the new type name, prefixed with `$` and a non-digit character. After that comes
the `=` token, and then follows a list of (at least one) a type definitions.

A type definitions begins with a list of (possibly zero) field labels. A field
label is a label in the context of the type declaration only, and is prefixed
with `.`.

Then, follows the underlying type, which is `$` prefixed. Finally, there is a
list of (possibly zero) type descriptors, separated by (at least one)
whitespace.

The `*` descriptor represents a pointer. `*<n>` is a nested pointer (pointer of
a pointer of a...) exactly `n` times, where `<n>` is the decimal representation
of a strictly positive integer. If `n` is not specified, it is assumed that
`n=1`.

Similarly, the `^` descriptor represents an array. `^<n>` is an array of size
`n`, where `<n>` is the decimal representation of a strictly positive integer
`n`.

Descriptors are applied from left to right, in order. e.g. `$8 ^100 *` is a
pointer to an array of 100 bytes, and `$8 * ^100` is an array of 100 pointers.
The number of descriptors si not bounded. e.g. `$8 * ^100 *` is a pointer to an
array of pointers.

```usm
type $bool = $1

type $str = $8 *

type $person =
    .name $str
    .age $32
    .isMale $bool

type $peopleArray = $person * ^100
```

> [!NOTE]
> The type definitions above will be used in examples throughout the
> specification.

### Function Pointer Type Definition<a name="function-pointer-type-definition"></a>

If a type definition contains the `@` token, it is treaded as a function pointer
alias. The (possibly empty) type list before the `@` token represents the
function return types, and the (possibly empty) type list after the `@` token
represents the function parameter types.

```usm
type $voidOp = @                ; no parameters, no returns
type $binaryOp = $32 @ $32 $32  ; two parameters, one return
```

## Functions<a name="functions"></a>

A function declaration always begins with the top level token `func`. Then
follows a list of (possibly zero) return types, than the function global name
(`@` prefixed), and finally follows a list of (possibly zero) type and register
pairs for each parameter that the function accepts.

It is possible to declare a function without providing an implementation. In
that case, the compiler should expect to find the implementation in another
object file that should be eventually linked.

```usm
func $32 @add $32 %a $32 %b
```

An implementation can be provided be appending the `=` token after the function
declaration. Then, a list of at least one instruction is expected, separated by
at least one newline between them. The function definition ends when an `EOF`
token is reached, or another top level token is encountered.

```usm
func $32 @add $32 %a $32 %b =
    %c = add %a %b
    ret %c
```

## Instructions<a name="instructions"></a>

An instruction consists of (possibly zero) target registers, and an expression.
The return types from an expression is always known, and should match the target
register types. If some (possibly all) registers are appearing for the first
time in function, their type should be inferred from the corresponding
expression return type.

```usm
%a %b %c ... =     ...
-----┬------   -----┬------
 target(s)      expression
```

Expression return values can be assigned to the *epsilon register* `%` if the
corresponding value should be ignored. If the expression returns more values
than the amount of target registers `n`, only the first `n` values from the
expression are assigned to the target registers, and the rest of the values are
ignored. If the expression does not return any values, or all returned values
are ignored, the `=` token should be emitted.

```usm
%q, %r = divmod $32 #7 $32 #3
%q, %r = divmod %a %b

; keep quotient, ignore reminder
%q % = divmod %a %b
%q = divmod %a %b

; ignore quotient, keep reminder
% %r = divmod %a %b

; ignore both
% % = divmod %a %b
% = divmod %a %b
divmod %a %b
```

### Expressions<a name="expressions"></a>

There are two distinct expression types: an *operator expression*, and an
*immediate values expression*.

An operator expression a operator name. It is an identifier which is *not*
prefixed with a special character. Then, follows the arguments to the operation,
which are operation specific, can can be immediate values, function labels, or
registers. An operation with specific parameter types should return a
deterministic set of (possibly zero) return types, which are then assigned to
the corresponding target registers. Valid operators and their implementation are
not part of this specification, and are implementation specific.

```usm
%a %b %c ... = dosomething %a %b $32 #1234 ...
-----┬------   -----┬----- ---------┬---------
 target(s)      operator     specific params
```

In addition, a list of (at least one) type and immediate initialization pairs
can be supplied as an expression to directly initialize the registers with
immediate values.

```usm
%0 %1 = $person ... $32 ...
        -----┬----- ---┬---
          imm #0     imm #1
```

## Immediate Values<a name="immediate-values"></a>

Immediate values are used to initialize registers and globals.

### Integer Immediate Value<a name="integer-immediate-value"></a>

Initialize an integer value using the syntax `#<n><b>` where `<b>` is replaced
with the immediate base (as described below), and `<n>` is replaced with a
possibly negative integer, in the provided base representation (as described
below).

| Base             | Allowed Suffix (`<b>`) | Allowed Digits (`<n>`)      |
| ---------------- | ---------------------- | --------------------------- |
| Hexadecimal (16) | `h` or `H`             | `0`-`9`, `a`-`f` or `A`-`F` |
| Decimal (10)     | empty, `d` or `D`      | `0`-`9`                     |
| Octal (8)        | `o` or `O`             | `0`-`7`                     |
| Binary (2)       | `b` or `B`             | `0`, `1`                    |

```usm
func @main =
    %0 = $32 #-1337
    %1 = $32 #4294967295
    %2 = $64 #DEADBEEFh
    %3 = $32 #-1234567o
    %4 = $8 #100b
```

### Character Immediate Value<a name="character-immediate-value"></a>

For convenience, an initialization of integers can be also done via a unicode
character. Using the syntax `#'<c>'`, where `<c>` is replaced by a unicode
character, the immediate value will be translated to the appropriate
[unicode code point](https://en.wikipedia.org/wiki/Code_point#In_Unicode).

### Pointer Immediate Value<a name="pointer-immediate-value"></a>

A pointer type can be only explicitly initialized to the zero immediate `#0` (or
to a global with the same type).

## Globals<a name="globals"></a>

It is possible to declare a global without initialization. In that case the
initial value of the global is undefined.

```usm
glob $person @author  ; undefined value
```

### Global Initialization<a name="global-initialization"></a>

Global initialization is done by initializing the global underlying standard
types, in order of declaration of the global type. The initialization should be
provided after the declaration of the global and the `=` token.

If not all fields of the global are initialized (possibly, none), the
uninitialized fields are implicitly initialized to zero.

```usm
glob $32 @authorAge = #1337

glob $8 ^5 @authorName = #'A' #'l' #'o' #'n'

glob $person @author =
    @authorName
    @authorAge            ; .isMale is implicitly initialized to #0
```

Using type labels, it is possible to start initialize fields from a different
starting position, and skip explicit initialization of fields to zero.

```usm
glob $person @author =  ; the .name field is initialized to #0 implicitly.
    .age @authorAge       ; initialization is started from .age field
    #1                    ; and continues to the .isMale field
```

Note that it is possible to implicitly initialize all of the fields to `#0`, but
simply appending the `=` token after the global declaration.

```usm
glob $person @author =
```

If a type field is initialized more than once, the value of the whole structure
is undefined (that is, including other fields).

[^1]: A unicode whitespace character is one that has the ["WSpace=Y" property].
    For reference, see [Go's unicode.IsSpace standard function].

["wspace=y" property]: https://en.wikipedia.org/wiki/Whitespace_character#Unicode
[ci]: https://github.com/RealA10N/usm/actions/workflows/ci.yml/badge.svg
[codecov]: https://codecov.io/gh/RealA10N/usm/graph/badge.svg?token=ZXVrTG9OxC
[go's unicode.isspace standard function]: https://pkg.go.dev/unicode#IsSpace
[pre-commit.ci status]: https://results.pre-commit.ci/badge/github/RealA10N/usm/main.svg
