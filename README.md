# TWReporter's Golang Backend API

## Environment 
### Development
Please make sure that you install [Glide
  package manager](https://github.com/Masterminds/glide) in the environment.

```
cd $GOPATH/src/twreporter.org/go-api
glide install                           # Install packages and dependencies

// use Makefile
make start 
// or 
go run main.go                          # Run without live-reloading
```

### Live Reloading
Note that `GOPATH/bin` should be in your `PATH`.
```
go get github.com/codegangsta/gin
gin                                     # Run with live-reloading
```

### Production
```
go build
./go-api
```

## Dependencies Setup and Configurations
There are two major dependencies of go-api, one is MySQL database, 
another is MongoDB. <br/>
MySQL DB stores membership data, which is related to users.<br/>
MongoDB stores news entities, which is the content that go-api provides.<br/>

### Install docker-compose
[docker-compose installation](https://docs.docker.com/compose/install/) 

### Start/Stop MySQL and MongoDB with default settings
```
// start MySQL and MongoDB
make env-up

// stop MySQL and MongoDB
make env-down
```

### Configure MySQL Connection
Copy `configs/config.example.json` and rename as `configs/config.json`.
Change `DBSettings` fields to connect to your own database, like following example.
```
  "DBSettings": {
    "Name":     "test_membership",
    "User":     "test_membership",
    "Password": "test_membership",
    "Address":  "127.0.0.1",
    "Port":     "3306"
  },
```

### Configure MongoDB Connection
Copy `configs/config.example.json` and rename as `configs/config.json`.
Change `MongoDBSettings` fields to connect to your own database, like following example.
```
  "MongoDBSettings": {
    "URL": "localhost",
    "DBName": "plate",
    "Timeout": 5
  },
```

### AWS SES Setup
Currently the source code sends email through AWS SES,

If you want to send email through your AWS SES, just put your AWS SES config under `~/.aws/credentials`
```
[default]
aws_access_key_id = ${AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${AWS_SECRET_ACCESS_KEY}
```

Otherwise, you have to change the `utils/mail.go` to integrate with your email service.

## Functional Testing
### Prerequisite
* Make sure the environment you run the test has a running `MySQL` server and `MongoDB` server<br/>

### How To Run Tests
```
// use Makefile
make test

// or

go test $(glide novendor)

// or print logs
go test -v $(glide novendor)
```

## RESTful API
`go-api` is a RESTful API built by golang.

It provides several RESTful web services, including
- [User login/oauth(facebook & google)](https://github.com/twreporter/go-api#users)
- [Read posts](https://github.com/twreporter/go-api#read-posts)
- [Read topics](https://github.com/twreporter/go-api#read-topics)
- [Read the combination of sections on index page](https://github.com/twreporter/go-api#read-posts-of-latest-editor-picked-latest-topic-reviews-topics-photography-and-infographic-sections-of-index-page)
- [Read the posts of multiple categories on index page](https://github.com/twreporter/go-api#read-posts-of-character-culture_movie-human_rights-international-land_environment-photo_audio-political_society-and-transformed_justice-categories)
- [Create/Read/Update/Delete bookmarks of a user](https://github.com/twreporter/go-api#bookmarks)
- [Create/Read/Update/Delete registration(s)](https://github.com/twreporter/go-api#registrations)
- [Create/Read/Update/Delete service(s)](https://github.com/twreporter/go-api#services)

## USERS
### Signin
- workflow: 
	1. user send `POST` request to `v1/signin` endpoint
	2. system will send activation email to user
	3. user click activation link(`<a>` link) in the email body
	4. go-api server verifies the token
	5. if verified, user will get a jwt(Json Web Token)
	6. user can use JWT to send personal requests(go-api server will verfiy the jwt).
	
### OAuth
Before Oauth signin, you have to setup the oauth config in `configs/config.json`
```
  "OauthSettings": {
    "FacebookSettings": {
      "ID": "${ID_YOU_GET_FROM_FACEBOOK_DEVELOPER}",
      "Secret": "${SECRECT_YOU_GET_FROM_FACEBOOK_DEVELOPER}",
      "URL": "http://${GO_API_SERVER_HOST_NAME}:8080/v1/auth/facebook/callback",
      "Statestr": "${THE_STATE_YOU_WANT_TO_USE_IN_AUTHORIZE_URL}"
    },
    "GoogleSettings": {
      "Id": "${ID_YOU_GET_FROM_GOOGLE_DEVELOPER}",
      "Secret": "${SECRECT_YOU_GET_FROM_FACEBOOK_DEVELOPER}",
      "Url": "http://${GO_API_SERVER_HOST_NAME}:8080/v1/auth/google/callback",
      "Statestr": "THE_STATE_YOU_WANT_TO_USE_IN_AUTHORIZE_URL"
    }
  },
  "ConsumerSettings": {
    "Domain": "${CONSUMER_DOMAIN_NAME}",
    "Protocol": "http",
    "Host": "${CONSUMER_HOST_NAME}",
    "Port": "3000"
  },
```
- workflow
	1. users click oauth login button, broswer send GET request to `/v1/oauth/goolge` or `/v1/oauth/facebook` endpoints(on go-api server)
	2. go-api server redirect users to google or facebook oauth confirmation page
	3. on goolge/facebook oauth page, user input account and password
	4. if verified by facebook/google, facebook/google will redirect user to `/v1/oauth/google/callback` or `/v1/oauth/facebook/callback` endpoints on go-api server.
	5. if verified by go-api server, go-api server will redirect user to customer page(here, will be `${ConsumerSettings.Protocol}://${ConsumerSettings.Host}:${ConsumerSettings.Port}/`).
	6. jwt(Json Web Token) will be set in the response header(`Set-Cookie: ${cookie}`), and user can get the jwt from browser `cookie`.
	7. user can use JWT to send personal requests(go-api server will verfiy the jwt).
### Example
[TWReporter main site](https://www.twreporter.org/signin) is using the above workflow, you can try to signin on our site.

### Signin Endpoint
- URL: `/v1/signin`
- Method: `POST`
- Data Params:
```
{
  "email": "nickhsine@twreporter.org",
  "destination": "https://www.twreporter.org"
}
```
	- Required: `email` 
	- Optional:`destination`
	- Explain: 
	
	`email` is the user email, the activation email will be sent to.
	
	`destionation` is the redirect URL after user signed in.
  
- Response: 
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "data": {
      	"email": "nickhsine@twreporter.org",
	"destination": "https://www.twreporter.org"
      },
        "status": "success"
    }
    ```
  * **Code:** 400 <br />
  **Content:** `{"status": "fail", "data": "{"email":"email is required", "destination":"destination is optional"}"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "error", "message": "Internal server error: Sending activation email occurs error"}`
   
### User Activation Endpoint
- URL: `/v1/activate`
- Method: `GET`
- URL Param:
  * Required: `email` and `token`
- Response:
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "status": "success",
      "id": "USER_ID",
      "privilege": "PRIVILEGE",
      "firstname": "Nick",
      "lastname": "Li",
      "email": "nickhsine@twreporter.org",
      "jwt": "JSON_WEB_TOKEN"
    }
    ```
  * **Code:** 401 <br />
  **Content:** `{"status": "error", "message": "ActivateToken is expired"}`
  * **Code:** 500 <br />
  **Content:** `{"status": "error", "message": "Generating JWT occurs error"}`
  

### Renew JWT Endpoint
- URL: `/v1/token/:userID`
  * example: `/v1/token/100`
- Method: `GET`
- Response:
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "status": "success",
      "data": {
      	"token": "NEW_JSON_WEB_TOKEN",
	"token_type": "Bearer"
      }
    }
    ```
  * **Code:** 401 <br />
  **Content:** `{"status": "error", "message": ""}`
  * **Code:** 500 <br />
  **Content:** `{"status": "error", "message": "Renewing JWT occurs error"}`

### OAuth Endpoints
- URL: `/v1/auth/google` | `/v1/auth/facebook`
- Method: `GET`
- Response:
  * **Code:** 302 <br />
    **Header:**
    ```	
    "Set-Cookie: auth_info={\"id\":100,\"privilege\":0,\"firstname\":\"\",\"lastname\":\"\",\"email\":\"nickhsine97753017@gmail.com\",\"jwt\":\"jwt_token_goes_here\"}; Domain=twreporter.org; Max-Age=100 HttpOnly"
    ```
    **Redirect URL:** `http://testtest.twreporter.org:3000/?login=google`
  * **Code:** 401 <br />
  **Content:** `{"status": "error", "message": ""}`
  * **Code:** 500 <br />
  **Content:** `{"status": "error", "message": "Renewing JWT occurs error"}`

## POSTS 
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

## TOPICS
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

## INDEX_PAGE
### Read posts of latest, editor picked, latest topic, reviews, topics, photography and infographic sections of index page 
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
        } ... ],
        "topics": [{
          // topic goes here
        }, {
          // topic goes here
        } ... ],
        "photos": [{
          // post goes here
        }], 
        "infographics": [{
          // post goes here
        },{
          // post goes here
        }]
      },
        "status": "ok"
    }
    ```
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

### Read posts of character, culture_movie, human_rights, international, land_environment, photo_audio, political_society and transformed_justice categories.
- URL: `/v1/index_page_categories`
- Method: `GET`
- Response:
  * **Code:** 200 <br />
    **Content:**
    ```
    {
      "records": {
        "character": [{
          // post goes here
        }, {
          // post goes here
        }, ... ], 
        "culture_movie": [{
          // post goes here
        }, {
          // post goes here
        }, ... ],
        "human_rights": [{
          // post goes here
        }, ... ],
        "international": [{
          // post goes here
        }, {
        } ... ],
        "land_environment": [{
          // post goes here
        }, {
          // post goes here
        } ... ],
        "photo_audio": [{
          // post goes here
        }], 
        "political_society": [{
          // post goes here
        } ... ],
        "transformed_justice": [{
          // post goes here
        } ... ],
      },
        "status": "ok"
    }
    ```
  * **Code:** 500 <br />
  **Content:** `{"status": "Internal server error", "error": "${here_goes_error_msg}"}`

## BOOKMARKS
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
        "slug": "about-us-footer",
        "host_name": "www.twreporter.org",
        "is_external": false,
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
   "slug": "about-us-footer",
   "host_name": "www.twreporter.org",
   "is_external": false,
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

## SERVICES
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

## REGISTRATIONS
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

## License
Go-api is [MIT licensed](https://github.com/twreporter/go-api/blob/master/LICENSE)
