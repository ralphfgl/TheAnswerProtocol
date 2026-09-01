install:
	cd frontend && npm install
run-server:
	cd server && go run main.go
run-client:
	cd cli && go run main.go
run-client-gui:
	cd frontend && npm run dev
lint:
	golangci-lint run
	cd frontend && npm run lint
clean:
	rm -rf frontend/node_modules