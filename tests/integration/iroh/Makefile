all: 

.PHONY: build
build:
	@sudo docker compose build
	
.PHONY: run
run: 
	@echo "Before this step you should run the connect and replicator commands manually (for now) from the output of the previous step." 
	@sudo docker compose up --force-recreate

.PHONY: graphql-a-add
graphql-a-add:
	@echo "NodeA add Article"
	@sudo docker compose exec nodeA defradb client --url localhost:9181 query 'mutation { create_Article(input: { content: "Hello world", published: true }) {_docID content published }}'

.PHONY: graphql-b-get
graphql-b-get:
	@echo "NodeB get Article"
	@sudo docker compose exec nodeB defradb client --url localhost:9182 query 'query { Article { _docID content published }}'