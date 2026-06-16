/* Minimal test harness for the resourcename C API (no external deps). */
#include "resourcename.h"

#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static int failures = 0;
static int checks = 0;

#define CHECK(cond, msg)                                             \
    do {                                                             \
        checks++;                                                    \
        if (!(cond)) {                                               \
            failures++;                                              \
            printf("FAIL: %s (%s:%d)\n", (msg), __FILE__, __LINE__); \
        }                                                            \
    } while (0)

static const char *value_for(const rn_pair *pairs, size_t n, const char *key) {
    for (size_t i = 0; i < n; i++) {
        if (strcmp(pairs[i].key, key) == 0) {
            return pairs[i].value;
        }
    }
    return NULL;
}

/* ---- key/value (pairs) API ----------------------------------------------- */

static void test_pairs_roundtrip(void) {
    rn_template *t = NULL;
    CHECK(rn_template_compile("//music.example.com/artists/{id}/{name}", &t) == RN_OK,
          "compile basic template");
    CHECK(rn_template_placeholder_count(t) == 2, "two placeholders");

    rn_pair in[] = {{(char *)"id", (char *)"ar-42"}, {(char *)"name", (char *)"Radiohead"}};
    char *rn = NULL;
    CHECK(rn_generate_pairs(t, in, 2, &rn) == RN_OK, "generate_pairs ok");
    CHECK(rn && strcmp(rn, "//music.example.com/artists/ar-42/Radiohead") == 0, "generate value");

    rn_pair *out = NULL;
    size_t n = 0;
    CHECK(rn_parse_pairs(t, rn, &out, &n) == RN_OK, "parse_pairs ok");
    CHECK(n == 2, "parse count");
    CHECK(out && strcmp(value_for(out, n, "id"), "ar-42") == 0, "parse id");
    CHECK(out && strcmp(value_for(out, n, "name"), "Radiohead") == 0, "parse name");

    free(rn);
    rn_pairs_free(out, n);
    rn_template_free(t);
}

static void test_parse_mismatch(void) {
    rn_template *t = NULL;
    rn_template_compile("//music.example.com/artists/{id}", &t);

    rn_pair *out = NULL;
    size_t n = 0;
    CHECK(rn_parse_pairs(t, "//music.example.com/artists/a/b", &out, &n) == RN_ERR_TEMPLATE_MISMATCH,
          "parse rejects extra segment");
    CHECK(rn_parse_pairs(t, "//other.example.com/x/y", &out, &n) == RN_ERR_TEMPLATE_MISMATCH,
          "parse rejects different prefix");

    rn_template_free(t);
}

static void test_generate_errors(void) {
    rn_template *t = NULL;
    rn_template_compile("//music.example.com/artists/{id}/{name}", &t);

    char *rn = NULL;
    rn_pair missing[] = {{(char *)"id", (char *)"ar-42"}};
    CHECK(rn_generate_pairs(t, missing, 1, &rn) == RN_ERR_MISSING_VALUE, "missing value rejected");

    rn_pair extra[] = {{(char *)"id", (char *)"ar-42"},
                       {(char *)"name", (char *)"Radiohead"},
                       {(char *)"other", (char *)"x"}};
    CHECK(rn_generate_pairs(t, extra, 3, &rn) == RN_ERR_UNEXPECTED_KEY, "unexpected key rejected");

    rn_pair slash[] = {{(char *)"id", (char *)"ar/42"}, {(char *)"name", (char *)"Radiohead"}};
    CHECK(rn_generate_pairs(t, slash, 2, &rn) == RN_ERR_SLASH_IN_VALUE, "slash in value rejected");

    rn_template_free(t);
}

static void test_template_validation(void) {
    rn_template *t = NULL;
    CHECK(rn_template_compile("", &t) == RN_ERR_EMPTY_TEMPLATE, "empty template rejected");
    CHECK(rn_template_compile("//music.example.com/artists", &t) == RN_ERR_NO_PLACEHOLDERS,
          "no-placeholder template rejected");
    CHECK(rn_template_compile("//x/{id}/y/{id}", &t) == RN_ERR_DUPLICATE_PLACEHOLDER,
          "duplicate placeholder rejected");
}

/* ---- struct API ----------------------------------------------------------- */

typedef struct {
    char *id;
    char *name;
} Artist;

static const rn_field ARTIST_FIELDS[] = {
    {"id", offsetof(Artist, id)},
    {"name", offsetof(Artist, name)},
};

static void test_struct_roundtrip(void) {
    rn_template *t = NULL;
    rn_template_compile("//music.example.com/artists/{id}/{name}", &t);

    Artist in = {.id = "ar-42", .name = "Radiohead"};
    char *rn = NULL;
    CHECK(rn_generate(t, &in, ARTIST_FIELDS, 2, &rn) == RN_OK, "struct generate ok");
    CHECK(rn && strcmp(rn, "//music.example.com/artists/ar-42/Radiohead") == 0, "struct generate value");

    Artist out = {0};
    CHECK(rn_parse(t, rn, &out, ARTIST_FIELDS, 2) == RN_OK, "struct parse ok");
    CHECK(out.id && strcmp(out.id, "ar-42") == 0, "struct parse id");
    CHECK(out.name && strcmp(out.name, "Radiohead") == 0, "struct parse name");

    rn_parse_free(&out, ARTIST_FIELDS, 2);
    CHECK(out.id == NULL && out.name == NULL, "parse_free nulls fields");

    Artist partial = {.id = "ar-42", .name = NULL};
    char *rn2 = NULL;
    CHECK(rn_generate(t, &partial, ARTIST_FIELDS, 2, &rn2) == RN_ERR_MISSING_VALUE,
          "struct generate reports missing field");

    free(rn);
    rn_template_free(t);
}

int main(void) {
    test_pairs_roundtrip();
    test_parse_mismatch();
    test_generate_errors();
    test_template_validation();
    test_struct_roundtrip();

    printf("%d checks, %d failures\n", checks, failures);
    return failures == 0 ? 0 : 1;
}
