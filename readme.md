# A Google Drive style API, using OpenFGA, written in Go

## Start your FGA Server
More info on [openfga.dev](https://openfga.dev)

The FGA Server needs to be running before this Go server, cause on start it will create a new model and populate it with the model in [`fga-model.json`](./fga-model.json).

```
docker run -p 8080:8080 -p 8081:8081 -p 3000:3000 openfga/openfga run
```
## Add your environment variables
```
cp .env.sample .env
```

## Start your Go server
```
go run main.go
```

## API Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET`     | `/documents/:id`       | Get a document |
| `POST`    | `/documents`           | Create a new document, the JSON body accepts `name` and `content` |
| `POST`    | `/documents/:id/share` | Share a document, the JSON body accepts `relation` and `user`. You can also use this endpoint to add a file to a parent folder by setting the `relation` to `parent` and the `user` to `folder:<id>` |
| `GET`     | `/folders/:id`         | Get a folder |
| `POST`    | `/folders`             | Create a new folder, the JSON body accepts `name` |
| `POST`    | `/folderws/:id/share`  | Share a folder, the JSON body accepts `relation` and `user`. You can also use this endpoint to add a folder to a parent folder by setting the `relation` to `parent` and the `user` to `folder:<id>` |
| `GET`    | `/documents` | Get all documents. This endpoint is included for debugging, no FGA check is done here |
| `GET`    | `/folders`   | Get all folder. This endpoint is included for debugging, no FGA check is done here |

## Test tokens (These don't work anywhere else).
The API uses the subject (`sub`) from the JWT formatted Access Token persented as Bearer token in the `Authorization` header. This is a demo, and therefor not all necessary checks are in place. Please don't use this in production!

### Sam
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJzYW1Ab2t0YS5jb20iLCJuYW1lIjoiU2FtIEJlbGxlbiIsImlhdCI6MTUxNjIzOTAyMn0.UgLEipGU-69_dKuhgCsV7mrBcCvRJBV880kuMJbLBy8
```

### Chiara
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJjaGlhcmFAb2t0YS5jb20iLCJuYW1lIjoiQ2hpYXJhIiwiaWF0IjoxNTE2MjM5MDIyfQ.-mbo6VBG1xZzK-T7bEuKqUQR1B-eu-ACRIKrtai1JEU
```

## Test scenario
### `1` Create a new document for user Sam
`POST /documents/`
```
{
  "name": "Test document",
  "content": "Test content"
}
```
### `2` User Sam should be able to see the document 
`GET /documents/:id`
### `3` User Chiara should NOT be able to see the document 
`GET /documents/:id`
### `4` Create a new folder for user Sam
`POST /folders`
```
{
  "name": "Test Folder",
}
```
### `5` User Sam should be able to see the folder 
`GET /folders/:id`
### `6` User Chiara should NOT be able to see the document
`GET /folders/:id`
### `7` Add the document to the folder for user Sam
`POST /documents/:id/share`
```
{
  "relation": "parent",
  "user": "folder:id
}
```
### `8` Share the folder with user Chiara
`POST /folders/:id/share`
```
{
  "relation": "viewer",
  "user": "user:chiara@okta.com
}
```
