# GO-OFFER

Go Offer is an open source project that aims to bring together several world-wide marketplaces and provide best offers and results in a single consolidated platform.
It's goal is to unify diverse market sources in order to connect people directly with products and services that they need through a global marketplace platform.
Connecting users to local and international marketplaces and enabling the consumer to receive informative insights of a product or service in a transparent manner.
Find out about prices, availability, features, reviews and more searching in one place and receiving the best results from multiple sources.
Join the project and help make the world a better palce for trade.

***
##### Current beta: [Search Prod](https://searchprod.com)
***
[Frontend Javascript/Angular](https://github.com/guilhebl/offer-web)
[Scala backend version](https://github.com/guilhebl/offer-backend)

## Getting started

in order to get started you must first create your own API keys at:

- Amazon Product Advertising API

- Ebay Product Search API

- Walmart API

- BestBuy API

after creating your api keys set the values in "app-config.properties" file


### building

go build

### running

after building run command

./go-offer


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

### testing

go test ./...


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