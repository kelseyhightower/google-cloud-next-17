# Next

The repo holds a simple demo app that connect to spanner and logs access events.

## Prep

Create a GKE cluster:

```
gcloud container clusters create us-west1 \
  --zone=us-west-1-a
```

### Create a Cloud Spanner Instance

```
gcloud beta spanner instances create google-cloud-next-demo \
  --config regional-us-central1 \
  --description "Google Cloud Next Demo" \
  --nodes 1
```

```
gcloud beta spanner databases create example \
  --instance google-cloud-next-demo
```

```
gcloud beta spanner databases ddl update example \
  --instance google-cloud-next-demo \
  --ddl='CREATE TABLE event (id STRING(MAX), message STRING(MAX), region STRING(MAX), timestamp TIMESTAMP) PRIMARY KEY (id)'
```

Store the Cloud Spanner instance ID in a config map:

```
DATABASE_ID=$(gcloud beta spanner instances describe google-cloud-next-demo \
  --format='value(name)')
```

```
kubectl create configmap spanner --from-literal=database-id=${DATABASE_ID}
```

### Create a Service Account to access Spanner

```
export PROJECT_ID=$(gcloud config get-value core/project)
```

```
export SERVICE_ACCOUNT_NAME="spanner-demo"
```

```
gcloud beta iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
  --display-name "spanner service account"
```

```
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role='roles/spanner.databaseUser'
```

```
gcloud beta iam service-accounts keys create \
  --iam-account "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  service-account.json
```

Store the service account in a Kubernetes secret:

```
kubectl create secret generic spanner --from-file service-account.json
```

### Deploy the example application

```
kubectl create -f deployments/next.yaml
```

```
kubectl create -f services/next.yaml
```
