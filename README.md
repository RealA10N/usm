# USM - The Universal Assembly Language

[![CI](https://github.com/RealA10N/usm/actions/workflows/ci.yml/badge.svg)](https://github.com/RealA10N/usm/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/RealA10N/usm/graph/badge.svg?token=ZXVrTG9OxC)](https://codecov.io/gh/RealA10N/usm)

One Universal assembly language to rule them all.

```mermaid
graph TD;
    LEX[Lexer]
    PRS[Parser]
    SSA[Static Single Assignment];
    CP[Constant Propagation];
    DCE[Dead Code Elimination];
    RA[Register Allocation];

    subgraph "Sparse Conditional Constant Propagation"
        CP
        DCE
    end

    LEX --> PRS;
    PRS --> SSA;
    SSA --> CP;
    CP --> DCE;
    DCE --> CP;
    DCE --> RA;
```
