# TWReporter's Golang Backend API

## Functional Test
### Prerequisite
* Make sure the environment you run the test has a running `MySQL` server and `MongoDB` server
* Execute the following commands after logining into MySQL server. 
```
CREATE USER 'gorm'@'localhost' IDENTIFIED BY 'gorm';
CREATE DATABASE gorm;
GRANT ALL ON gorm.* TO 'gorm'@'localhost';
```

### How To Run Tests
```
go test $(glide novendor)

// or print logs
go test -v $(glide novendor)
```

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
- Read posts
- Read topics
- Read the combination of sections of index page
- Create/Read/Update/Delete bookmarks of a user
- Create/Read/Update/Delete registration(s)
- Create/Read/Update/Delete service(s)

### Read posts
- URL: `/v1/posts`
- Method: `GET`
- URL param:
  * Optional:
  `
  where=[string]
  offset=[integer]
  limit=[integer]
  sort=[string] 
  full=[boolean]
  `
  * Explain:
  `offset`: the number you want to skip
  `limit`: the number you want server to return
  `sort`: the field to sort by in the returned records
  `full`: if true, each record in the returued records will have all the embedded assets

  * example:
  `?where={"tags":{"$in":"57bab17eab5c6c0f00db77d1"}}&offset=10&limit=10&sort=-publishedDate&full=true` <br />
  this example will get 10 full records tagged by 57bab17eab5c6c0f00db77d1 and sorted by publishedDate ascendingly.

- Response:
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "records": [{
        // post data structure goes here
      }],
        "status": "ok"
    }
    ```
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Read topics
- URL: `/v1/topics`
- Method: `GET`
- URL param:
  * Optional:
  `
  where=[string]
  offset=[integer]
  limit=[integer]
  sort=[string] 
  full=[boolean]
  `
  * Explain:
  `offset`: the number you want to skip
  `limit`: the number you want server to return
  `sort`: the field to sort by in the returned records
  `full`: if true, each record in the returued records will have all the embedded assets

  * example:
  `?where={"slug":"far-sea-fishing-investigative-report"}&full=true` <br />
  this example will get 1 full topic.

- Response:
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "records": [{
        // topic goes here
      }],
        "status": "ok"
    }
    ```
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Read content of index page in the first screen
- URL: `/v1/index_page`
- Method: `GET`
- Response:
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "records": {
        "latest": [{
          // post goes here
        }, {
          // post goes here
        }, {
          // post goes here
        }, ... ], 
        "editor_picks": [{
          // post goes here
        }, {
          // post goes here
        }, {
          // post goes here
        }, ... ],
        "latest_topic": [{
          // topic goes here
        }],
        "reviews": [{
          // post goes here
        }, {
          // post goes here
        }, {
          // post goes here
        }, ... ]
      },
        "status": "ok"
    }
    ```
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Get bookmarks
- URL: `/v1/users/:userID/bookmarks`
  * example: `/v1/users/1/bookmarks`
- Authorization of Header: `Bearer ${JWT_TOKEN}`
- Method: `GET`

- Response: 
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "records": [{
        "id": bookmarkID_1,
        "created_at": "2017-05-09T11:42:50.084994666+08:00",
        "updated_at": "2017-05-09T11:42:50.084994666+08:00",
        "deleted_at": null,
        "href": "https://www.twreporter.org/a/about-us-footer",
	"title": "關於我們",
	"desc": "《報導者》是「財團法人報導者文化基金會」成立的非營利網路媒體...",
	"thumbnail": "https://www.twreporter.org/asset/logo-desk.svg"
      }, ... ],
        "status": "ok"
    }
    ```
  * **Code:** 401 <br />  
  * **Code:** 403 <br />
  * **Code:** 404 <br />
  **Content:** `{"status": "Record not found", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Create a bookmark
- URL: /users/:userID/bookmarks
- Authorization of Header: `Bearer ${JWT_TOKEN}`
- Content-Type of Header: `application/json`
- Method: `POST`
- Data Params:
```
{
   "href": "https://www.twreporter.org/a/about-us-footer",
   "title": "關於我們",
   "desc": "《報導者》是「財團法人報導者文化基金會」成立的非營利網路媒體...",
   "thumbnail": "https://www.twreporter.org/asset/logo-desk.svg"
}
```

- Response: 
  * **Code:** 201 <br />
    **Content:**
    ```
    {
        "status": "ok"
    }
    ```
  * **Code:** 400 <br />
  * **Code:** 401 <br />
  * **Code:** 403 <br />
  **Content:** `{"status": "Bad request", "error": "${here_goes_error_msg}"}`
  * **Code:** 404 <br />
  **Content:** `{"status": "Record not found", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`
  
### Delete a bookmark
- URL: /users/:userID/bookmarks/:bookmarkID
- Authorization of Header: `Bearer ${JWT_TOKEN}`
- Method: `DELETE`

- Response: 
  * **Code:** 204 <br />
  * **Code:** 401 <br />
  * **Code:** 403 <br />
  * **Code:** 404 <br />
  **Content:** `{"status": "Record not found", "error": "${here_goes_error_msg}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

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
