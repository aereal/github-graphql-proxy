GITHUB_GRAPHQL_SCHEMA = ./github.gql
APOLLO_ROUTER = router
ROVER = rover
SUPERGRAPH = ./supergraph.gql
SCHEMA = ./schema.gql

.PHONY: setup
setup: $(GITHUB_GRAPHQL_SCHEMA) $(SUPERGRAPH)

$(GITHUB_GRAPHQL_SCHEMA):
	curl -sSL https://docs.github.com/public/schema.docs.graphql > $@

.PHONY: supergraph
supergraph: $(SUPERGRAPH)

$(SUPERGRAPH): $(GITHUB_GRAPHQL_SCHEMA) $(SCHEMA)
	$(ROVER) supergraph compose --config supergraph.yml > $@
