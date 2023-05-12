# Changelog

## [0.0.4](https://github.com/Excoriate/aws-secrets-rotation-lambda/compare/v0.0.3...v0.0.4) (2023-05-12)


### Other

* add status badges fixed ([da3afe5](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/da3afe584a02ed95b20c6fe346b92d21a2523898))

## [0.0.3](https://github.com/Excoriate/aws-secrets-rotation-lambda/compare/v0.0.2...v0.0.3) (2023-05-11)


### Other

* remove old logo ([6fe3001](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/6fe3001540d5b51c43cc58ac71c3474df058a589))

## [0.0.2](https://github.com/Excoriate/aws-secrets-rotation-lambda/compare/v0.0.1...v0.0.2) (2023-05-11)


### Features

* add dagger pipeline first steps ([44ac651](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/44ac6514006922f91535ce09dfb5835c18365689))
* Add dagger task to push file to S3 ([d3e4e36](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/d3e4e369d17562e6f633318721e0f522af00949c))
* add finalize rotation step, add feature flag ([4927326](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/4927326448b67fee96d41e05525a381e5933505d))
* Add first structure ([e06acdd](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/e06acddca8db48163d70354b932c3e3bbec5fc7d))
* add infratructure tasks in dagger ([ab44b04](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/ab44b042e7dc84f1976e30fc4e4cea069c24fd4d))
* add package logic ([23184ab](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/23184ab142b31e8ce8131bfe8f62905a3f739951))
* complete create secret step ([c40d016](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/c40d0166163786e628c55012764d8c73cb493e42))


### Bug Fixes

* Add correct type checking for ResourceNotFound error while finalizing a rotation ([fcb450f](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/fcb450f05898eba9ca6455d10029306be06ff0ce))
* fix put secret value operation ([3dc1488](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/3dc148882157faa5af021b959f4534e5c08623b8))
* invalid validation looking for the secret version with AWSCURRENT label ([f11c51e](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/f11c51ed8bce5634452dd9ecb707bca4041b757c))


### Refactoring

* add working version of lambda iac module ([65a81bf](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/65a81bf76238c3a3c331f293c6ff8901e12910a6))
* adjust lambda compile task in dagger ([9b256b1](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/9b256b1c1fd2b3b72fdcc22263a1e8e2f82f6f16))
* change name for package zip file command ([93082da](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/93082da31ba8c453073bd9c674e2eb29a0744186))
* Enhance existing code, add data module ([48df38c](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/48df38ca070e74062fa2c53e572814fc910e6351))
* fix random bugs, add rotator in infra pipeilne, enhance name for demo resources ([d822cad](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/d822cad5131706a3aa799959631dad413742d89f))


### Other

* add docs ([b311a51](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/b311a517f72af01c4643c15aff993c100f687748))
* add missing gitignore statements ([369ced6](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/369ced619315068d3c67c76006e9ff1e99fa1181))
* add release.yml ([30ef859](https://github.com/Excoriate/aws-secrets-rotation-lambda/commit/30ef859584a0a08637487c8aeac61c857f53125c))
