test:
	@./go.test.sh
coverage:
	@./go.coverage.sh
update_schema:
	go run github.com/gotd/getdoc/cmd/getdoc -out-dir _schema -out-file latest.json -pretty true
	# Download to ${layer_version}.json
	go run github.com/gotd/getdoc/cmd/getdoc -out-dir _schema -pretty true
