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

It helps in matching the user entered partial text to tokens at any position in the indexed content.


For example, a match query for “jur”  or “pa” will match all of the below titles.

-Jurassic Park

-Jurassic Park III

-The Lost World: Jurassic Park

This approach has the benefits of,

Matches text from terms positioned at any places of the searched field.
This approach is generally fast for queries as it boils down to a normal term search.
Spelling mistakes can be taken care by the match query through the fuzzy option.

Note: This sample code is purely for demo purposes and no optimisations/clean ups are taken care.
