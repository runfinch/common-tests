# Changelog

## [0.7.22](https://github.com/runfinch/common-tests/compare/v0.7.21...v0.7.22) (2024-04-19)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/gomega from 1.32.0 to 1.33.0 ([#146](https://github.com/runfinch/common-tests/issues/146)) ([63d640b](https://github.com/runfinch/common-tests/commit/63d640bfee39c0b0547b6c0231ace5bdf2695156))

## [0.7.21](https://github.com/runfinch/common-tests/compare/v0.7.20...v0.7.21) (2024-03-28)


### Bug Fixes

* add a delay for system event monitoring to start before pull completes as the run commands are async ([#144](https://github.com/runfinch/common-tests/issues/144)) ([5de585f](https://github.com/runfinch/common-tests/commit/5de585f5bece7ed358928444cbde4cfe4426ff57))
* add custom wait for retry logic ([#141](https://github.com/runfinch/common-tests/issues/141)) ([3b69319](https://github.com/runfinch/common-tests/commit/3b693190773706dde9b6f8dd11171e26cc0df404))

## [0.7.20](https://github.com/runfinch/common-tests/compare/v0.7.19...v0.7.20) (2024-03-27)


### Bug Fixes

* image pull retry logic ([#139](https://github.com/runfinch/common-tests/issues/139)) ([4c30166](https://github.com/runfinch/common-tests/commit/4c30166cc5f7e1af73c7730aeea5dc72fea8d8d5))

## [0.7.19](https://github.com/runfinch/common-tests/compare/v0.7.18...v0.7.19) (2024-03-27)


### Bug Fixes

* Retry image pull for 3 times and then fail ([#137](https://github.com/runfinch/common-tests/issues/137)) ([3f4765f](https://github.com/runfinch/common-tests/commit/3f4765f82d255f710dae6aaf9a889b58e18f08ca))

## [0.7.18](https://github.com/runfinch/common-tests/compare/v0.7.17...v0.7.18) (2024-03-26)


### Bug Fixes

* use of uninitialized global variable ([#135](https://github.com/runfinch/common-tests/issues/135)) ([02a134d](https://github.com/runfinch/common-tests/commit/02a134d69a1eecf6406c50d3088b23645c2726ff))

## [0.7.17](https://github.com/runfinch/common-tests/compare/v0.7.16...v0.7.17) (2024-03-25)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.17.0 to 2.17.1 ([#132](https://github.com/runfinch/common-tests/issues/132)) ([c55e33b](https://github.com/runfinch/common-tests/commit/c55e33bf70ecce4f44d02d17a44a681739764abe))


### Bug Fixes

* track localImages in a new map to enable proper cleanup ([#133](https://github.com/runfinch/common-tests/issues/133)) ([c8a5e72](https://github.com/runfinch/common-tests/commit/c8a5e7222feb9fd476070727d8edec91068c1280))

## [0.7.16](https://github.com/runfinch/common-tests/compare/v0.7.15...v0.7.16) (2024-03-22)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.15.0 to 2.16.0 ([#123](https://github.com/runfinch/common-tests/issues/123)) ([146e16f](https://github.com/runfinch/common-tests/commit/146e16fe020f6ef94139d473b78103e374180647))
* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.16.0 to 2.17.0 ([#129](https://github.com/runfinch/common-tests/issues/129)) ([860b13e](https://github.com/runfinch/common-tests/commit/860b13edb022200db5462fdb74e00ec259d4f353))
* **deps:** Bump github.com/onsi/gomega from 1.31.1 to 1.32.0 ([#130](https://github.com/runfinch/common-tests/issues/130)) ([c61629e](https://github.com/runfinch/common-tests/commit/c61629e0221b6ce0ebe5c8ad83406f358f6aeff2))

## [0.7.15](https://github.com/runfinch/common-tests/compare/v0.7.14...v0.7.15) (2024-03-11)


### Bug Fixes

* Add a wait for server to start before doing a curl, to avoid sync issue ([#124](https://github.com/runfinch/common-tests/issues/124)) ([cb3138e](https://github.com/runfinch/common-tests/commit/cb3138e72fe7284e2c27a34859ce22e24ae440e9))

## [0.7.14](https://github.com/runfinch/common-tests/compare/v0.7.13...v0.7.14) (2024-02-27)


### Bug Fixes

* use new values of session ([#121](https://github.com/runfinch/common-tests/issues/121)) ([a09ae51](https://github.com/runfinch/common-tests/commit/a09ae519d21f7a6886a54f3c1a36f4d3ea5e8309))

## [0.7.13](https://github.com/runfinch/common-tests/compare/v0.7.12...v0.7.13) (2024-02-12)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.14.0 to 2.15.0 ([#115](https://github.com/runfinch/common-tests/issues/115)) ([70f9539](https://github.com/runfinch/common-tests/commit/70f953921afa4c81090a63a218613946200b3f19))
* **deps:** Bump github.com/onsi/gomega from 1.30.0 to 1.31.1 ([#117](https://github.com/runfinch/common-tests/issues/117)) ([70cc410](https://github.com/runfinch/common-tests/commit/70cc410777cb65616e80b30684be4491cf7be2de))


### Bug Fixes

* increase event timeout ([#118](https://github.com/runfinch/common-tests/issues/118)) ([22ca9e4](https://github.com/runfinch/common-tests/commit/22ca9e496d21e2a2b88cb149b4e8123375697518))

## [0.7.12](https://github.com/runfinch/common-tests/compare/v0.7.11...v0.7.12) (2024-01-13)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.13.2 to 2.14.0 ([#112](https://github.com/runfinch/common-tests/issues/112)) ([ece5ec6](https://github.com/runfinch/common-tests/commit/ece5ec6eab9870208c693bdf34ee8371373b3501))


### Bug Fixes

* allow propagation time for async conditions ([#111](https://github.com/runfinch/common-tests/issues/111)) ([8aeb41a](https://github.com/runfinch/common-tests/commit/8aeb41a875c62a0aedfe7042861c424d42ab7bb1))
* increase nginx pull timeout ([#114](https://github.com/runfinch/common-tests/issues/114)) ([33308d0](https://github.com/runfinch/common-tests/commit/33308d0ea9235454783bd88cc40de3d84852974a))

## [0.7.11](https://github.com/runfinch/common-tests/compare/v0.7.10...v0.7.11) (2024-01-06)


### Bug Fixes

* fix panic in HTTPGetAndAssert ([#109](https://github.com/runfinch/common-tests/issues/109)) ([b572343](https://github.com/runfinch/common-tests/commit/b5723431c20c68df7f93bbb99a04f6683368c8bc))

## [0.7.10](https://github.com/runfinch/common-tests/compare/v0.7.9...v0.7.10) (2024-01-05)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.13.1 to 2.13.2 ([#102](https://github.com/runfinch/common-tests/issues/102)) ([91f0e82](https://github.com/runfinch/common-tests/commit/91f0e82764480a5be385d3cdcee91d6b38e4d6be))


### Bug Fixes

* disable build --ssh test on Windows ([#106](https://github.com/runfinch/common-tests/issues/106)) ([e7fe1eb](https://github.com/runfinch/common-tests/commit/e7fe1ebc083f5ab66c7426a1d464e78c941faa12))

## [0.7.9](https://github.com/runfinch/common-tests/compare/v0.7.8...v0.7.9) (2023-11-21)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.13.0 to 2.13.1 ([#100](https://github.com/runfinch/common-tests/issues/100)) ([9aae0d5](https://github.com/runfinch/common-tests/commit/9aae0d51e3f52a4b52445a367291e9d9b6401bb1))
* **deps:** Bump github.com/onsi/gomega from 1.28.0 to 1.29.0 ([#98](https://github.com/runfinch/common-tests/issues/98)) ([1cb5cc6](https://github.com/runfinch/common-tests/commit/1cb5cc6d4f3ae20ca617f82c083caa5bc56f5531))
* **deps:** Bump github.com/onsi/gomega from 1.29.0 to 1.30.0 ([#99](https://github.com/runfinch/common-tests/issues/99)) ([fffae02](https://github.com/runfinch/common-tests/commit/fffae0218ec7836bc98d7392781918fed2fefa68))
* **deps:** Bump golang.org/x/net from 0.14.0 to 0.17.0 ([#95](https://github.com/runfinch/common-tests/issues/95)) ([4e8bcda](https://github.com/runfinch/common-tests/commit/4e8bcdae827a129fb08761594114561949793990))

## [0.7.8](https://github.com/runfinch/common-tests/compare/v0.7.7...v0.7.8) (2023-10-10)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.12.1 to 2.13.0 ([#93](https://github.com/runfinch/common-tests/issues/93)) ([dcd9dee](https://github.com/runfinch/common-tests/commit/dcd9dee430a5bdc8690472b34c675143ff56ec4c))

## [0.7.7](https://github.com/runfinch/common-tests/compare/v0.7.6...v0.7.7) (2023-10-05)


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/gomega from 1.27.10 to 1.28.0 ([#90](https://github.com/runfinch/common-tests/issues/90)) ([d8a87bb](https://github.com/runfinch/common-tests/commit/d8a87bb07ca00770c75fdd0dac7914c3304fbd37))

## [0.7.6](https://github.com/runfinch/common-tests/compare/v0.7.5...v0.7.6) (2023-09-21)


### Bug Fixes

* add --all to volume prune tests ([#87](https://github.com/runfinch/common-tests/issues/87)) ([9248bec](https://github.com/runfinch/common-tests/commit/9248bec81bbbd68b588a746bc409cd7b2c41ae03))

## [0.7.5](https://github.com/runfinch/common-tests/compare/v0.7.4...v0.7.5) (2023-09-21)


### Bug Fixes

* adds --all to the volume prune command to prune named volumes ([#86](https://github.com/runfinch/common-tests/issues/86)) ([4973e9f](https://github.com/runfinch/common-tests/commit/4973e9fa956b1339fa282d065576b417acfe2c52))
* Update logs test args ([#83](https://github.com/runfinch/common-tests/issues/83)) ([011c2e3](https://github.com/runfinch/common-tests/commit/011c2e335c4da40842bfacccab24779ce63aaa04))


### Build System or External Dependencies

* **deps:** Bump github.com/onsi/ginkgo/v2 from 2.12.0 to 2.12.1 ([#85](https://github.com/runfinch/common-tests/issues/85)) ([89d408f](https://github.com/runfinch/common-tests/commit/89d408f34ceb0be386cc0ff780aaa52638b267d5))

## [0.7.4](https://github.com/runfinch/common-tests/compare/v0.7.3...v0.7.4) (2023-09-20)


### Bug Fixes

* Fix container filepath to make it platform independent ([#80](https://github.com/runfinch/common-tests/issues/80)) ([5496e94](https://github.com/runfinch/common-tests/commit/5496e94a7ec5db81f58e787a7d6dcf29efab7e37))
* increase acceptable time deviation for stop tests with -t ([#81](https://github.com/runfinch/common-tests/issues/81)) ([c292f6d](https://github.com/runfinch/common-tests/commit/c292f6d8f79cf51b0ec2b5ec42db186e9e5661df))

## [0.7.3](https://github.com/runfinch/common-tests/compare/v0.7.2...v0.7.3) (2023-08-25)


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.11.0 to 2.12.0 ([#78](https://github.com/runfinch/common-tests/issues/78)) ([7bae5c1](https://github.com/runfinch/common-tests/commit/7bae5c16524336711c8258ab27c59c54ebf399cd))
* **deps:** bump github.com/onsi/gomega from 1.27.8 to 1.27.10 ([#74](https://github.com/runfinch/common-tests/issues/74)) ([f0f6fa1](https://github.com/runfinch/common-tests/commit/f0f6fa1a053db57b741068b208d463cb729a274f))

## [0.7.2](https://github.com/runfinch/common-tests/compare/v0.7.1...v0.7.2) (2023-08-08)


### Bug Fixes

* make tests compatible with nerdctlv1.5 ([#75](https://github.com/runfinch/common-tests/issues/75)) ([6876cd0](https://github.com/runfinch/common-tests/commit/6876cd046728c28f527b56770fd04735f7dc7067))

## [0.7.1](https://github.com/runfinch/common-tests/compare/v0.7.0...v0.7.1) (2023-06-27)


### Bug Fixes

* add retry to assert containers do not exist for compose down ([#73](https://github.com/runfinch/common-tests/issues/73)) ([88f732f](https://github.com/runfinch/common-tests/commit/88f732f12979b0064b852812db4b48affedf5e4c))


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.10.0 to 2.11.0 ([#71](https://github.com/runfinch/common-tests/issues/71)) ([45e9414](https://github.com/runfinch/common-tests/commit/45e9414dba27581a286784e16ff0ab54301220b2))
* **deps:** bump github.com/onsi/ginkgo/v2 from 2.9.5 to 2.10.0 ([#69](https://github.com/runfinch/common-tests/issues/69)) ([a6ad55d](https://github.com/runfinch/common-tests/commit/a6ad55dd08ee0d3316f51891a795e4f4f5dc9dcd))
* **deps:** bump github.com/onsi/gomega from 1.27.7 to 1.27.8 ([#68](https://github.com/runfinch/common-tests/issues/68)) ([6c72750](https://github.com/runfinch/common-tests/commit/6c7275007bf34fb6ddecc4013c16f1d79ff6d1d0))

## [0.7.0](https://github.com/runfinch/common-tests/compare/v0.6.5...v0.7.0) (2023-05-26)


### Features

* Tests for bind mounts ([#66](https://github.com/runfinch/common-tests/issues/66)) ([22a7f7e](https://github.com/runfinch/common-tests/commit/22a7f7e7bd917e443aa47aaa9eaa5dac03a5a10b))
* verify the result of finch inspect has State.Status and State.Error ([#64](https://github.com/runfinch/common-tests/issues/64)) ([b761a7a](https://github.com/runfinch/common-tests/commit/b761a7ab19fe15e0d0bf34441fad1248ac6b3e83))


### Build System or External Dependencies

* **deps:** bump github.com/onsi/ginkgo/v2 from 2.9.2 to 2.9.5 ([#62](https://github.com/runfinch/common-tests/issues/62)) ([0bd0901](https://github.com/runfinch/common-tests/commit/0bd090128548cdeb8cf381c8c53b2177fe009ab6))
* **deps:** bump github.com/onsi/gomega from 1.27.5 to 1.27.6 ([#54](https://github.com/runfinch/common-tests/issues/54)) ([72120b5](https://github.com/runfinch/common-tests/commit/72120b57b4c70945df307a1aea80d609e7c27a95))
* **deps:** bump github.com/onsi/gomega from 1.27.6 to 1.27.7 ([#65](https://github.com/runfinch/common-tests/issues/65)) ([590a984](https://github.com/runfinch/common-tests/commit/590a9845b46218c1c8d669a5e5a9269dfc86a232))

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
