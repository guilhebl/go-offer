# GO-OFFER

Go Offer is an open source project that aims to bring together several world-wide marketplaces and provide best offers and results in a single consolidated platform.
It's goal is to unify diverse market sources in order to connect people directly with products and services that they need through a global marketplace platform.
Connecting users to local and international marketplaces and enabling the consumer to receive informative insights of a product or service in a transparent manner.
Find out about prices, availability, features, reviews and more searching in one place and receiving the best results from multiple sources.
Join the project and help make the world a better place for trade.


##### Checkout [Current LIVE beta](https://searchprod.com)

&nbsp;&nbsp;

##### Current stack:

. [Frontend](https://github.com/guilhebl/offer-web) -> Angular4, Typescript, NgRx, Gulp

. [Backend](https://github.com/guilhebl/offer-java) -> Java

##### Future stack:

. [Frontend](https://github.com/guilhebl/offer-web) -> Angular5, Typescript, NgRx, Webpack

. [Backend](https://github.com/guilhebl/go-offer) -> Go

## Getting started

in order to get started you must first create your own API keys at:

- Amazon Product Advertising API

- Ebay Product Search API

- Walmart API

- BestBuy API

After creating your api keys set the values in "app-config.properties" file replacing proper entries that have the string "TEST123456789"
with the appropriate key values created in previous step.


### building

go build

### running

after building run command

./go-offer

### static folder

Static files such as HTML,CSS,JS files are located inside the `static` folder


make sure cache is disabled when running tests.


### setup DB

Setup and install [apache cassandra](https://linode.com/docs/databases/cassandra/deploy-scalable-cassandra/)

[Go Client](https://academy.datastax.com/resources/getting-started-apache-cassandra-and-go)


cqlsh -u cassandra -p cassandra
 CREATE ROLE [new_superuser] WITH PASSWORD = '[secure_password]' AND SUPERUSER = true AND LOGIN = true;
 ALTER ROLE cassandra WITH PASSWORD = 'cassandra' AND SUPERUSER = false AND LOGIN = false;
 REVOKE ALL PERMISSIONS ON ALL KEYSPACES FROM cassandra;
 GRANT ALL PERMISSIONS ON ALL KEYSPACES TO [superuser];

create 2 keyspace, one for production and one for test:

`
CREATE KEYSPACE atlanteus
  WITH REPLICATION = {
   'class' : 'SimpleStrategy',
   'replication_factor' : 1
  };
`

create test keyspace:

`
CREATE KEYSPACE test
  WITH REPLICATION = {
   'class' : 'SimpleStrategy',
   'replication_factor' : 1
  };
`

check if keyspaces were correctly created using `DESCRIBE keyspaces;`

Run:

`CREATE ROLE admin WITH PASSWORD = '[strong password]' AND SUPERUSER = true AND LOGIN = true;`


### setup cache (optional)

Setup and install REDIS cache as described [here](http://www.geekpills.com/operating-system/linux/install-configure-redis-ubuntu-17-10)

to flush cache: `redis-cli FLUSHDB`

make sure cache is disabled when running tests.


### REST API methods

1. Search Trending offers:

```
GET localhost:8080/offers
```

2. Search offers by keyword:

```
curl -H "Content-Type: application/json" -X POST -d '{ "searchColumns":[ { "name":"name", "value":"skyrim" } ], "sortOrder":"asc", "page":1, "rowsPerPage":10 }' http://localhost:8080/offers
```

3. Get Product Detail

```
GET localhost:8080/offers/887276234465?idType=upc&source=walmart.com
```

4. Search offers from datastore (cassandra)

```
GET localhost:8080/offerlist
```


5. Add offers to datastore (cassandra)

```
curl -H "Content-Type: application/json" -X POST -d '{"upc":"upc1","name":"test record","partyName":"amazon.com","semanticName":"http:/item01","mainImageFileUrl":"http:/item01.jpg","partyImageFileUrl":"amazon-logo.jpg","productCategory":"laptops","price":500,"rating":3.88,"numReviews":120}' http://localhost:8080/offerlist
```

### testing

to run main functional tests 
stack must be running without REDIS cache, but with Cassandra DB on as it runs some tests against the database, 
from main folder type:

1. `go test`

to run all tests
2. `go test ./...`

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to the project, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Pull requests need to be based on and opened against the `master` branch.

## License

The offergo source code is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html)
