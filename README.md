# fdc-ingest
Provides a utility for loading the [USDA FDC](https://fdc.nal.usda.gov/download-datasets.html) csv files into a Couchbase datastore.  

# What's in the repo    
/ingest -- source for the ingest implemenations      
/couchbase -- shell scripts for initializing Couchbase

# Building   
The steps below outline how to go about building and running the utility using Couchbase.  

The build requires [Go](https://golang.org/dl/) version 11 or greater.  If you are using [Couchbase](https://www.couchbase.com/downloads), then version 6 or greater is preferred, version 5 will work as well as the community or enterprise edition.

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

### Step 3: Install and Configure [Couchbase](https://www.couchbase.com)     
If you do not already have access to a CouchBase instance then you will need to install at least version 5 or greater of the Community edition.  There are a number of easy deployment [options](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp) as well -- from a local workstation, docker or the public cloud.  Checkout the latter from [Google](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp), [Amazon](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp) and [Azure](https://resources.couchbase.com/cloud-partner-gcp/docs-deploy-gcp).    
1. From your Couchbase console or REST API using the scripts provided in this [repo](https://github.com/littlebunch/fdc-ingest/tree/master/couchbase), create: 1) a bucket, e.g. fdcdata, 2) a user, e.g. fdcadmin, with the Application Access role to the bucket and 3) the indexes.  Note, if you are using the provided script to create the indexes, then be sure to run the build command before loading data.   
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

### Step 4:  Load the USDA csv data

1. [Download](https://fdc.nal.usda.gov/download-datasets.html) and unzip the supporting data, BFPD, FNDDS and SR csv files into a location of your choice.   
2. Load the supporting data files.   (Note: the NUT and DERV types must be loaded first.)
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_Supporting_Data_csv/nutrient.csv -t NUT 
```
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_Supporting_Data_csv/food_nutrient_derivation.csv -t DERV
```
If you are loading Standard Release Legacy data then it's recommended that you load food category data before loading the food data:
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_Supporting_Data_csv/food_category.csv -t FGSR   
```
And,if you are loading Food Survey data then it's recommend that you load WWEIA food category data before loading the food data:
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_Supporting_Data_csv/wweia_food_category.csv -t FGSR   
```
3. Load the data files.    
You can, of course, pick and choose which food datasets to load.  Each may be loaded into the same bucket or can be loaded into separate buckets depending on your application needs.   If loading into separate buckets then be sure to change your config.yml or environment variables before each run.   

Here's how to load Branded Food Products: 
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_branded_food_csv/ -t BFPD    
```
Or, Food Survey:
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_survey_food_csv/ -t FNDDS  
``` 
Or, Standard Release Legacy:
```
$GOBIN/fdcloader -c /path/to/config.yml -i /path/to/FoodData_Central_sr_csv_2019-04-02/ -t SR
``` 
Loading is not particulary fast, especially for Branded Food Products which can take over an hour on a standard MacBook Pro with 8GB RAM and about 40 minutes on a Linux machine with 16GB RAM and 4 cores.  There's probably a fair amount of room for improvement, especially in taking advantage of Go concurrancy.   

