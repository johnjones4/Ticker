
bin/ticker:
	go build -o bin/ticker .

ticker/provider/coingecko/swagger.json:
	curl https://raw.githubusercontent.com/coingecko/coingecko-api-oas/refs/heads/main/coingecko-public-api-v3.json > ticker/provider/coingecko/swagger.json

ticker/provider/noaa/swagger.json:
	curl https://api.weather.gov/openapi.json > ticker/provider/noaa/swagger.json

clients: ticker/provider/coingecko/swagger.json ticker/provider/noaa/swagger.json
	oapi-codegen -generate client,types,skip-prune -package yahoofinance ticker/provider/yahoofinance/swagger.yml > ticker/provider/yahoofinance/client.go
	oapi-codegen -generate client,types,skip-prune -package coingecko ticker/provider/coingecko/swagger.json > ticker/provider/coingecko/client.go
	oapi-codegen -generate client,types,skip-prune -package noaa ticker/provider/noaa/swagger.json > ticker/provider/noaa/client.go

install: bin/ticker
	mv bin/ticker /usr/bin/ticker
	cp ./res/env /etc/default/ticker
	cp ./res/ticker.service /etc/systemd/system/ticker.service
	cp config.json /etc/ticker.json
	systemctl daemon-reexec
	systemctl daemon-reload
	systemctl enable ticker.service
	systemctl start ticker.service
