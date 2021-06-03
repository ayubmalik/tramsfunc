# tramsfunc

A Google Cloud Plaform (GCP) function to provide Transport for Greater Manchester (TFGM) Metrolink tram information.

This function is not intended for end users. Rather it provides the backend API, so that _clients_ can call it without
having to register their own [TFGM API key](https://developer.tfgm.com/).

To deploy this function into GCP a TFGM API key **is** needed.

## Running locally

This project uses [Function Frameworks](https://cloud.google.com/functions/docs/running/function-frameworks) to run the
function locally.

```
# start local server
TFGM_API_KEY='some api key' go run cmd/local-server/main.go
```

```
# call the API
curl http://localhost:8080/tramsfunc
```
