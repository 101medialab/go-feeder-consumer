# go-feeder-consumer
Asynchroniously retrieve and consume data simultanously in multiple threads in Golang

Conside you need to read data from database, after fetching data you need to transform them into another structure then feeding them to another process, maybe calling webhook or indexing data into ElasticSearch.
You can get the job done for sure, but it can be very slow if you are not writing it properly. Of course you can user WatiGroup, channel, go routine to speed things up. However when you are going to add more and more steps into your workflow it can be troublesome and get complicated.

This package aims to create a component that is reusable for any workflow that need to do something first in parallel and then feeding data from this process to the next step which the next step need to further process these data in parallel again.
