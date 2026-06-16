/* Template lifecycle (free) and read-only accessors. Compilation is in compile.c. */
#define _POSIX_C_SOURCE 200809L

#include "internal.h"

#include <stdlib.h>

void rn_template_free(rn_template *t) {
    if (!t) {
        return;
    }
    if (t->re_valid) {
        regfree(&t->re);
    }
    for (size_t i = 0; i < t->nnames; i++) {
        free(t->names[i]);
    }
    free(t->names);
    for (size_t i = 0; i < t->nsegments; i++) {
        if (!t->segments[i].is_placeholder) {
            free(t->segments[i].literal);
        }
    }
    free(t->segments);
    free(t->tmpl);
    free(t);
}

const char *rn_template_string(const rn_template *t) {
    return t ? t->tmpl : NULL;
}

size_t rn_template_placeholder_count(const rn_template *t) {
    return t ? t->nnames : 0;
}

const char *rn_template_placeholder_at(const rn_template *t, size_t index) {
    return (t && index < t->nnames) ? t->names[index] : NULL;
}
