CC ?= cc
CFLAGS ?= -std=c11 -Wall -Wextra -O2
CPPFLAGS += -Ilib
AR ?= ar

LIB := libresourcename.a
LIB_SRC := $(wildcard lib/*.c)
LIB_OBJ := $(LIB_SRC:.c=.o)

.PHONY: all lib run test clean
all: lib example

lib: $(LIB)

$(LIB): $(LIB_OBJ)
	$(AR) rcs $@ $^

lib/%.o: lib/%.c lib/resourcename.h lib/internal.h
	$(CC) $(CFLAGS) $(CPPFLAGS) -c $< -o $@

example: example.c $(LIB)
	$(CC) $(CFLAGS) $(CPPFLAGS) $< $(LIB) -o $@

test_runner: test.c $(LIB)
	$(CC) $(CFLAGS) $(CPPFLAGS) $< $(LIB) -o $@

run: example
	./example

test: test_runner
	./test_runner

clean:
	rm -f $(LIB_OBJ) $(LIB) example test_runner
