### 1.0.11
- Bug fix. Order authors

### 1.0.10
- Bug Fix. Only add Cache-Control in the Response Header when Request Header contains Origin directive'

### 1.0.9 
- Hot Fix. Add hard coded  Access-Control-Allow-Origin in response header

### 1.0.8
- Upgrade github.com/gin-contrib/cors to the latest commit to fix the cors problem

### 1.0.7
- Set Access-Control-Allow-Origin: https://www.twreporter.org 

### 1.0.6
- Add Access-Control-Allow-Origin in response header

### 1.0.5
- Fix Typo. 

### 1.0.4
- Bug fix. Avoid fetching all the records if ids is an empty array instead of nil value

### 1.0.3
- Allow any request in development environment
- Update categories 

### 1.0.2
- Return leading_image_portrait field in post model.
- Add AmazonSES service to send mails.
- Integrate with circleci.
- Add Cache-Control in response header and simply cors setting
- Add CORS in response header for production environment
- Only allow to query `published` posts and topics in production environment

### 1.0.1
- Fix typo. Agolia to Algolia.

### 1.0.0
- initialization
