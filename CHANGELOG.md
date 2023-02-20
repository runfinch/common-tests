# Changelog

## [0.6.0](https://github.com/runfinch/common-tests/compare/v0.5.0...v0.6.0) (2023-02-20)


### Features

* add tests for host-gateway speciap IP with equal sign ([#33](https://github.com/runfinch/common-tests/issues/33)) ([1296336](https://github.com/runfinch/common-tests/commit/1296336de63e3ac88c7d22acde97cc100d42b075))

## [0.5.0](https://github.com/runfinch/common-tests/compare/v0.4.0...v0.5.0) (2023-02-13)


### Features

* add tests for special IP in --add-host flag ([#29](https://github.com/runfinch/common-tests/issues/29)) ([1fecd9f](https://github.com/runfinch/common-tests/commit/1fecd9f5cb00982c88f2367eebdf4a78ad918c9c))

## [0.4.0](https://github.com/runfinch/common-tests/compare/v0.3.1...v0.4.0) (2023-02-01)


### Features

* add additional tests for env vars ([#26](https://github.com/runfinch/common-tests/issues/26)) ([d3b48e2](https://github.com/runfinch/common-tests/commit/d3b48e238cbb43e790d29bf33cb6d1adb39a2e16))


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.7.0 to 2.8.0 ([#27](https://github.com/runfinch/common-tests/issues/27)) ([723b70e](https://github.com/runfinch/common-tests/commit/723b70ed612517d279b1e851b755965b9d76bc27))
* **deps:** bump github.com/onsi/gomega from 1.24.2 to 1.26.0 ([#24](https://github.com/runfinch/common-tests/issues/24)) ([33e2c83](https://github.com/runfinch/common-tests/commit/33e2c8358089ad58edc5715909215196a18fb410))

## [0.3.1](https://github.com/runfinch/common-tests/compare/v0.3.0...v0.3.1) (2023-01-17)


### Bug Fixes

* Fix run -e/--env tests and add missing variable redefinition ([#22](https://github.com/runfinch/common-tests/issues/22)) ([84960f8](https://github.com/runfinch/common-tests/commit/84960f89215881460c3b6c462e02cd1f53f74878))

## [0.3.0](https://github.com/runfinch/common-tests/compare/v0.2.0...v0.3.0) (2023-01-12)


### ⚠ BREAKING CHANGES

* StdOut,StdErr -> Stdout,Stderr ([#20](https://github.com/runfinch/common-tests/issues/20))

### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.6.0 to 2.6.1 ([#15](https://github.com/runfinch/common-tests/issues/15)) ([ab4e024](https://github.com/runfinch/common-tests/commit/ab4e024075b03b34bd125d96d21c8361c6851f4f))
* **deps:** bump github.com/onsi/ginkgo/v2 from 2.6.1 to 2.7.0 ([#19](https://github.com/runfinch/common-tests/issues/19)) ([e695dc5](https://github.com/runfinch/common-tests/commit/e695dc51448406c809adb6395f8fa2db7cc0a9bd))
* **deps:** bump github.com/onsi/gomega from 1.24.1 to 1.24.2 ([#14](https://github.com/runfinch/common-tests/issues/14)) ([b4a7aa2](https://github.com/runfinch/common-tests/commit/b4a7aa2474ecc97bdb1a2283b04ea43ca2c769f7))


### refactor

* StdOut,StdErr -&gt; Stdout,Stderr ([#20](https://github.com/runfinch/common-tests/issues/20)) ([92fab5a](https://github.com/runfinch/common-tests/commit/92fab5a416075c802c2aaeef00e4ae263ff36aed))

## [0.2.0](https://github.com/runfinch/common-tests/compare/v0.1.1...v0.2.0) (2022-12-13)


### Features

* add e2e tests for resource and user flags ([#5](https://github.com/runfinch/common-tests/issues/5)) ([1d5ec0d](https://github.com/runfinch/common-tests/commit/1d5ec0db09b523f47f9825ef7237ed1d9c51401a))


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.5.1 to 2.6.0 ([#12](https://github.com/runfinch/common-tests/issues/12)) ([a676453](https://github.com/runfinch/common-tests/commit/a676453d03acf86b361202fb3d7e5414b66beb0d))

## [0.1.1](https://github.com/runfinch/common-tests/compare/v0.1.0...v0.1.1) (2022-11-30)


### Bug Fixes

* --pid=host tests ([#8](https://github.com/runfinch/common-tests/issues/8)) ([77342d8](https://github.com/runfinch/common-tests/commit/77342d8745bbb185bea2445cc150c0ff2dca0056))


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.5.0 to 2.5.1 ([#3](https://github.com/runfinch/common-tests/issues/3)) ([abf1f07](https://github.com/runfinch/common-tests/commit/abf1f07985e32a173032d7f49d9c4e621576ff47))
