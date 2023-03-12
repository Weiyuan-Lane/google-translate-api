
# google-translate-api

This is an exploratory self project, written in go, exploring the v2 (simple) and v3 (advanced) versions of the Google Translate API

---

### Running this project

To run this project, you need to have `docker` installed, as well as having an API key and service account created for using V2 and V3 of the translate APIs. Place the values in the `.env` file.

To start running the local server at port `8080`, simply run:
```
docker-compose up
```


---

### Cloud Run

The V3 of the Translate API works without an API key, as long as a service account with the right permissions is assigned. It works because of  [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials)

To build this project and testing in Cloud Run for V3 of the translate API (without using API Key), follow the following steps:

(make sure you replace the `%GCP-PROJECT-ID%` value to your project id)

1. Build docker image
```
docker build -t gcr.io/%GCP-PROJECT-ID%/google-translate-api:0.0.1 --platform linux/amd64 -f tools/Dockerfile.cloudrun .
```

2. Push the image to your image registry in your project 
```
docker push gcr.io/%GCP-PROJECT-ID%/google-translate-api:0.0.2
```

3. Visit [Cloud Run](https://console.cloud.google.com/run), and create a new service using your new image and service account (which has the translate permissions assigned)

Viola! You're done!
