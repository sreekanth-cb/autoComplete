# autoComplete

Sample project to demo type ahead or autocomplete using FTS/n-gram approach


Autocomplete suggestions are delivered using the Edge N gram based approach.
User has to define the n-gram (eg: min-2/max-10) token filter for the field on 
which the type ahead or suggestion is intended.

As a quick reminder, an edge-n-gram toekn filter will tokenize the input like below.

An edge-n-gram tokenizer of min length 2 and max length 6 would tokenize given 
the text “jurassic park” like below.

ju, jur, jura, juras, jurass, pa, par, park. 


The approach followed here involves the use of different index and query time analyzers.

During indexing, autocomplete intended fields has to be analyzed with a custom analyzer 
making use of edge-n-gram token filter. (should be chained with to_lower filter too if needed). 
And later during query (match query) time, simple analyzer should be used to prevent 
the unnecessary splitting of the query text.

Along with the match query, the client is requesting for the actual auto completed 
field contents using `Fields`  option in the search request and this value will be 
used as the auto completed text for the user.


For example, a match query for “jur”  or “pa” will match all of the below titles.

-Jurassic Park

-Jurassic Park III

-The Lost World: Jurassic Park

This approach has the benefits of,

Matches text from terms positioned at any places of the searched field.
This approach is generally fast for queries as it boils down to a normal term search.
Spelling mistakes can be taken care by the match query through the fuzzy option.

Note: This sample code is purely for demo purposes and no optimisations/clean ups are taken care.

## Steps to create sample auto complete app for the movie title field.

1. Install Couchbase server and create a couchbase bucket named `movies`.
2. Import the movies data from the movies_metadata.csv file to the movies bucket.
    eg: /cbimport csv -c http://127.0.0.1:8091 -u Username -p password -b movies --d file:////<path>/movies_metadata.csv -g %title% -t 4
3. Create an FTS index for the movie_title field on bucket movie.
    ```
    curl -XPUT -H "Content-Type: application/json" \
    -u <username>:<password> http://[::1]:8094/api/index/FTS -d \
    '{
    "type": "fulltext-index",
    "name": "FTS",
    "sourceType": "couchbase",
    "sourceName": "movies",
    "sourceUUID": "<<bucket_UUID_to_be_filled",
    "planParams": {
        "maxPartitionsPerPIndex": 171,
        "indexPartitions": 6
    },
    "params": {
        "doc_config": {
        "docid_prefix_delim": "",
        "docid_regexp": "",
        "mode": "type_field",
        "type_field": "type"
        },
        "mapping": {
        "analysis": {
            "analyzers": {
            "custom": {
                "char_filters": [
                "asciifolding"
                ],
                "token_filters": [
                "edgengram",
                "to_lower"
                ],
                "tokenizer": "unicode",
                "type": "custom"
            }
            },
            "token_filters": {
            "edgengram": {
                "back": "false",
                "max": 10,
                "min": 2,
                "type": "edge_ngram"
            }
            }
        },
        "default_analyzer": "standard",
        "default_datetime_parser": "dateTimeOptional",
        "default_field": "_all",
        "default_mapping": {
            "dynamic": false,
            "enabled": true,
            "properties": {
            "title": {
                "dynamic": false,
                "enabled": true,
                "fields": [
                {
                    "analyzer": "custom",
                    "include_in_all": true,
                    "index": true,
                    "name": "movie_title",
                    "store": true,
                    "type": "text"
                }
                ]
            }
            }
        },
        "default_type": "_default",
        "docvalues_dynamic": true,
        "index_dynamic": true,
        "store_dynamic": false,
        "type_field": "_type"
        },
        "store": {
        "indexType": "scorch",
        "kvStoreName": ""
        }
    },
    "sourceParams": {}
    }'

4. Just run the local client using -> `go run sample.go` (expected $GOPATH to be set)
5. Goto the url - http://localhost:12345/static/ in browser and start searching the movie titles.
