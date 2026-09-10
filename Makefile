install:
	cd frontend && npm install
run-server:
	cd server && go run server.go parsing.go commands.go handlers.go events.go combat.go logging.go quest.go
run-client:
	cd cli-client && go run cli_client.go
run-client-gui-dev:
	cd frontend && npm run dev
run-client-gui:
	cd frontend && npm run build && npm run preview
run-proxy:
	cd proxy && go run main.go
lint:
	golangci-lint run
	cd frontend && npm run lint
clean:
	rm -rf frontend/node_modules
