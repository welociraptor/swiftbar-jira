all: install

.PHONY: build
build:
	go build -ldflags="-s -w" -o jira.5m.bin
	xattr -w "com.ameba.SwiftBar" "$(shell cat metadata.txt | base64)" jira.5m.bin

.PHONY: install
install: build
	cp jira.5m.bin swiftbar-jira.yaml ${SWIFTBAR_PLUGIN_DIR}
