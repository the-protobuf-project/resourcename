/*
 * Internal declarations shared across the resourcename implementation files.
 * Not part of the public API.
 */
#ifndef RESOURCENAME_INTERNAL_H
#define RESOURCENAME_INTERNAL_H

#include "resourcename.h"

#include <regex.h>
#include <stddef.h>

/* ---- string helpers (util.c) --------------------------------------------- */

/* Duplicate the first `n` bytes of `src` as a NUL-terminated string. */
char *rn_dup_n(const char *src, size_t n);
/* Duplicate a NUL-terminated string. */
char *rn_dup_str(const char *src);

/* A growable, always-NUL-terminated string buffer (zero-initialize to use). */
typedef struct {
    char *buf;
    size_t len;
    size_t cap;
} rn_strbuf;

int rn_sb_putc(rn_strbuf *s, char c);            /* append one char; 0 ok, -1 nomem */
int rn_sb_putn(rn_strbuf *s, const char *p, size_t n); /* append n bytes */
int rn_sb_puts(rn_strbuf *s, const char *str);   /* append a C string */

/* True if `c` is an ERE metacharacter that must be escaped in a literal. */
int rn_is_ere_meta(char c);

/* Index of `name` in `names`, or -1 if absent. */
long rn_name_index(char *const *names, size_t nnames, const char *name);

/* ---- template internals (template.c) ------------------------------------- */

/* One piece of a template: literal text, or a reference to a placeholder name. */
typedef struct {
    int is_placeholder;
    char *literal;     /* owned; valid when !is_placeholder */
    size_t name_index; /* index into names[]; valid when is_placeholder */
} rn_segment;

struct rn_template {
    char *tmpl;           /* original template string */
    char **names;         /* placeholder names, in order */
    size_t nnames;
    rn_segment *segments; /* literal/placeholder pieces for generation */
    size_t nsegments;
    regex_t re;           /* compiled matcher for parsing */
    int re_valid;
};

#endif /* RESOURCENAME_INTERNAL_H */
