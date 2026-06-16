# resourcename (C)

Google [AIP-122](https://aip.dev/122)-style **resource name** templates with
`{placeholder}` segments: compile a template, **parse** a full name into its
components, and **generate** a name from components. This is the C port of the
Go `Template` core.

Each `{placeholder}` matches exactly one path segment (`[^/]+`), so generated
values may not contain `/`.

## Layout

```text
c/
├── lib/                 # the library (compile with -Ilib)
│   ├── resourcename.h   # public API
│   ├── internal.h       # private shared declarations
│   ├── util.c           # string helpers, status messages
│   ├── compile.c        # rn_template_compile
│   ├── template.c       # rn_template_free + accessors
│   ├── generate.c       # rn_generate / rn_generate_pairs
│   └── parse.c          # rn_parse / rn_parse_pairs
├── example.c
├── test.c
└── Makefile
```

## Structs (recommended)

C has no reflection, so a struct is paired with a table of `rn_field`
descriptors (the C analog of Go's struct tags). Mapped fields must be `char *`;
fill the offset with `offsetof`.

```c
#include "resourcename.h"
#include <stddef.h>

typedef struct { char *id; char *name; } Artist;

static const rn_field ARTIST_FIELDS[] = {
    { "id",   offsetof(Artist, id) },
    { "name", offsetof(Artist, name) },
};

rn_template *t = NULL;
rn_template_compile("//music.example.com/artists/{id}/{name}", &t);

/* generate: struct -> resource name */
Artist a = { .id = "ar-42", .name = "Radiohead" };
char *name = NULL;
rn_generate(t, &a, ARTIST_FIELDS, 2, &name);   /* "//music.example.com/artists/ar-42/Radiohead" */
free(name);

/* parse: resource name -> struct (fields are heap-allocated) */
Artist parsed = {0};
rn_parse(t, "//music.example.com/artists/ar-42/Radiohead", &parsed, ARTIST_FIELDS, 2);
rn_parse_free(&parsed, ARTIST_FIELDS, 2);

rn_template_free(t);
```

## Key/value API (primitive)

`rn_generate` / `rn_parse` are built on lower-level pair calls you can use
directly when you don't have a struct:

```c
rn_pair in[] = { {"id", "ar-42"}, {"name", "Radiohead"} };
char *name = NULL;
rn_generate_pairs(t, in, 2, &name);
free(name);

rn_pair *out = NULL;
size_t n = 0;
rn_parse_pairs(t, name, &out, &n);
rn_pairs_free(out, n);
```

Every call returns an `rn_status`; use `rn_status_message()` for a description.
Callers own the strings from `rn_generate` / `rn_generate_pairs` (free with
`free()`), `rn_parse` (free with `rn_parse_free()`), and `rn_parse_pairs` (free
with `rn_pairs_free()`).

## API Reference

### Compilation

- `rn_status rn_template_compile(const char *template, rn_template **out)` —
  Compile a template. On success returns `RN_OK` and stores a template in `*out`;
  otherwise returns an error and sets `*out=NULL`.
- `void rn_template_free(rn_template *t)` — Free a template (NULL safe).
- `const char *rn_template_string(const rn_template *t)` — Get the original template string.
- `size_t rn_template_placeholder_count(const rn_template *t)` — Number of placeholders.
- `const char *rn_template_placeholder_at(const rn_template *t, size_t index)` —
  Get placeholder name at index (declaration order).

### Struct API

- `rn_status rn_generate(const rn_template *t, const void *record, const rn_field *fields, size_t nfields, char **out)` —
  Build a resource name from struct fields. Stores a heap string in `*out` (free with `free()`).
- `rn_status rn_parse(const rn_template *t, const char *resource, void *record, const rn_field *fields, size_t nfields)` —
  Parse resource name into struct fields. Heap copies are written to `*record`.
- `void rn_parse_free(void *record, const rn_field *fields, size_t nfields)` —
  Free strings written by `rn_parse`.

### Pairs API

- `rn_status rn_generate_pairs(const rn_template *t, const rn_pair *pairs, size_t count, char **out)` —
  Build a resource name from key/value pairs. Stores a heap string in `*out`.
- `rn_status rn_parse_pairs(const rn_template *t, const char *resource, rn_pair **out_pairs, size_t *out_count)` —
  Parse resource name into key/value pairs. Stores a heap array in `*out_pairs`.
- `void rn_pairs_free(rn_pair *pairs, size_t count)` — Free pairs array (NULL safe).

### Status Codes

- `RN_OK` — Success
- `RN_ERR_EMPTY_TEMPLATE` — Template string was empty
- `RN_ERR_NO_PLACEHOLDERS` — Template had no `{placeholder}` segments
- `RN_ERR_DUPLICATE_PLACEHOLDER` — A placeholder name appeared more than once
- `RN_ERR_INVALID_TEMPLATE` — Template could not be compiled or bad argument
- `RN_ERR_TEMPLATE_MISMATCH` — Resource name did not match the template
- `RN_ERR_MISSING_VALUE` — A required placeholder value was not provided
- `RN_ERR_UNEXPECTED_KEY` — A provided key is not a placeholder in the template
- `RN_ERR_SLASH_IN_VALUE` — A value contained `/`, invalid for one segment
- `RN_ERR_NOMEM` — Allocation failed

Use `const char *rn_status_message(rn_status status)` to get a human-readable message.

## Build, run, test

Uses only the C standard library and POSIX `<regex.h>` (macOS and Linux). No
external dependencies.

```bash
make            # builds libresourcename.a and the example
make run        # runs the example
make test       # builds and runs the test suite
```

## Notes / limitations

- Placeholder values are single path segments (`[^/]+`); they may not contain `/`.
- `rn_generate` / `rn_generate_pairs` error on a missing value, an unexpected key,
  or a `/` in a value.
- Templates with no placeholders, an empty string, or duplicate placeholder names
  are rejected at compile time.
