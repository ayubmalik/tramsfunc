CHECKVARS=$(and $(TFGM_API_KEY),OK)
GCP_PROJECT=tramsfunc
GCP_REGION=europe-west2
TFGM_API_URL=https://api.tfgm.com/odata

deploy: check-env test deploy-function

check-env:
	@if [ "$(CHECKVARS)" != "OK" ]; then echo "Please set TFGM_API_KEY env var before running deploy target!"; exit 1; fi

test: clean
	go test ./

deploy-function:
	@gcloud functions deploy tramsfunc --project $(GCP_PROJECT) --region $(GCP_REGION) --runtime=go113 --entry-point API \
  	--trigger-http --allow-unauthenticated --set-env-vars TFGM_API_URL=$(TFGM_API_URL),TFGM_API_KEY=$(TFGM_API_KEY)

clean:
	@go clean -testcache
