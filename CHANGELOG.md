# Changelog

## v0.1.0-rc.2 (2022-06-18)

 * [3abef45](https://github.com/tensorchord/envd/commit/3abef452fb45bdcdbe4291caeae1ebd1a12589e4) fix: Fix the bug about uid (#447)
 * [eff6ffa](https://github.com/tensorchord/envd/commit/eff6ffac3dd6d7f0ffd0313dc4eec06eb753d4de) fix: Fix typo (#445)

### Contributors

 * Ce Gao

## v0.1.0-rc.1 (2022-06-18)

 * [6a35a57](https://github.com/tensorchord/envd/commit/6a35a579847163fe255cf981c85636fb2e4f3e5d) chore(README): Add who should use section (#442)
 * [6e1cf05](https://github.com/tensorchord/envd/commit/6e1cf0509844060040bbd85d4f19e29410fb7a6f) fix: replace useless .editorconfig (#440)
 * [1c23cea](https://github.com/tensorchord/envd/commit/1c23cea84bfb37f2cd5ee0df63bf625448147994) release: Separate alpha and stable release in Homebrew (#439)
 * [274e183](https://github.com/tensorchord/envd/commit/274e18317597bb7a7a6413f59a52e7d5274ac85c) Update PyTorch installation CMD in examples (#435)
 * [cfda1be](https://github.com/tensorchord/envd/commit/cfda1bed6f9861f198d0a37ebc8f46ed4bbb51ab) chore(deps): bump pypa/cibuildwheel from 2.6.1 to 2.7.0 (#428)
 * [d25157e](https://github.com/tensorchord/envd/commit/d25157ef8169f2aa7c7169dfa18b389041a036a0) chore(deps): bump github.com/spf13/viper from 1.4.0 to 1.12.0 (#430)
 * [d394410](https://github.com/tensorchord/envd/commit/d394410331eb25d20a33b6adef68506cbc0b6602) chore(deps): bump goreleaser/goreleaser-action from 2 to 3 (#427)
 * [7070506](https://github.com/tensorchord/envd/commit/7070506f6c2acb44692dd06090ec9d614927ef65) chore(deps): bump github.com/gliderlabs/ssh from 0.3.3 to 0.3.4 (#429)
 * [2e287c9](https://github.com/tensorchord/envd/commit/2e287c9b1ea4d16c822532e56a39ae67bdf1982c) chore(destroy): add current path `.` as the default path (#422)
 * [750db5a](https://github.com/tensorchord/envd/commit/750db5a20328f5cb10a118ae348b5545dfef5a1f) chore(Makefile): add `help` target (#421)
 * [0c00005](https://github.com/tensorchord/envd/commit/0c0000571b7ffb1a4a9af418b1c91a1c28253ce8) Bootstrap gets error if the envd_buildkitd was stopped before (#417)
 * [8a1bd1e](https://github.com/tensorchord/envd/commit/8a1bd1e3cd6d8372d5181fa2c873726583eb9029) feat #383 (#416)

### Contributors

 * Aaron Sun
 * Ce Gao
 * Kevin Su
 * Yuchen Cheng
 * Zhenzhen Zhao
 * dependabot[bot]
 * kenwoodjw

## v0.1.0-alpha.12 (2022-06-17)

 * [8531491](https://github.com/tensorchord/envd/commit/853149189d88e6ecf6a5a924e3aa19d8f7993f7e) fix: Fix default ssh shell (#411)
 * [5f3b16b](https://github.com/tensorchord/envd/commit/5f3b16bf579fe372ec2063bf3ce8904a4e5d2e2b) feat: Support configuring CRAN mirror for R environment (#405)
 * [8e27e99](https://github.com/tensorchord/envd/commit/8e27e9962588bd1b14268eaed6f6e067a9c5d908) fix: Only configure conda for Python environment (#406)
 * [6160899](https://github.com/tensorchord/envd/commit/6160899a10ec941f178dd821bee42eb615b932e3) feat(cli): support `envd build --output` (#402)

### Contributors

 * Ce Gao
 * Yuan Tang
 * Yuchen Cheng

## v0.1.0-alpha.11 (2022-06-17)

 * [d3fda6d](https://github.com/tensorchord/envd/commit/d3fda6db2783c93f4f5ea7954627500431059558) fix: Hack the gid (#399)
 * [e478c1a](https://github.com/tensorchord/envd/commit/e478c1a51191ab74f5c733d79106520032240e7f) bug: Fix source is released twice for macos and linux (#394)

### Contributors

 * Ce Gao
 * Jinjing Zhou

## v0.1.0-alpha.10 (2022-06-16)

 * [a55ad88](https://github.com/tensorchord/envd/commit/a55ad8808e5a23dba628f5a55c335603771083e9) feat(lang): Set default user to current (#390)
 * [df5bde3](https://github.com/tensorchord/envd/commit/df5bde376bcdd7db84ccd50b132cf05c3546deb9) feat(release): Support Homebrew in goreleaser (#389)
 * [7c10b71](https://github.com/tensorchord/envd/commit/7c10b71266c1b4ed250dc1a850603e930f531cc2) fix: cannot assign requested address (#386)
 * [d48e3ab](https://github.com/tensorchord/envd/commit/d48e3abf0b2443e1740dd0c2ddd5ebcf81e6fa6f) fix: Output error details when debug flag is enabled (#385)
 * [e2c8adb](https://github.com/tensorchord/envd/commit/e2c8adb41eb394d7fed56d5494f1ed5fc0832356) fix: use python3 explicitly to avoid type hints error (#379)
 * [fc2afe9](https://github.com/tensorchord/envd/commit/fc2afe9b6fe2bf9e9430039a0270de45981ba5ac) fix: add classifiers in setup.py (#380)
 * [33bdd7a](https://github.com/tensorchord/envd/commit/33bdd7a3313d18a9e781db98b033f5a7ceffe58b) doc: Add universe api doc (#374)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Jun
 * Keming
 * Manjusaka
 * Yuchen Cheng

## v0.1.0-alpha.9 (2022-06-16)

 * [3b3945a](https://github.com/tensorchord/envd/commit/3b3945aea9297fc69a8f9c787ef27812f1d0efb9) fix: Add v before tags (#371)

### Contributors

 * Ce Gao

## v0.1.0-alpha.8 (2022-06-16)


### Contributors


## v0.1.0-alpha.7 (2022-06-16)

 * [f60a976](https://github.com/tensorchord/envd/commit/f60a9766b2744c93e026d9be1094847dd0e9949a) enhancement(CLI): Update the description of envd (#364)
 * [9640bf1](https://github.com/tensorchord/envd/commit/9640bf1ce0a2cb8ed663aa1093a309bd60fac627) fix: config pip source speed up in china (#354)
 * [1b56ce2](https://github.com/tensorchord/envd/commit/1b56ce2db0deede211b3db5f783df07f9e94530d) add cpu example (#338)
 * [a410552](https://github.com/tensorchord/envd/commit/a41055285eaa08f96b97bd4d6c2f88ac74506c76) fix: remove py wrapper traceback information (#341)
 * [ae629bb](https://github.com/tensorchord/envd/commit/ae629bbb868ca4c6bf129817a48a59191f0e9605) feat: Support specifying number of GPUs (#336)
 * [7d577f7](https://github.com/tensorchord/envd/commit/7d577f72f6810c6c9a244c6af750ac783a0d0064) feat: Suport conda env (#335)
 * [6fe2ae0](https://github.com/tensorchord/envd/commit/6fe2ae0d329abe4a34b0c232357b78df8a4bf6a9) manually use docker distribution 2.8.1 (#333)
 * [2327ffd](https://github.com/tensorchord/envd/commit/2327ffd84048a7d8b6befc65ee2e72796a4603d0) fix: Disable unit test in macOS (#328)
 * [c80082b](https://github.com/tensorchord/envd/commit/c80082b708c7666db27ee8d6abf98f82c1234e54) workflow: enable macOS in CI without conditions (#327)
 * [35ef36d](https://github.com/tensorchord/envd/commit/35ef36dc24b4d83937120a7643eebd58f123fc47) fix: pypi sdist (#318)
 * [2b81df6](https://github.com/tensorchord/envd/commit/2b81df67fbdd8b4eb799c890c68e2683c6edb6b7) fix: typo in readme (#325)
 * [85123b6](https://github.com/tensorchord/envd/commit/85123b6d427363f595f22ae0410bec6de7a092ef) fix: fix typo (#324)
 * [559e143](https://github.com/tensorchord/envd/commit/559e1435a9585fa32cee3eff73a11a577bcec111) chore(CI): Enable code coverage (#323)
 * [a4cb9dc](https://github.com/tensorchord/envd/commit/a4cb9dc0bf93970d604ac94dd88074f641bcff9b) fix(release): Change docker user (#321)
 * [a5c3427](https://github.com/tensorchord/envd/commit/a5c3427c3feb38dbe5f0d7fdd633a8c978129837) chore(deps): bump github.com/moby/buildkit from 0.10.1 to 0.10.3 (#313)
 * [43ad124](https://github.com/tensorchord/envd/commit/43ad1249ff938277b65bb91d7d8cf6a128380ad3) chore(deps): bump github.com/pkg/sftp from 1.13.4 to 1.13.5 (#309)
 * [5a9c947](https://github.com/tensorchord/envd/commit/5a9c947edf1763745f1523ba7e4f2acf5f476990) chore(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.2 (#310)
 * [4fb34e2](https://github.com/tensorchord/envd/commit/4fb34e29471b0a5625d075af02854f784aceb8a7) chore(deps): bump github.com/urfave/cli/v2 from 2.4.0 to 2.8.1 (#312)
 * [0c63064](https://github.com/tensorchord/envd/commit/0c6306476922f0abcf30771a7d724b150e2188b6) fix: add api/__init__.py (#317)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming
 * Ling Jin
 * Xu Jin
 * Yuan Tang
 * Yuchen Cheng
 * Zhizhen He
 * dependabot[bot]
 * kenwoodjw

## v0.1.0-alpha.6 (2022-06-13)

 * [12cf334](https://github.com/tensorchord/envd/commit/12cf3345b6c09106508271dd84bd41bc03ceedbc) fix: Fix twine (#301)

### Contributors

 * Ce Gao

## v0.1.0-alpha.5 (2022-06-13)

 * [f42e162](https://github.com/tensorchord/envd/commit/f42e1625fd11332411b821931d0664494bfc1927) fix: Instal twine (#300)

### Contributors

 * Ce Gao

## v0.1.0-alpha.4 (2022-06-13)

 * [7720529](https://github.com/tensorchord/envd/commit/7720529500b58d854c09c817485f0edc2a1198dc) feat(lang): Support config.conda_channel and install.conda_packages (#293)
 * [452f3dc](https://github.com/tensorchord/envd/commit/452f3dc8033d1d2163aed7bc61bcd8d54ad81aec) feat: add pypi sdist (#298)
 * [cfe65fe](https://github.com/tensorchord/envd/commit/cfe65fe690f050a8e759063c3d6a3f71aa051f05) fix: py27 subprocess (#296)
 * [5cf52ec](https://github.com/tensorchord/envd/commit/5cf52ec7f03f94394aa17140986cb85545dfe942) fix: Update readme about installation (#295)

### Contributors

 * Ce Gao
 * Keming

## v0.1.0-alpha.3 (2022-06-13)

 * [3bf2710](https://github.com/tensorchord/envd/commit/3bf27107eb3bf01916ec703b9ae697dd87a92ad7) action: Add pypi release pipeline (#277)
 * [35e6e1b](https://github.com/tensorchord/envd/commit/35e6e1baa1af2c8b25ba47f88949f049b996d5d9) workflow: Enable macOS in CI (#287)
 * [2bc13df](https://github.com/tensorchord/envd/commit/2bc13df54f4ce9150db0519986344e781c3e5f32) bug: fix version without tag (#288)
 * [ef3886c](https://github.com/tensorchord/envd/commit/ef3886c7730540104737e207970afcb0b3876c2a) Revert "workflow: enable macOS in CI (#280)" (#286)
 * [02f83aa](https://github.com/tensorchord/envd/commit/02f83aa6bb7d1944988ab2e62e328d6cfcd3ff77) workflow: enable macOS in CI (#280)
 * [2cd9a0e](https://github.com/tensorchord/envd/commit/2cd9a0e3fa0191d801e8c8c3d6e4a64244c58861) fix: Update contributing (#284)
 * [f166cf8](https://github.com/tensorchord/envd/commit/f166cf8059cf658c59c51f155402ddf5eef9a922) feat: Support destroy environment by name (#281)
 * [9c237e4](https://github.com/tensorchord/envd/commit/9c237e4cad996be5baf1116d8794708bf598c893) fix: Bump version and fix base image (#279)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Yuan Tang
 * Yuchen Cheng

## v0.1.0-alpha.2 (2022-06-12)

 * [e048fc0](https://github.com/tensorchord/envd/commit/e048fc06c3f1b5dd2d0f69aded389a67a4608ced) fix: Use 127.0.0.1 instead of containerIP in ssh (#276)
 * [ae16402](https://github.com/tensorchord/envd/commit/ae16402016f27505d100d820d6a2aeba1aa9838a) fix: Hard code OS (#270)

### Contributors

 * Ce Gao

## v0.1.0-alpha.1 (2022-06-11)

 * [0bf757f](https://github.com/tensorchord/envd/commit/0bf757f3a371a5546ffc223505b5c6839b5c459f) fix: Fix typo in the file name (#266)
 * [846dc0e](https://github.com/tensorchord/envd/commit/846dc0efc67c0eb72bc45d93b055ec4edb49bfd1) feat: Support only print the version number (#265)
 * [5e82ccb](https://github.com/tensorchord/envd/commit/5e82ccb281d399ed7196f1de0fab02759900c60b) fix: Typo (#264)
 * [45e4562](https://github.com/tensorchord/envd/commit/45e456257ba23513a77f9bf27f0392cc5e003376) add api doc (#262)
 * [723b32f](https://github.com/tensorchord/envd/commit/723b32fe3c2ba40d691ae518077ac21255609014) fix: Set default value to GUID (#260)
 * [c4f525a](https://github.com/tensorchord/envd/commit/c4f525ae8cfd377cb8b0de4b7a026f2098627ef9) feat: add pypi package (#258)
 * [66d83d6](https://github.com/tensorchord/envd/commit/66d83d699651cf5b20b4acafe0c2369f1945c561) feat: Initial support for R language (#257)
 * [ba64556](https://github.com/tensorchord/envd/commit/ba64556c7effa0d687089ec5a7987f0d7874a991) fix: Fix summary (#256)
 * [473da34](https://github.com/tensorchord/envd/commit/473da34cc754524b46d4135dcceb1f56f3398e99) fix: Remove default assignee (#254)
 * [08a75eb](https://github.com/tensorchord/envd/commit/08a75ebdab560ca2c1989bcc98ec4ec60ede46ce) feat: Move cmd to pkg/app (#250)
 * [15d9c51](https://github.com/tensorchord/envd/commit/15d9c51f3b656d61f70370b67bfa9e97ca5e075e) fix: Incorrect cache ID directory (#251)
 * [8f8c81e](https://github.com/tensorchord/envd/commit/8f8c81eeb88599dffcf3baa99c6e1539b067dd78) Readme: Fix Readme (#247)
 * [3886502](https://github.com/tensorchord/envd/commit/388650253f152b00bcfc07273d2aaf414824cdd7) feat(lang): Fix pip cache (#243)
 * [1d930e1](https://github.com/tensorchord/envd/commit/1d930e1b31843633aa36f562d6333253299d42b2) feat: Refactor syntax (#238)
 * [289cb07](https://github.com/tensorchord/envd/commit/289cb07ff0b7a3688a9c7032762917efa0766173) feat(CLI): Add dockerhub mirror flag (#242)
 * [1a50052](https://github.com/tensorchord/envd/commit/1a50052216cb418be6f57d4a369d40c597961fa9) fix(zsh): Remnant characters when tab (#239)
 * [ced22ea](https://github.com/tensorchord/envd/commit/ced22ea6899f51fa6f97ea2852dd4ff97a7ff747) feat: Add version command with enriched information (#236)
 * [f14cd6a](https://github.com/tensorchord/envd/commit/f14cd6aed69eb40d4a6d520372c4ee327888113b) feat(lang): Support git_config rule (#235)
 * [8555bc6](https://github.com/tensorchord/envd/commit/8555bc6e24b2fb16f8e41fde358ce7c9515589ed) fix(README): Add build in readme (#234)
 * [ebf3197](https://github.com/tensorchord/envd/commit/ebf3197ec9b44004dbbf3a46539ab90024d51a9c) fix(README): Fix pip index doc (#233)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming
 * Yuan Tang
 * Yuchen Cheng

## v0.0.1-rc.1 (2022-06-02)

 * [2cac96f](https://github.com/tensorchord/envd/commit/2cac96ff3a80ea23140abae03870002556c4e2d8) fix: Add missing release binaries for Darwin (#231)
 * [e0194f9](https://github.com/tensorchord/envd/commit/e0194f9b8650427546557becef3bafb0b4f218ba) fix: Connecting timeout to wait for buildkitd is ignored (#230)
 * [8c0c98d](https://github.com/tensorchord/envd/commit/8c0c98d6741b70a0afb17bd93b500df826d1a2a9) feat: Support extra PyPI index (#229)
 * [31514df](https://github.com/tensorchord/envd/commit/31514dfe7b24a7830a8a5c4c845ad56f640112c8) feat(CLI): Add pause and resume (#228)
 * [3184df9](https://github.com/tensorchord/envd/commit/3184df9e4007a9247572827e5933651ee626601d) feat: Use build function as the default target (#226)
 * [6a4b7b5](https://github.com/tensorchord/envd/commit/6a4b7b569bfec867657d222a8bd3f0079584727a) feat(CLI): List deps (#223)
 * [2d8f6fb](https://github.com/tensorchord/envd/commit/2d8f6fbe5b919ab8459769f712a067635bd537c4) fix(README): Add why and how does it work (#225)
 * [a2ffdda](https://github.com/tensorchord/envd/commit/a2ffdda36144727cc642bc4ce9df274498a55efc) fix: Optimize err when jupyter port is already allocated (#224)
 * [f62355f](https://github.com/tensorchord/envd/commit/f62355f4d2e57dc303ee34af177205647117441f) fix(workflow): Request CI build only when review is required (#227)
 * [855a5f5](https://github.com/tensorchord/envd/commit/855a5f5c480297ed1d80d6e7056d0cc730f5ddeb) fix: Add gpu error message (#217)
 * [93face5](https://github.com/tensorchord/envd/commit/93face5fa3b20a36735930d144b3a1cc47773eff) fix: Add a asciinema demo (#222)
 * [cd2c0db](https://github.com/tensorchord/envd/commit/cd2c0db5b04dce639eff2144ae4fbbe04291020e) fix(home): Remove afterall to avoid flaky tests (#215)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Yuan Tang

## v0.0.1-alpha.6 (2022-05-29)

 * [b8a5516](https://github.com/tensorchord/envd/commit/b8a551676c61e27ddea2d35e53a697f10717bdea) fix(ssh): Do not create the ssh key if the key exists (#211)
 * [40c7d06](https://github.com/tensorchord/envd/commit/40c7d062f0130aeaac25b2828ae7798e50b0180d) feat: Auth ssh with key (#205)
 * [e0984bb](https://github.com/tensorchord/envd/commit/e0984bb112d21c4aab36a48cd3a04b9cb27f59a6) feat: Add prefix for cache id (#204)
 * [ee31696](https://github.com/tensorchord/envd/commit/ee3169648d645fdff0c26ac04a88a9ac9a6bfa72) fix(Makefile): Fix addlicense for more general use (#207)
 * [432f497](https://github.com/tensorchord/envd/commit/432f497ddef84206a54651e3ad88437a8c003f44) Fix: Fix (#199)
 * [728e3a9](https://github.com/tensorchord/envd/commit/728e3a9a7383cc6f258040c8a1fec9b84ec0a7d3) fix: Poll instead of err (#197)
 * [6ed91e5](https://github.com/tensorchord/envd/commit/6ed91e5b401073cce9d6442081532bad48426bda) feat(lang): Support pip cache with uid (#198)
 * [a589700](https://github.com/tensorchord/envd/commit/a589700c356d24735be5e37f42a3ac383e7f931b) feat(cli): Add ls command to list all envs (#177)
 * [5d9b7a2](https://github.com/tensorchord/envd/commit/5d9b7a2976b6154a1ca150ed407ac0d31a3a9574) fix(Makefile): use `$$` to represent `$` for shell command (#196)
 * [0964149](https://github.com/tensorchord/envd/commit/09641491bd5f3f81db6d3fd7980997d0d283a26d) fix(CI): Enable check (#186)
 * [4ce2890](https://github.com/tensorchord/envd/commit/4ce2890611d3dc67bdaa721bddd68ecdac1bc9a4) feat(ssh): Add prefix (#182)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Yuchen Cheng
 * Zhenzhen Zhao

## v0.0.1-alpha.5 (2022-05-20)

 * [a76b2aa](https://github.com/tensorchord/envd/commit/a76b2aa8ebb528a019a363e89d9a22b36308add5) feat(vscode): Use openvsx (#174)
 * [0968f83](https://github.com/tensorchord/envd/commit/0968f83502b4595477cb2bea602f0be1711f1baa) chore(README): Update badge (#175)
 * [3b71950](https://github.com/tensorchord/envd/commit/3b7195079892e5e3f33e20642de7995fa9c3276a) feat: Use UID 1000 to build (#173)
 * [161927c](https://github.com/tensorchord/envd/commit/161927c29b50c64093c39a2fc77984e38b46d4a8) feat(ssh): Add/Remove SSH config entry (#172)
 * [858b2fd](https://github.com/tensorchord/envd/commit/858b2fd4019c6336d848b18408d84c04c956d16d) feat(vscode): fine-grained cache management (#164)
 * [d600431](https://github.com/tensorchord/envd/commit/d60043174eca672b840b7de80a1b01a68000acd1) fix(lang): Fix user-defined packages (#168)
 * [74c8e20](https://github.com/tensorchord/envd/commit/74c8e205c9cee57a650be7b59bf837fa1be2dc97) chore(test): Add test cases for shell and home manager (#163)
 * [32ea318](https://github.com/tensorchord/envd/commit/32ea31892f4c08fceff4703b2dffacad8a6a2fcf) feat(runtime): Use tini to manage process (#162)
 * [0ca21d2](https://github.com/tensorchord/envd/commit/0ca21d2cbfbe8b8483a33680bce3c18f83de63aa) chore(CLI): e2e build test (#152)
 * [a7735f1](https://github.com/tensorchord/envd/commit/a7735f1c28a242546dad31961055e1f2282e8438) feat(CLI): fine-grained cache management (#155)
 * [33cdee2](https://github.com/tensorchord/envd/commit/33cdee26f4b2812919d6266ce9ff46211c0332a1) feat(lang): Change user (#153)
 * [0804dc7](https://github.com/tensorchord/envd/commit/0804dc775871f184d4acb3469835cef28bf68a03) feat(lang): Support git credentials (#145)
 * [9d9d282](https://github.com/tensorchord/envd/commit/9d9d282640de68a4d9b68d4a1ac4c34ee5ac83c6) fix: Add discord (#146)
 * [0b187da](https://github.com/tensorchord/envd/commit/0b187da03900fb765ef9269bc736d498d63cdbaa) feat(CLI): Support multiple environments (#142)
 * [b35c287](https://github.com/tensorchord/envd/commit/b35c287858b471387f2b3f3479b1df6845b9ce12) feat(cli): beautify the image pulling progress in bootstrap (#144)
 * [e027a8b](https://github.com/tensorchord/envd/commit/e027a8bf44f2029970c195459896b307c4a23bc4) feat(CLI): Support build context (#138)
 * [62d5049](https://github.com/tensorchord/envd/commit/62d504945f3ac9ebc03565442f35481cc0194df2) fix: Fix progress display reorder problem (#139)
 * [7ef5c71](https://github.com/tensorchord/envd/commit/7ef5c7101d6c67a4fbf2ef79459ea13649eeac9e) chore(make) golangci-lint install (#141)
 * [567673a](https://github.com/tensorchord/envd/commit/567673ab5e58411cf0cadd3c6a62dce03f440d64) chore(license) copyright (#137)

### Contributors

 * Ce Gao
 * Jian Zeng
 * Jinjing Zhou
 * Keming

## v0.0.1-alpha.4 (2022-05-11)

 * [9795d2b](https://github.com/tensorchord/envd/commit/9795d2b934a02b46d0fad8b3979d3fc62932ff4c) chore(readme) fix name (#136)
 * [8aa2ea8](https://github.com/tensorchord/envd/commit/8aa2ea8a9d50e615d78c5b6506f53b6cf8f24402) fix: Add detailed output for compile stage (#131)
 * [72e6dd5](https://github.com/tensorchord/envd/commit/72e6dd52bbfbdd228e522edf04a08ae65af22f39) chore(make,readme) minor fix (#132)
 * [20054a7](https://github.com/tensorchord/envd/commit/20054a706e14aa981ddf727c9ada9319d4bb9c42) feat(CLI): Replace .midi with XDG base dir (#129)
 * [d6aec6f](https://github.com/tensorchord/envd/commit/d6aec6f27fb66f4d7b21efe911b2c85cbeb31c44) fix(lang): Fix vscode plugin parse logic (#123)
 * [22f8fb9](https://github.com/tensorchord/envd/commit/22f8fb9419769b896d28ac858798337ef8211c68) chore(workflow): Run linter in CI (#126)
 * [6d505a2](https://github.com/tensorchord/envd/commit/6d505a2cb9c2fb5f3fb20b41f8638c5853abff43) feat(lang): Support command exec (#119)
 * [b95aee9](https://github.com/tensorchord/envd/commit/b95aee9d56a00558a65ec348d828e4040a9329b7) chore(ssh): Remove hard coded keys (#120)
 * [870cebd](https://github.com/tensorchord/envd/commit/870cebde01a3a2b423c23020b9b9dfa2fedbd028) fix(docker): Negotiate api version with server (#122)
 * [da4d1a0](https://github.com/tensorchord/envd/commit/da4d1a09140fe973665676b48df652d98fc0f761) feat: Add MNIST example (#115)
 * [7d15da4](https://github.com/tensorchord/envd/commit/7d15da431fee8402339d35092a20178a54c6d631) support runtime mount directory (#118)
 * [aeb846e](https://github.com/tensorchord/envd/commit/aeb846ec0ac30f1614e2273f6dc1fdddfe818daa) feat(docker): Mount base dir into the container (#113)
 * [3c28ffb](https://github.com/tensorchord/envd/commit/3c28ffbab95e3106a124d87b71293435daf8faaf) feat(CLI): Support destroy subcommand (#110)
 * [b8fbce1](https://github.com/tensorchord/envd/commit/b8fbce10d4f46ba521e1613cfe04a6883d7d1562) feat(lang): Support jupyter (#107)
 * [d3256c6](https://github.com/tensorchord/envd/commit/d3256c64b8c20071aa6676b9cd169e7f8cf421fa) feat(lang): Support zsh (#85)
 * [721e0fb](https://github.com/tensorchord/envd/commit/721e0fb28ff7a75c552557bb043e83e6102de187) fix: Fix name (#104)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming

## v0.0.1-alpha.3 (2022-05-05)

 * [5d39458](https://github.com/tensorchord/envd/commit/5d39458f15833277afdca0d9b2271e57d279a47d) fix: Fix template (#103)
 * [b5773f5](https://github.com/tensorchord/envd/commit/b5773f5623f92545ecc5e573c88a2914a61f1df6) fix(workflow): Add docker login (#102)

### Contributors

 * Ce Gao

## v0.0.1-alpha.2 (2022-05-03)

 * [b9750ba](https://github.com/tensorchord/envd/commit/b9750badf718b70842be8779f1c3a26fdf817283) fix: typo (#100)
 * [d847350](https://github.com/tensorchord/envd/commit/d8473507c962b47a04c061ba33e3b0789893a1c4) chore(lint): make linter happy (#97)
 * [3d3677e](https://github.com/tensorchord/envd/commit/3d3677e27f5429099199f8c0d8fda4b3c1fbfd84) fix: cannot bootstrap buildkitd if the container does not exist (#87)
 * [62ce39e](https://github.com/tensorchord/envd/commit/62ce39ec6c3346f6f9f840b0414640e0f31f27e7) chore(buildkit): Add test cases for builder (#83)
 * [9534253](https://github.com/tensorchord/envd/commit/953425348896fa6b6a8304fced06948442b25323) fix: Add a waiting before ssh (#82)
 * [ac81b41](https://github.com/tensorchord/envd/commit/ac81b419e0e8f4dd8f78f29c4d8e56fc146fae6e) feat: Use docker to keep ssh (#80)
 * [44dc07a](https://github.com/tensorchord/envd/commit/44dc07a3c6b2730048b673dbb3e77209815157c0) fix: Fix no such file bug (#81)
 * [08aa611](https://github.com/tensorchord/envd/commit/08aa611da89f42e57af09d3d186163984a819836) feat(CLI): Download midi-ssh instead of the built-in one (#73)
 * [2ad4d94](https://github.com/tensorchord/envd/commit/2ad4d949955cde550d6fe6c71bb078bd18c5432c) fix: Use docker-container as the buildkitd socket (#72)
 * [758bc0f](https://github.com/tensorchord/envd/commit/758bc0ff6b0db04593ca9fafb2fbb78805bc57d2) chore(error): Replace pkg/errors with crdb/errors (#71)
 * [c4e58ae](https://github.com/tensorchord/envd/commit/c4e58ae0aaba59b97874a94b1c1bb52af7610c07) refactor: Support connection retry (#70)

### Contributors

 * Ce Gao
 * Keming
 * Siyuan Wang

## v0.0.1-alpha.1 (2022-04-26)


### Contributors


