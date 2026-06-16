/* Extracting placeholder values from a resource name into pairs or a struct. */
#define _POSIX_C_SOURCE 200809L

#include "internal.h"

#include <regex.h>
#include <stdlib.h>
#include <string.h>

rn_status rn_parse_pairs(const rn_template *t, const char *resource,
                         rn_pair **out_pairs, size_t *out_count) {
    if (out_pairs) {
        *out_pairs = NULL;
    }
    if (out_count) {
        *out_count = 0;
    }
    if (!t || !resource || !out_pairs || !out_count) {
        return RN_ERR_INVALID_TEMPLATE;
    }

    size_t ngroups = t->nnames + 1;
    regmatch_t *matches = (regmatch_t *)calloc(ngroups, sizeof(regmatch_t));
    if (!matches) {
        return RN_ERR_NOMEM;
    }
    if (regexec(&t->re, resource, ngroups, matches, 0) != 0) {
        free(matches);
        return RN_ERR_TEMPLATE_MISMATCH;
    }

    rn_pair *pairs = (rn_pair *)calloc(t->nnames, sizeof(rn_pair));
    if (!pairs) {
        free(matches);
        return RN_ERR_NOMEM;
    }
    for (size_t i = 0; i < t->nnames; i++) {
        regmatch_t m = matches[i + 1];
        size_t len = (m.rm_so >= 0 && m.rm_eo >= m.rm_so) ? (size_t)(m.rm_eo - m.rm_so) : 0;
        pairs[i].key = rn_dup_str(t->names[i]);
        pairs[i].value = rn_dup_n(resource + (m.rm_so >= 0 ? m.rm_so : 0), len);
        if (!pairs[i].key || !pairs[i].value) {
            rn_pairs_free(pairs, t->nnames);
            free(matches);
            return RN_ERR_NOMEM;
        }
    }
    free(matches);
    *out_pairs = pairs;
    *out_count = t->nnames;
    return RN_OK;
}

void rn_pairs_free(rn_pair *pairs, size_t count) {
    if (!pairs) {
        return;
    }
    for (size_t i = 0; i < count; i++) {
        free(pairs[i].key);
        free(pairs[i].value);
    }
    free(pairs);
}

rn_status rn_parse(const rn_template *t, const char *resource,
                   void *record, const rn_field *fields, size_t nfields) {
    if (!t || !resource || !record || (nfields > 0 && !fields)) {
        return RN_ERR_INVALID_TEMPLATE;
    }

    rn_pair *parsed = NULL;
    size_t n = 0;
    rn_status status = rn_parse_pairs(t, resource, &parsed, &n);
    if (status != RN_OK) {
        return status;
    }

    for (size_t i = 0; i < nfields; i++) {
        const char *value = NULL;
        for (size_t k = 0; k < n; k++) {
            if (strcmp(parsed[k].key, fields[i].placeholder) == 0) {
                value = parsed[k].value;
                break;
            }
        }
        if (!value) {
            continue; /* this field's placeholder is not in the template */
        }
        char *copy = rn_dup_str(value);
        if (!copy) {
            rn_pairs_free(parsed, n);
            return RN_ERR_NOMEM;
        }
        *(char **)((char *)record + fields[i].offset) = copy;
    }

    rn_pairs_free(parsed, n);
    return RN_OK;
}

void rn_parse_free(void *record, const rn_field *fields, size_t nfields) {
    if (!record || !fields) {
        return;
    }
    for (size_t i = 0; i < nfields; i++) {
        char **slot = (char **)((char *)record + fields[i].offset);
        free(*slot);
        *slot = NULL;
    }
}
