CC = clang
CLANG_FORMAT = clang-format

SRC = array.c tests/test_array.c
CFLAGS = -I. -Itests
BUILD_DIR = build
TARGET = $(BUILD_DIR)/test_runner

VALGRIND_FLAGS = --leak-check=full --track-origins=yes --show-leak-kinds=all
VALGRIND = valgrind $(VALGRIND_FLAGS)

FORMAT_FILES = $(shell find tests -name '*.c' -o -name '*.h' | grep -v -F -f .clang-format-ignore)

.PHONY: all test clean download format

all: $(TARGET)

$(TARGET): $(SRC)
	@mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) $^ -o $@

test: $(TARGET)
	./$(TARGET)

clean:
	rm -rf $(BUILD_DIR)


download:
	curl -o tests/greatest.h https://raw.githubusercontent.com/silentbicycle/greatest/master/greatest.h

format:
	$(CLANG_FORMAT) -i $(FORMAT_FILES)

lint:
	clang-format --dry-run --Werror $(FORMAT_FILES)

valgrind: $(TARGET)
	$(VALGRIND) --error-exitcode=1 $(TARGET)

