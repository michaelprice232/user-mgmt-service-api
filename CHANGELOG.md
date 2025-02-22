# Changelog

## [1.3.1](https://github.com/michaelprice232/user-mgmt-service-api/compare/v1.3.0...v1.3.1) (2025-01-16)


### Bug Fixes

* allow pipeline to get JWT token ([#22](https://github.com/michaelprice232/user-mgmt-service-api/issues/22)) ([f0bb71f](https://github.com/michaelprice232/user-mgmt-service-api/commit/f0bb71f7767e1ef8af8340aa9424a98c3910347c))

## [1.3.0](https://github.com/michaelprice232/user-mgmt-service-api/compare/v1.2.0...v1.3.0) (2025-01-16)


### Features

* build final SemVer Docker images ([#20](https://github.com/michaelprice232/user-mgmt-service-api/issues/20)) ([ff1f3a3](https://github.com/michaelprice232/user-mgmt-service-api/commit/ff1f3a3ac7446bfed7cbd114d609aea4badf116f))

## [1.2.0](https://github.com/michaelprice232/user-mgmt-service-api/compare/v1.1.0...v1.2.0) (2025-01-15)


### Features

* added a --version CLI flag for displaying build version, GOOS & GOARCH ([6abe30b](https://github.com/michaelprice232/user-mgmt-service-api/commit/6abe30b03eaaef91d38ed7b961c004050768c47b))
* added DELETE /users/&lt;logon_name&gt; endpoint ([4fce427](https://github.com/michaelprice232/user-mgmt-service-api/commit/4fce42766cdd793ce7c4491f4c3cf6dd64eb293b))
* added functionality for reading from database instead of memory ([0392f75](https://github.com/michaelprice232/user-mgmt-service-api/commit/0392f75e5138c0e9f2dbaa5152a3f0b43c28ee66))
* added healthcheck endpoint on GET /health ([1af6d8a](https://github.com/michaelprice232/user-mgmt-service-api/commit/1af6d8a57070687a32d83eac7a6df89e54004a83))
* added initial set of unit tests covering happy paths ([94a91bb](https://github.com/michaelprice232/user-mgmt-service-api/commit/94a91bb361abbed23d8a80f0d8d0c45c7aa53eaa))
* added LogonName to User struct type to hold the unique system logon name for each user ([b2ae5b2](https://github.com/michaelprice232/user-mgmt-service-api/commit/b2ae5b27fca8a2f2e957ed56ba2d3d92ad58db02))
* added POST /users endpoint for adding new users ([32f7da1](https://github.com/michaelprice232/user-mgmt-service-api/commit/32f7da1d3b5677a4ffad55f24a5f09956965bbfd))
* added PUT /users/&lt;logon_name&gt; endpoint inc. tests ([11adc29](https://github.com/michaelprice232/user-mgmt-service-api/commit/11adc291504c9240509becb78658ee0366c2d931))
* added UserID field to User struct which will be used for managing objects by other CRUD endpoints ([7ac93f6](https://github.com/michaelprice232/user-mgmt-service-api/commit/7ac93f653c3bbf4596b858484c9a9d72e7ef7339))
* added validation for request payload in the POST /users handler ([94db8cf](https://github.com/michaelprice232/user-mgmt-service-api/commit/94db8cf4da922222457759f623642760d56c0fee))
* enable graceful shutdown of web server based on INT & TERM OS signals ([18ed2ae](https://github.com/michaelprice232/user-mgmt-service-api/commit/18ed2ae12bcc233a4798a038095c08347a6bc682))
* pull DB config from environment variables during app startup rather than hardcoded within the functions ([3cd48a1](https://github.com/michaelprice232/user-mgmt-service-api/commit/3cd48a19df599103988b2e2a89f720f772b67b26))
* Read pages on demand from the Database rather than loading the entire dataset into memory ([f966297](https://github.com/michaelprice232/user-mgmt-service-api/commit/f966297a11b681801ee564823a6d12dd0e770117))


### Bug Fixes

* added validation for checking that user_id has not been passed in the request payload ([1c76eda](https://github.com/michaelprice232/user-mgmt-service-api/commit/1c76eda40da0945abfeeebf9f833a3c0f9c35f57))
* ensure UsersResponse TotalPages/MorePages are valid when using name filters ([ac95e76](https://github.com/michaelprice232/user-mgmt-service-api/commit/ac95e764b0435473e842ac56e6886b9ad035b203))
* ensure validation failures of query strings are passed back as errors and processing stopped ([c77fc4c](https://github.com/michaelprice232/user-mgmt-service-api/commit/c77fc4cdb6bec45ba0a505a4bb3e12c64f8c7d4f))

## Changelog
