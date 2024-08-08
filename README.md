# USM - The Universal Assembly Language

[![CI](https://github.com/RealA10N/usm/actions/workflows/ci.yml/badge.svg)](https://github.com/RealA10N/usm/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/RealA10N/usm/graph/badge.svg?token=ZXVrTG9OxC)](https://codecov.io/gh/RealA10N/usm)

One Universal assembly language to rule them all.

## Registers

A register is a location that can store values of a certain type.
Registers are defined and bounded to the context of a single function.
The first assignment of the register to a value (possibly, as a function parameter)
defines the register type, and the type of the register cannot be changed
afterwards. A register type can be any valid type, and the size of the register
(in bits) is unbounded.

Unlike in other, machine specific, assembly languages, the number of available
registers are not bounded by USM, and their names can be any sequence of non
whitespace[^1] unicode characters, prefixed by `%`.

[^1]: A unicode whitespace character is one that has the ["WSpace=Y" property](https://en.wikipedia.org/wiki/Whitespace_character#Unicode). For reference, see [Go's unicode.IsSpace standard function](https://pkg.go.dev/unicode#IsSpace).

Registers are not necessarily stored in memory, and thus can't be directly
dereferenced.

## Functions

A function declaration always begins with the top level token `func`.
Then follows a list of (possibly zero) return types, than the function global
name (`@` prefixed), and finally follows a list of (possibly zero) type and
register pairs for each parameter that the function accepts.

It is possible to declare a function without providing an implementation.
In that case, the compiler should expect to find the implementation in another
object file that should be eventually linked.

```usm
func $32 @add $32 %a $32 %b
```

An implementation can be provided be appending the `=` token after the function
declaration. Then, a list of at least one instruction is expected, separated
by at least one newline between them. The function definition ends when an `EOF`
token is reached, or another top level token is encountered.

```usm
func $32 @add $32 %a $32 %b =
    %c = add %a %b
    ret %c
```

## Instructions

An instruction consists of (possibly zero) target registers, and an expression.
The return types from an expression is always known, and should match the target
register types. If some (possibly all) registers are appearing for the first
time in function, their type should be inferred from the corresponding
expression return type.

```
%a %b %c ... =     ...
-----┬------   ------------
 target(s)      expression
```

Expression return values can be assigned to the *epsilon register* `%` if the
corresponding value should be ignored. If the expression returns more values
than the amount of target registers `n`, only the first `n` values from the
expression are assigned to the target registers, and the rest of the values
are ignored. If the expression does not return any values, or all returned
values are ignored, the `=` token should be emitted.

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

### Expressions

There are two distinct expression types: an *operator expression*, and an
*immediate values expression*.

An operator expression a operator name. It is an identifier which is *not*
prefixed with a special character. Then, follows the arguments to the operation,
which are operation specific, can can be immediate values, function labels,
or registers. An operation with specific parameter types should return a
deterministic set of (possibly zero) return types, which are then assigned to
the corresponding target registers.

```usm
%a %b %c ... = dosomething %a %b $32 #1234 ...
-----┬------   -----┬----- ---------┬---------
 target(s)        opr id     specific params
```

In addition, a list of (at least one) type and immediate initialization pairs
can be supplied as an expression to directly initialize the registers with
immediate values.

```usm
%0 %1 = $person ... $32 ...
        ----------- -------
          imm #0     imm #1
```

## Immediate Values

Immediate values are used to initialize registers and globals.

### Integer Immediate Value

Initialize an integer value using the syntax `#<n><b>` where `<b>` is replaced
with the immediate base (as described below), and `<n>` is replaced with a
possibly negative integer, in the provided base representation (as described
below).

| Base             | Allowed Suffix (`<b>`) | Allowed Digits (`<n>`)      |
|------------------|------------------------|-----------------------------|
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

### Character Immediate Value

For convenience, an initialization of integers can be also done via a unicode
character. Using the syntax `#'<c>'`, where `<c>` is replaced by a unicode
character, the immediate value will be translated to the appropriate [unicode
code point](https://en.wikipedia.org/wiki/Code_point#In_Unicode).

### Pointer Immediate Value

A pointer type can be only explicitly initialized to the zero immediate `#0`
(or to a global with the same type).

## Globals

It is possible to declare a global without initialization. In that case the
initial value of the global is undefined.

```usm
glob $person @author  ; undefined value
```

### Global Initialization

Global initialization is done by initializing the global underlying standard types,
in order of declaration of the global type. The initialization should be provided after the declaration of the global and the `=` token.

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

Note that it is possible to implicitly initialize all of the fields to `#0`,
but simply appending the `=` token after the global declaration.

```usm
glob $person @author =
```

If a type field is initialized more than once, the value of the whole structure
is undefined (that is, including other fields).
