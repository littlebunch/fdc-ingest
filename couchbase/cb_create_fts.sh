#!/bin/sh
###########################################################################
# creates a full text index required for API search                       #
# run this after data is fully loaded. change -u values and url as needed # 
###########################################################################
curl -u administrator:administrator -XPUT http://localhost:8094/api/index/fts_srlegacy -H 'cache-control: no-cache' -H 'content-type:application/json' -d '
{
  "type": "fulltext-index",
  "name": "fts_srlegacy",
  "sourceType": "couchbase",
  "sourceName": "srlegacy",
  "planParams": {
    "maxPartitionsPerPIndex": 171
  },
  "params": {
    "doc_config": {
      "docid_prefix_delim": "",
      "docid_regexp": "",
      "mode": "type_field",
      "type_field": "type"
    },
    "mapping": {
      "analysis": {},
      "default_analyzer": "standard",
      "default_datetime_parser": "dateTimeOptional",
      "default_field": "_all",
      "default_mapping": {
        "dynamic": true,
        "enabled": false
      },
      "default_type": "_default",
      "docvalues_dynamic": true,
      "index_dynamic": true,
      "store_dynamic": false,
      "type_field": "_type",
      "types": {
        "FOOD": {
          "dynamic": false,
          "enabled": true,
          "properties": {
            "company": {
              "dynamic": false,
              "enabled": true,
              "fields": [
                {
                  "analyzer": "keyword",
                  "include_term_vectors": true,
                  "index": true,
                  "name": "company_kw",
                  "type": "text"
                },
                {
                  "include_in_all": true,
                  "include_term_vectors": true,
                  "index": true,
                  "name": "company",
                  "store": true,
                  "type": "text"
                }
              ]
            },
            "dataSource": {
              "dynamic": false,
              "enabled": true,
              "fields": [
                {
                  "include_in_all": true,
                  "include_term_vectors": true,
                  "index": true,
                  "name": "dataSource",
                  "store": true,
                  "type": "text"
                }
              ]
            },
            "fdcId": {
              "dynamic": false,
              "enabled": true,
              "fields": [
                {
                  "include_in_all": true,
                  "include_term_vectors": true,
                  "index": true,
                  "name": "fdcId",
                  "store": true,
                  "type": "text"
                }
              ]
            },
            "foodDescription": {
              "dynamic": false,
              "enabled": true,
              "fields": [
                {
                  "include_in_all": true,
                  "include_term_vectors": true,
                  "index": true,
                  "name": "foodDescription",
                  "store": true,
                  "type": "text"
                },
                {
                  "analyzer": "keyword",
                  "include_term_vectors": true,
                  "index": true,
                  "name": "foodDescription_kw",
                  "type": "text"
                }
              ]
            },
            "ingredients": {
              "dynamic": false,
              "enabled": true,
              "fields": [
                {
                  "analyzer": "keyword",
                  "include_term_vectors": true,
                  "index": true,
                  "name": "ingredients_kw",
                  "type": "text"
                },
                {
                  "include_in_all": true,
                  "include_term_vectors": true,
                  "index": true,
                  "name": "ingredients",
                  "store": true,
                  "type": "text"
                }
              ]
            },
            "type": {
              "dynamic": false,
              "enabled": true,
              "fields": [
                {
                  "include_in_all": true,
                  "include_term_vectors": true,
                  "index": true,
                  "name": "type",
                  "store": true,
                  "type": "text"
                }
              ]
            },
            "upc": {
              "dynamic": false,
              "enabled": true,
              "fields": [
                {
                  "analyzer": "keyword",
                  "include_term_vectors": true,
                  "index": true,
                  "name": "upc_kw",
                  "type": "text"
                },
                {
                  "include_in_all": true,
                  "include_term_vectors": true,
                  "index": true,
                  "name": "upc",
                  "store": true,
                  "type": "text"
                }
              ]
            }
          }
        }
      }
    },
    "store": {
      "indexType": "scorch",
      "kvStoreName": ""
    }
  },
  "sourceParams": {}
}'
