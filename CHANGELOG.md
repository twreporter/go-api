## Unreleased

### 7.2.4-rc.13

### Notable Changes

- fix
  - Replacement keyword issue

### Commits

- [[`7ddc215`](https://github.com/twreporter/go-api/commit/7ddc215)] - fix: Replacement keyword issue

### 7.2.4-rc.12

### Notable Changes

- fix
  - Typo in code comments
  - Removal of roleCheck in setUser function
  - Sending Explorer role mail only if the user is an Explorer

### Commits

- [[`f4e3f74`](https://github.com/twreporter/go-api/commit/f4e3f74)] - fix: typo
- [[`61a527c`](https://github.com/twreporter/go-api/commit/61a527c)] - fix: remove roleCheck in setUser
- [[`4293783`](https://github.com/twreporter/go-api/commit/4293783)] - fix: only send Explorer role mail if is Explorer

### 7.2.4-rc.11

### Notable Changes

- fix
  - Role activation and checks
  - Role-related email handling
  - Unused string removal from templates

### Commits

- [[`0ca38e6`](https://github.com/twreporter/go-api/commit/0ca38e6)] - fix: only role check if not activated
- [[`fbffc83`](https://github.com/twreporter/go-api/commit/fbffc83)] - fix: only role check if not activated
- [[`fd2f6d8`](https://github.com/twreporter/go-api/commit/fd2f6d8)] - fix: move role mail from signin to setUser
- [[`2c2444e`](https://github.com/twreporter/go-api/commit/2c2444e)] - fix: remove unused string from template

### 7.2.4-rc.10

### Notable Changes

- fix
  - Change IsTrailblazer SQL condition

### Commits

- [[`22ed4f5`](https://github.com/twreporter/go-api/commit/22ed4f5)] - fix: change IsTrailblazer SQL condition

### 7.2.4-rc.9

### Notable Changes

- fix
  - Role assign mail templates loader
  - Add detail logs

### Commits

- [[`9e94c34`](https://github.com/twreporter/go-api/commit/9e94c34)] - fix: add detail logs
- [[`c62abc8`](https://github.com/twreporter/go-api/commit/c62abc8)] - fix: role assign mail templates loader

### 7.2.4-rc.8

### Notable Changes

- fix
  - Typo in code
  - Trailblazer donate downgrade bug
  - Fine-tuned roleCheck logic
  - Added role check before sending role change email
  - Sending explorer email only when user is activated
  - Refactored sendRoleMail and updated trailblazer email handling
  - Updated variables for better code clarity
  - Changed role condition regarding periodic donations in the trailblazerCheck feature
  - Added a 2 months condition to the trailblazer check
- feat
  - Added role assign related emails
  - Updated trailblazer check to accommodate multiple periodic donations

### Commits

- [[`ece2bc9`](https://github.com/twreporter/go-api/commit/ece2bc9)] - fix: typo
- [[`ddadb1e`](https://github.com/twreporter/go-api/commit/ddadb1e)] - fix: trailblazer donate downgrade bug
- [[`79f93b4`](https://github.com/twreporter/go-api/commit/79f93b4)] - fix: fine-tune roleCheck logic
- [[`62a07b0`](https://github.com/twreporter/go-api/commit/62a07b0)] - fix: Add role check before sending role change mail
- [[`9afb442`](https://github.com/twreporter/go-api/commit/9afb442)] - fix: sending explorer email only when user is activated
- [[`7c351a4`](https://github.com/twreporter/go-api/commit/7c351a4)] - fix: refactor sendRoleMail and update trailblazer mail
- [[`12e3837`](https://github.com/twreporter/go-api/commit/12e3837)] - fix: variables
- [[`2875321`](https://github.com/twreporter/go-api/commit/2875321)] - feat: Add role assign related mails
- [[`9098bbd`](https://github.com/twreporter/go-api/commit/9098bbd)] - fix: change role condition about periodic_donations
- [[`0212acc`](https://github.com/twreporter/go-api/commit/0212acc)] - feat: Update trailblazer check for multiple periodic donations
- [[`9a22964`](https://github.com/twreporter/go-api/commit/9a22964)] - fix: Add 2 months condition to trailblazer check

### 7.2.4-rc.7

### Notable Changes

- fix
  - Mysql query

### Commits

- [[`6cac8d8`](https://github.com/twreporter/go-api/commit/6cac8d8)] - Merge pull request #563 from howar31/master 
- [[`2c9f53d`](https://github.com/twreporter/go-api/commit/2c9f53d)] - fix: Fix role replacement query

### 7.2.4-rc.6

### Notable Changes

- role
  - Add donation role flows
  - Add role assign function
  - Change role related operation from ID to Key
- fix
  - Replace Clear and Append with Replace
  - Fix typo

### Commits

- [[`27389bf`](https://github.com/twreporter/go-api/commit/27389bf)] - Merge pull request #556 from howar31/master
- [[`8d417cf`](https://github.com/twreporter/go-api/commit/8d417cf)] - fix: replace Clear and Append with Replace
- [[`991f4bf`](https://github.com/twreporter/go-api/commit/991f4bf)] - fix: typo
- [[`bfd17fd`](https://github.com/twreporter/go-api/commit/bfd17fd)] - fix: Change role related operation from ID to Key
- [[`6b0d7cf`](https://github.com/twreporter/go-api/commit/6b0d7cf)] - feat: Add donation role flows
- [[`6c26368`](https://github.com/twreporter/go-api/commit/6c26368)] - feat: Add role assign function

### 7.2.4-rc.5

### Notable Changes

- api/auth, api/oauth
  - redirect to onboarding url if user deactivated after signin
- doc
  - add `onboarding` parameter for auth & oauth api

### Commits

- [[`a6906021cc`](https://github.com/twreporter/go-api/commit/a6906021cc)] - **fix**: add updated `index.html` (Aylie Chou)
- [[`65a0d7eee9`](https://github.com/twreporter/go-api/commit/65a0d7eee9)] - **fix**: add onboarding url on oauth destination (Aylie Chou)
- [[`28e98c6e97`](https://github.com/twreporter/go-api/commit/28e98c6e97)] - **fix**: redirect to onboarding url if user deactivated (Aylie Chou)

### 7.2.4-rc.4

### Notable Changes

- migration
  - Add migration for role table key column
- api/user
  - Change getUser's response preferences and maillist from string to array
- doc
  - Add role key column

### Commits

- [[`c33ff42`](https://github.com/twreporter/go-api/commit/c33ff42)] - Merge pull request #555 from liruchen32/fix/fix-user-api-response

### 7.2.4-rc.3

### Notable Changes

- migration
  - Add migration for user table activated column
- api/user
  - New function for setting user's preferences and maillist related
  - Add activated field to getUser API
  - Fix first/last name for user
  - Add goroutine for activated update
- doc
  - Add new User section
- tests
  - Fix tests
- typo
  - Fix typo in test
  - Fix typo for setUser route

### Commits

- [[`da6cbfe`](https://github.com/twreporter/go-api/commit/da6cbfe)] - Merge pull request #552 from howar31/feat-getuser
- [[`c18209a`](https://github.com/twreporter/go-api/commit/c18209a)] - fix: reviews
- [[`2b1f789`](https://github.com/twreporter/go-api/commit/2b1f789)] - Merge branch 'master' into feat-getuser
- [[`3285298`](https://github.com/twreporter/go-api/commit/3285298)] - feat: add maillist and read_preference to getUser
- [[`65e9165`](https://github.com/twreporter/go-api/commit/65e9165)] - fix: first/last name for user; add goroutine for activated update
- [[`ca1de82`](https://github.com/twreporter/go-api/commit/ca1de82)] - Merge pull request #551 from howar31/fix-setuserRoute
- [[`283fa53`](https://github.com/twreporter/go-api/commit/283fa53)] - fix: tests
- [[`627656a`](https://github.com/twreporter/go-api/commit/627656a)] - feat: add activated to related flow
- [[`8216cd1`](https://github.com/twreporter/go-api/commit/8216cd1)] - feat: add activated to getUser
- [[`4c51786`](https://github.com/twreporter/go-api/commit/4c51786)] - feat: Add migration for user table activated column
- [[`4873ba3`](https://github.com/twreporter/go-api/commit/4873ba3)] - feat: implement GetUser API
- [[`b4792b1`](https://github.com/twreporter/go-api/commit/b4792b1)] - fix: typo in test
- [[`8ac57a6`](https://github.com/twreporter/go-api/commit/8ac57a6)] - fix: typo for setUser route

### 7.2.4-rc.2

### Notable Changes

- migration
  - Add migration for job queue
- api/user
  - New function for setting user's preferences and maillist related
- doc
  - Add new User section

### Commits

- [[`60ae353`](https://github.com/twreporter/go-api/commit/60ae353)] - feat: Add Mailgroup, SetUser related functions
- [[`77d708b`](https://github.com/twreporter/go-api/commit/77d708b)] - Merge branch 'twreporter:master' into master
- [[`26f0957`](https://github.com/twreporter/go-api/commit/26f0957)] - fix: refactor naming for read preference
- [[`abf8304`](https://github.com/twreporter/go-api/commit/abf8304)] - feat: add SetUser test case
- [[`82cfdcf`](https://github.com/twreporter/go-api/commit/82cfdcf)] - feat: add mailchimp sync related functions
- [[`3caa514`](https://github.com/twreporter/go-api/commit/3caa514)] - fix: reviews
- [[`e5538bf`](https://github.com/twreporter/go-api/commit/e5538bf)] - feat: Complete SetUser by adding ReadPreference implementation
- [[`ccfd7dd`](https://github.com/twreporter/go-api/commit/ccfd7dd)] - fix: reviews, add readable name for interestIDs, transaction for setUser
- [[`6767394`](https://github.com/twreporter/go-api/commit/6767394)] - fix: test case
- [[`74fbf33`](https://github.com/twreporter/go-api/commit/74fbf33)] - fix: move interest ids to config
- [[`4b8ce95`](https://github.com/twreporter/go-api/commit/4b8ce95)] - fix: improve read_preferences transaction
- [[`1698c91`](https://github.com/twreporter/go-api/commit/1698c91)] - fix: update configs
- [[`d2ce95f`](https://github.com/twreporter/go-api/commit/d2ce95f)] - feat: Add User API related docs
- [[`9fa30e9`](https://github.com/twreporter/go-api/commit/9fa30e9)] - fix: reviews and html
- [[`7bdb079`](https://github.com/twreporter/go-api/commit/7bdb079)] - fix: reviews
- [[`13b5c70`](https://github.com/twreporter/go-api/commit/13b5c70)] - fix: reviews, seperate to user.apib
- [[`61f668a`](https://github.com/twreporter/go-api/commit/61f668a)] - fix: reveal read_preferences
- [[`c8c52d1`](https://github.com/twreporter/go-api/commit/c8c52d1)] - fix: resolve migration sequence

### 7.2.4-rc.1

### Notable Changes

- migrations
  - create roles and users_roles table [#540](https://github.com/twreporter/go-api/pull/540)

### Commits

- [[`2198d0a`](https://github.com/twreporter/go-api/commit/2198d0a)] - migrations: migrations for roles table (Lucien Lu)
- [[`b6c58d9`](https://github.com/twreporter/go-api/commit/b6c58d9)] - migrations: use current timestamp for updated_at and created_at (Lucien Lu)
- [[`40ca155`](https://github.com/twreporter/go-api/commit/40ca155)] - migrations: migrations for users_roles table (Lucien Lu)

### 7.2.4-rc.0

### Notable Changes

- bump version to trigger CI/CD for applying new k8s config

### Commits

## Released

### 7.2.3 (Current)

### Notable Changes

- api/posts
  - fix post without category_set should not return 404

### Commits

- [[`ebc9c88e2a`](https://github.com/twreporter/go-api/commit/ebc9c88e2a)] - **fix**: /posts endpoint return 404 if post.caetgory\_set is empty (Aylie Chou)

### 7.2.2

### Notable Changes

- api/donation
  - publish data to neticrm on endpoint /line-notify
- api/posts
  - fix category_set related sorting issue
  - add category_set related functions
- chore
  - fix changelog typo
  - revert #510

### Commits

- [[`873056179f`](https://github.com/twreporter/go-api/commit/873056179f)] - **fix**: CHANGELOG repo name typo (Aylie Chou)
* [[`237abe1ae3`](https://github.com/twreporter/go-api/commit/237abe1ae3)] - **fix**: typo (Aylie Chou)
- [[`d5e6a8499c`](https://github.com/twreporter/go-api/commit/d5e6a8499c)] - **fix**: endpoint /line-notify publish to neticrm topic (Aylie Chou)
- [[`633c49a169`](https://github.com/twreporter/go-api/commit/633c49a169)] - **fix**: getMetaOfPosts sorting issue (Howard Ou)
- [[`3853a0e600`](https://github.com/twreporter/go-api/commit/3853a0e600)] - **fix**: typo (Howard Ou)
- [[`68fae18b38`](https://github.com/twreporter/go-api/commit/68fae18b38)] - **fix**: typo (Howard Ou)
- [[`1ecbf5ad5a`](https://github.com/twreporter/go-api/commit/1ecbf5ad5a)] - Revert **fix**: sorting large data (Howard Ou)
- [[`9bd62e938e`](https://github.com/twreporter/go-api/commit/9bd62e938e)] - **fix**: sorting large data (Howard Ou)
- [[`81bed66e21`](https://github.com/twreporter/go-api/commit/81bed66e21)] - **fix**: sorting & paging for posts api (Howard Ou)
- [[`3985e4442b`](https://github.com/twreporter/go-api/commit/3985e4442b)] - **fix**: wording (Howard Ou)
- [[`5b45c30d70`](https://github.com/twreporter/go-api/commit/5b45c30d70)] - **fix**: update logic (Howard Ou)
- [[`9b6f222027`](https://github.com/twreporter/go-api/commit/9b6f222027)] - **fix**: update category_set check (Howard Ou)
- [[`f74f09689a`](https://github.com/twreporter/go-api/commit/f74f09689a)] - **fix**: add cases for category_set check (Howard Ou)
- [[`bd2231a618`](https://github.com/twreporter/go-api/commit/bd2231a618)] - **fix**: fine tune category_set checking function (Howard Ou)
- [[`e6a3749c4b`](https://github.com/twreporter/go-api/commit/e6a3749c4b)] - **fix**: change category_set consistancy test (Howard Ou)
- [[`f860033af4`](https://github.com/twreporter/go-api/commit/f860033af4)] - **fix**: join new and old categories/categoryset search result (Howard Ou)
- [[`d18214a94e`](https://github.com/twreporter/go-api/commit/d18214a94e)] - **fix**: category_set for full post (Howard Ou)
- [[`540616faa9`](https://github.com/twreporter/go-api/commit/540616faa9)] - **fix**: test case for category_set (Howard Ou)

### 7.2.1

### Notable Changes

- circleci
  - use gke-gcloud-auth-plugin
  - use twreporter-gcp context for gcp related credential
- api/doc
  - update api docs for new info structure

### Commits

- [[`0105fa91d9`](https://github.com/twreporter/go-api/commit/0105fa91d9)] - **fix**: use twreporter-gcp context for gcp related credential (Aylie Chou)
- [[`0d6dba4f20`](https://github.com/twreporter/go-api/commit/0d6dba4f20)] - **fix**: use gke-gcloud-auth-plugin (Aylie Chou)
- [[`acac5b7678`](https://github.com/twreporter/go-api/commit/acac5b7678)] - **docs**: add new category info to /v2/index-page (Aylie Chou)
- [[`53dba2da89`](https://github.com/twreporter/go-api/commit/53dba2da89)] - **docs**: update api doc (Aylie Chou)
- [[`3379bb8c7e`](https://github.com/twreporter/go-api/commit/3379bb8c7e)] - **docs**: update api doc for parameter type change (Aylie Chou)
- [[`2a848b6195`](https://github.com/twreporter/go-api/commit/2a848b6195)] - **docs**: update doc/index.html for github page (Aylie Chou)
* [[`42df484e9f`](https://github.com/twreporter/go-api/commit/42df484e9f)] - **docs**: add docs for category sets (Aylie Chou)

### 7.2.0

### Notable Changes

- api/donation
  - add slack notify when publish to cloud pub/sub fail
  - fix periodic donation do not store pay\_method field
- api/doc
  - fix typo for `auto\_tax\_deduction`
- api/bookmark
  - add `added\_at` in GetBookmarksOfAUser response

### Commits

- [[`63389f7753`](https://github.com/twreporter/go-api/commit/63389f7753)] - **fix**: add `added\_at` in GetBookmarksOfAUser response (Aylie Chou)
- [[`6a8c3d6311`](https://github.com/twreporter/go-api/commit/6a8c3d6311)] - **fix**: fix test case & donation model declaration (Aylie Chou)
- [[`9f757de4ba`](https://github.com/twreporter/go-api/commit/9f757de4ba)] - **fix**: fix periodic donation do not store pay\_method field (Aylie Chou)
- [[`6ba0b9d25d`](https://github.com/twreporter/go-api/commit/6ba0b9d25d)] - **feat**: add slack notify when publish to cloud pub/sub fail (Aylie Chou)
- [[`6d95cf8513`](https://github.com/twreporter/go-api/commit/6d95cf8513)] - api/doc: fix typo for `auto\_tax\_deduction` (#480) (Tai-Jiun Fang)

### 7.1.3

### Notable Changes
- .circleci
  - add pubsub credencials

### Commits

* [[`b201454e52`](https://github.com/twreporter/go-api/commit/b201454e52)] - **fix**: set env GOOGLE\_APPLICATION\_CREDENTIALS to pubsub\_credentials path (Aylie Chou)
* [[`a96f92e69b`](https://github.com/twreporter/go-api/commit/a96f92e69b)] - **fix**: add pubsub credentials (Aylie Chou)

### 7.1.2

### Notable Changes
- api/donation
  - publish donation/user update data to cloud pub/sub

### Commits

* [[`3b83dbf8bd`](https://github.com/twreporter/go-api/commit/3b83dbf8bd)] - **fix**: update go-mod-lib vresion & publish to pub.sub topic when user update (Aylie Chou)
* [[`778e281c69`](https://github.com/twreporter/go-api/commit/778e281c69)] - **fix**: publish to cloud pub/sub topic when user update (Aylie Chou)
* [[`a92b576558`](https://github.com/twreporter/go-api/commit/a92b576558)] - **fix**: update go-mod-lib version to fix test error (Aylie Chou)
* [[`8421ff7a49`](https://github.com/twreporter/go-api/commit/8421ff7a49)] - **feat**: publish upsert donation info to cloud pub/sub (Aylie Chou)

### 7.1.1

### Notable Changes
- model/donation
  - upsert fields for donation redesign [#463](https://github.com/twreporter/go-api/pull/463)
- api/donation
  - address updated fields data flow [#463](https://github.com/twreporter/go-api/pull/463)
  - add UpdateUserDataByCardholder func [#463](https://github.com/twreporter/go-api/pull/463)
- doc/prime-donation, doc/periodic-donation
  - update donation docs due to spec change [#462](https://github.com/twreporter/go-api/pull/462), [#464](https://github.com/twreporter/go-api/pull/464), [#465](https://github.com/twreporter/go-api/pull/465)
- api/mail
  - update success donation email content [#470](https://github.com/twreporter/go-api/pull/470)

### Commits
* [[`b1cd841d9c`](https://github.com/twreporter/go-api/commit/b1cd841d9c)] - api/mail: update success donation email content (#470) (Tai-Jiun Fang)
* [[`43d852583b`](https://github.com/twreporter/go-api/commit/43d852583b)] - **fix**: fix sprintf type (Aylie Chou)
* [[`4d9571758d`](https://github.com/twreporter/go-api/commit/4d9571758d)] - **fix**: print out err in UpdateUserDataByCarholder func (Aylie Chou)
* [[`7680c04eba`](https://github.com/twreporter/go-api/commit/7680c04eba)] - **doc**: add development guide about OAuth setting to README.md (#466) (Tai-Jiun Fang)
* [[`1831ab39a8`](https://github.com/twreporter/go-api/commit/1831ab39a8)] - **doc**: add migration usage with Makefile in README.md (Aylie Chou)
* [[`419571eca5`](https://github.com/twreporter/go-api/commit/419571eca5)] - **fix**: fix typo & add migrate cmd in Makefile (Aylie Chou)
* [[`7bd3b89181`](https://github.com/twreporter/go-api/commit/7bd3b89181)] - **fix**: fix typo & sync cardholder to user table asynchronously (Aylie Chou)
* [[`73ecddf94e`](https://github.com/twreporter/go-api/commit/73ecddf94e)] - api/donation: add updateUserDataFromCardholder function (Aylie Chou)
* [[`ef6b28361d`](https://github.com/twreporter/go-api/commit/ef6b28361d)] - model/donation: add address related fields to cardholder schema (Aylie Chou)
* [[`a09c6cb32a`](https://github.com/twreporter/go-api/commit/a09c6cb32a)] - **chore**: update aglio doc index.html (Aylie Chou)
* [[`385c2b097c`](https://github.com/twreporter/go-api/commit/385c2b097c)] - **doc**: update prime-donation & periodic-donation doc (Aylie Chou)
* [[`2f9c449b05`](https://github.com/twreporter/go-api/commit/2f9c449b05)] - **doc**: update type of `read\_preference` (#464) (Tai-Jiun Fang)
* [[`cdb79a21ec`](https://github.com/twreporter/go-api/commit/cdb79a21ec)] - **feat**: add data flow for updated db schema (Aylie Chou)
* [[`c60442832d`](https://github.com/twreporter/go-api/commit/c60442832d)] - **fix**: add field update & preserve to-be-drop field fornew-old compatibility (Aylie Chou)
* [[`46197c5195`](https://github.com/twreporter/go-api/commit/46197c5195)] - **doc**: update donation docs due to spec change (#462) (Tai-Jiun Fang)
* [[`14aac66ef7`](https://github.com/twreporter/go-api/commit/14aac66ef7)] - **fix**: fix syntax error (Aylie Chou)
* [[`3a924a4d02`](https://github.com/twreporter/go-api/commit/3a924a4d02)] - **fix**: fix table name (user -\> users) (Aylie Chou)
* [[`39e8035321`](https://github.com/twreporter/go-api/commit/39e8035321)] - **chore**: update migrate file (000007/000008) (Aylie Chou)
* [[`020804f829`](https://github.com/twreporter/go-api/commit/020804f829)] - **chore**: add migrate up/down file (Aylie Chou)

### 7.1.0

#### Notable Changes

- api/account:
  - add authenticate endpoint ([#456](https://github.com/twreporter/go-api/pull/456))
  - redirect to the error redirection URL when error occurs ([#459](https://github.com/twreporter/go-api/pull/459))
- api/mail:
  - update authenticate and signin email ([#456](https://github.com/twreporter/go-api/pull/456)) 

#### Commits
* [[`9d29f0a162`](https://github.com/twreporter/go-api/commit/9d29f0a162)] - **doc**: update signin body (Tai-Jiun Fang)
* [[`2a973c116e`](https://github.com/twreporter/go-api/commit/2a973c116e)] - api/account: redirect to the error redirection URL when error occurs (Tai-Jiun Fang)
* [[`66691f4005`](https://github.com/twreporter/go-api/commit/66691f4005)] - **doc**: update CHANGELOG.md (#457) (Tai-Jiun Fang)
* [[`d680a128c0`](https://github.com/twreporter/go-api/commit/d680a128c0)] - api/mail: update authenticate and signin email (Tai-Jiun Fang)
* [[`ebbbf188e1`](https://github.com/twreporter/go-api/commit/ebbbf188e1)] - api/account: add authenticate endpoint (Tai-Jiun Fang)
* [[`fdef35f7e9`](https://github.com/twreporter/go-api/commit/fdef35f7e9)] - **doc**: add /v2/authenticate (Tai-Jiun Fang)

## Released

### 7.0.2 (Current), 2021-10-13

#### Notable Changes

- api/donation:
  - fix patch of Line POINT only transaction ([#451](https://github.com/twreporter/go-api/pull/451)) 

#### Commits
- [[32866db](https://github.com/twreporter/go-api/commit/32866dbb29d7bcf93153217c49681daa9684ff43)] - api/donation: fix patch of Line POINT only transaction(Ching-Yang, Tseng)
- [[c229bec](https://github.com/twreporter/go-api/commit/c229bec14c403a7946f85aaf1a6d764d2801a4a0)] - api/donation: return error detail(Ching-Yang, Tseng)

### 7.0.1, 2021-09-23

#### Notable Changes

- api/donation:
  - adjust backend notification status ([#447](https://github.com/twreporter/go-api/pull/447))

#### Commits
- [[7e754db](https://github.com/twreporter/go-api/commit/7e754db82083171f8e0b15dade0df6183ef569ee)] - doc: update CHANGELOG(Ching-Yang, Tseng)
- [[d03d6fb](https://github.com/twreporter/go-api/commit/d03d6fb761a374896bdaeaa15c75ca5fadeb401b)] - api/donation: adjust backend notification status(Ching-Yang, Tseng)
- [[b841258](https://github.com/twreporter/go-api/commit/b841258f012bbaa9bf9b2c9c9e72320ea0984e37)] - doc: adjust backend notify success status code(Ching-Yang, Tseng)

### 7.0.0, 2021-07-27

#### Notable Changes

- api/news:
  - deprecate /v1/posts, /v1/topics/, /v1/index_page, /v1/index_page_categories endpoints([#432](https://github.com/twreporter/go-api/pull/432))

- api/donation:
  - record api update ([#433](https://github.com/twreporter/go-api/pull/433), [#440](https://github.com/twreporter/go-api/pull/440))
  - log donor info if the update occurs error ([#436](https://github.com/twreporter/go-api/pull/436))

- api/auth:
  - fix invalid mgo session ([#434](https://github.com/twreporter/go-api/pull/434))
  - log error if the session cannot be saved ([#435](https://github.com/twreporter/go-api/pull/435))

#### Commits
- [[458d643](https://github.com/twreporter/go-api/commit/458d643f4417172d3e8bae0a214a9f584ba82b5d)] - doc: update CHANGELOG(Ching-Yang, Tseng)
- [[bb6017b](https://github.com/twreporter/go-api/commit/bb6017b92c49bcd216b1b877d897910885b6140b)] - api/donation: fix time filter unit(Ching-Yang, Tseng)
- [[4805a5c](https://github.com/twreporter/go-api/commit/4805a5c4c92a9fd959755a0a04bfd42ac179faac)] - doc: update CHANGELOG(Ching-Yang, Tseng)
- [[bf53737](https://github.com/twreporter/go-api/commit/bf5373701373700817ef8210364c0d0cbf81f60d)] - api/donation: log donor info if the update occurs error(Ching-Yang, Tseng)
- [[6584064](https://github.com/twreporter/go-api/commit/658406499bb7f6b7ec803d2f2934db11929fa6b8)] - api/auth: fix invalid mgo session(Ching-Yang, Tseng)
- [[a6275bd](https://github.com/twreporter/go-api/commit/a6275bddb189b91132f1bb645b9e1a98b6778d3b)] - api/donation: update the endpoint according API change(Ching-Yang, Tseng)
- [[f78ebac](https://github.com/twreporter/go-api/commit/f78ebac6635b415039623fcc2775d35ea386b8c4)] - api/auth: log error if the session cannot be saved(Ching-Yang, Tseng)
- [[5935451](https://github.com/twreporter/go-api/commit/593545154f8be646c9d8e407cf15f9fa77b2f850)] - api/donation: update test according to API change(Ching-Yang, Tseng)
- [[37a540e](https://github.com/twreporter/go-api/commit/37a540e8ddffcb79e997d13f2ce2b6b3b4798b33)] - doc: update /v1/tappay_record endpoint(Ching-Yang, Tseng)
- [[29e3316](https://github.com/twreporter/go-api/commit/29e331651ad04cdaf62a67873e2e947df4e4dc13)] - doc: remove deprecated api on readme(Ching-Yang, Tseng)
- [[d814162](https://github.com/twreporter/go-api/commit/d8141621f617bd43bc663f504a75a49112af14eb)] - api/news: remove v1 news model files(Ching-Yang, Tseng)
- [[cb87625](https://github.com/twreporter/go-api/commit/cb876252ab595d8e4786a14683b8d49fcee885b5)] - api/news: remove NewsController(Ching-Yang, Tseng)
- [[d7904b6](https://github.com/twreporter/go-api/commit/d7904b6f53cee1ea4bedaa75d3e374a4b679b138)] - api/news: remove /v1/search/authors, /v1/search/posts(Ching-Yang, Tseng)
- [[7ef2b47](https://github.com/twreporter/go-api/commit/7ef2b47cbce4a4107238bf5d0c6ad63687614c03)] - api/news: remove /v1/authors endpoint(Ching-Yang, Tseng)
- [[5875add](https://github.com/twreporter/go-api/commit/5875addeca498ffb81e2518cecbbc68072f6fdc2)] - api/news: remove /v1/topics, /v1/topics/:slug(Ching-Yang, Tseng)
- [[1cefd32](https://github.com/twreporter/go-api/commit/1cefd32542ada626df534622b087a2eb41cc9afc)] - api/news: remove /v1/posts, /v1/posts/:slug(Ching-Yang, Tseng)
- [[f2f32fd](https://github.com/twreporter/go-api/commit/f2f32fdd2c379d7b9f4d8cd5701a4be02bde4310)] - api/news: deprecate /v1/index_page(Ching-Yang, Tseng)

### 6.2.1, 2021-04-19

#### Notable Changes

- api/news:
  - api/news: add state of the topic field in a post([#428](https://github.com/twreporter/go-api/pull/428))
  
#### Commits
- [[ad4a6a4](https://github.com/twreporter/go-api/commit/ad4a6a4aa63ea920331c31873ed8c88785758dd8)] - api/news: add state of the topic in a post(Ching-Yang, Tseng)

### 6.2.0, 2021-04-16

#### Notable Changes

- core:
  - rename module to github host([#416](https://github.com/twreporter/go-api/pull/416))
  
- chore:
  - fix configGenerator issue during CI([#423](https://github.com/twreporter/go-api/pull/423))
  
- api/donation:
  - return error if patch record failed([#418](https://github.com/twreporter/go-api/pull/418))
  
- api/news:
  - add /v2/authors, /v2/authors/:author_id, /v2/authors/:author_id/posts route([#419](https://github.com/twreporter/go-api/pull/419), [#420](https://github.com/twreporter/go-api/pull/420))

#### Commits
- [[3671230](https://github.com/twreporter/go-api/commit/367123017c72e3847bc83cc746496e01848cfed9)] - chore: workaround configMapGenerator issue(Ching-Yang, Tseng)
- [[d47794f](https://github.com/twreporter/go-api/commit/d47794fffb1de5c2e49688566dae70c17db6adc8)] - chore: bump version for demo the author page(Ching-Yang, Tseng)
- [[0bacc08](https://github.com/twreporter/go-api/commit/0bacc082757847b055453a28967bd05bf8d447a1)] - api/news: assemble the authors data in mongo(Ching-Yang, Tseng)
- [[55fa742](https://github.com/twreporter/go-api/commit/55fa742d31e7a5480d2149135e6d1effa04c679b)] - api/news: add timeout during algolia search(Ching-Yang, Tseng)
- [[32fc662](https://github.com/twreporter/go-api/commit/32fc662ceea2f57a43855ee3cd2c5f0b49387388)] - api/news: add timeout for author page(Ching-Yang, Tseng)
- [[971d33a](https://github.com/twreporter/go-api/commit/971d33a16479beb5174157c0460970212fc6708d)] - api/news: return 204 if no further content available(Ching-Yang, Tseng)
- [[b3faadf](https://github.com/twreporter/go-api/commit/b3faadf2c2583b0ea6d657567670ab85715fe0b5)] - api/news: fix author id unmarshaler(Ching-Yang, Tseng)
- [[6e6c562](https://github.com/twreporter/go-api/commit/6e6c562f048c255fe5b79ab6d457b9687de59b9b)] - api/news: get authors from algolia in primary(Ching-Yang, Tseng)
- [[5aa2c5e](https://github.com/twreporter/go-api/commit/5aa2c5e6f33a4c1e47185b37d5793ef6a1614bcd)] - api/news: fix incorrect thumbnail and update time(Ching-Yang, Tseng)
- [[3e3b610](https://github.com/twreporter/go-api/commit/3e3b610e54617c84270a22f20c98cc547ecfd602)] - api/news: fix GetAuthorByID should return 404 if Invalid ID is provided(Ching-Yang, Tseng)
- [[ddbb553](https://github.com/twreporter/go-api/commit/ddbb5532866d3867a212dba710d32a57804b8b2b)] - api/news: fix TestGetAuthorByID with InvalidID(Ching-Yang, Tseng)
- [[8c47673](https://github.com/twreporter/go-api/commit/8c476733dddd4d5e93d68ed497ca942a7ebc602b)] - api/news: add test for getting post by an author (writer)(Ching-Yang, Tseng)
- [[e1ff366](https://github.com/twreporter/go-api/commit/e1ff3667137553f7213823a85fea2e319c50428f)] - api/news: add test for getting posts by an author (photographer)(Ching-Yang, Tseng)
- [[32bc2b3](https://github.com/twreporter/go-api/commit/32bc2b3a9fde7f4e5f6d900e12b889ad0e9a680f)] - api/news: add test for getting posts by an author (designer)(Ching-Yang, Tseng)
- [[e0c9343](https://github.com/twreporter/go-api/commit/e0c934326fb654d06722d3bd6de92aff38bba3b5)] - api/news: Implement `/v2/authors/:author_id/posts`(Ching-Yang, Tseng)
- [[0771fda](https://github.com/twreporter/go-api/commit/0771fda8308b4cbb1fa045944dd12639b2eeb251)] - api/news: add test for getting posts by an author (engineer)(Ching-Yang, Tseng)
- [[5ed40bc](https://github.com/twreporter/go-api/commit/5ed40bc14c758adbc6831c52bf41b0e80e7fe308)] - api/news: return 404 if author_id not exists(Ching-Yang, Tseng)
- [[187105f](https://github.com/twreporter/go-api/commit/187105f5aa35650c75e08edbb460997e073495af)] - api/news: add test for unable to get an author(Ching-Yang, Tseng)
- [[5d39b1e](https://github.com/twreporter/go-api/commit/5d39b1e60fecd42b0f0fd99fabf4cee85d7ec8cc)] - api/news: implement `/v2/author/:author_id`(Ching-Yang, Tseng)
- [[2b527c6](https://github.com/twreporter/go-api/commit/2b527c6b157e19f90a7a41aa566ab771a6b0c0a9)] - api/news: add test for get an author by ID(Ching-Yang, Tseng)
- [[620fa7a](https://github.com/twreporter/go-api/commit/620fa7a98d0f78101f541655c171e8656053d362)] - api/news: return 204 if no author match(Ching-Yang, Tseng)
- [[c88715a](https://github.com/twreporter/go-api/commit/c88715a25901f05485c67bca7a20c71f0d15f462)] - api/news: add test for search returns empty(Ching-Yang, Tseng)
- [[655e74e](https://github.com/twreporter/go-api/commit/655e74e833afb74cea8b961f13837743e9bcf070)] - api/news: implement /v2/authors endpoint(Ching-Yang, Tseng)
- [[151e5ca](https://github.com/twreporter/go-api/commit/151e5ca9a5a8ae49de65e59b74d6528dc546533a)] - api/news: add test for author keyword search(Ching-Yang, Tseng)
- [[7cf094c](https://github.com/twreporter/go-api/commit/7cf094c8df9e36a146756f966b01ca401e77a957)] - chore: add another mongo testdb(Ching-Yang, Tseng)
- [[02f366a](https://github.com/twreporter/go-api/commit/02f366aa717eee02d2c946951afd2216c14be517)] - doc: add api document of /v2/authors/{author_id}/post(Ching-Yang, Tseng)
- [[8d4257b](https://github.com/twreporter/go-api/commit/8d4257b5e98f8d55f922c6b6c121228e78f99478)] - doc: add api document for /v2/authors/{author_id}(Ching-Yang, Tseng)
- [[8c01fc3](https://github.com/twreporter/go-api/commit/8c01fc3cd01712e8218bd64425ee57c1cbe4a48d)] - doc: add api document to /v2/authors endpoint(Ching-Yang, Tseng)
- [[f0cf6a1](https://github.com/twreporter/go-api/commit/f0cf6a13f0298743f2dd15d76a0d604a27117ad0)] - api/donations: return error if patch record failed(Ching-Yang, Tseng)
- [[4c816a4](https://github.com/twreporter/go-api/commit/4c816a4f17e25b1fa0c59cf34476ba9861e71101)] - all: rename module to github host(Ching-Yang, Tseng)

### 6.1.3, 2020-10-21

#### Notable Changes

- api/donation:
   - add receipt_header to pay_by_prime_donation, pay_by_other_method_donation, periodic_donation tables ([#406](https://github.com/twreporter/go-api/pull/406))

- api/mail:
  - update personal info reveal in donation success email ([#409](https://github.com/twreporter/go-api/pull/409))

#### Commits
- [[e2a186d](https://github.com/twreporter/go-api/commit/e2a186dd6286c5e04dd43870ac2b23084f513e9e)] - doc: update CHANGELOG(Ching-Yang, Tseng)
- [[640b3eb](https://github.com/twreporter/go-api/commit/640b3ebb290b66c771f2596b4ba35f56d766d28d)] - api/mail: update donation success email(Ching-Yang, Tseng)
- [[9f60841](https://github.com/twreporter/go-api/commit/9f60841732a443bd965c4308e8a69e27c6a871e5)] - fix: update membershio_user.sql to catch up the latest migrations(nickhsine)
- [[61e0359](https://github.com/twreporter/go-api/commit/61e03593389677ed09360764b6d5f4a2495f091b)] - api/donations: validate receipt_header if not zero value(Ching-Yang, Tseng)
- [[1b2fb74](https://github.com/twreporter/go-api/commit/1b2fb7405b93528204ee53182ad39d47f3876591)] - doc: include PR in notable changes section(Ching-Yang, Tseng)
- [[3d3c492](https://github.com/twreporter/go-api/commit/3d3c49279dc380583794e58454ac6da71056d06a)] - api/donations: implement receipt_header deletion(Ching-Yang, Tseng)
- [[03232fc](https://github.com/twreporter/go-api/commit/03232fc757e392df6355f5bc4104bdb2f3ac0d24)] - api/donations: add tests for patch empty fields(Ching-Yang, Tseng)
- [[7a23f93](https://github.com/twreporter/go-api/commit/7a23f93749d416b9dfc0338cfe48c54449e2b833)] - api/donation: refactor patch prime/periodic donations(Ching-Yang, Tseng)
- [[899cff7](https://github.com/twreporter/go-api/commit/899cff70e320d6bb86483d735a0c7f561c0e8404)] - api/donation: support GET `receipt_header` field(Ching-Yang, Tseng)
- [[81c1fa9](https://github.com/twreporter/go-api/commit/81c1fa91dcc75e18516d041908e947b8301c0a9b)] - api/donation: support patch receipt_header field(Ching-Yang, Tseng)
- [[fb204bb](https://github.com/twreporter/go-api/commit/fb204bbcfd8bd1c994313aeb35dfa256b9f1d40b)] - doc: update CHANGELOG(Ching-Yang, Tseng)
- [[e76f5fb](https://github.com/twreporter/go-api/commit/e76f5fb137ee2e58aa81e68a642d90f2fc3bd71f)] - api/donation: add receipt_header field(Ching-Yang, Tseng)

### 6.1.2, 2020-09-18

#### Notable Changes

- api/news:
   - fix out of order related documents
   - fix out of order designers/engineers/photographers/writers

#### Commits
- [[4b17fd5](https://github.com/twreporter/go-api/commit/4b17fd534cef4ed8c4515b32e6abd8619ebfb62e)] - doc: Update CHANGELOG(Ching-Yang, Tseng)
- [[17c59c5](https://github.com/twreporter/go-api/commit/17c59c59af0ffd7c674ded59037c9cc3837c1f76)] - api/news: fix out-of-order authors(Ching-Yang, Tseng)
- [[72d5300](https://github.com/twreporter/go-api/commit/72d5300e82462da2d26d459e5e51ab7462373131)] - doc: Update CHANGELOG(Ching-Yang, Tseng)
- [[4ce93c8](https://github.com/twreporter/go-api/commit/4ce93c82fe5119493eff41fbccaa4bd826082715)] - api/news: fix out of order related document(Ching-Yang, Tseng)

### 6.1.1, 2020-09-03

#### Notable Changes

- api/news:
   - fix `leading_video` field decoder

#### Commits
- [[e837034](https://github.com/twreporter/go-api/commit/e83703416b9af31d40cbca3a8d987b9f6e8f4595)] - doc: update the CHANGELOG(Ching-Yang, Tseng)
- [[8c3a164](https://github.com/twreporter/go-api/commit/8c3a164785e9d920992e6525cd2fd47bf172489c)] - api/news: fix video bson document decoder(Ching-Yang, Tseng)

### 6.1.0 (Current), 2020-08-28

#### Notable Changes

- api/news:
   - add /v2/posts, /v2/posts/SLUG endpoints
   - add /v2/topics, /v2/topics/SLUG endpoints
   - add /v2/index_page endpoint (combine the records of /v1/index_page and /v1/index_page_categories)

#### Commits
- [[7aaee3b](https://github.com/twreporter/go-api/commit/7aaee3bbb97b0d0ded47926ad88a142b5d05e4ec)] - api/news: improve filter performance(Ching-Yang, Tseng)
- [[189d4bb](https://github.com/twreporter/go-api/commit/189d4bbaf6a7c582e6d67eea1c29dc77860db6a6)] - api/news: filter draft related documents(Ching-Yang, Tseng)
- [[fe693fe](https://github.com/twreporter/go-api/commit/fe693fefbd51bad4046f05ed7bd6c192f8ea670f)] - doc: Update CHANGELOG(Ching-Yang, Tseng)
- [[af42cae](https://github.com/twreporter/go-api/commit/af42caeb6227e89c6ff3ccb192a6175c790ad62b)] - api/news: fix index page response format(Ching-Yang, Tseng)
- [[a4a3a1b](https://github.com/twreporter/go-api/commit/a4a3a1bf2a7d81d3e8ecc0d0cc8533485af20eb5)] - api/news: fix post list query parameter parser(Ching-Yang, Tseng)
- [[6d01667](https://github.com/twreporter/go-api/commit/6d016671a54f861780e72a9724cad475c04e1f7d)] - api/news: parameterize the timeout value(Ching-Yang, Tseng)
- [[57df91e](https://github.com/twreporter/go-api/commit/57df91ea2e22e9aaad41a4c9ec56aa88a8c5e3e1)] - api/news: fix post/topic query count command(Ching-Yang, Tseng)
- [[15e78d2](https://github.com/twreporter/go-api/commit/15e78d2a64d7451e666a699c9702a28c56671cf7)] - api/news: do not build sort stage for single query(Ching-Yang, Tseng)
- [[838a908](https://github.com/twreporter/go-api/commit/838a9087a25640c74a3b949ed7d0d690bc3d8d5c)] - api/news: simplify index page jobs pipeline(Ching-Yang, Tseng)
- [[9c32bde](https://github.com/twreporter/go-api/commit/9c32bde06ed08ebae75dae894802920f3e7e35e3)] - api/news: rewrite index jobs preparation(Ching-Yang, Tseng)
- [[cf0e390](https://github.com/twreporter/go-api/commit/cf0e39070918214e22bd4a3ea1914c3500f20f32)] - api/news: append full flag on FullPosts(Ching-Yang, Tseng)
- [[3852f86](https://github.com/twreporter/go-api/commit/3852f86ab32763821fc60cc41ba725b727a513ff)] - api/news: group Option functions(Ching-Yang, Tseng)
- [[5ede79c](https://github.com/twreporter/go-api/commit/5ede79c86c94917e06b734d32837c8fda0c5292b)] - api/news: bail out error and handle cursor error(Ching-Yang, Tseng)
- [[1cea5b4](https://github.com/twreporter/go-api/commit/1cea5b4a047b9b46b050ae9886d7f413396ee5e0)] - api/news: fix test fail from default query changed(Ching-Yang, Tseng)
- [[e3f7467](https://github.com/twreporter/go-api/commit/e3f74673eb18d4d8b21a5e8dc303bc877857ddfb)] - chore: update mongo image for testing to 3.6.18(Ching-Yang, Tseng)
- [[ce374ab](https://github.com/twreporter/go-api/commit/ce374abbd9585db5c70a96363d76326c517329a3)] - api/news: remove draft implementation(Ching-Yang, Tseng)
- [[3cc2212](https://github.com/twreporter/go-api/commit/3cc221298e5bbf71d13d2a7c40107d20f4dd101f)] - api/news: only listing published post(Ching-Yang, Tseng)
- [[df6fd57](https://github.com/twreporter/go-api/commit/df6fd57d9e0b684c7c4c817235b0db05cb73898e)] - api/news: fix index page sorting issue(Ching-Yang, Tseng)
- [[2d31c5a](https://github.com/twreporter/go-api/commit/2d31c5a938aba324ab91ea901a0850fdd04fc1d1)] - api/news: move index page endpoint into controller(Ching-Yang, Tseng)
- [[a6a4434](https://github.com/twreporter/go-api/commit/a6a4434f6be616265bb6ebb92576f384211861b3)] - api/news: move post section and category into internal/news(Ching-Yang, Tseng)
- [[10ae823](https://github.com/twreporter/go-api/commit/10ae823db5d8ebdf99888ef8284d25c8febe002b)] - api/news: add news query builder with default(Ching-Yang, Tseng)
- [[4676928](https://github.com/twreporter/go-api/commit/4676928391de22017453ca7bc16deba4cc772442)] - api/news: move GetTopics endpoint to controller(Ching-Yang, Tseng)
- [[a30e27b](https://github.com/twreporter/go-api/commit/a30e27b42a9d857a7f6e9624bbea3df06cd0c39e)] - api/news: refactor server side error handler(Ching-Yang, Tseng)
- [[e0a075c](https://github.com/twreporter/go-api/commit/e0a075c6e4e16587e81a8d3204148c5d2a1f9968)] - api/news: move GetATopic endpoint to controller(Ching-Yang, Tseng)
- [[004ff52](https://github.com/twreporter/go-api/commit/004ff52f93e2dbe63672cdeb2ee6a7a98bddc712)] - api/news: moves GetPosts endpoint to controller(Ching-Yang, Tseng)
- [[5288844](https://github.com/twreporter/go-api/commit/5288844afc74163f673193c0b9b624d9a7883ec9)] - core: fix test fail from function signature change(Ching-Yang, Tseng)
- [[fa1c2a4](https://github.com/twreporter/go-api/commit/fa1c2a4aba34967c9af3f46b816a48dc1b697928)] - core: upgrade golang to the 1.14.4(Ching-Yang, Tseng)
- [[f7255f2](https://github.com/twreporter/go-api/commit/f7255f21f89d24ffa29c7a12a7ba8736fb9f4d3c)] - api/news: moves GetAPost endpoint into controller(Ching-Yang, Tseng)
- [[f7fb0b3](https://github.com/twreporter/go-api/commit/f7fb0b36fe5e7e74700cc9008cd76aac564f5d9b)] - api/news: adapt to current directory layout(Ching-Yang, Tseng)
- [[6e715db](https://github.com/twreporter/go-api/commit/6e715db91d7f7d643caac8ddac128576b7cc3ab7)] - api/news: move model into internal/models(Ching-Yang, Tseng)
- [[63993c4](https://github.com/twreporter/go-api/commit/63993c4c6724ef3217fc72b0a5bd4f193b57a3c7)] - api/news: refactor lookup stages(Ching-Yang, Tseng)
- [[83298ea](https://github.com/twreporter/go-api/commit/83298ea7d162fb8a5bd2f8e7f18d09288521c610)] - api/news: convert mongoQuery to query documents(Ching-Yang, Tseng)
- [[1b9a1f5](https://github.com/twreporter/go-api/commit/1b9a1f53fa72ccf23f7cd1cb0a521813d3042fe8)] - api/news: convert Query to mongoQuery object(Ching-Yang, Tseng)
- [[7396dcf](https://github.com/twreporter/go-api/commit/7396dcf37253873a2bf4a02694fc006545da4f3d)] - api/news: refactor topic list query parser(Ching-Yang, Tseng)
- [[79fe61a](https://github.com/twreporter/go-api/commit/79fe61a2aafc87c12a7c6a0d8dd499a137d07289)] - api/news: refactor single topic query parser(Ching-Yang, Tseng)
- [[7714683](https://github.com/twreporter/go-api/commit/7714683c92f723ba2ee763af3c8abb6cf7677c16)] - api/news: define constant string variables(Ching-Yang, Tseng)
- [[8287863](https://github.com/twreporter/go-api/commit/82878636d76d4e3ac6cb43d6277b2327798e7b65)] - api/news: refactor post list query(Ching-Yang, Tseng)
- [[1851f35](https://github.com/twreporter/go-api/commit/1851f357875a7f31a516fb3f8ebbf55902556ffa)] - api/news: refactor query for single post retrieval(Ching-Yang, Tseng)
- [[8d37028](https://github.com/twreporter/go-api/commit/8d37028158aaa4a85ed3a065c5dece8680263900)] - api/news: implement posts query filter(Ching-Yang, Tseng)
- [[9163756](https://github.com/twreporter/go-api/commit/916375667356a03eb41837f6072e087a227656dc)] - api/news: remove filter during index page fetch(Ching-Yang, Tseng)
- [[14d04df](https://github.com/twreporter/go-api/commit/14d04df242f4904ae0bfb5eaa3f80a7d5c91ea8c)] - api/news: implement posts/topics list(Ching-Yang, Tseng)
- [[dc4a5b6](https://github.com/twreporter/go-api/commit/dc4a5b6abab18c3573f186af864c34c296b8ec87)] - api/news: implement index page fetch(Ching-Yang, Tseng)
- [[d5bd96d](https://github.com/twreporter/go-api/commit/d5bd96d47729bb03232b399a774f49b00ba14a6a)] - api/news: implement GetTopics in storage layer(Ching-Yang, Tseng)
- [[40421bb](https://github.com/twreporter/go-api/commit/40421bbf1cc397f4ec0a1552d67b37524847f681)] - api/news: implement GetPosts in storage layer(Ching-Yang, Tseng)
- [[dbd641b](https://github.com/twreporter/go-api/commit/dbd641bcf3fe68a17c0f0113bbc9d2361eeaee51)] - api/news: add storage layer function signature(Ching-Yang, Tseng)
- [[19fa0d8](https://github.com/twreporter/go-api/commit/19fa0d8539c1e6b9a30f6603ea990223b796ca31)] - api/news: adjust storage interface w.r.t query(Ching-Yang, Tseng)
- [[b91b5c6](https://github.com/twreporter/go-api/commit/b91b5c627ea54ef5d93bb40daadecdf99855d077)] - api/news: prototype the post/topic query model(Ching-Yang, Tseng)
- [[c864622](https://github.com/twreporter/go-api/commit/c86462242cf5b61f44b62e56cd5806d79ee66b2d)] - api/news: prototype v2 controller(Ching-Yang, Tseng)
- [[3f6b2f1](https://github.com/twreporter/go-api/commit/3f6b2f11f1af8a1d121ec6cf6ada09cc0ae735fa)] - core: add mongo db connection with new driver(Ching-Yang, Tseng)
- [[fea7d18](https://github.com/twreporter/go-api/commit/fea7d18078ba95b8a895e3ee13cb7cfddde2bcd0)] - api/news: add new model for posts(babygoat)
- [[c14a339](https://github.com/twreporter/go-api/commit/c14a33993395ac316eeb6f485f5b0e3b670d6d9d)] - doc: fix `full` field type in Topic(Ching-Yang, Tseng)
- [[6533e85](https://github.com/twreporter/go-api/commit/6533e85511d6f7efca67b59a22841e13b1bff332)] - doc: fix `full` field type in Post group(Ching-Yang, Tseng)
- [[f19d8a5](https://github.com/twreporter/go-api/commit/f19d8a5ea6cc0fb71366f2731fb0c6e5ba1f49a4)] - api/news: replace field writters with writers(Ching-Yang, Tseng)
- [[2637e1b](https://github.com/twreporter/go-api/commit/2637e1b114012f3c89d492da0d15402f6130d563)] - doc: return meta instead of empty content during post/topic list(Ching-Yang, Tseng)
- [[aeb4a20](https://github.com/twreporter/go-api/commit/aeb4a20d8c6d7132b47a1dd9df2f5b5883743507)] - doc: add missing required fields(Ching-Yang, Tseng)
- [[8d2ba67](https://github.com/twreporter/go-api/commit/8d2ba673bfd88997a04d19168d2121e300382892)] - doc: singularize topics field(Ching-Yang, Tseng)
- [[7f177d2](https://github.com/twreporter/go-api/commit/7f177d2606ec80c2aec7d19e973ddce84d9be589)] - doc: fix nested array field schema type missing(Ching-Yang, Tseng)
- [[e9f7f76](https://github.com/twreporter/go-api/commit/e9f7f769345a0d89d457d07dfc9e4bb3ba8ed63c)] - doc: give relateds sample value instead of default(Ching-Yang, Tseng)
- [[7c74107](https://github.com/twreporter/go-api/commit/7c7410760cc63ca96bc2871b51d1a6188b74fd4e)] - doc: add client side error for invalid slug(Ching-Yang, Tseng)
- [[93af752](https://github.com/twreporter/go-api/commit/93af75242aed38975050cd1fa928cbaf00d589e0)] - doc: describe different response w.r.t parameter(Ching-Yang, Tseng)
- [[1590853](https://github.com/twreporter/go-api/commit/1590853f052c397777ca3057c9dd602047b27951)] - doc: adjust /v2/topics endpoints(Ching-Yang, Tseng)
- [[dcd8975](https://github.com/twreporter/go-api/commit/dcd8975adba21eaad4dba69093f9395d595cb79f)] - doc: adjust /v2/posts endpoints(Ching-Yang, Tseng)
- [[f6aba81](https://github.com/twreporter/go-api/commit/f6aba810246cee78d04f9e06b7feb91e098c54f5)] - doc: add /v2/index_page endpoint(Ching-Yang, Tseng)
- [[9ffbfa3](https://github.com/twreporter/go-api/commit/9ffbfa32bcef9895f2cc667ada8a2638d15e5150)] - doc: generate result document(Ching-Yang, Tseng)
- [[c74559f](https://github.com/twreporter/go-api/commit/c74559fa7463cac0a2a48123bdd7a447f86aa1a4)] - doc: add the v2 endpoint of listing topic(Ching-Yang, Tseng)
- [[5f876e9](https://github.com/twreporter/go-api/commit/5f876e9fa075447fc628e52a7e24ffb729db93a3)] - doc: add the v2 endpoint to get a single topic(Ching-Yang, Tseng)
- [[0bb1f95](https://github.com/twreporter/go-api/commit/0bb1f95d8d9420240a21a81ee7eeb403e2b6fad8)] - doc: add v2 endpoint for listing posts(Ching-Yang, Tseng)
- [[c5729ab](https://github.com/twreporter/go-api/commit/c5729ab46903dffb14d4e2f189465438971dbde4)] - doc: adds v2 endpoint of getting a post(Ching-Yang, Tseng)

### 6.0.4, 2020-06-10

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
* [[`95830c578d`](https://github.com/twreporter/go-api/commit/95830c578d)] - **doc**: update CHANGELOG.md (Taylor Fang)
* [[`73dad89691`](https://github.com/twreporter/go-api/commit/73dad89691)] - **chore**: donations table schema change (Taylor Fang)
* [[`90541811e7`](https://github.com/twreporter/go-api/commit/90541811e7)] - **doc**: fix json format (Taylor Fang)
* [[`9098549594`](https://github.com/twreporter/go-api/commit/9098549594)] - **doc**: update docs according to review comment (Taylor Fang)
* [[`03da19a4ce`](https://github.com/twreporter/go-api/commit/03da19a4ce)] - **doc**: update docs to add `receipt\_header` field (Taylor Fang)
* [[`83b7f799c8`](https://github.com/twreporter/go-api/commit/83b7f799c8)] - **doc**: update CHANGELOG.md (nickhsine)
* [[`1be3f110dd`](https://github.com/twreporter/go-api/commit/1be3f110dd)] - api/mail: update success donation email template (nickhsine)
* [[`95c0fba4d2`](https://github.com/twreporter/go-api/commit/95c0fba4d2)] - **doc**: Update Changelog (Ching-Yang, Tseng)
* [[`219646de12`](https://github.com/twreporter/go-api/commit/219646de12)] - api/news: expand throught by new connections (Ching-Yang, Tseng)
* [[`47685da57c`](https://github.com/twreporter/go-api/commit/47685da57c)] - **core**: change mongo query mode (Ching-Yang, Tseng)
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
