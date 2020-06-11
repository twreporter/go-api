## Unreleased

## Released
### 6.0.4 (Current), 2020-06-10

#### Notable Changes

- api/donation
  - Append utm tag to donation link
  - Add `receipt_header` column

- api/mail
  - Update footer of the email template
  - generate client id for tracking
  - update success donation email template

- api/news
  - expand throughput by new connections

- core
  - change mongo query mode

#### Commits
* [[`a25b664fe8`](https://github.com/twreporter/go-api/commit/a25b664fe8)] - Merge pull request #380 from taylrj/update-schema (Tai-Jiun Fang)
* [[`95830c578d`](https://github.com/twreporter/go-api/commit/95830c578d)] - **doc**: update CHANGELOG.md (Taylor Fang)
* [[`73dad89691`](https://github.com/twreporter/go-api/commit/73dad89691)] - **chore**: donations table schema change (Taylor Fang)
* [[`6405beec18`](https://github.com/twreporter/go-api/commit/6405beec18)] - Merge pull request #378 from taylrj/add-receipt-title (Tai-Jiun Fang)
* [[`90541811e7`](https://github.com/twreporter/go-api/commit/90541811e7)] - **doc**: fix json format (Taylor Fang)
* [[`9098549594`](https://github.com/twreporter/go-api/commit/9098549594)] - **doc**: update docs according to review comment (Taylor Fang)
* [[`03da19a4ce`](https://github.com/twreporter/go-api/commit/03da19a4ce)] - **doc**: update docs to add `receipt\_header` field (Taylor Fang)
* [[`815b123a33`](https://github.com/twreporter/go-api/commit/815b123a33)] - Merge pull request #377 from nickhsine/donation-email-temp (nick)
* [[`83b7f799c8`](https://github.com/twreporter/go-api/commit/83b7f799c8)] - **doc**: update CHANGELOG.md (nickhsine)
* [[`1be3f110dd`](https://github.com/twreporter/go-api/commit/1be3f110dd)] - api/mail: update success donation email template (nickhsine)
* [[`88d0999641`](https://github.com/twreporter/go-api/commit/88d0999641)] - Merge pull request #375 from babygoat/bump-6.0.4 (babygoat)
* [[`95c0fba4d2`](https://github.com/twreporter/go-api/commit/95c0fba4d2)] - **doc**: Update Changelog (Ching-Yang, Tseng)
* [[`e08fa70398`](https://github.com/twreporter/go-api/commit/e08fa70398)] - Merge pull request #374 from babygoat/mongo-read-skew (babygoat)
* [[`219646de12`](https://github.com/twreporter/go-api/commit/219646de12)] - api/news: expand throught by new connections (Ching-Yang, Tseng)
* [[`47685da57c`](https://github.com/twreporter/go-api/commit/47685da57c)] - **core**: change mongo query mode (Ching-Yang, Tseng)
* [[`879e07821a`](https://github.com/twreporter/go-api/commit/879e07821a)] - Merge pull request #373 from babygoat/success-donation-email-template-update (babygoat)
* [[`68a2811f95`](https://github.com/twreporter/go-api/commit/68a2811f95)] - api/mail: update footer of the email template (Ching-Yang, Tseng)
* [[`ebd52aedb9`](https://github.com/twreporter/go-api/commit/ebd52aedb9)] - api/mail: generate client id for tracking (Ching-Yang, Tseng)
* [[`845d3d696a`](https://github.com/twreporter/go-api/commit/845d3d696a)] - api/donation: append utm tag to donation link (Ching-Yang, Tseng)

### 6.0.3

#### Notable Changes

- api/user
  - Prevent user from retrieving the bookmarks of others

#### Commits
- [[a01296b](https://github.com/twreporter/go-api/commit/a01296b9c9f433daa5aadbe1a2e70d896ac60a92)] - Prevent a user from retrieving bookmark of others(babygoat)
- [[67fd87f](https://github.com/twreporter/go-api/commit/67fd87f2e2d7971e219e2cf983bda12a2e1c8b0f)] - Refactor tests of bookmark(babygoat)

### 6.0.2, 2020-03-05

#### Notable Changes

- api/donation:
  - Fix incorrect linepay notification format

#### Commits
- [[8c5d196](https://github.com/twreporter/go-api/commit/8c5d196945e75035a12a47f14c828e1932146870)] - Fix incorrect linepay notification format(babygoat)

### 6.0.1, 2020-03-04

#### Notable Changes

- api/donation:
  - Prior to use proxy for tappay request if configured

#### Commits
- [[33573be](https://github.com/twreporter/go-api/commit/33573be1e1bbb9a2fdb7d55eae9693a9e991a91a)] - Dynamically configure donation proxy(babygoat)

### 6.0.0, 2020-02-21

#### Notable Changes

- core:
  - Rewrite error handle with pkg/errors
  - Integrate log formatter on staging/production for stackdriver

#### Commits
- [[493272f](https://github.com/twreporter/go-api/commit/493272f53ffbf44b2058f414f98f3daf994b2e05)] - Update logformatter for the gin format fix(babygoat)
- [[4def9e9](https://github.com/twreporter/go-api/commit/4def9e956929bde9efc6bd2d0ad4f71601b4f453)] - Setup logger(babygoat)
- [[53e87d7](https://github.com/twreporter/go-api/commit/53e87d79b955dd10175011212944fe00516d6232)] - Fix gin 1.4.0 import and logrus module typo(babygoat)
- [[fb2e6de](https://github.com/twreporter/go-api/commit/fb2e6de66116528bd9458b066c182ae4fff99968)] - Add the recovery middleware in production(babygoat)
- [[7764efd](https://github.com/twreporter/go-api/commit/7764efd65672b0f92e8b3e50d098da16d8e6be96)] - Remove vague bookmark update(babygoat)
- [[8292b52](https://github.com/twreporter/go-api/commit/8292b52b52ec58c2a4740929ca80275c46765d78)] - Adjust log severity(babygoat)
- [[812a76d](https://github.com/twreporter/go-api/commit/812a76d60f2eb065742ecaa72870c9903ee886f4)] - Remove unnecessary error log(babygoat)
- [[8472775](https://github.com/twreporter/go-api/commit/8472775a2b04875bd123cdba5cda870cbf6ebab0)] - Remove AppError(babygoat)
- [[b206e53](https://github.com/twreporter/go-api/commit/b206e53d1022e27528bd05e9943c342d864de88e)] - Rewrite errors in utils/service/configs(babygoat)
- [[56a8890](https://github.com/twreporter/go-api/commit/56a88906c88ff143df31d58a29a03e07ccf93351)] - Rewrite controller layer error(babygoat)
- [[dc1d814](https://github.com/twreporter/go-api/commit/dc1d814f259ece66224c1b2ac8d1e76a82b1c8cb)] - Add utility for transfer error to http response(babygoat)
- [[8e15949](https://github.com/twreporter/go-api/commit/8e1594941edb44e73948bec496fe6e16e2d28aa2)] - Add storage errors utilities(babygoat)
- [[2c9aaf3](https://github.com/twreporter/go-api/commit/2c9aaf3d5ccfacf7f9fb9e1db6c91cff3fbc625e)] - Wrap storage error with pkg/errors(Ching-Yang, Tseng)
- [[66b08d0](https://github.com/twreporter/go-api/commit/66b08d077adcda53a84c412389417ae087115314)] - Remove deprecated routes(Ching-Yang, Tseng)

### 5.1.3, 2020-02-06

#### Notable Changes
- api/news:
  - return empty records if there is no query result
  - handle edge case: `?where={categories:{"in": null}}` query string 

#### Commits
- [[f9ae74a](https://github.com/twreporter/go-api/commit/f9ae74ab9960027b0bfadc1ceaca2adebe6a9b0d)] - fix: handle url query parsing failure 
- [[626b694](https://github.com/twreporter/go-api/commit/626b6943cf2bfdac4d564383173190a5d79aa190)] - fix: make (posts|topics) records be empty array rather than null
- [[66aeaee](https://github.com/twreporter/go-api/commit/66aeaeeec1e8a7ea6ba356c2ed7d1d81517be553)] - fix: handle NilObjectId query

### 5.1.2, 2020-02-04
#### Notable Changes
- api/donation:
  - Config frontend host of linepay in runtime

#### Commits
- [[23df36d](https://github.com/twreporter/go-api/commit/23df36df3d85a4de180de14ce1817928d574b4d0)] - Config frontend host of linepay in runtime(babygoat)
- [[3a7948e](https://github.com/twreporter/go-api/commit/3a7948eade39583544189a59ddb799e62b7acac0)] - bug: show latest review and photo articles

### 5.1.1, 2019-11-26
#### Notable Changes
- api/donation:
  - add `line_pay_product_image_url` linepay icon
  - increase size of `bank_result_msg` column
- chore:
  - Include the kubernetes config during deployment
  - update circleci config for new cluster
- api/auth:
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
- api/donation:
  - add `/v1/tappay_query` endpoint for querying TapPay Record API 
  - add `refunded` value into `status` field of donation related tables
- api/news:
  - fix `/v1/index_page` endpoint returning old photography posts
  
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
