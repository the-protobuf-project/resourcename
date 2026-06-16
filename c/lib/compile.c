/* Compiling a template string into placeholders, segments, and a matcher. */
#define _POSIX_C_SOURCE 200809L

#include "internal.h"

#include <stdlib.h>
#include <string.h>

/* Append a segment to a growable rn_segment array. Returns 0 / -1 (nomem). */
static int push_segment(rn_segment **segs, size_t *n, size_t *cap, rn_segment seg) {
    if (*n == *cap) {
        size_t ncap = *cap ? *cap * 2 : 8;
        rn_segment *grown = (rn_segment *)realloc(*segs, ncap * sizeof(rn_segment));
        if (!grown) {
            return -1;
        }
        *segs = grown;
        *cap = ncap;
    }
    (*segs)[(*n)++] = seg;
    return 0;
}

/* Flush accumulated literal text (if any) as a literal segment. */
static int flush_literal(rn_segment **segs, size_t *n, size_t *cap, rn_strbuf *lit) {
    if (lit->len == 0) {
        return 0;
    }
    char *owned = rn_dup_n(lit->buf, lit->len);
    if (!owned) {
        return -1;
    }
    rn_segment seg = {0, owned, 0};
    if (push_segment(segs, n, cap, seg) != 0) {
        free(owned);
        return -1;
    }
    lit->len = 0;
    lit->buf[0] = '\0';
    return 0;
}

rn_status rn_template_compile(const char *template, rn_template **out) {
    if (out) {
        *out = NULL;
    }
    if (!template || !out) {
        return RN_ERR_INVALID_TEMPLATE;
    }
    if (template[0] == '\0') {
        return RN_ERR_EMPTY_TEMPLATE;
    }

    char **names = NULL;
    size_t nnames = 0, names_cap = 0;
    rn_segment *segs = NULL;
    size_t nsegs = 0, segs_cap = 0;
    rn_strbuf lit = {0}, pat = {0};
    rn_status status = RN_ERR_NOMEM;

    if (rn_sb_putc(&pat, '^') != 0) {
        goto fail;
    }
    for (size_t i = 0; template[i] != '\0';) {
        size_t j = i + 1;
        if (template[i] == '{') {
            while (template[j] != '\0' && template[j] != '}') {
                j++;
            }
        }
        /* A "{name}" with non-empty content is a placeholder; anything else is literal. */
        if (template[i] == '{' && template[j] == '}' && j > i + 1) {
            char *name = rn_dup_n(template + i + 1, j - (i + 1));
            if (!name) {
                goto fail;
            }
            if (rn_name_index(names, nnames, name) >= 0) {
                free(name);
                status = RN_ERR_DUPLICATE_PLACEHOLDER;
                goto fail;
            }
            if (flush_literal(&segs, &nsegs, &segs_cap, &lit) != 0) {
                free(name);
                goto fail;
            }
            if (nnames == names_cap) {
                size_t ncap = names_cap ? names_cap * 2 : 8;
                char **grown = (char **)realloc(names, ncap * sizeof(char *));
                if (!grown) {
                    free(name);
                    goto fail;
                }
                names = grown;
                names_cap = ncap;
            }
            names[nnames] = name;
            rn_segment seg = {1, NULL, nnames};
            if (push_segment(&segs, &nsegs, &segs_cap, seg) != 0) {
                goto fail;
            }
            nnames++;
            if (rn_sb_puts(&pat, "([^/]+)") != 0) {
                goto fail;
            }
            i = j + 1;
            continue;
        }
        if (rn_sb_putc(&lit, template[i]) != 0) {
            goto fail;
        }
        if (rn_is_ere_meta(template[i]) && rn_sb_putc(&pat, '\\') != 0) {
            goto fail;
        }
        if (rn_sb_putc(&pat, template[i]) != 0) {
            goto fail;
        }
        i++;
    }
    if (flush_literal(&segs, &nsegs, &segs_cap, &lit) != 0 || rn_sb_putc(&pat, '$') != 0) {
        goto fail;
    }
    if (nnames == 0) {
        status = RN_ERR_NO_PLACEHOLDERS;
        goto fail;
    }

    rn_template *t = (rn_template *)calloc(1, sizeof(rn_template));
    if (!t) {
        goto fail;
    }
    if (regcomp(&t->re, pat.buf, REG_EXTENDED) != 0) {
        free(t);
        status = RN_ERR_INVALID_TEMPLATE;
        goto fail;
    }
    t->re_valid = 1;
    t->tmpl = rn_dup_str(template);
    if (!t->tmpl) {
        rn_template_free(t);
        goto fail;
    }
    t->names = names;
    t->nnames = nnames;
    t->segments = segs;
    t->nsegments = nsegs;
    free(lit.buf);
    free(pat.buf);
    *out = t;
    return RN_OK;

fail:
    free(lit.buf);
    free(pat.buf);
    for (size_t i = 0; i < nnames; i++) {
        free(names[i]);
    }
    free(names);
    for (size_t i = 0; i < nsegs; i++) {
        if (!segs[i].is_placeholder) {
            free(segs[i].literal);
        }
    }
    free(segs);
    return status;
}
