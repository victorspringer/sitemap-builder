test:
	cd smb && go test -v -race $(go list ./...)