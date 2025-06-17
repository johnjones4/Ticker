ticker/provider/coingecko/swagger.json:
	curl https://raw.githubusercontent.com/coingecko/coingecko-api-oas/refs/heads/main/coingecko-public-api-v3.json > ticker/provider/coingecko/swagger.json

clients: ticker/provider/coingecko/swagger.json
	oapi-codegen -generate client,types,skip-prune -package yahoofinance ticker/provider/yahoofinance/swagger.yml > ticker/provider/yahoofinance/client.go
	oapi-codegen -generate client,types,skip-prune -package coingecko ticker/provider/coingecko/swagger.json > ticker/provider/coingecko/client.go