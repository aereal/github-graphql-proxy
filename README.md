# github-graphql-proxy

## setup

```sh
make setup
```

### run [Apollo Router][]

```sh
./_tools/router -c router.yml -s supergraph.gql
```

### update [supergraph][] schema

```sh
make supergraph
```

[Apollo Router]: https://www.apollographql.com/docs/router/quickstart/
[supergraph]: https://www.apollographql.com/docs/federation/federated-types/overview#supergraph-schema
