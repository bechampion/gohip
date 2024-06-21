all: delete_latest_tag recreate_tag push_tag

delete_latest_tag:
	@echo "Deleting latest local tag: latest"
	@git tag -d latest

recreate_tag:
	@echo "Recreating tag: latest"
	@git tag latest

push_tag:
	@echo "Force pushing tag: latest"
	@git push --force origin latest

.PHONY: all delete_latest_tag recreate_tag push_tag

build:
	go build -o gohip-$(GOOS)-$(GOARCH)

install: build
	mkdir -p $(DESTDIR)/usr/bin
	cp gohip-$(GOOS)-$(GOARCH) $(DESTDIR)/usr/bin/gohip

debian-pkg: install
	mkdir -p $(DESTDIR)/DEBIAN
	cp build-aux/debian/control $(DESTDIR)/DEBIAN/
	echo "Version: $(RELEASE_VERSION)" >> $(DESTDIR)/DEBIAN/control
	cp build-aux/debian/postinst $(DESTDIR)/DEBIAN/
	chmod 775 $(DESTDIR)/DEBIAN/postinst
	dpkg-deb --build $(DESTDIR) gohip-$(RELEASE_VERSION)-x86_64.deb
	md5sum gohip-$(RELEASE_VERSION)-x86_64.deb > gohip-$(RELEASE_VERSION)-x86_64.deb.md5sum
