
all: minesweeper

build-dir:
	mkdir -p $(PWD)/build/

minesweeper: build-dir
	CGO_ENABLED=0 go build $(BUILD_FLAGS) $(BUILD_CONSTANTS) -o $(PWD)/build/$@ $(PWD)/cmd/$@.go

lint:
	golangci-lint run
