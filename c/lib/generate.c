/* Building resource names from key/value pairs or from struct fields. */
#include "internal.h"

#include <stdlib.h>
#include <string.h>

rn_status rn_generate_pairs(const rn_template *t, const rn_pair *pairs, size_t count,
                            char **out) {
    if (out) {
        *out = NULL;
    }
    if (!t || !out || (count > 0 && !pairs)) {
        return RN_ERR_INVALID_TEMPLATE;
    }

    /* Reject keys that are not placeholders in this template. */
    for (size_t k = 0; k < count; k++) {
        if (!pairs[k].key || rn_name_index(t->names, t->nnames, pairs[k].key) < 0) {
            return RN_ERR_UNEXPECTED_KEY;
        }
    }

    rn_strbuf result = {0};
    for (size_t s = 0; s < t->nsegments; s++) {
        const rn_segment *seg = &t->segments[s];
        if (!seg->is_placeholder) {
            if (rn_sb_puts(&result, seg->literal) != 0) {
                goto nomem;
            }
            continue;
        }
        const char *name = t->names[seg->name_index];
        const char *value = NULL;
        for (size_t k = 0; k < count; k++) {
            if (pairs[k].key && strcmp(pairs[k].key, name) == 0) {
                value = pairs[k].value;
                break;
            }
        }
        if (!value) {
            free(result.buf);
            return RN_ERR_MISSING_VALUE;
        }
        if (strchr(value, '/') != NULL) {
            free(result.buf);
            return RN_ERR_SLASH_IN_VALUE;
        }
        if (rn_sb_puts(&result, value) != 0) {
            goto nomem;
        }
    }
    if (!result.buf) {
        result.buf = rn_dup_str("");
        if (!result.buf) {
            return RN_ERR_NOMEM;
        }
    }
    *out = result.buf;
    return RN_OK;

nomem:
    free(result.buf);
    return RN_ERR_NOMEM;
}

rn_status rn_generate(const rn_template *t, const void *record,
                      const rn_field *fields, size_t nfields, char **out) {
    if (out) {
        *out = NULL;
    }
    if (!t || !record || (nfields > 0 && !fields) || !out) {
        return RN_ERR_INVALID_TEMPLATE;
    }

    rn_pair *pairs = NULL;
    if (nfields > 0) {
        pairs = (rn_pair *)calloc(nfields, sizeof(rn_pair));
        if (!pairs) {
            return RN_ERR_NOMEM;
        }
    }
    size_t npairs = 0;
    for (size_t i = 0; i < nfields; i++) {
        char *value = *(char *const *)((const char *)record + fields[i].offset);
        if (!value) {
            continue; /* leave rn_generate_pairs to report the missing value */
        }
        pairs[npairs].key = (char *)fields[i].placeholder;
        pairs[npairs].value = value;
        npairs++;
    }

    rn_status status = rn_generate_pairs(t, pairs, npairs, out);
    free(pairs);
    return status;
}
