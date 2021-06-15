CHECKVARS=$(and $(GCP_PROJECT),$(TFGM_API_URL),$(TFGM_API_KEY),OK)

deploy: check-env test

check-env:
	@if [ "$(CHECKVARS)" != "OK" ]; then echo "Please set GCP_PROJECT, TFGM_URL and TFGM_AP_KEY env vars before running deploy target!"; exit 1; fi

test: clean
	go test ./

clean:
	@go clean -testcache
	@rm -rf dist
	@mkdir dist
