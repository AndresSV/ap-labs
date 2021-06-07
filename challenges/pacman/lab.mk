# build & test automation

APP_NAME=pacman

build:
	go build ${APP_NAME}.go

test: build
	@echo Test 1 - Launch Pacman
	./${APP_NAME}

clean:
	rm -rf pacman