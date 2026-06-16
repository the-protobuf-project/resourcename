/*
 * Demonstrates the resourcename C API working with structs.
 *
 * C has no reflection, so each struct is paired with a table of rn_field
 * descriptors (the C analog of the Go example's struct tags). rn_generate reads
 * the struct's fields into a resource name; rn_parse reads a name back into a
 * struct.
 */
#include "resourcename.h"

#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>

/* A basic resource: //music.example.com/artists/{id}/{name} */
typedef struct {
    char *id;
    char *name;
} Artist;

static const rn_field ARTIST_FIELDS[] = {
    {"id", offsetof(Artist, id)},
    {"name", offsetof(Artist, name)},
};
#define ARTIST_NFIELDS (sizeof(ARTIST_FIELDS) / sizeof(ARTIST_FIELDS[0]))

/* A multi-segment resource: //music.example.com/artists/{id}/albums/{title}/{year} */
typedef struct {
    char *id;
    char *title;
    char *year;
} ArtistAlbum;

static const rn_field ALBUM_FIELDS[] = {
    {"id", offsetof(ArtistAlbum, id)},
    {"title", offsetof(ArtistAlbum, title)},
    {"year", offsetof(ArtistAlbum, year)},
};
#define ALBUM_NFIELDS (sizeof(ALBUM_FIELDS) / sizeof(ALBUM_FIELDS[0]))

static void demo_basic(void) {
    printf("\n1. Basic struct (generate / parse):\n");

    rn_template *t = NULL;
    rn_status st = rn_template_compile("//music.example.com/artists/{id}/{name}", &t);
    if (st != RN_OK) {
        fprintf(stderr, "   compile error: %s\n", rn_status_message(st));
        return;
    }

    Artist artist = {.id = "ar-42", .name = "Radiohead"};
    char *rn = NULL;
    st = rn_generate(t, &artist, ARTIST_FIELDS, ARTIST_NFIELDS, &rn);
    if (st == RN_OK) {
        printf("   Generated:  %s\n", rn);
    }

    Artist parsed = {0};
    st = rn_parse(t, rn, &parsed, ARTIST_FIELDS, ARTIST_NFIELDS);
    if (st == RN_OK) {
        printf("   Parsed:     id=%s, name=%s\n", parsed.id, parsed.name);
    }

    rn_parse_free(&parsed, ARTIST_FIELDS, ARTIST_NFIELDS);
    free(rn);
    rn_template_free(t);
}

static void demo_nested(void) {
    printf("\n2. Multi-segment struct (id / album title / year):\n");

    rn_template *t = NULL;
    rn_status st =
        rn_template_compile("//music.example.com/artists/{id}/albums/{title}/{year}", &t);
    if (st != RN_OK) {
        fprintf(stderr, "   compile error: %s\n", rn_status_message(st));
        return;
    }

    ArtistAlbum album = {.id = "radiohead", .title = "in-rainbows", .year = "2007"};
    char *rn = NULL;
    st = rn_generate(t, &album, ALBUM_FIELDS, ALBUM_NFIELDS, &rn);
    if (st == RN_OK) {
        printf("   Generated:  %s\n", rn);
    }

    ArtistAlbum parsed = {0};
    st = rn_parse(t, rn, &parsed, ALBUM_FIELDS, ALBUM_NFIELDS);
    if (st == RN_OK) {
        printf("   Parsed:     id=%s, title=%s, year=%s\n", parsed.id, parsed.title, parsed.year);
    }

    rn_parse_free(&parsed, ALBUM_FIELDS, ALBUM_NFIELDS);
    free(rn);
    rn_template_free(t);
}

int main(void) {
    printf("=== Resource Name Demo (C) ===\n");
    demo_basic();
    demo_nested();
    return 0;
}
