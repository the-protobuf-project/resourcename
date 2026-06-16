/*
 * resourcename — Google AIP-122 style resource-name templates for C.
 *
 * A template looks like "//music.example.com/artists/{id}/{name}". Each
 * {placeholder} matches exactly one path segment ([^/]+), so values may not
 * contain '/'. Compile a template once, then:
 *
 *   - rn_generate / rn_parse        work with your own structs (recommended)
 *   - rn_generate_pairs / rn_parse_pairs  work with key/value arrays (primitive)
 *
 * This is the C port of the Go `Template`.
 */
#ifndef RESOURCENAME_H
#define RESOURCENAME_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Opaque compiled template. Create with rn_template_compile, free with rn_template_free. */
typedef struct rn_template rn_template;

/* Status codes returned by every fallible call. See rn_status_message(). */
typedef enum {
    RN_OK = 0,
    RN_ERR_EMPTY_TEMPLATE,        /* template string was empty */
    RN_ERR_NO_PLACEHOLDERS,       /* template had no {placeholder} segments */
    RN_ERR_DUPLICATE_PLACEHOLDER, /* a placeholder name appeared more than once */
    RN_ERR_INVALID_TEMPLATE,      /* template could not be compiled / bad argument */
    RN_ERR_TEMPLATE_MISMATCH,     /* resource name did not match the template */
    RN_ERR_MISSING_VALUE,         /* a required placeholder value was not provided */
    RN_ERR_UNEXPECTED_KEY,        /* a provided key is not a placeholder in the template */
    RN_ERR_SLASH_IN_VALUE,        /* a value contained '/', invalid for one segment */
    RN_ERR_NOMEM                  /* allocation failed */
} rn_status;

/* ---- compile -------------------------------------------------------------- */

/* Compile a template. On success returns RN_OK and stores a template in *out
 * (free with rn_template_free); otherwise returns an error and sets *out=NULL. */
rn_status rn_template_compile(const char *template, rn_template **out);

/* Free a template created by rn_template_compile. NULL is allowed. */
void rn_template_free(rn_template *t);

/* The original template string. */
const char *rn_template_string(const rn_template *t);

/* Number of placeholders, and the placeholder name at `index` (declaration order). */
size_t rn_template_placeholder_count(const rn_template *t);
const char *rn_template_placeholder_at(const rn_template *t, size_t index);

/* ---- struct mapping (recommended) ----------------------------------------- *
 *
 * C has no reflection, so a struct's layout is described once with a table of
 * rn_field entries (the C analog of Go's struct tags). Every mapped field must
 * be a `char *`; fill the offset with offsetof:
 *
 *     typedef struct { char *id; char *name; } Artist;
 *     static const rn_field ARTIST_FIELDS[] = {
 *         { "id",   offsetof(Artist, id) },
 *         { "name", offsetof(Artist, name) },
 *     };
 */
typedef struct {
    const char *placeholder; /* placeholder name in the template */
    size_t offset;           /* offsetof(YourStruct, field); field must be char* */
} rn_field;

/* Build a resource name from the `char *` fields of *record. A NULL field counts
 * as missing. On success stores a heap string in *out (free with free()). */
rn_status rn_generate(const rn_template *t, const void *record,
                      const rn_field *fields, size_t nfields, char **out);

/* Parse `resource` and store heap copies of each captured segment into the
 * matching `char *` fields of *record (zero-initialize *record first). Free the
 * written strings with rn_parse_free. */
rn_status rn_parse(const rn_template *t, const char *resource,
                   void *record, const rn_field *fields, size_t nfields);

/* Free (and NULL out) the strings written into *record by rn_parse. */
void rn_parse_free(void *record, const rn_field *fields, size_t nfields);

/* ---- key/value mapping (primitive) ---------------------------------------- */

/* A placeholder key/value pair: input to rn_generate_pairs, output of rn_parse_pairs. */
typedef struct {
    char *key;
    char *value;
} rn_pair;

/* Build a resource name from `count` pairs (one per placeholder). On success
 * stores a heap string in *out (free with free()). */
rn_status rn_generate_pairs(const rn_template *t, const rn_pair *pairs, size_t count,
                            char **out);

/* Parse `resource` into placeholder values. On success stores a heap array of
 * *out_count pairs in *out_pairs (free with rn_pairs_free). */
rn_status rn_parse_pairs(const rn_template *t, const char *resource,
                         rn_pair **out_pairs, size_t *out_count);

/* Free a pairs array returned by rn_parse_pairs. NULL is allowed. */
void rn_pairs_free(rn_pair *pairs, size_t count);

/* Human-readable message for a status code (never NULL). */
const char *rn_status_message(rn_status status);

#ifdef __cplusplus
}
#endif

#endif /* RESOURCENAME_H */
