# Changelog

## [0.6.4](https://github.com/runfinch/common-tests/compare/v0.6.3...v0.6.4) (2023-04-06)


### Bug Fixes

* better handling of concurrent http servers ([#57](https://github.com/runfinch/common-tests/issues/57)) ([0ae6182](https://github.com/runfinch/common-tests/commit/0ae6182b53beffdc279ffc1df7d539fe8205fd56))

## [0.6.3](https://github.com/runfinch/common-tests/compare/v0.6.2...v0.6.3) (2023-03-27)


### Build System or External Dependencies

* **deps:** bump github.com/onsi/gomega from 1.27.2 to 1.27.4 ([#48](https://github.com/runfinch/common-tests/issues/48)) ([f113e7b](https://github.com/runfinch/common-tests/commit/f113e7b2f65a66982773079a30dee06dd5e6e326))
* **deps:** bump github.com/onsi/gomega from 1.27.4 to 1.27.5 ([#52](https://github.com/runfinch/common-tests/issues/52)) ([bd056e7](https://github.com/runfinch/common-tests/commit/bd056e7d947a2432611ca923585422c399c56395))

## [0.6.2](https://github.com/runfinch/common-tests/compare/v0.6.1...v0.6.2) (2023-03-16)


### Bug Fixes

* Fix tests to match nerdctl 1.2.1 outputs ([#50](https://github.com/runfinch/common-tests/issues/50)) ([3d9b4f4](https://github.com/runfinch/common-tests/commit/3d9b4f4794d8df965dd2d611b2bed59aabff7dc2))


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.8.3 to 2.8.4 ([#41](https://github.com/runfinch/common-tests/issues/41)) ([a9476c1](https://github.com/runfinch/common-tests/commit/a9476c13bc4febd40a4f98cc8e6f8eebc04cfb5e))
* **deps:** bump github.com/onsi/gomega from 1.27.1 to 1.27.2 ([#40](https://github.com/runfinch/common-tests/issues/40)) ([e8fc71a](https://github.com/runfinch/common-tests/commit/e8fc71a9c94afe2084bfdb129de5f5828adfa8b8))

## [0.6.1](https://github.com/runfinch/common-tests/compare/v0.6.0...v0.6.1) (2023-02-28)


### Bug Fixes

* Switch from `nc -l` to `nginx` in `run -p/--publish` test ([7a6a6c3](https://github.com/runfinch/common-tests/commit/7a6a6c36d11796b2048d90353f06d25013b132a8))


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.8.0 to 2.8.3 ([#37](https://github.com/runfinch/common-tests/issues/37)) ([7b76f03](https://github.com/runfinch/common-tests/commit/7b76f03b77bb7a39b0a68aa6ad75942e67998e29))
* **deps:** bump github.com/onsi/gomega from 1.26.0 to 1.27.1 ([#36](https://github.com/runfinch/common-tests/issues/36)) ([e5a684e](https://github.com/runfinch/common-tests/commit/e5a684eada0303629645d600cf94cc49e8fbdba2))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#34](https://github.com/runfinch/common-tests/issues/34)) ([f218705](https://github.com/runfinch/common-tests/commit/f218705a28f93d8ae6463b75662c3ff108433e7b))

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
