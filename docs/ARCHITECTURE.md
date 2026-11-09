# Architecture

## 1. The Judge
This library holds the judgement for a module definition. It ensures that an allowed action declares its rule and refuses declarations where `Rules` and `Allow` disagree.

## 2. Who Calls It
- **definitionc:** Calls this library when it generates code.
- **crud:** Calls this library when it mounts the module.

## 3. Separation of Concerns
The `model` library holds only the data a module declares (`Rules`, `MFA`, `ModuleDefinition`). This library is purely the logic/checker (the judge).

## 4. Design Characteristics
- Isomorphic: works everywhere (no build tags).
- No maps are used.
- Returns translated webtyp/fmt errors.
