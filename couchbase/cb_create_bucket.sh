#!/bin/sh
######################################################
# Creates a bucket named 'gnutdata' on local cluster #
# change the -u and URL values as required           #
######################################################
curl -u Administrator:password -X POST http://10.0.0.183:8091/pools/default/buckets -d name=bfpddata -d ramQuotaMB=2000 -d bucketType=couchbase -d authType=none
# Creates a bucket_admin user for the bucket
# API call is only available on the Couchbase EE edition
#curl -u Administrator:password -X PUT --data "name=Admin&password=gnutadmin&roles=Application_Access[gnutbfpd]" -H "Content-Type: application/x-www-form-urlencoded" http://localhost:8091/settings/rbac/users/local/gnutadmin
