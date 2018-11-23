#********************************************************************************
#*  @file
#*  @copyright defined in go-seele/LICENSE
#*********************************************************************************

# Makefile to build the command lines and tests in Seele project.
# This Makefile doesn't consider Windows Environment. If you use it in Windows, please be careful.

all: run
run:
	go build -o ./build/run ./e2e/
	@echo "Done e2e building"

.PHONY: run
