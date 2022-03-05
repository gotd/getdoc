test:
	@./go.test.sh
coverage:
	@./go.coverage.sh
update_schema:
	@go run github.com/gotd/getdoc/cmd/getdoc --out-dir _schema --host corefork.telegram.org --out-file latest.json --pretty true
	@go run github.com/gotd/getdoc/cmd/getdoc --out-dir _schema --host corefork.telegram.org --pretty true
