GITHUB_GRAPHQL_SCHEMA = ./github.gql
TOOLS_DIR = ./_tools
APOLLO_ROUTER = $(TOOLS_DIR)/router
ROVER = $(TOOLS_DIR)/rover
APOLLO_ROUTER_TMP_DIR := $(shell mktemp -d)
ROVER_TMP_DIR := $(shell mktemp -d)
SUPERGRAPH = ./supergraph.gql
SCHEMA = ./schema.gql

.PHONY: setup
setup: $(GITHUB_GRAPHQL_SCHEMA) $(APOLLO_ROUTER) $(ROVER) $(SUPERGRAPH)

$(GITHUB_GRAPHQL_SCHEMA):
	curl -sSL https://docs.github.com/public/schema.docs.graphql > $@

$(APOLLO_ROUTER):
	curl -sSL https://github.com/apollographql/router/releases/download/v1.0.0-alpha.3/router-v1.0.0-alpha.3-x86_64-apple-darwin.tar.gz | tar xzf - --strip-components=1 -C $(APOLLO_ROUTER_TMP_DIR)
	mkdir -p $(dir $@)
	cp $(APOLLO_ROUTER_TMP_DIR)/router $@

$(ROVER):
	curl -sSL https://github.com/apollographql/rover/releases/download/v0.8.2/rover-v0.8.2-x86_64-apple-darwin.tar.gz | tar xzf - --strip-components=1 -C $(ROVER_TMP_DIR)
	mkdir -p $(dir $@)
	cp $(ROVER_TMP_DIR)/rover $@

.PHONY: supergraph
supergraph: $(SUPERGRAPH)

$(SUPERGRAPH): $(ROVER) $(GITHUB_GRAPHQL_SCHEMA) $(SCHEMA)
	$(ROVER) supergraph compose --config supergraph.yml > $@
