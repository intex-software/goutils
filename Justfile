[private]
default:
	@just --list

update:
	go get -u ./...
	go mod tidy

version := replace(replace_regex(`git describe --tags --always --match=v*`, "v|-g.*", ""), "-", ".")

version:
	@echo {{version}}

create-version newVersion:
	git tag -a -m "Neu Version v{{newVersion}}" "v{{newVersion}}"
