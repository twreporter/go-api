## Unreleased

## Released
### 5.1.3(Current), 2020-02-05
#### Commits
- [[f9ae74a](https://github.com/twreporter/go-api/commit/f9ae74ab9960027b0bfadc1ceaca2adebe6a9b0d)] - fix: handle url query parsing failure 

### 5.1.2, 2020-02-04
#### Notable Changes
- donation:
  - Config frontend host of linepay in runtime

#### Commits
- [[23df36d](https://github.com/twreporter/go-api/commit/23df36df3d85a4de180de14ce1817928d574b4d0)] - Config frontend host of linepay in runtime(babygoat)
- [[3a7948e](https://github.com/twreporter/go-api/commit/3a7948eade39583544189a59ddb799e62b7acac0)] - bug: show latest review and photo articles

### 5.1.1, 2019-11-26
#### Notable Changes
- donation:
  - add `line_pay_product_image_url` linepay icon
- db_schema:
  - increase size of `bank_result_msg` column
- build:
  - Include the kubernetes config during deployment
  - update circleci config for new cluster
- oauth:
  - remove v1 oauth endpoints
  - upgrade Facebook Graph API: v2.8 -> v3.2

#### Commits
- [[b9796f8](https://github.com/twreporter/go-api/commit/b9796f8d13987e9ba22b7f6252e03e31b495925f)] - Update config for release environment(babygoat)
- [[986c9f8](https://github.com/twreporter/go-api/commit/986c9f8d96cee3f0e21a94c755569b85cc55a51a)] - Fix incorrect environment setup(babygoat)
- [[e13e89a](https://github.com/twreporter/go-api/commit/e131e89aebb8c187bd68cc17417bd67bc5d14648)] - Do not overwrite the default image name(babygoat)
- [[4087761](https://github.com/twreporter/go-api/commit/40877613d48baa6cdec2c63f8ce02042186a81c9)] - Increase bank_result_msg column(Ching-Yang, Tseng)
- [[c17bc98](https://github.com/twreporter/go-api/commit/c17bc986967f70900e219ba248c1c2c8c4e56d56)] - Fix incorrect kustomize PATH(babygoat)
- [[483ece5](https://github.com/twreporter/go-api/commit/483ece556f8ad77a2e858e864832368159e5a89a)] - Fix incorrect context injection(babygoat)
- [[cfcbd5e](https://github.com/twreporter/go-api/commit/cfcbd5ed5d71a87d55fb518325e320ccca12f95b)] - Fix missing package version file(babygoat)
- [[105e34c](https://github.com/twreporter/go-api/commit/105e34c9fc5d0cf71195b5a1f0450a6b82ac3252)] - Include kubernetes config during deployment(babygoat)
- [[d873dd3](https://github.com/twreporter/go-api/commit/d873dd34ec4b84f42264016c8c509c0195613a7c)] - Only send linepay logo url during linepay trx(babygoat)
- [[fc3cbbc](https://github.com/twreporter/go-api/commit/fc3cbbc36fbe36a249f714bea9db073357585e03)] - Update linepay merchant logo(babygoat)
- [[a74f8ef](https://github.com/twreporter/go-api/commit/a74f8efc218dd7bc79cf400061d52021abb3b9ec)] - Provide valid icon image link(babygoat)
g for next branch(babygoat)
- [[35000c3](https://github.com/twreporter/go-api/commit/35000c3847846ba6440ca63dd1c635f973d220b7)] - Add linepay icon during transaction(babygoat)
- [[cddb19d](https://github.com/twreporter/go-api/commit/cddb19de136f61960da09857ca461eabfb13a4ad)] - remove /v1/auth/faceboook and /v1/auth/google oauth endpoints(nickhsine)
- [[06e4cbf](https://github.com/twreporter/go-api/commit/06e4cbf6da206c4180aded3cf084129621e7f94a)] - update controllers/oauth.go: upgrade facebook graph API from v2.8 to v3.2(nickhsine)
- [[5a88338](https://github.com/twreporter/go-api/commit/5a88338fdd6570b03b65e8c5d38d61d24d48ef6a)] - update circleci config due to k8s cluster change(nickhsine) 

### 5.1.0, 2019-10-15
#### Notable Changes
- feature:
  - add `/v1/tappay_query` endpoint for querying TapPay Record API 
- bug:
  - fix `/v1/index_page` endpoint returning old photography posts
- db-schema:
  - add `refunded` value into `status` field of donation related tables
  
#### Commits
- [[d0712d4](https://github.com/twreporter/go-api/commit/d0712d4fde8e3a4b1c55012b6be875421cb4cd4d)] - bug fix: /v1/index_page endpoint returns old photography post(nickhsine)
- [[fd26655](https://github.com/twreporter/go-api/commit/fd2665539a8ada5923c2e52ff4b78f23103fa4db)] - Add refunded payment status(babygoat)
- [[09dbbe3](https://github.com/twreporter/go-api/commit/09dbbe3559f65cf29cf6abd36effd2056b313021)] - Clean up test users after each test(babygoat)
- [[9c3934e](https://github.com/twreporter/go-api/commit/9c3934ea83f188401bca1f47e35214ce61eeec70)] - Implement transaction query on tappay server (babygoat)
- [[82d00c4](https://github.com/twreporter/go-api/commit/82d00c4fc1f2eae5af4a0634153a37d220288e04)] - Add tests for query transaction record endpoint(babygoat)
- [[cf06678](https://github.com/twreporter/go-api/commit/cf066785262f0716266736d21325a8557639020b)] - Filter out the secret transaction info(babygoat)
- [[936bbf8](https://github.com/twreporter/go-api/commit/936bbf8e7b7bbf17fff77b767868b5944940d0c9)] - Add the endpoint document for tappay query(babygoat)

### 5.0.4

#### Improvement
- Fix missing transaction time error when the transaction fails

### 5.0.3

#### Models
- models/post.go: add `is_external` field in `Post` struct

### 5.0.2

#### Login
- Append login_time query param on login redirect destionation url

### 5.0.1

#### Donations
- Change the default value of `send_receipt` column of `periodic_donations` and `pay_by_prime_donations` table to `no`.

#### Miscellaneous
- Revise login email template

### 5.0.0

#### Breaking Change
- Dependency management migration: `glide` -> `go module`
- Go version upgrade: 1.10 -> 1.12.6

#### New Features
- Line pay support

#### CircleCI refactoring
- Update dockerfile
- Update circleci script
- Add mysql health check in circleci script

#### Code refactoring
- Refactor the auth setup: this patch extracts the auth tokens creation( authorization header,
cookie) into helper function.

- Refactor donation patch error: this patch refactor the donation patch errors into table-driven tests.

- Refactor donation get error: this patch refactors the donation get error into table-driven tests.

- Improve error records: this patch improves the error record fields by recordind extra
`rec_trade_id`, `bank_result_code` and `bank_result_msg`.

#### Bug fix
- Fix wrong address in the success email

#### Miscellaneous
- Update globals/constants.go: change field names of categories

### 4.0.1
#### Bug Fix
- Fix wrong address in the success email

### 4.0.0
#### New Feature  
  * Enforce the donation through forward-proxy

#### Breaking Change
  * Deprecate the following donation endpoints
     - /v1/periodic-donations/:id GET method
     - /v1/periodic-donations/:id PATCH method
     - /v1/donations/prime/:id GET method
     - /v1/donations/prime/:id PATCH method
     - /v1/donations/others/:id GET method
  * Add the follwoing donation endpoints in replace of the above deprecated ones
     - /v1/periodic-donations/orders/:order GET method
     - /v1/periodic-donations/orders/:order PATCH method
     - /v1/donations/prime/orders/:order GET method
     - /v1/donations/prime/orders/:order PATCH method
     - /v1/donations/others/orders/:order GET method
  * Change the donation information update link from /contribute/{frequency}/:id -> /contribute/{frequency}/:order

#### Miscellaneous
  * Improve the CI build flow
  * Add `is_anonymous` field to prime/periodic donation.
  * Fine tune donation success email context.

### 3.0.3
#### Thank you mail refinement
- Format template/success-donation.tmpl
- Update template/success-donation.tmpl. Add do-not-reply message
- Email sender name changed: `no-reply@twreporter.org` -> `報導者 The Reporter <no-reply@twreporter.org>`
- Thank-you mail wording revised: 捐款 -> 贊助

### 3.0.2
#### Schema Change 
- Correct `send_receipt` enumeration value of `periodic_donations`, `pay_by_prime_donations` and `pay_by_other_method_donations` tables

### 3.0.1
#### Bug fixed
- `leading_image_portrait` is missing in full post object.

### 3.0.0
#### Improve authentication and authorization protocol 
  1. A user signs in through the login form or social account
  2. After authentication,  /v2/auth/activate or /v2/auth/{google,
  facebook}/callback will set `id_token` cookie in jwt format.
  3. Frontend server will then launch another request to
  /v2/auth/token along with the bear token in Authorization header
  from `id_token` to get the `access_token`.
  4. After validating the `id_token`, go-api returns the `access_token` in
     response payload.
  5. When users want to sign out, the frontend server should redirect users to
  /v2/auth/logout endpoint, which will unset `id_token` cookie.
  
  Besides protocol improvement, there are some refactors as well,
  - refactor the token generation utilties for backward compatibility.
  - enable sessions while doing google|facebook oAuth.

#### New Feature:
##### Donation endpoints
  - /v1/periodic_donations endpoint with POST method
  - /v1/donations/prime endpoint with POST method
  The above endpoints allow users to contribute monthly(the upper one) or one-time(the lower one).

  - /v1/periodic_donations/:id (PATCH method)
  - /v1/donations/prime/:id (PATCH method)
  The above endpoints allow users to patch detailed information to the certain donation record

  - /v1/periodic-donations/:id?user_id=:userID (GET method)
  - /v1/donations/prime/:id?user_id=:userID (GET method)
  The above endpoints allow users to get the certain donation record
 
##### Mail endpoints
  - /v1/mail/send_activation (POST method)
  - /v1/mail/send_success_donation (POST method)

#### Configuration refactoring
  - use `viper` to load the config
  - change config file format from json to yaml
  - add controllers/mail.go to handle HTTP request/response
  - replace utils/mail.go by services/mail.go
  - use template/signin.tmpl to generate activation mail HTML
  - use template/success-donation.tmpl to generate success donation mail HTML
  - send HTTP POST request to mail endpoints after signin and donation success

#### Miscellaneous
- Use Allow-Origins to constrain access from different sites with respect to the environments
- Send thank-you donation email after success donation
- Api documents for donation, mail, version 2 auth/oauth endpoints
- `form` request body is not supported on donation and mail controllers anymore
- Only development environment can return draft posts or topics

### 2.1.4
- Update /v1/search/posts and /v1/search/authors to use new algolia indices

### 2.1.3
- [feature] /v1/authors{?limit,offset,sort}[GET] endpoint for fetching authors

### 2.1.2
- [Performance] Replace md5 hash function by crc32 on subscriptions table
  - rename web_push_subscriptions to web_push_subs
  - add UNIQUE KEY on `endpoint` field, and set `endpoint` to varchar(500)
  - rename `hash_endpoint` to `crc32_endpoint`, and remove UNIQUE KEY from `crc32_endpoint`

### 2.1.1
- Update membership_user.sql. Remove soft delete on web_push_subscriptions table
 
### 2.1.0
- New endpoint for subscribing webpush notification

### 2.0.3
- Add a new resized target option: w400

### 2.0.2
- Set Cache-Control: no-store in the response header for oauth endpoints 
- Sort EditorPicksSection by updated_at field in controllers/index_page.go

### 2.0.1
- [PR#102](https://github.com/twreporter/go-api/pull/102)
- Use userID, email and standard claims to generate JWT.
- Code refactors. Fix typo, add error check.
- Still redirect to destination even if oauth fails.
- New endpoint "/v1/token/:userID", which is used to renew JWT for clients.
- Set Cache-Control: no-store for those endpoints related to users

### 2.0.0
**Major Change**
- Drop password and signup process, only send activation email every time user want to sign in.
- Dedup the clients accounts. Connect the client who signs in by oauth or by email to the existed record.
- Move oauth controllers from subfolders to root controllers folder.
  - controllers/oauth/goolge/google.go -> controllers/google-oauth.go
  - controllers/oauth/facebook/facebook.go -> controllers/goog-oauth.go
- Update controller/google-oauth.go. Set jwt in the cookies
- Update middlewares/jwt.go. Add SetEmailClaim function.
- Update membership_user.sql
- Change email content wording and styling
- Add GinResponseWrapper function, which deliver the response to the client
- Update controllers/account.go
  - code refactor since the return value of each function is wrapped by
  GinReponseWrapper
- Function test refactor

### 1.1.8
- Bug fix. Output `html` field in ContentBody.

### 1.1.7
- Embed Theme field in post model 
- Make activation email more stylish

### 1.1.6
- Fetch sections asynchronously for index_page controller

### 1.1.5
- Sort returned bookmarks of a user by `updated_at` in `users_bookmarks` table

### 1.1.4
- Sort relateds according to the order set by editors

### 1.1.3 
- Add job title of authors
- Check JWT expiration time
- Allow DELETE method and Authorization Header
- Add endpoint /users/:userID/bookmarks/:bookmarkSlug to get a bookmark of a user
- Refine create/delete/get bookmark

### 1.1.2
- Update models/post.go. Add LeadingImageDescription field

### 1.1.1
- Fetch posts without is_feature: true in photos_section

### 1.1.0
- Update Bookmark model. Replace Href by Slug, Style and External.

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
