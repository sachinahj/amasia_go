SELECT lg.*
FROM yelp.Log lg
WHERE lg.IsDone=false
;

SELECT DISTINCT ZipCode
FROM amasia.ZipCodeCategoriesLevel4
;

SELECT DISTINCT ZipCode
FROM yelp.Log
;

SELECT count(*)
FROM yelp.Log
WHERE type="businesses_search"
AND IsDone=true
AND CreatedAt > "2017-02-28 00:00:00"
AND CreatedAt < "2017-02-29 00:00:00"
;
