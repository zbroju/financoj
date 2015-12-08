SHELL = /bin/bash

CC = gcc
CFLAGS = -O3
LIBRARIES = -lsqlite3 -lconfig
EXECUTABLE = mm

SRC_DIR = src
BUILD_DIR = build
TARGET_DIR = target
DIST_DIR = mymoney
SOURCES = $(wildcard $(SRC_DIR)/*.c)
HEADERS = $(wildcard $(SRC_DIR)/*.h)
OBJECTS = $(patsubst $(SRC_DIR)/%.c, $(BUILD_DIR)/%.o, $(SOURCES))

all: build
	mkdir -p $(TARGET_DIR)/$(DIST_DIR)
	cp LICENSE README.md example.mmrc $(BUILD_DIR)/$(EXECUTABLE) $(TARGET_DIR)/$(DIST_DIR)
	tar -czvf $(TARGET_DIR)/mymoney.tar.gz -C $(TARGET_DIR) $(DIST_DIR)

build: pre_compile $(EXECUTABLE)

debug: pre_debug

$(EXECUTABLE):$(OBJECTS)
	$(CC) $(CFLAGS) -o $(BUILD_DIR)/$@ $(OBJECTS) $(LIBRARIES)

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c $(SRC_DIR)/%.h
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -r $(BUILD_DIR) $(TARGET_DIR)

pre_compile:
	mkdir -p $(BUILD_DIR)
	mkdir -p $(TARGET_DIR)

pre_debug: CFLAGS += -Wall -g
pre_debug: pre_compile $(EXECUTABLE)
