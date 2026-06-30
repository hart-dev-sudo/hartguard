BINARY  = port-scan-detector
PREFIX  = /usr/local
CFGDIR  = /etc/hartguard
LOGDIR  = /var/log/hartguard

build:
	go build -o $(BINARY) .

test:
	go test ./internal/detector/...

clean:
	rm -f $(BINARY)

install: build
	install -m 0755 $(BINARY) $(PREFIX)/bin/$(BINARY)
	install -m 0644 $(BINARY).service /etc/systemd/system/$(BINARY).service
	install -d $(CFGDIR)
	install -d $(LOGDIR)
	@if [ ! -f $(CFGDIR)/$(BINARY).yaml ]; then \
		install -m 0644 config.yaml $(CFGDIR)/$(BINARY).yaml; \
		echo "Config installed to $(CFGDIR)/$(BINARY).yaml — edit before starting service"; \
	else \
		echo "Config already exists at $(CFGDIR)/$(BINARY).yaml — skipping"; \
	fi
	systemctl daemon-reload
	@echo "Run: systemctl enable --now $(BINARY)"

uninstall:
	systemctl stop $(BINARY) 2>/dev/null || true
	systemctl disable $(BINARY) 2>/dev/null || true
	rm -f $(PREFIX)/bin/$(BINARY)
	rm -f /etc/systemd/system/$(BINARY).service
	systemctl daemon-reload

.PHONY: build test clean install uninstall
