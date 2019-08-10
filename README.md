# fdc-ingest
Provides a utility for loading the [USDA FDC](https://fdc.nal.usda.gov/download-datasets.html) csv files into a Couchbase datastore.  

# What's in the repo    
/ingest -- source for the ingest implemenations      
/couchbase -- shell scripts for initializing Couchbase

# Building   
The steps below outline how to go about building and running the utility using Couchbase.  

The build requires go version 12.  If you are using [Couchbase](https://www.couchbase.com/downloads), then version 6 or greater is preferred, version 5 will work as well as the community or enterprise edition.

### Step 1: Clone this repo
Clone this repo into any location other than your $GOPATH:
```
git clone git@github.com:littlebunch/fdc-ingest.git
```
and cd to the repo root, e.g.:
```
cd ~/fdc-ingest
```
      
### Step 2: Build the binary

The repo contains go.mod and supporting files so a build will automatically install and version all needed libraries.  If you don't want to use go mod then rm go.mod and have at it the old-fashioned way.   

Build data loader utility:   
```
go build -o $GOBIN/fdcloader loader.go
```
You're free to choose different names for -o binaries as you like. Or, just issue a go run.

### Step 3: Install [Couchbase](https://www.couchbase.com)     
If you do not already have access to a CouchBase instance then you will need to install at least version 5 or greater of the Community edition.  There are a number of easy deployment [options](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp) as well -- from a local workstation, docker or the public cloud.  Checkout the latter from [Google](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp), [Amazon](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp) and [Azure](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp).     

### Step 4:  Load the USDA csv data
1. From your Couchbase console or REST API using the scripts provided in this [repo](https://github.com/littlebunch/fdc-ingest/tree/master/couchbase), create a bucket, e.g. gnutdata, a user, e.g. gnutadmin, with the Application Access role to the bucket and indexes.  
2. Configure config.yml (see below) for host, bucket and user id/pw values you have selected.  A template is provided to get you started.

```
couchdb:   
  url:  localhost   
  bucket: gnutdata   //default  bucket    
  fts: fd_food  // default full-text index   
  user: <your_user>    
  pwd: <your_password>    

```
You can also set Environment variables which will override config.yml values:  
```
COUCHBASE_URL=localhost   
COUCHBASE_BUCKET=gnutdata   
COUCHBASE_FTSINDEX=fd_food   
COUCHBASE_USER=user_name   
COUCHBASE_PWD=user_password   
```
3. [Download](https://fdc.nal.usda.gov/download-datasets.html) and unzip the supporting data, BFPD, FNDDS and SR csv files into a location of your choice.   
4. Load the data files. (Note: the NUT and DERV types must be loaded first.)
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_Supporting_Data_csv/ -t NUT 
```
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_Supporting_Data_csv/ -t DERV
```
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_branded_food_csv/ -t BFPD    
```
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_survey_food_csv/ -t FNDDS  
```    
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_sr_csv_2019-04-02/ -t SR
``` 


