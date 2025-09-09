
all: minesweeper

build-dir:
	mkdir -p $(PWD)/build/

minesweeper: build-dir
	go build $(BUILD_FLAGS) $(BUILD_CONSTANTS) -o $(PWD)/build/$@ $(PWD)/cmd/$@.go
