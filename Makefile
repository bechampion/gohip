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

test:
	go test -v ./systemd ./others ./osdata ./types .

build: build-dynamic build-static

build-dynamic:
	go build -o gohip-$(GOOS)-$(GOARCH)

build-static:
	CGO_ENABLED=0 go build -o gohip-static-$(GOOS)-$(GOARCH)

install: build
	mkdir -p $(DESTDIR)/usr/bin
	cp gohip-$(GOOS)-$(GOARCH) $(DESTDIR)/usr/bin/gohip
	cp gohip-static-$(GOOS)-$(GOARCH) $(DESTDIR)/usr/bin/gohip-static

debian-pkg: install
	mkdir -p $(DESTDIR)/DEBIAN
	mkdir -p $(DESTDIR)/etc/vpnc/post-connect.d/

	cp build-aux/scripts/split.sh $(DESTDIR)/etc/vpnc/post-connect.d/split.sh
	chmod 755 $(DESTDIR)/etc/vpnc/post-connect.d/split.sh

	cp build-aux/debian/control $(DESTDIR)/DEBIAN/

	echo "Version: $(RELEASE_VERSION)" >> $(DESTDIR)/DEBIAN/control
	cp build-aux/debian/postinst $(DESTDIR)/DEBIAN/
	chmod 775 $(DESTDIR)/DEBIAN/postinst
	dpkg-deb --build $(DESTDIR) gohip-$(RELEASE_VERSION)-x86_64.deb
	md5sum gohip-$(RELEASE_VERSION)-x86_64.deb > gohip-$(RELEASE_VERSION)-x86_64.deb.md5sum
