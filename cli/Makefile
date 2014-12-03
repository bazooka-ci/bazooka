default: install

install: godep-restore
	cd bzk && \
    go install

godep-restore:
	godep restore ./...

godep-save:
	godep save ./...
