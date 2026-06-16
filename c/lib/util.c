/* String helpers, ERE-metachar test, name lookup, and status messages. */
#define _POSIX_C_SOURCE 200809L

#include "internal.h"

#include <stdlib.h>
#include <string.h>

char *rn_dup_n(const char *src, size_t n) {
    char *out = (char *)malloc(n + 1);
    if (!out) {
        return NULL;
    }
    memcpy(out, src, n);
    out[n] = '\0';
    return out;
}

char *rn_dup_str(const char *src) {
    return rn_dup_n(src, strlen(src));
}

/* Grow the buffer so it can hold `extra` more bytes plus a NUL terminator. */
static int sb_reserve(rn_strbuf *s, size_t extra) {
    if (s->len + extra + 1 <= s->cap) {
        return 0;
    }
    size_t cap = s->cap ? s->cap : 16;
    while (s->len + extra + 1 > cap) {
        cap *= 2;
    }
    char *buf = (char *)realloc(s->buf, cap);
    if (!buf) {
        return -1;
    }
    s->buf = buf;
    s->cap = cap;
    return 0;
}

int rn_sb_putc(rn_strbuf *s, char c) {
    if (sb_reserve(s, 1) != 0) {
        return -1;
    }
    s->buf[s->len++] = c;
    s->buf[s->len] = '\0';
    return 0;
}

int rn_sb_putn(rn_strbuf *s, const char *p, size_t n) {
    if (sb_reserve(s, n) != 0) {
        return -1;
    }
    memcpy(s->buf + s->len, p, n);
    s->len += n;
    s->buf[s->len] = '\0';
    return 0;
}

int rn_sb_puts(rn_strbuf *s, const char *str) {
    return rn_sb_putn(s, str, strlen(str));
}

int rn_is_ere_meta(char c) {
    switch (c) {
    case '.':
    case '\\':
    case '+':
    case '*':
    case '?':
    case '[':
    case ']':
    case '^':
    case '$':
    case '(':
    case ')':
    case '{':
    case '}':
    case '|':
        return 1;
    default:
        return 0;
    }
}

long rn_name_index(char *const *names, size_t nnames, const char *name) {
    for (size_t i = 0; i < nnames; i++) {
        if (strcmp(names[i], name) == 0) {
            return (long)i;
        }
    }
    return -1;
}

const char *rn_status_message(rn_status status) {
    switch (status) {
    case RN_OK:
        return "ok";
    case RN_ERR_EMPTY_TEMPLATE:
        return "template cannot be empty";
    case RN_ERR_NO_PLACEHOLDERS:
        return "template must contain at least one placeholder";
    case RN_ERR_DUPLICATE_PLACEHOLDER:
        return "template contains a duplicate placeholder";
    case RN_ERR_INVALID_TEMPLATE:
        return "invalid template";
    case RN_ERR_TEMPLATE_MISMATCH:
        return "resource name does not match template";
    case RN_ERR_MISSING_VALUE:
        return "missing value for placeholder";
    case RN_ERR_UNEXPECTED_KEY:
        return "unexpected key not present in template";
    case RN_ERR_SLASH_IN_VALUE:
        return "value must not contain '/'";
    case RN_ERR_NOMEM:
        return "out of memory";
    }
    return "unknown error";
}
