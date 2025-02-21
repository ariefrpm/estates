# Backend Engineering Interview Assignment (Golang)

## API Test Adjustments (Important)

I have adjusted the API test cases to align with the documentation:

- Removed the API test for the hello endpoint, as it was intended only as a guideline and can be safely deleted.
- Corrected the response status codes for Create Estate and Create Tree APIs to HTTP 201 (Created) as per requirements (previously HTTP 200).

API tests passed on my end.

```sh
go clean -testcache
go test ./tests/...
ok      github.com/SawitProRecruitment/EstateService/tests      0.313s
```

## Code Structure

The project follows a **hexagonal architecture** (similar to Clean Architecture) by separating business logic from external dependencies.

- Core (`domain`, `usecase`, `port/interfaces`) is independent of infrastructure.
- Adapters (`handler`, `storage`) depend on the core — not the other way around.

## Clear separation of concerns:

- `core/domain` → Pure business objects (independent of any external dependencies).
- `core/usecase` → Business logic, coordinating with repository interfaces.
- `core/interfaces` → Interfaces for inbound (usecase) and outbound (repo) dependencies.
- `storage/postgres` → Outbound adapter, implementing the repo interface.
- `handler` → Inbound adapter, consuming the usecase interface.

## Requirements

To run this project you need to have the following installed:

1. [Go](https://golang.org/doc/install) version 1.21
2. [GNU Make](https://www.gnu.org/software/make/)
3. [oapi-codegen](https://github.com/deepmap/oapi-codegen)

    Install the latest version with:
    ```
    go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
    ```
4. [mock](https://github.com/uber-go/mock)

    Install the latest version with:
    ```
    go install go.uber.org/mock/mockgen@latest
    ```

5. [Docker](https://docs.docker.com/get-docker/) version 20
   
   We will use this for testing your API.

6. [Docker Compose](https://docs.docker.com/compose/install/) version 1.29

7. [Node](https://nodejs.org/en) v20

   We will use this for testing your API

8. [NPM](https://www.npmjs.com/) v10

    We will use this for testing your API.

## Initiate The Project

To start working, execute

```
make init
```

## Running

You should be able to run using the script `run.sh`:

```bash
./run.sh
```

You may see some errors since you have not created the API yet.

However for testing, you can use Docker run the project, run the following command:

```
docke -compose up --build
```

You should be able to access the API at http://localhost:8080

If you change `database.sql` file, you need to reinitate the database by running:

```
docker compose down --volumes
```

## Testing

To run test, run the following command:

```
make test
```
