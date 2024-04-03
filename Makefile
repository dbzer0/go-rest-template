build:
	$(MAKE) -C app build

build-cross-platform:
	$(MAKE) -C app bin-cross-platform

clean:
	$(MAKE) -C app clean
	if [ -f coverage.html ] ; then rm coverage.html ; fi
	if [ -d .cover ] ; then rm -rf .cover ; fi
	docker-compose down --rmi all -v 2>/dev/null || true
	docker-compose stop >/dev/null
	docker-compose rm >/dev/null

rebuild:
	docker-compose build techno-sso
	docker-compose build unit

unit:
	docker-compose run --rm unit

update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif

coverage:
	docker-compose run --rm unit && [ -f ./coverage.html ] && xdg-open coverage.html

.PHONY: all build clean unit rebuild coverage
