SHELL = /bin/bash

CC = gcc
CFLAG = -Wall
LIBRARIES = -lsqlite3 -lconfig
EXECUTABLE = mm

SRC_DIR = src
BUILD_DIR = build
TARGET_DIR = target
SOURCES = $(wildcard $(SRC_DIR)/*.c)
HEADERS = $(wildcard $(SRC_DIR)/*.h)
OBJECTS = $(patsubst $(SRC_DIR)/%.c, $(BUILD_DIR)/%.o, $(SOURCES))

all: pre_compile $(EXECUTABLE)

$(EXECUTABLE):$(OBJECTS)
	$(CC) $(CFLAGS) -o $(TARGET_DIR)/$@ $(OBJECTS) $(LIBRARIES)

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c $(SRC_DIR)/%.h
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -r $(BUILD_DIR) $(TARGET_DIR)

pre_compile:
	mkdir -p $(BUILD_DIR)
	mkdir -p $(TARGET_DIR)
