# TWReporter's Golang Backend API

## Configurations

#### MySQL connection
Copy `configs/config.example.json` and rename as `configs/config.json`. Change its content to connect to your own database.

## Development
Please make sure that you install [Glide
  package manager](https://github.com/Masterminds/glide) in the environment.

```
cd $GOPATH/src/twreporter.org/go-api
glide install                           # Install packages and dependencies
go run main.go                          # Run without live-reloading
```

#### Live Reloading
Note that `GOPATH/bin` should be in your `PATH`.
```
go get github.com/codegangsta/gin
gin                                     # Run with live-reloading
```


## Production
```
go build
./go-api
```


## Testing
```
$ go test $(glide novendor)             # run go test over all directories of the project except the vendor directory
```

## RESTful API
`go-api` is a RESTful API built by golang.

It provides several RESTful web services, including
- User login/logout/signup
- Create/Read/Update/Delete bookmarks of a user
- Create/Read/Update/Delete registration(s)
- Create/Read/Update/Delete service(s)

### Create a service
- URL: `/v1/services/`
- Content-Type of Header: `application/json`
- Method: `POST`
- Data Params:
```
{
  "name": "news_letter"
}
```
- Response: 
  * **Code:** 201 <br />
    **Content:**
    ```
    {
      "record": {
        "ID": 1,
        "CreatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "UpdatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "DeletedAt": null,
        "Name": "news_letter"
      },
        "status": "ok"
    }
    ```
  * **Code:** 400 <br />
  **Content:** `{"status": "Bad request", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Read a service
- URL: `/v1/services/:id`
- Method: `GET`
- Response: 
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "record": {
        "ID": 1,
        "CreatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "UpdatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "DeletedAt": null,
        "Name": "news_letter"
      },
        "status": "ok"
    }
    ```
  * **Code:** 404 <br />
  **Content:** `{"status": "Resource not found", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Update a service
Update a service or create a service if not existed
- URL: `/v1/services/:id`
- Method: `PUT`
- Response: 
- Data Params:
```
{
  "name": "news_letter"
}
```
- Response: 
  * **Code:** 200<br />
    **Content:**
    ```
    {
      "record": {
        "ID": 1,
        "CreatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "UpdatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "DeletedAt": null,
        "Name": "news_letter"
      },
        "status": "ok"
    }
    ```
  * **Code:** 201<br />
    **Content:**
    ```
    {
      "record": {
        "ID": 1,
        "CreatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "UpdatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "DeletedAt": null,
        "Name": "news_letter"
      },
        "status": "ok"
    }
    ```
  * **Code:** 400 <br />
  **Content:** `{"status": "Bad request", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Delete a service
- URL: `/v1/services/:id`
- Method: `DELETE`
- Response: 
  * **Code:** 204 <br />
  * **Code:** 404 <br />
  **Content:** `{"status": "Resource not found", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Create a registration
- URL: `/v1/registrations/:service/`
  * example: `/v1/registrations/news_letter/`
- Content-Type of Header: `application/json`
- Method: `POST`
- Data Params:
```
{
  "email": "nickhsine@twreporter.org"
}
```
- Response: 
  * **Code:** 201 <br />
    **Content:**
    ```
    {
      "record": {
        "CreatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "UpdatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "DeletedAt": null,
        "Email": "nickhsine@twreporter.org",
        "Service": "news_letter",
        "Active": false,
        "ActivateToken": ""
      },
        "status": "ok"
    }
    ```
  * **Code:** 400 <br />
  **Content:** `{"status": "Bad request", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Read a registration
- URL: `/v1/registrations/:service/:email`
  * example: `/v1/registrations/news_letter/nickhsine%40twreporter.org`
- Method: `GET`
- Response: 
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "record": {
        "CreatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "UpdatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "DeletedAt": null,
        "Email": "nickhsine@twreporter.org",
        "Service": "news_letter",
        "Active": false,
        "ActivateToken": ""
      },
        "status": "ok"
    }
    ```
  * **Code:** 404 <br />
  **Content:** `{"status": "Resource not found", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Read registrations
- URL: `/v1/registrations/:service`
  * example: `/v1/registrations/news_letter`
- Method: `GET`
- URL param: 
  * Optional: 
    `
    offset=[integer] 
    limit=[integer]
    order_by=[string]
		active_code=[integer]
    `
  * example: 
    `?offset=10&limit=10&order=updated_at&active_code=2`
- Response: 
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "record": {
        "CreatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "UpdatedAt": "2017-05-09T11:42:50.084994666+08:00",
        "DeletedAt": null,
        "Email": "nickhsine@twreporter.org",
        "Service": "news_letter",
        "Active": false,
        "ActivateToken": ""
      },
        "status": "ok"
    }
    ```
  * **Code:** 400 <br />
  **Content:** `{"status": "Bad request", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Delete a registration
- URL: `/v1/registrations/:service/:email`
  * example: `/v1/registrations/news_letter/nickhsine%40twreporter.org`
- Method: `DELETE`
- Response: 
  * **Code:** 204 <br />
  * **Code:** 404 <br />
  **Content:** `{"status": "Resource not found", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Activate(Update) a registration
After register a service such as news_letter, the system will send a email to the user for activating the registration.<br />
The email will contain a link like `go-api.twreporter.org/v1/activation/news_letter/nickhsine%40twreporter.org?activateToken=${here_goes_the_token}`.<br />
When user clicks the link, the registration will be activated.

- URL: `/v1/activation/:service/:userEmail`
	* example: `/v1/registrations/news_letter/nickhsine%40twreporter.org`
- Method: `GET`
- Response:
	* **Code: ** 307 <br /> 
	Redirect to the front-end website. If the activation fails, the redirect url will have error URL param.<br />
	For example https://www.twreporter.org/?error=Account+token+is+not+correct&error_code=403
