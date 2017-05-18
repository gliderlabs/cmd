.PHONY: build dev stage clobber publisher publish-local

BUILD_DIR := ../build/

build: stage
	cd $(BUILD_DIR)www && hugo

dev: stage
	cd $(BUILD_DIR)www && hugo server

stage: clobber $(BUILD_DIR)www

clobber:
	rm -rf $(BUILD_DIR)www

publish-local:
	cd .. && tar \
		--exclude=vendor \
		--exclude=.git \
		--exclude=build \
		-cf - . | ssh alpha.cmd.io progrium/cmd-publish master

publish-ci:
	@cd .. && tar \
		--exclude=vendor \
		--exclude=.git \
		--exclude=build \
		-cf - . | ssh $(PUBLISH_TOKEN)@alpha.cmd.io progrium/cmd-publish $(CIRCLE_BRANCH)

../ui/static/semantic:
	$(MAKE) -C ../ui static/semantic

$(BUILD_DIR)www: ../ui/static/semantic
	mkdir -p $(BUILD_DIR)
	cd $(BUILD_DIR) && hugo new site www > /dev/null
	cp -r ../ui/static/semantic $(BUILD_DIR)www/static/semantic
	rm -rf $(BUILD_DIR)www/content
	rm $(BUILD_DIR)www/config.toml
	# remember, the target path is relative to the symlink file
	ln -s ../../docs $(BUILD_DIR)www/content
	ln -s ../../../www $(BUILD_DIR)www/themes/www
	ln -s ../../www/config.toml $(BUILD_DIR)www/config.toml
