# Changelog

## v0.3.10 (2023-01-25)

 * [efdd05f](https://github.com/tensorchord/envd/commit/efdd05f5219c92a5cbae14a57fb55199535b9cd4) fix: rm remote cache for v1 (#1447)
 * [f940354](https://github.com/tensorchord/envd/commit/f9403540e1efe08b8512a728f9cd0eb3232dfb76) feat: add shm size to start options (#1445)
 * [b0320cb](https://github.com/tensorchord/envd/commit/b0320cb2b5f82bc164210c2995caa4ca8dbf50a2) chore(deps): bump github.com/onsi/gomega from 1.24.2 to 1.25.0 (#1442)
 * [0a2d9af](https://github.com/tensorchord/envd/commit/0a2d9af669c2b6d8128fc1b83087f05e88b4d70e) chore(deps): bump golang.org/x/crypto from 0.2.0 to 0.5.0 (#1443)
 * [34c3150](https://github.com/tensorchord/envd/commit/34c3150db5e23dcc25b3270c396c0b1e62a0ab74) chore(deps): bump pypa/cibuildwheel from 2.11.4 to 2.12.0 (#1438)
 * [83c63c7](https://github.com/tensorchord/envd/commit/83c63c7dbce98cd2db53fc317a5366d78ecb5bfd) fix: pip install requirements.txt (#1434)
 * [f47f7a3](https://github.com/tensorchord/envd/commit/f47f7a353cd84dc03f577116bac33090ab8f0eac) doc: add init and daemon debug guide (#1435)
 * [b6d49c6](https://github.com/tensorchord/envd/commit/b6d49c64290e32b4c92be009f40371b296899573) feat: add `make` as a default system package (#1433)
 * [686d029](https://github.com/tensorchord/envd/commit/686d0295b36b4b8227c65b8e8aee60bd062006f6) fix: Add docker auth (#1432)
 * [2d03de9](https://github.com/tensorchord/envd/commit/2d03de9872610b1177f63b613c0f38d36a9eed19) feat: :sparkles: init nerdctl support (#1378)

### Contributors

 * Ce Gao
 * Keming
 * Wei Zhang
 * dependabot[bot]

## v0.3.9 (2023-01-18)


### Contributors


## v0.3.8 (2023-01-20)

 * [2d03de9](https://github.com/tensorchord/envd/commit/2d03de9872610b1177f63b613c0f38d36a9eed19) feat: :sparkles: init nerdctl support (#1378)

### Contributors

 * Wei Zhang

## v0.3.7 (2023-01-18)

 * [0f556d0](https://github.com/tensorchord/envd/commit/0f556d0b630df70c17b98f5444ba3dcd67abd0bd) fix: change default channel to `defaults` (#1427)
 * [949f188](https://github.com/tensorchord/envd/commit/949f188a074fffef300cc4489806317357931d0b) LLM inference example (#1425)
 * [46f5996](https://github.com/tensorchord/envd/commit/46f59963dfe10c369b21320f754b20721cd085c1) chore(deps): bump github.com/tensorchord/envd-server from 0.0.23 to 0.0.24 (#1388)
 * [b93cd48](https://github.com/tensorchord/envd/commit/b93cd48b1e87642677d1a34120427b7b28192d90) feat: add checksum check for Julia archive file (#1419)
 * [c87846e](https://github.com/tensorchord/envd/commit/c87846e68f2eb0220e0c8698e94894c82fb04491) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.6.1 to 2.7.0 (#1418)
 * [64b5a0a](https://github.com/tensorchord/envd/commit/64b5a0ab5f7799ba3210bf0afa1c5521d078d963) Optimize the logic of adding entry in wsl (#1414)
 * [e1bbd56](https://github.com/tensorchord/envd/commit/e1bbd565b266e785ebb05b0e6227161f12bb8b14) feat: support trust the pypi index in v1, add api doc (#1412)
 * [d1179a6](https://github.com/tensorchord/envd/commit/d1179a6cbaece9d5627fa76dd8a1d464cb253f0e) fix: v1 gid on darwin (#1411)
 * [c0b88a9](https://github.com/tensorchord/envd/commit/c0b88a969bebde216d4cc63cbafd46772b258749) fix: create a legal container env name from cur dir (#1406)
 * [47b4dba](https://github.com/tensorchord/envd/commit/47b4dba83d036d3669afb9c44ed8a9f9f5028dd1) fix: julia v1 test (#1408)

### Contributors

 * Jinjing Zhou
 * Keming
 * dependabot[bot]
 * li mengyang
 * x0oo0x

## v0.3.6 (2023-01-11)

 * [52d94bb](https://github.com/tensorchord/envd/commit/52d94bb686f88111b678c8b5e39119780f190986) Intial e2e test for julia environment (#1401)
 * [7916dbb](https://github.com/tensorchord/envd/commit/7916dbb7f83fab6b4f4f2233ddbc7072234e92bb) Support Julia environment and Julia package installation (#1400)
 * [e7cabe6](https://github.com/tensorchord/envd/commit/e7cabe6e01615ba49d18acba3884c4371d22e782) feat(lang/py): add trust for pip_index (#1395)
 * [c269b89](https://github.com/tensorchord/envd/commit/c269b8966325734dde3edbb3f9f0b93c759609dc) fix: mv rlang test to doc_extra test (#1393)
 * [b4c08c8](https://github.com/tensorchord/envd/commit/b4c08c835b427a426d20681c4e6093e4ce3c786d) Initial e2e test for R language environment (#1342)
 * [3f7d699](https://github.com/tensorchord/envd/commit/3f7d699c9dfc36cb85433a8796d745774a91f1db) Refactor R package installation (#1381)
 * [aab7f53](https://github.com/tensorchord/envd/commit/aab7f53d2cf57549598c73040d98165fb2ecf5cb) feat: `envd exec` get runtime graph from container labels (#1392)
 * [ede2013](https://github.com/tensorchord/envd/commit/ede201335b95d5071224f9aeeb03d81830472582) chore(deps): bump golang.org/x/term from 0.3.0 to 0.4.0 (#1389)
 * [7034c5d](https://github.com/tensorchord/envd/commit/7034c5dade81293572fdb604d0c323461a13772b) Revert "fix: increase the default buildkit cache limit" (#1386)
 * [f56d95a](https://github.com/tensorchord/envd/commit/f56d95a2eb66bef6044d5653983f714cd0df6da3) fix: increase the default buildkit cache limit (#1382)
 * [3c38b5e](https://github.com/tensorchord/envd/commit/3c38b5e62753d51ec470586b44ada1b14af6eb64) fix: v0 user passwd, add test for root permission (#1379)
 * [51b586e](https://github.com/tensorchord/envd/commit/51b586e4a143fb90562e81c4271d6abdb922bce6) fix: add conda&python path for non-conda mode (#1376)
 * [754b6c3](https://github.com/tensorchord/envd/commit/754b6c3e773437ee59b7eadfdcd850c8610ef9eb) fix: change to dynamic PATH (#1375)
 * [a218dbc](https://github.com/tensorchord/envd/commit/a218dbc2ea98fc50599b5260daa73234be424334) fix(CLI): Remove image list in destroy (#1373)

### Contributors

 * Ce Gao
 * Keming
 * Weixiao Huang
 * dependabot[bot]
 * x0oo0x

## v0.3.5 (2023-01-03)

 * [a18c880](https://github.com/tensorchord/envd/commit/a18c8801961f7aa44a8f534ee1323b3588fd8be3) fix: Update rstudio server (#1370)

### Contributors

 * Ce Gao

## v0.3.4 (2023-01-03)

 * [38990c3](https://github.com/tensorchord/envd/commit/38990c31baece96fba182f4913686332de95e309) fix: Remove libssl1.1 (#1368)

### Contributors

 * Ce Gao

## v0.3.3 (2023-01-03)


### Contributors


## v0.3.2 (2023-01-03)

 * [095237d](https://github.com/tensorchord/envd/commit/095237dee602eb76dedf9041370173eee367d467) fix: fix nil pointer dereference when exiting shell (#1364)
 * [aae5e27](https://github.com/tensorchord/envd/commit/aae5e27feaa06436fa81c3dc414b28b6b7e21a3b) chore(deps): bump github.com/mattn/go-isatty from 0.0.16 to 0.0.17 (#1361)
 * [b0b5d71](https://github.com/tensorchord/envd/commit/b0b5d71ccb463a7e26e6f3f0cf2ce5d094731367) fix: remove duplicated config and fix Dockerfile pattern (#1362)
 * [bb0fe7c](https://github.com/tensorchord/envd/commit/bb0fe7cac7d6194f2fe095a5aa61bf68a96a564d) Fix: Remove unnecessary user switching code when installing R package (#1351)
 * [4f406ce](https://github.com/tensorchord/envd/commit/4f406ce057c36e1c45ee1e5d2b4c27a8120796ca) fix: fix broken link (#1349)
 * [3093e1a](https://github.com/tensorchord/envd/commit/3093e1a4c64ff1cbff58f656124951d561445a73) fix: v1 change default user group to 1000:1000 (#1347)
 * [571abfe](https://github.com/tensorchord/envd/commit/571abfede9e2166ba5fbefaaf046a4ff926309ad) feat: add --format for commands to output json (#1323)
 * [296680c](https://github.com/tensorchord/envd/commit/296680ce6a02d5a0605936a49de250ba236baaa5) fix: install black[jupyter] for CI check (#1343)
 * [cb33a65](https://github.com/tensorchord/envd/commit/cb33a656801a2f6b6cd04f4b4631ff17b3a2d3a8) chore(deps): bump github.com/tensorchord/envd-server from 0.0.21 to 0.0.23 (#1341)
 * [722fca1](https://github.com/tensorchord/envd/commit/722fca1d11a5519a79f6be29ca48ef67a2414b55) chore(deps): bump pypa/cibuildwheel from 2.11.3 to 2.11.4 (#1340)
 * [7426973](https://github.com/tensorchord/envd/commit/7426973cc33dd56a3967a4c2163575a92359a963) doc: add pytorch2 example (#1339)
 * [e1cb495](https://github.com/tensorchord/envd/commit/e1cb4953d008eb3ed0816e667a59595da9204a3c) fix: doc url in envd github issue page (#1338)
 * [122021d](https://github.com/tensorchord/envd/commit/122021da173e62ddbada683ddba4dfd8eb0432cc) refact: change Shlex(fmt.Sprintf to Shlexf, use ` instead of escape \" (#1337)

### Contributors

 * Bingtan Lu
 * Keming
 * Zhizhen He
 * cutecutecat
 * dependabot[bot]
 * x0oo0x
 * xing0821
 * xxchan

## v0.3.1 (2022-12-23)

 * [8cc3cc4](https://github.com/tensorchord/envd/commit/8cc3cc4064449025b0aeb764c80568796df79a32) feat: Support installing R environment in DEB822 format using internal standardized structure (#1329)
 * [49b3236](https://github.com/tensorchord/envd/commit/49b323609c52d6550b150e5cb44dc941d480d0d0) feat(login): Refactor logic with the new key API (#1333)
 * [1cd917d](https://github.com/tensorchord/envd/commit/1cd917dbe3b1fede67cfe9dd216e86622790b007) fix: upgrade horust to support arm (#1330)
 * [d39adc8](https://github.com/tensorchord/envd/commit/d39adc851a8e967e56221f4ac09775fc09f686af) refactor: combine code with the same logic (#1326)
 * [d1990ec](https://github.com/tensorchord/envd/commit/d1990ec5b8677d9c588919c887c5a304b01b3703) chore(deps): bump github.com/onsi/gomega from 1.24.1 to 1.24.2 (#1320)
 * [c598834](https://github.com/tensorchord/envd/commit/c598834ee34abc8675e19d2fb31bacbcdaa1acdb) fix: conda meta permission for envd user (#1325)
 * [0a48acd](https://github.com/tensorchord/envd/commit/0a48acd355bfe6c62380f7729991993063ca2991) feat: produce ABI-agnostic wheels for python (#1324)
 * [8042997](https://github.com/tensorchord/envd/commit/80429972639cf57d8127a7d5bd8d71ad41834d89) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.6.0 to 2.6.1 (#1321)

### Contributors

 * Ce Gao
 * Frost Ming
 * Keming
 * dependabot[bot]
 * rrain7
 * x0oo0x

## v0.3.0 (2022-12-16)

 * [30ff84d](https://github.com/tensorchord/envd/commit/30ff84d71155bd5bcf6d09df2d1d8effd95b130e) fix: v1 gpu base image name and envs (#1313)

### Contributors

 * Keming

## v0.2.6 (2022-12-15)

 * [31eec50](https://github.com/tensorchord/envd/commit/31eec501f898e35ca17f3a8399eaa99e3a2e9705) fix: v0 bashrc file (#1307)
 * [1fec54e](https://github.com/tensorchord/envd/commit/1fec54e3ecbec248bcdd9b8d9f0bb484f270bf3a) doc: v1 (#1284)
 * [47361c9](https://github.com/tensorchord/envd/commit/47361c99499f38550f13334dad004ad9f933d9ec) chore(vendor): Update envd server (#1305)
 * [d90dbf8](https://github.com/tensorchord/envd/commit/d90dbf8c5f246fe7f0eac27dc8338ccde797457b) fix: Add more error info in env ls (#1301)
 * [fc29fdd](https://github.com/tensorchord/envd/commit/fc29fdd5602fdcf6132bb58df50fad1ab3ecba94) feat: add option to output the completion script to terminal (#1303)
 * [1cc8c48](https://github.com/tensorchord/envd/commit/1cc8c4899245b0a35f29fc7351e7406f68bcc5bf) fix: flatten pypi packages in the label (#1302)
 * [0a416d1](https://github.com/tensorchord/envd/commit/0a416d1a5ad2fef5abecc44beac526ea33511540) doc: add more installation methods (#1300)
 * [d0a803e](https://github.com/tensorchord/envd/commit/d0a803e99c8cb346ef1b68724ab6e6f241390f18) fix: init uid/gid in Graph (#1299)
 * [db0396d](https://github.com/tensorchord/envd/commit/db0396d410f6f1e05f215269a63455bca73b3a0b) fix: make starship prompt pure text (#1296)
 * [a24c4cb](https://github.com/tensorchord/envd/commit/a24c4cba6c11960bce2c174d61df22be150a3ccd) fix: Update svg and simplify readme (#1297)
 * [935e55d](https://github.com/tensorchord/envd/commit/935e55d5907378ab967d4c767bff7ec3b0c0909a) Interactive cli (#1294)
 * [3800ea2](https://github.com/tensorchord/envd/commit/3800ea2971ea260431fd8fc5ab22ad9b069d8251) fix: add PATH when switch to envd user (#1295)
 * [382a6c1](https://github.com/tensorchord/envd/commit/382a6c150e30ea883c3229b37d851d85ead737bc) feat: install pypi packages in different groups (#1289)
 * [c66f201](https://github.com/tensorchord/envd/commit/c66f201c90a3ffee039a13679c41e792a1ebeac5) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.5.1 to 2.6.0 (#1293)
 * [d15c3ca](https://github.com/tensorchord/envd/commit/d15c3caf450009c4e94172b45474d8e268ac76e8) chore(deps): bump golang.org/x/term from 0.2.0 to 0.3.0 (#1292)
 * [50a8bb6](https://github.com/tensorchord/envd/commit/50a8bb620c178adaa80ce59edea364759afbd5bc) chore(deps): bump github.com/urfave/cli/v2 from 2.23.6 to 2.23.7 (#1291)
 * [99e8a15](https://github.com/tensorchord/envd/commit/99e8a15fd240d83810e2f4414bcb66d3addcb355) chore(deps): bump pypa/cibuildwheel from 2.11.2 to 2.11.3 (#1290)
 * [f7f2b34](https://github.com/tensorchord/envd/commit/f7f2b34a80f616ae4db1c0fdb1c01f6583f4821f) Add config.owner to v1 syntax for setting UID and GID (#1286)
 * [5c4aae9](https://github.com/tensorchord/envd/commit/5c4aae93f797a6ce504f3d09f6395e5492627703) fix: Remove graph access in ssh (#1281)
 * [56d25a0](https://github.com/tensorchord/envd/commit/56d25a01dd37b2c069ae9b4492811c8571a1e714) feat: Support specify resource requirement when creating the environment (#1268)
 * [4f82290](https://github.com/tensorchord/envd/commit/4f8229033995cb5919f08c6889515927942b96ff) feat: support no default language in v1 (#1278)
 * [8e9cb2e](https://github.com/tensorchord/envd/commit/8e9cb2e7bcbaef7c86996d7376ca27af8e741ed0) feat: support serving (#1228)
 * [8651444](https://github.com/tensorchord/envd/commit/865144482b380cf5d3d6d30dff55a45087aeabc4) Remove incorrect `sudo` command when installing system package (#1273)
 * [235dde9](https://github.com/tensorchord/envd/commit/235dde905643e2eed7ba91e10423b92598d6a12a) fix: Remove ir direct access in ssh (#1271)
 * [06b4880](https://github.com/tensorchord/envd/commit/06b488035a5c21a74c8727fc7d8e85f404a2bb00) feat: Add JWT and user/pwd in envd server context (#1243)
 * [7badc67](https://github.com/tensorchord/envd/commit/7badc67fcf9be813d2124aed7f18acd855847e63) chore(deps): bump github.com/tensorchord/envd-server from 0.0.11 to 0.0.12 (#1266)
 * [cefd1df](https://github.com/tensorchord/envd/commit/cefd1dfa8101a363a762921ec7439f64d1463e2a) chore(deps): bump github.com/urfave/cli/v2 from 2.23.5 to 2.23.6 (#1264)

### Contributors

 * Alex Xi
 * Ce Gao
 * Frost Ming
 * Jinjing Zhou
 * Keming
 * dependabot[bot]
 * nullday
 * x0oo0x

## v0.2.5-rc.7 (2022-12-04)


### Contributors


## v0.2.5-rc.6 (2022-12-04)

 * [8594d10](https://github.com/tensorchord/envd/commit/8594d1044874ed2baa3bc74ef9b5e35f3766a6f7) fix: Update python setuptools to python3 (#1261)

### Contributors

 * Ce Gao

## v0.2.5-rc.5 (2022-12-02)

 * [0042b24](https://github.com/tensorchord/envd/commit/0042b2481da98c6810e858069f769f3b7ef4808a) feat: add completion command (#1258)

### Contributors

 * tison

## v0.2.5-rc.4 (2022-12-02)


### Contributors


## v0.2.5-rc.3 (2022-12-02)


### Contributors


## v0.2.5-rc.2 (2022-12-02)

 * [5cbed88](https://github.com/tensorchord/envd/commit/5cbed887e0aa6b0c766f48c0cf8effe85f49d20b) fix: put the binary under bin directly (#1254)
 * [c7fd12b](https://github.com/tensorchord/envd/commit/c7fd12b9047b5c59e555a3d083cdf7cdef1449d6) bug: Update envd-server version to 0.0.11 (#1245)

### Contributors

 * Frost Ming
 * Jinjing Zhou

## v0.2.5-rc.1 (2022-11-28)

 * [a4dad9b](https://github.com/tensorchord/envd/commit/a4dad9b212c41e07175f201e067f891c76277843) chore(deps): bump github.com/tensorchord/envd-server from 0.0.9 to 0.0.10 (#1235)

### Contributors

 * dependabot[bot]

## v0.2.5-alpha.8 (2022-11-23)

 * [2f37901](https://github.com/tensorchord/envd/commit/2f379016002e58c3a87b2a8b2e2afd6ffc6ab9b9) feat: support add new dir to runtime PATH (#1218)
 * [e85f062](https://github.com/tensorchord/envd/commit/e85f0620ad76379e22fdf36f70697b1f8d71444b) feat: merge coverage file and report to goverall (#1211)
 * [7f10440](https://github.com/tensorchord/envd/commit/7f10440a8163f4d23964783fdf83c2e03c66e4f9) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.5.0 to 2.5.1 (#1213)
 * [6db3f15](https://github.com/tensorchord/envd/commit/6db3f152325b674ea67e24fc4ee694f99ea206c3) chore(deps): bump github.com/tensorchord/envd-server from 0.0.8 to 0.0.9 (#1212)
 * [5efb0b1](https://github.com/tensorchord/envd/commit/5efb0b12b98e3000d0ddc6ef015734e65411ac80) fix: only mkdir /home/envd/* when image==nil (#1208)
 * [8c8c3f3](https://github.com/tensorchord/envd/commit/8c8c3f33ada10013b350071a1a36529d993053c4) feat: add name to `envd envs describe` (#1201)

### Contributors

 * Alex Xi
 * Keming
 * cutecutecat
 * dependabot[bot]

## v0.2.5-alpha.7 (2022-11-15)

 * [3fa3070](https://github.com/tensorchord/envd/commit/3fa307098a585b524804bd0beff6697e68440b42) feat: support build owner (#1202)

### Contributors

 * Keming

## v0.2.5-alpha.6 (2022-11-14)

 * [5462d30](https://github.com/tensorchord/envd/commit/5462d30a556f91741b770e4b4b24c1fd639f5299) fix: release CI add runneradmin to docker group (#1199)

### Contributors

 * Keming

## v0.2.5-alpha.5 (2022-11-14)

 * [731b6f8](https://github.com/tensorchord/envd/commit/731b6f8271a77096d9a29a03e353e0388f825485) fix: github action release cache command (#1197)

### Contributors

 * Keming

## v0.2.5-alpha.4 (2022-11-14)

 * [4a7d6cc](https://github.com/tensorchord/envd/commit/4a7d6cca18d9a033c02c1ef235570f64ef987300) fix: daemon command (#1185)

### Contributors

 * Keming

## v0.2.5-alpha.3 (2022-11-14)

 * [212f996](https://github.com/tensorchord/envd/commit/212f9966d12f2322fe8c97d14db093448399119e) chore(deps): bump github.com/urfave/cli/v2 from 2.23.4 to 2.23.5 (#1188)
 * [c4854e8](https://github.com/tensorchord/envd/commit/c4854e842a7cee0f2600f5c6eb8d19ca3510369d) chore(deps): bump github.com/onsi/gomega from 1.24.0 to 1.24.1 (#1189)
 * [62b7b2f](https://github.com/tensorchord/envd/commit/62b7b2f204aa5b53688a2b0356d8a746a5fac1e4) chore(deps): bump github.com/moby/buildkit from 0.10.5 to 0.10.6 (#1190)
 * [b1eb44e](https://github.com/tensorchord/envd/commit/b1eb44e4cc80b3422122429d3f29e5d7ffcaf3d6) fix: pre-create the workdir (#1174)
 * [66d9cb0](https://github.com/tensorchord/envd/commit/66d9cb061a033220814813a349fc12495927aef2) fix: Move user to image config (#1173)
 * [dd9dd30](https://github.com/tensorchord/envd/commit/dd9dd3083b363d1d698060255892e6a132b8182a) refact: move env to server (#1172)
 * [2d37b48](https://github.com/tensorchord/envd/commit/2d37b48736c5ba34ea0800c925c31c0c3ff40f13) fix stable diffusion demo via b/c huggingface api change (#1170)
 * [c364958](https://github.com/tensorchord/envd/commit/c364958c1d21d25a7cb782413f6666bd666f7dd1) feat: support envd-server image (#1150)
 * [f9e7c67](https://github.com/tensorchord/envd/commit/f9e7c67a0f16c802fa2aa75a562f77e6f17baef9) chore(deps): bump github.com/spf13/viper from 1.13.0 to 1.14.0 (#1164)
 * [5c9ac69](https://github.com/tensorchord/envd/commit/5c9ac69a584463cf906b27d0ccd200403d7f4974) chore(deps): bump github.com/onsi/gomega from 1.23.0 to 1.24.0 (#1163)
 * [d1de890](https://github.com/tensorchord/envd/commit/d1de8906326ee10a1021dbff80892d8f0c74e637) chore(deps): bump github.com/urfave/cli/v2 from 2.23.0 to 2.23.4 (#1162)
 * [ad5ce11](https://github.com/tensorchord/envd/commit/ad5ce117bdc4433f8894c4b813748d33c6edbd96) chore(deps): bump dependabot/fetch-metadata from 1.3.4 to 1.3.5 (#1161)
 * [6b7a0e9](https://github.com/tensorchord/envd/commit/6b7a0e94937d5e0532139d519e424ccb16bd9de1) fix: cuda tag in image cache (#1156)
 * [1af3449](https://github.com/tensorchord/envd/commit/1af3449ce72e218350c3d4fa19aee64f98ced970) fix: pre mkdir for runtime mount (#1153)

### Contributors

 * Ce Gao
 * Keming
 * dependabot[bot]
 * xieydd

## v0.2.5-alpha.2 (2022-11-03)

 * [46d24fd](https://github.com/tensorchord/envd/commit/46d24fd1331b1b0e223cb2d2223e17d648080bf1) fix: horust cache (#1151)

### Contributors

 * Keming

## v0.2.5-alpha.1 (2022-11-03)


### Contributors


## v0.2.5 (2022-12-06)

 * [235dde9](https://github.com/tensorchord/envd/commit/235dde905643e2eed7ba91e10423b92598d6a12a) fix: Remove ir direct access in ssh (#1271)
 * [06b4880](https://github.com/tensorchord/envd/commit/06b488035a5c21a74c8727fc7d8e85f404a2bb00) feat: Add JWT and user/pwd in envd server context (#1243)
 * [7badc67](https://github.com/tensorchord/envd/commit/7badc67fcf9be813d2124aed7f18acd855847e63) chore(deps): bump github.com/tensorchord/envd-server from 0.0.11 to 0.0.12 (#1266)
 * [cefd1df](https://github.com/tensorchord/envd/commit/cefd1dfa8101a363a762921ec7439f64d1463e2a) chore(deps): bump github.com/urfave/cli/v2 from 2.23.5 to 2.23.6 (#1264)
 * [8594d10](https://github.com/tensorchord/envd/commit/8594d1044874ed2baa3bc74ef9b5e35f3766a6f7) fix: Update python setuptools to python3 (#1261)
 * [0042b24](https://github.com/tensorchord/envd/commit/0042b2481da98c6810e858069f769f3b7ef4808a) feat: add completion command (#1258)
 * [5cbed88](https://github.com/tensorchord/envd/commit/5cbed887e0aa6b0c766f48c0cf8effe85f49d20b) fix: put the binary under bin directly (#1254)
 * [c7fd12b](https://github.com/tensorchord/envd/commit/c7fd12b9047b5c59e555a3d083cdf7cdef1449d6) bug: Update envd-server version to 0.0.11 (#1245)
 * [a4dad9b](https://github.com/tensorchord/envd/commit/a4dad9b212c41e07175f201e067f891c76277843) chore(deps): bump github.com/tensorchord/envd-server from 0.0.9 to 0.0.10 (#1235)
 * [2f37901](https://github.com/tensorchord/envd/commit/2f379016002e58c3a87b2a8b2e2afd6ffc6ab9b9) feat: support add new dir to runtime PATH (#1218)
 * [e85f062](https://github.com/tensorchord/envd/commit/e85f0620ad76379e22fdf36f70697b1f8d71444b) feat: merge coverage file and report to goverall (#1211)
 * [7f10440](https://github.com/tensorchord/envd/commit/7f10440a8163f4d23964783fdf83c2e03c66e4f9) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.5.0 to 2.5.1 (#1213)
 * [6db3f15](https://github.com/tensorchord/envd/commit/6db3f152325b674ea67e24fc4ee694f99ea206c3) chore(deps): bump github.com/tensorchord/envd-server from 0.0.8 to 0.0.9 (#1212)
 * [5efb0b1](https://github.com/tensorchord/envd/commit/5efb0b12b98e3000d0ddc6ef015734e65411ac80) fix: only mkdir /home/envd/* when image==nil (#1208)
 * [8c8c3f3](https://github.com/tensorchord/envd/commit/8c8c3f33ada10013b350071a1a36529d993053c4) feat: add name to `envd envs describe` (#1201)
 * [3fa3070](https://github.com/tensorchord/envd/commit/3fa307098a585b524804bd0beff6697e68440b42) feat: support build owner (#1202)
 * [5462d30](https://github.com/tensorchord/envd/commit/5462d30a556f91741b770e4b4b24c1fd639f5299) fix: release CI add runneradmin to docker group (#1199)
 * [731b6f8](https://github.com/tensorchord/envd/commit/731b6f8271a77096d9a29a03e353e0388f825485) fix: github action release cache command (#1197)
 * [4a7d6cc](https://github.com/tensorchord/envd/commit/4a7d6cca18d9a033c02c1ef235570f64ef987300) fix: daemon command (#1185)
 * [212f996](https://github.com/tensorchord/envd/commit/212f9966d12f2322fe8c97d14db093448399119e) chore(deps): bump github.com/urfave/cli/v2 from 2.23.4 to 2.23.5 (#1188)
 * [c4854e8](https://github.com/tensorchord/envd/commit/c4854e842a7cee0f2600f5c6eb8d19ca3510369d) chore(deps): bump github.com/onsi/gomega from 1.24.0 to 1.24.1 (#1189)
 * [62b7b2f](https://github.com/tensorchord/envd/commit/62b7b2f204aa5b53688a2b0356d8a746a5fac1e4) chore(deps): bump github.com/moby/buildkit from 0.10.5 to 0.10.6 (#1190)
 * [b1eb44e](https://github.com/tensorchord/envd/commit/b1eb44e4cc80b3422122429d3f29e5d7ffcaf3d6) fix: pre-create the workdir (#1174)
 * [66d9cb0](https://github.com/tensorchord/envd/commit/66d9cb061a033220814813a349fc12495927aef2) fix: Move user to image config (#1173)
 * [dd9dd30](https://github.com/tensorchord/envd/commit/dd9dd3083b363d1d698060255892e6a132b8182a) refact: move env to server (#1172)
 * [2d37b48](https://github.com/tensorchord/envd/commit/2d37b48736c5ba34ea0800c925c31c0c3ff40f13) fix stable diffusion demo via b/c huggingface api change (#1170)
 * [c364958](https://github.com/tensorchord/envd/commit/c364958c1d21d25a7cb782413f6666bd666f7dd1) feat: support envd-server image (#1150)
 * [f9e7c67](https://github.com/tensorchord/envd/commit/f9e7c67a0f16c802fa2aa75a562f77e6f17baef9) chore(deps): bump github.com/spf13/viper from 1.13.0 to 1.14.0 (#1164)
 * [5c9ac69](https://github.com/tensorchord/envd/commit/5c9ac69a584463cf906b27d0ccd200403d7f4974) chore(deps): bump github.com/onsi/gomega from 1.23.0 to 1.24.0 (#1163)
 * [d1de890](https://github.com/tensorchord/envd/commit/d1de8906326ee10a1021dbff80892d8f0c74e637) chore(deps): bump github.com/urfave/cli/v2 from 2.23.0 to 2.23.4 (#1162)
 * [ad5ce11](https://github.com/tensorchord/envd/commit/ad5ce117bdc4433f8894c4b813748d33c6edbd96) chore(deps): bump dependabot/fetch-metadata from 1.3.4 to 1.3.5 (#1161)
 * [6b7a0e9](https://github.com/tensorchord/envd/commit/6b7a0e94937d5e0532139d519e424ccb16bd9de1) fix: cuda tag in image cache (#1156)
 * [1af3449](https://github.com/tensorchord/envd/commit/1af3449ce72e218350c3d4fa19aee64f98ced970) fix: pre mkdir for runtime mount (#1153)
 * [46d24fd](https://github.com/tensorchord/envd/commit/46d24fd1331b1b0e223cb2d2223e17d648080bf1) fix: horust cache (#1151)
 * [28d73f9](https://github.com/tensorchord/envd/commit/28d73f9c1633f74363adcb20c41519863e993f4b) feat: Fetch base image metadata (#1148)
 * [a2aa9fe](https://github.com/tensorchord/envd/commit/a2aa9fe6ecfd3ddeb1e8073b4768898e2952a812) feat: add environment name as label (#1135)
 * [bc1c513](https://github.com/tensorchord/envd/commit/bc1c51356b5e94ef28e818acc65aef87833ef705) Bug/add support for wsl2 (#1134)
 * [5053054](https://github.com/tensorchord/envd/commit/505305407f7375ed9ddf30014f72dee82f33e600) feat(cli): Add alias rm to envd image remove (#1136)
 * [7bf8c00](https://github.com/tensorchord/envd/commit/7bf8c00dfd8c34cc5552590b67b1e69eecaf5a51) feat: expose support listen addr (#1128)
 * [8776a25](https://github.com/tensorchord/envd/commit/8776a259de0b15d1a8bc7e76473b9a49e44e7890) fix: Fix zsh tab issue (#1129)
 * [cff88a9](https://github.com/tensorchord/envd/commit/cff88a9ac759eee48508eab65802f1adc77b9ecb) fix: Fix the order to make conda cache always work (#1126)
 * [07768b7](https://github.com/tensorchord/envd/commit/07768b7af8d9f39639a91ba156e3274889840c50) feat: record cmd duration in telemetry (#1122)
 * [0bdca3c](https://github.com/tensorchord/envd/commit/0bdca3c2d54d0eea8ee6b3a14641ae5563971b0e) fix: Remove cwd in conda to enable remote cache (#1125)
 * [535e26a](https://github.com/tensorchord/envd/commit/535e26a0a0d9d3855317d34dcf209e4d8d08efe9) feat(debug): Add llb export (#1124)
 * [d693254](https://github.com/tensorchord/envd/commit/d693254c1a8b19e02aca85934818a4e7b6dde661) fix: Disable telemetry in CI (#1116)
 * [3cba45e](https://github.com/tensorchord/envd/commit/3cba45e36496cd6010d9b07acbb0c7d34826f828) chore(deps): bump github.com/urfave/cli/v2 from 2.20.3 to 2.23.0 (#1119)
 * [57a63e0](https://github.com/tensorchord/envd/commit/57a63e015c813a0a446768b00c83039923ce0df3) chore(deps): bump pypa/cibuildwheel from 2.11.1 to 2.11.2 (#1115)
 * [77b3d84](https://github.com/tensorchord/envd/commit/77b3d84741f990e9cf0b3a8f9bd56abe40fcd895) chore(deps): bump github.com/onsi/gomega from 1.22.1 to 1.23.0 (#1117)

### Contributors

 * Alex Xi
 * Ce Gao
 * Frost Ming
 * Isaac
 * Jinjing Zhou
 * Keming
 * cutecutecat
 * dependabot[bot]
 * nullday
 * tison
 * xieydd

## v0.2.4-alpha.17 (2022-10-29)

 * [197b6c6](https://github.com/tensorchord/envd/commit/197b6c612891f0b7bc95b4b1f12a7daacbe7e51f) feat: Add telemetry with the help of segment.io (#1113)
 * [dedd731](https://github.com/tensorchord/envd/commit/dedd73113214d68fd6f8421446fcf2ea8895252e) feat: Add listening_addr to expose fun (#1110)
 * [f672c8f](https://github.com/tensorchord/envd/commit/f672c8f67c9fd793427fb9ec556e767fdf1ef50b) feat: support build time run without mount host (#1109)
 * [214f7c8](https://github.com/tensorchord/envd/commit/214f7c88c64c8b9152e0abb3acfb67e57cc1bf68) fix: panic if the user specify entrypoint for non-costom image (#1108)
 * [8f89ba6](https://github.com/tensorchord/envd/commit/8f89ba6c57795f76f38e0d613ab62517d7b03205) fix: auth with the same name (#1106)
 * [e831fe9](https://github.com/tensorchord/envd/commit/e831fe9c46a4f4085e95395722ffd836bd305c8e) fix: cycle import detection in the interpreter (#1104)
 * [356a707](https://github.com/tensorchord/envd/commit/356a707bfc7e4f9945868fbbcd94a03277d3e601) feat: add config.repo in envd (#1101)
 * [4429bae](https://github.com/tensorchord/envd/commit/4429baebf0bfb6a4f496609c90117881975b4054) fix(CLI): Fix format for app (#1100)
 * [0a975de](https://github.com/tensorchord/envd/commit/0a975de0103d84841252b54b549cddaa6d82d8a6) feat: support envd server destroy env (#1096)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming
 * Yilong Li
 * nullday

## v0.2.4-alpha.16 (2022-10-26)

 * [ac87c76](https://github.com/tensorchord/envd/commit/ac87c76f78e8d6ff8c156365310ed5740e8e054b) feature: support runtime init script (#1085)
 * [3416e70](https://github.com/tensorchord/envd/commit/3416e70ad5753441b109e1f6d755d9a02cd26e5c) feat: enable convenient vscode debugging (#1089)
 * [753fde3](https://github.com/tensorchord/envd/commit/753fde329e7da22837af27cddb48969947b9125b) fix:#1090 (#1092)
 * [1ddd513](https://github.com/tensorchord/envd/commit/1ddd513c03b78b8382d309c3a9fb162ae39973b9) feat(CLI): Support envd-server in up command (#1087)
 * [ab9fe2c](https://github.com/tensorchord/envd/commit/ab9fe2ce255ea18780310f37903b216853386589) feat: Support forward (#1014)
 * [4cd6fef](https://github.com/tensorchord/envd/commit/4cd6fef5d6284538ea523011aac47feec8dcf36a) chore(deps): bump github.com/moby/buildkit from 0.10.4 to 0.10.5 (#1081)
 * [a81bbb0](https://github.com/tensorchord/envd/commit/a81bbb0aab9b1a7d1bf37a0374884f34af386c3b) chore(deps): bump github.com/urfave/cli/v2 from 2.20.2 to 2.20.3 (#1078)
 * [a89f466](https://github.com/tensorchord/envd/commit/a89f466517aec054c0e0c5a295f09b5e7b3c33b7) chore(deps): bump github.com/opencontainers/image-spec from 1.1.0-rc1 to 1.1.0-rc2 (#1077)
 * [a209f79](https://github.com/tensorchord/envd/commit/a209f79d07b945c2a1a87972f191be7f4428d6c1) chore(deps): bump github.com/onsi/gomega from 1.21.1 to 1.22.1 (#1080)
 * [a8a6339](https://github.com/tensorchord/envd/commit/a8a633928a5ec70f5deccbebf81619bb1dbe33da) chore(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (#1079)

### Contributors

 * Ce Gao
 * Friends A
 * Jinjing Zhou
 * Yilong Li
 * dependabot[bot]

## v0.2.4-alpha.15 (2022-10-21)

 * [9faeebe](https://github.com/tensorchord/envd/commit/9faeebe27f47ce2ee1d99b04325eb41745390099) fix: Fix build.sh (#1072)
 * [ff30b3a](https://github.com/tensorchord/envd/commit/ff30b3a5a4b1ccdc089b601c1d4c2bfdd81ce61c) fix: Fix the file name typo (#1071)

### Contributors

 * Ce Gao

## v0.2.4-alpha.14 (2022-10-21)

 * [51b00fb](https://github.com/tensorchord/envd/commit/51b00fb189422b4791ef2c53cc80ac3a5067b67a) feat(context): Support unix context and daemonless (#1062)
 * [d41b674](https://github.com/tensorchord/envd/commit/d41b67430da64f19b67de4710ad3df1f56f9cc94) bug: Fix sdist didn't include go files (#1068)
 * [181ee71](https://github.com/tensorchord/envd/commit/181ee71a73f3723ca68cd26fd40f53ca96231429) chore(CLI): Move destroy logic from docker to envd engine. (#1050)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Tumushimire Yves

## v0.2.4-alpha.13 (2022-10-20)

 * [672bcd7](https://github.com/tensorchord/envd/commit/672bcd7c04f093f41d113749e39da3f193aea406) bug: fix install from pypi source release (#1064)

### Contributors

 * Jinjing Zhou

## v0.2.4-alpha.12 (2022-10-20)

 * [cd5cf52](https://github.com/tensorchord/envd/commit/cd5cf522e29ba54113a22867acc396e6fcd128e1) feat: use horust as the dev container supervisor (#1051)

### Contributors

 * Keming

## v0.2.4-alpha.11 (2022-10-20)

 * [044d89e](https://github.com/tensorchord/envd/commit/044d89e0d0901b3cf75d62c8505b2455b5b91dc5) feat: Update env client to support multiple envs (#1052)
 * [711dfb2](https://github.com/tensorchord/envd/commit/711dfb2ac9f08e467d315427b894191aa598e1d5) feat(lang): Add proposal for custom base image (#567)
 * [b601304](https://github.com/tensorchord/envd/commit/b601304a72bc586c78149a3dc9cec3bd1b1d1c35) chroe: make python package after apt package (#1048)
 * [cf78dd1](https://github.com/tensorchord/envd/commit/cf78dd131e869ad20fe67a4a1159d7bf69e6fb63) feat: Implement env client in envd engine (#1049)
 * [bc6255b](https://github.com/tensorchord/envd/commit/bc6255bc3d494be85afc7c132f2e48f72fffd0fa) example: Add torch profiler example (#1026)
 * [ded3fce](https://github.com/tensorchord/envd/commit/ded3fce1b481aab9c51c2f13dbb2d92288532da0) feat: add `envd down` as an alias for `envd destroy` (#1047)
 * [9b8582d](https://github.com/tensorchord/envd/commit/9b8582d910e83288ef90ba45b599221341ad29f6) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.2.0 to 2.3.1 (#1039)
 * [0ee19ed](https://github.com/tensorchord/envd/commit/0ee19ed79d3332a25082878ca19c1bc032072aef) chore(deps): bump pypa/cibuildwheel from 2.10.2 to 2.11.1 (#1036)
 * [6ba7633](https://github.com/tensorchord/envd/commit/6ba763368d81035b0f13e5eb1549856ecaaa30d1) chore(deps): bump k8s.io/api from 0.25.2 to 0.25.3 (#1038)
 * [0500dca](https://github.com/tensorchord/envd/commit/0500dca9b5174510c1c741ec1d9cfbf97d09a2e9) chore(deps): bump github.com/urfave/cli/v2 from 2.19.2 to 2.20.2 (#1037)
 * [847f187](https://github.com/tensorchord/envd/commit/847f187b45983a56093dd962d5638dae5f260743) chore(deps): bump github.com/onsi/gomega from 1.21.1 to 1.22.1 (#1040)
 * [1fb2721](https://github.com/tensorchord/envd/commit/1fb272197f405c5ee9dcca46c82079095afa406d) bug: Fix flaky auth test (#1034)
 * [d0134d3](https://github.com/tensorchord/envd/commit/d0134d301eba153dcda0c79a3a23babf75a14434) fix(zsh): ignore inserting zsh-completion if system don't have zsh shell (#1025)
 * [ec92a29](https://github.com/tensorchord/envd/commit/ec92a295b835664d559f201ba2542a6a609a560d) feat: support using current directory name as env name in envd describe (#1033)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Zhenzhen Zhao
 * dependabot[bot]
 * wangxiaolei

## v0.2.4-alpha.10 (2022-10-14)

 * [d3f93a8](https://github.com/tensorchord/envd/commit/d3f93a8389cbf198ce217e2001b7ae0c7df62d2b) bug: fix explicit channel setting with conda env.yaml (#1008)
 * [bcd8029](https://github.com/tensorchord/envd/commit/bcd802976597b3ef07906aa26682fc8ffa0e891a) bug: Fix install from source will lose version information (#1024)
 * [a58197c](https://github.com/tensorchord/envd/commit/a58197c476745bbe568e5712e54912936ca3a14e) fix: Add two debug logs (#1023)
 * [0fc4d8d](https://github.com/tensorchord/envd/commit/0fc4d8d12e9a13c00e760ee105765b0270a2c3ee) feat(context): Get ssh hostname from context, instead of hard-coded string (#1020)
 * [262a07b](https://github.com/tensorchord/envd/commit/262a07b69f976317379d42a8583b1d084c41e73b) feat(CLI): Support envd images prune #976 (#1012)
 * [88dc908](https://github.com/tensorchord/envd/commit/88dc908bc09b99b009c93a72999951096c984daa) fix(CLI): Fix a typo (#1015)
 * [931e1a0](https://github.com/tensorchord/envd/commit/931e1a09d0e3eda869d80e9b73e730c8c8768514) feat(cli): add `envd images remove` command (#1007)
 * [f26fc4c](https://github.com/tensorchord/envd/commit/f26fc4ce7e42ed0e567cb492c41551f6c2f2c3b4) feat(CLI): Support create command (#1001)
 * [f288ff3](https://github.com/tensorchord/envd/commit/f288ff3b8f2cdb2f63cbe313635f3596b3a2ace7) feat(CI): add nightly build & test tasks (#1006)
 * [46168da](https://github.com/tensorchord/envd/commit/46168dac70fb8cd9b841a32e35a414ad29ddfaa2) chore(test): Add test cases for pkg/home/auth.go (#1009)
 * [d87afda](https://github.com/tensorchord/envd/commit/d87afda9213fcb334ca55d247ba3c0663b7b180b) feat: Add envd-server runtime proposal (#303)
 * [fec7ada](https://github.com/tensorchord/envd/commit/fec7ada4efc2186fa5f858d5abe211a2beb2dfdd) fix: Remove hard code docker in envd engine init (#1000)
 * [cf56f0d](https://github.com/tensorchord/envd/commit/cf56f0de7bfb05be47d4704c14f9bbb892a86d28) chore(deps): bump github.com/onsi/gomega from 1.20.2 to 1.21.1 (#997)
 * [a55470e](https://github.com/tensorchord/envd/commit/a55470e363b35b689649596b449b996ada701630) feat(build): detect if the current environment is running before building (#892) (#989)
 * [d0219f2](https://github.com/tensorchord/envd/commit/d0219f2e3e999e27f334a39180031390a4554af9) chore(deps): bump github.com/urfave/cli/v2 from 2.17.1 to 2.19.2 (#998)

### Contributors

 * Ce Gao
 * JasonZhu
 * Jinjing Zhou
 * Tumushimire Yves
 * XRW
 * Yijiang Liu
 * Zhenzhen Zhao
 * dependabot[bot]

## v0.2.4-alpha.9 (2022-10-09)

 * [b2a9018](https://github.com/tensorchord/envd/commit/b2a90188c012fa186de2ecb7b1aa534062681020) feat: Support cli argument for host key (#992)
 * [91ea50c](https://github.com/tensorchord/envd/commit/91ea50cd74e36f2d0d966d2b6a1ddea0ad15c43b) bug: fix pypi package information (#959)

### Contributors

 * Ce Gao
 * Jinjing Zhou

## v0.2.4-alpha.8 (2022-10-08)

 * [0be9969](https://github.com/tensorchord/envd/commit/0be9969471f8358f1cc283215ddf8e4cfcecebad) fix: Remove dump checkout and remove pre-commit (#982)
 * [f4ebb02](https://github.com/tensorchord/envd/commit/f4ebb026c45ab29a84dfa57782ef5fec8c11cf14) Chore(test): Add test cases for pkg/util/starlarkutil (#979)
 * [06d3cfa](https://github.com/tensorchord/envd/commit/06d3cfa1e171a409893d59494754dc47fb886a52) *: upgrade golangci-linter and add noloopclosure (#980)
 * [17fedb8](https://github.com/tensorchord/envd/commit/17fedb8d2ae41318e2f705687c334b7e42413ab3) fix(ir): make sure default value won't be replaced with empty value (#970)
 * [c1ae887](https://github.com/tensorchord/envd/commit/c1ae887f12b23802363ac417c198e19983c5605f) feat(CLI): Support runner in context (#961)
 * [bf993e2](https://github.com/tensorchord/envd/commit/bf993e2cab44502abb74a79e732b74a5567b6194) refact: conda/mamba create/update env, fix user permissions (#933)
 * [0e79fb9](https://github.com/tensorchord/envd/commit/0e79fb91b9f3770f2e5356fc8a10182dd31e2bdf) chore(deps): bump github.com/urfave/cli/v2 from 2.16.3 to 2.17.1 (#968)
 * [c9045d2](https://github.com/tensorchord/envd/commit/c9045d24d412099a1ac0bf9b9acf78135157351d) chore(deps): bump dependabot/fetch-metadata from 1.3.3 to 1.3.4 (#967)

### Contributors

 * Ce Gao
 * Keming
 * Tumushimire Yves
 * Weizhen Wang
 * dependabot[bot]

## v0.2.4-alpha.7 (2022-10-01)

 * [c22408c](https://github.com/tensorchord/envd/commit/c22408c8f3b87ef959bda7d7681203ffd8d6212c) fix(ir): `apt install` and `conda env create` cache (#962)
 * [006f653](https://github.com/tensorchord/envd/commit/006f6538396b47b1bcab2a0d53ec9dae0221c8f5) feat: envd-sshd can read public key path from environment variable (#954)

### Contributors

 * Jinjing Zhou
 * Keming

## v0.2.4-alpha.6 (2022-09-29)

 * [0b73548](https://github.com/tensorchord/envd/commit/0b7354868943f9e4ea189be18eeea534959b5d6c) example: Use torch in mnist example (#927)

### Contributors

 * Jinjing Zhou

## v0.2.4-alpha.5 (2022-09-28)

 * [d16a975](https://github.com/tensorchord/envd/commit/d16a975d2ec1c8310e02d0e00841cf82900d49cf) fix: cuda tag (#949)

### Contributors

 * Keming

## v0.2.4-alpha.4 (2022-09-28)

 * [236cd0b](https://github.com/tensorchord/envd/commit/236cd0be8cdfa9e4731d3b69cec7027001c8770b) fix: version tag in build.sh (#947)

### Contributors

 * Keming

## v0.2.4-alpha.3 (2022-09-28)

 * [2eef587](https://github.com/tensorchord/envd/commit/2eef587b44dfe77828570bac8a3e2a7def61c5e1) fix: r & julia sshd image (#945)

### Contributors

 * Keming

## v0.2.4-alpha.2 (2022-09-28)

 * [b704029](https://github.com/tensorchord/envd/commit/b70402917edf874a0a8e630664b637fa8c22cd53) feat(ir): all in llb (#941)
 * [f5f70e0](https://github.com/tensorchord/envd/commit/f5f70e0de304b8ff4767cd935ab3d307ed5599a2) chore(deps): bump pypa/cibuildwheel from 2.10.1 to 2.10.2 (#936)
 * [bd69c3d](https://github.com/tensorchord/envd/commit/bd69c3df326f1f213472960c099dff0c6d35e41c) feat: Support envd-server (#932)

### Contributors

 * Ce Gao
 * Keming
 * dependabot[bot]

## v0.2.4-alpha.1 (2022-09-21)


### Contributors


## v0.2.4 (2022-10-31)

 * [d693254](https://github.com/tensorchord/envd/commit/d693254c1a8b19e02aca85934818a4e7b6dde661) fix: Disable telemetry in CI (#1116)
 * [3cba45e](https://github.com/tensorchord/envd/commit/3cba45e36496cd6010d9b07acbb0c7d34826f828) chore(deps): bump github.com/urfave/cli/v2 from 2.20.3 to 2.23.0 (#1119)
 * [57a63e0](https://github.com/tensorchord/envd/commit/57a63e015c813a0a446768b00c83039923ce0df3) chore(deps): bump pypa/cibuildwheel from 2.11.1 to 2.11.2 (#1115)
 * [77b3d84](https://github.com/tensorchord/envd/commit/77b3d84741f990e9cf0b3a8f9bd56abe40fcd895) chore(deps): bump github.com/onsi/gomega from 1.22.1 to 1.23.0 (#1117)
 * [197b6c6](https://github.com/tensorchord/envd/commit/197b6c612891f0b7bc95b4b1f12a7daacbe7e51f) feat: Add telemetry with the help of segment.io (#1113)
 * [dedd731](https://github.com/tensorchord/envd/commit/dedd73113214d68fd6f8421446fcf2ea8895252e) feat: Add listening_addr to expose fun (#1110)
 * [f672c8f](https://github.com/tensorchord/envd/commit/f672c8f67c9fd793427fb9ec556e767fdf1ef50b) feat: support build time run without mount host (#1109)
 * [214f7c8](https://github.com/tensorchord/envd/commit/214f7c88c64c8b9152e0abb3acfb67e57cc1bf68) fix: panic if the user specify entrypoint for non-costom image (#1108)
 * [8f89ba6](https://github.com/tensorchord/envd/commit/8f89ba6c57795f76f38e0d613ab62517d7b03205) fix: auth with the same name (#1106)
 * [e831fe9](https://github.com/tensorchord/envd/commit/e831fe9c46a4f4085e95395722ffd836bd305c8e) fix: cycle import detection in the interpreter (#1104)
 * [356a707](https://github.com/tensorchord/envd/commit/356a707bfc7e4f9945868fbbcd94a03277d3e601) feat: add config.repo in envd (#1101)
 * [4429bae](https://github.com/tensorchord/envd/commit/4429baebf0bfb6a4f496609c90117881975b4054) fix(CLI): Fix format for app (#1100)
 * [0a975de](https://github.com/tensorchord/envd/commit/0a975de0103d84841252b54b549cddaa6d82d8a6) feat: support envd server destroy env (#1096)
 * [ac87c76](https://github.com/tensorchord/envd/commit/ac87c76f78e8d6ff8c156365310ed5740e8e054b) feature: support runtime init script (#1085)
 * [3416e70](https://github.com/tensorchord/envd/commit/3416e70ad5753441b109e1f6d755d9a02cd26e5c) feat: enable convenient vscode debugging (#1089)
 * [753fde3](https://github.com/tensorchord/envd/commit/753fde329e7da22837af27cddb48969947b9125b) fix:#1090 (#1092)
 * [1ddd513](https://github.com/tensorchord/envd/commit/1ddd513c03b78b8382d309c3a9fb162ae39973b9) feat(CLI): Support envd-server in up command (#1087)
 * [ab9fe2c](https://github.com/tensorchord/envd/commit/ab9fe2ce255ea18780310f37903b216853386589) feat: Support forward (#1014)
 * [4cd6fef](https://github.com/tensorchord/envd/commit/4cd6fef5d6284538ea523011aac47feec8dcf36a) chore(deps): bump github.com/moby/buildkit from 0.10.4 to 0.10.5 (#1081)
 * [a81bbb0](https://github.com/tensorchord/envd/commit/a81bbb0aab9b1a7d1bf37a0374884f34af386c3b) chore(deps): bump github.com/urfave/cli/v2 from 2.20.2 to 2.20.3 (#1078)
 * [a89f466](https://github.com/tensorchord/envd/commit/a89f466517aec054c0e0c5a295f09b5e7b3c33b7) chore(deps): bump github.com/opencontainers/image-spec from 1.1.0-rc1 to 1.1.0-rc2 (#1077)
 * [a209f79](https://github.com/tensorchord/envd/commit/a209f79d07b945c2a1a87972f191be7f4428d6c1) chore(deps): bump github.com/onsi/gomega from 1.21.1 to 1.22.1 (#1080)
 * [a8a6339](https://github.com/tensorchord/envd/commit/a8a633928a5ec70f5deccbebf81619bb1dbe33da) chore(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (#1079)
 * [9faeebe](https://github.com/tensorchord/envd/commit/9faeebe27f47ce2ee1d99b04325eb41745390099) fix: Fix build.sh (#1072)
 * [ff30b3a](https://github.com/tensorchord/envd/commit/ff30b3a5a4b1ccdc089b601c1d4c2bfdd81ce61c) fix: Fix the file name typo (#1071)
 * [51b00fb](https://github.com/tensorchord/envd/commit/51b00fb189422b4791ef2c53cc80ac3a5067b67a) feat(context): Support unix context and daemonless (#1062)
 * [d41b674](https://github.com/tensorchord/envd/commit/d41b67430da64f19b67de4710ad3df1f56f9cc94) bug: Fix sdist didn't include go files (#1068)
 * [181ee71](https://github.com/tensorchord/envd/commit/181ee71a73f3723ca68cd26fd40f53ca96231429) chore(CLI): Move destroy logic from docker to envd engine. (#1050)
 * [672bcd7](https://github.com/tensorchord/envd/commit/672bcd7c04f093f41d113749e39da3f193aea406) bug: fix install from pypi source release (#1064)
 * [cd5cf52](https://github.com/tensorchord/envd/commit/cd5cf522e29ba54113a22867acc396e6fcd128e1) feat: use horust as the dev container supervisor (#1051)
 * [044d89e](https://github.com/tensorchord/envd/commit/044d89e0d0901b3cf75d62c8505b2455b5b91dc5) feat: Update env client to support multiple envs (#1052)
 * [711dfb2](https://github.com/tensorchord/envd/commit/711dfb2ac9f08e467d315427b894191aa598e1d5) feat(lang): Add proposal for custom base image (#567)
 * [b601304](https://github.com/tensorchord/envd/commit/b601304a72bc586c78149a3dc9cec3bd1b1d1c35) chroe: make python package after apt package (#1048)
 * [cf78dd1](https://github.com/tensorchord/envd/commit/cf78dd131e869ad20fe67a4a1159d7bf69e6fb63) feat: Implement env client in envd engine (#1049)
 * [bc6255b](https://github.com/tensorchord/envd/commit/bc6255bc3d494be85afc7c132f2e48f72fffd0fa) example: Add torch profiler example (#1026)
 * [ded3fce](https://github.com/tensorchord/envd/commit/ded3fce1b481aab9c51c2f13dbb2d92288532da0) feat: add `envd down` as an alias for `envd destroy` (#1047)
 * [9b8582d](https://github.com/tensorchord/envd/commit/9b8582d910e83288ef90ba45b599221341ad29f6) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.2.0 to 2.3.1 (#1039)
 * [0ee19ed](https://github.com/tensorchord/envd/commit/0ee19ed79d3332a25082878ca19c1bc032072aef) chore(deps): bump pypa/cibuildwheel from 2.10.2 to 2.11.1 (#1036)
 * [6ba7633](https://github.com/tensorchord/envd/commit/6ba763368d81035b0f13e5eb1549856ecaaa30d1) chore(deps): bump k8s.io/api from 0.25.2 to 0.25.3 (#1038)
 * [0500dca](https://github.com/tensorchord/envd/commit/0500dca9b5174510c1c741ec1d9cfbf97d09a2e9) chore(deps): bump github.com/urfave/cli/v2 from 2.19.2 to 2.20.2 (#1037)
 * [847f187](https://github.com/tensorchord/envd/commit/847f187b45983a56093dd962d5638dae5f260743) chore(deps): bump github.com/onsi/gomega from 1.21.1 to 1.22.1 (#1040)
 * [1fb2721](https://github.com/tensorchord/envd/commit/1fb272197f405c5ee9dcca46c82079095afa406d) bug: Fix flaky auth test (#1034)
 * [d0134d3](https://github.com/tensorchord/envd/commit/d0134d301eba153dcda0c79a3a23babf75a14434) fix(zsh): ignore inserting zsh-completion if system don't have zsh shell (#1025)
 * [ec92a29](https://github.com/tensorchord/envd/commit/ec92a295b835664d559f201ba2542a6a609a560d) feat: support using current directory name as env name in envd describe (#1033)
 * [d3f93a8](https://github.com/tensorchord/envd/commit/d3f93a8389cbf198ce217e2001b7ae0c7df62d2b) bug: fix explicit channel setting with conda env.yaml (#1008)
 * [bcd8029](https://github.com/tensorchord/envd/commit/bcd802976597b3ef07906aa26682fc8ffa0e891a) bug: Fix install from source will lose version information (#1024)
 * [a58197c](https://github.com/tensorchord/envd/commit/a58197c476745bbe568e5712e54912936ca3a14e) fix: Add two debug logs (#1023)
 * [0fc4d8d](https://github.com/tensorchord/envd/commit/0fc4d8d12e9a13c00e760ee105765b0270a2c3ee) feat(context): Get ssh hostname from context, instead of hard-coded string (#1020)
 * [262a07b](https://github.com/tensorchord/envd/commit/262a07b69f976317379d42a8583b1d084c41e73b) feat(CLI): Support envd images prune #976 (#1012)
 * [88dc908](https://github.com/tensorchord/envd/commit/88dc908bc09b99b009c93a72999951096c984daa) fix(CLI): Fix a typo (#1015)
 * [931e1a0](https://github.com/tensorchord/envd/commit/931e1a09d0e3eda869d80e9b73e730c8c8768514) feat(cli): add `envd images remove` command (#1007)
 * [f26fc4c](https://github.com/tensorchord/envd/commit/f26fc4ce7e42ed0e567cb492c41551f6c2f2c3b4) feat(CLI): Support create command (#1001)
 * [f288ff3](https://github.com/tensorchord/envd/commit/f288ff3b8f2cdb2f63cbe313635f3596b3a2ace7) feat(CI): add nightly build & test tasks (#1006)
 * [46168da](https://github.com/tensorchord/envd/commit/46168dac70fb8cd9b841a32e35a414ad29ddfaa2) chore(test): Add test cases for pkg/home/auth.go (#1009)
 * [d87afda](https://github.com/tensorchord/envd/commit/d87afda9213fcb334ca55d247ba3c0663b7b180b) feat: Add envd-server runtime proposal (#303)
 * [fec7ada](https://github.com/tensorchord/envd/commit/fec7ada4efc2186fa5f858d5abe211a2beb2dfdd) fix: Remove hard code docker in envd engine init (#1000)
 * [cf56f0d](https://github.com/tensorchord/envd/commit/cf56f0de7bfb05be47d4704c14f9bbb892a86d28) chore(deps): bump github.com/onsi/gomega from 1.20.2 to 1.21.1 (#997)
 * [a55470e](https://github.com/tensorchord/envd/commit/a55470e363b35b689649596b449b996ada701630) feat(build): detect if the current environment is running before building (#892) (#989)
 * [d0219f2](https://github.com/tensorchord/envd/commit/d0219f2e3e999e27f334a39180031390a4554af9) chore(deps): bump github.com/urfave/cli/v2 from 2.17.1 to 2.19.2 (#998)
 * [b2a9018](https://github.com/tensorchord/envd/commit/b2a90188c012fa186de2ecb7b1aa534062681020) feat: Support cli argument for host key (#992)
 * [91ea50c](https://github.com/tensorchord/envd/commit/91ea50cd74e36f2d0d966d2b6a1ddea0ad15c43b) bug: fix pypi package information (#959)
 * [0be9969](https://github.com/tensorchord/envd/commit/0be9969471f8358f1cc283215ddf8e4cfcecebad) fix: Remove dump checkout and remove pre-commit (#982)
 * [f4ebb02](https://github.com/tensorchord/envd/commit/f4ebb026c45ab29a84dfa57782ef5fec8c11cf14) Chore(test): Add test cases for pkg/util/starlarkutil (#979)
 * [06d3cfa](https://github.com/tensorchord/envd/commit/06d3cfa1e171a409893d59494754dc47fb886a52) *: upgrade golangci-linter and add noloopclosure (#980)
 * [17fedb8](https://github.com/tensorchord/envd/commit/17fedb8d2ae41318e2f705687c334b7e42413ab3) fix(ir): make sure default value won't be replaced with empty value (#970)
 * [c1ae887](https://github.com/tensorchord/envd/commit/c1ae887f12b23802363ac417c198e19983c5605f) feat(CLI): Support runner in context (#961)
 * [bf993e2](https://github.com/tensorchord/envd/commit/bf993e2cab44502abb74a79e732b74a5567b6194) refact: conda/mamba create/update env, fix user permissions (#933)
 * [0e79fb9](https://github.com/tensorchord/envd/commit/0e79fb91b9f3770f2e5356fc8a10182dd31e2bdf) chore(deps): bump github.com/urfave/cli/v2 from 2.16.3 to 2.17.1 (#968)
 * [c9045d2](https://github.com/tensorchord/envd/commit/c9045d24d412099a1ac0bf9b9acf78135157351d) chore(deps): bump dependabot/fetch-metadata from 1.3.3 to 1.3.4 (#967)
 * [c22408c](https://github.com/tensorchord/envd/commit/c22408c8f3b87ef959bda7d7681203ffd8d6212c) fix(ir): `apt install` and `conda env create` cache (#962)
 * [006f653](https://github.com/tensorchord/envd/commit/006f6538396b47b1bcab2a0d53ec9dae0221c8f5) feat: envd-sshd can read public key path from environment variable (#954)
 * [0b73548](https://github.com/tensorchord/envd/commit/0b7354868943f9e4ea189be18eeea534959b5d6c) example: Use torch in mnist example (#927)
 * [d16a975](https://github.com/tensorchord/envd/commit/d16a975d2ec1c8310e02d0e00841cf82900d49cf) fix: cuda tag (#949)
 * [236cd0b](https://github.com/tensorchord/envd/commit/236cd0be8cdfa9e4731d3b69cec7027001c8770b) fix: version tag in build.sh (#947)
 * [2eef587](https://github.com/tensorchord/envd/commit/2eef587b44dfe77828570bac8a3e2a7def61c5e1) fix: r & julia sshd image (#945)
 * [b704029](https://github.com/tensorchord/envd/commit/b70402917edf874a0a8e630664b637fa8c22cd53) feat(ir): all in llb (#941)
 * [f5f70e0](https://github.com/tensorchord/envd/commit/f5f70e0de304b8ff4767cd935ab3d307ed5599a2) chore(deps): bump pypa/cibuildwheel from 2.10.1 to 2.10.2 (#936)
 * [bd69c3d](https://github.com/tensorchord/envd/commit/bd69c3df326f1f213472960c099dff0c6d35e41c) feat: Support envd-server (#932)
 * [171b82f](https://github.com/tensorchord/envd/commit/171b82fcbac1da9543b9acda1911d110c322d802) fix: e2e doc test (#926)
 * [ce36545](https://github.com/tensorchord/envd/commit/ce36545d4865a6062ab4b6332a0e26e0d7030db2) feat(cli): rm image when destroying the env (#925)
 * [eecc7cf](https://github.com/tensorchord/envd/commit/eecc7cf48cf354ef3d72f17dfcf1d6f19bf90b85) feat: informative error message (#859)
 * [a3fd464](https://github.com/tensorchord/envd/commit/a3fd4649f83f95d4e349348e4367bf564447780f) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.1.6 to 2.2.0 (#919)
 * [756c92d](https://github.com/tensorchord/envd/commit/756c92d48cb819149be5520948de122b687a25aa) chore(deps): bump github.com/urfave/cli/v2 from 2.16.2 to 2.16.3 (#918)
 * [b5acc08](https://github.com/tensorchord/envd/commit/b5acc08d46b6b84a0c57d19a2e1b62e247785bd4) chore(deps): bump pypa/cibuildwheel from 2.9.0 to 2.10.1 (#917)

### Contributors

 * Ce Gao
 * Friends A
 * JasonZhu
 * Jinjing Zhou
 * Keming
 * Tumushimire Yves
 * Weizhen Wang
 * XRW
 * Yijiang Liu
 * Yilong Li
 * Zhenzhen Zhao
 * dependabot[bot]
 * nullday
 * wangxiaolei

## v0.2.3 (2022-09-16)

 * [c3c0b4e](https://github.com/tensorchord/envd/commit/c3c0b4e33ab696c3863632d0d4179f8211813fcb) fix: Use macos 11 (#912)

### Contributors

 * Ce Gao

## v0.2.2 (2022-09-16)

 * [6721a88](https://github.com/tensorchord/envd/commit/6721a88aa1e2a97bfd3060e3168ab15b61f25609) chore(goreleaser): Skip homebrew (#909)

### Contributors

 * Ce Gao

## v0.2.1 (2022-09-16)


### Contributors


## v0.2.0-beta.1 (2022-09-27)


### Contributors


## v0.2.0-beta6 (2022-09-28)

 * [9046f3b](https://github.com/tensorchord/envd/commit/9046f3b1e0bd6e990e0b651584d042a585fd17c8) fix gorlease v
 * [9fea9e8](https://github.com/tensorchord/envd/commit/9fea9e8b1944acf82ac81d07874d7b2d22b4ce97) cache gpu
 * [c7917a6](https://github.com/tensorchord/envd/commit/c7917a607553338b58a4d40af1a3262cd5618bb8) use git tag
 * [c5aeec6](https://github.com/tensorchord/envd/commit/c5aeec6ee82fdeb2d6da2ecc0f5d3baf1093e918) fix ref bug

### Contributors

 * Keming

## v0.2.0-beta5 (2022-09-27)

 * [86f0b17](https://github.com/tensorchord/envd/commit/86f0b170902275d9320de819cfd708d6fd797015) fix ci secret

### Contributors

 * Keming

## v0.2.0-beta4 (2022-09-27)

 * [5cb30f3](https://github.com/tensorchord/envd/commit/5cb30f33067c3334435ade69a5c082a870a63a70) only run some ci in tensorchord
 * [cfefb15](https://github.com/tensorchord/envd/commit/cfefb15c7120b04bb711ee4a7f84e0803fbae487) fix ci yml
 * [84c80be](https://github.com/tensorchord/envd/commit/84c80be8cb4f753a2131fa671569925bd9c1bf7f) fix ci
 * [cea390d](https://github.com/tensorchord/envd/commit/cea390ddf1a61358ade9cebbe2f4fdc71a7233fe) del python dockerfile
 * [2a45a8f](https://github.com/tensorchord/envd/commit/2a45a8f2f918288c3d862597c4c94566f9138773) cuda version
 * [5e407d0](https://github.com/tensorchord/envd/commit/5e407d0eb1a73cac86721b66c2bed290d35ac76e) delete outdated comments
 * [16ca22f](https://github.com/tensorchord/envd/commit/16ca22f7d2e1f02060f299c4a69c74db647ef843) add llb log
 * [7327c13](https://github.com/tensorchord/envd/commit/7327c13c32ff897160b40e49eb70b884589b5bd1) copy envd-sshd
 * [f58e55b](https://github.com/tensorchord/envd/commit/f58e55ba4b52b7b9efa1a4669e5f5b11b364b7a0) change to docker hub
 * [0863ac5](https://github.com/tensorchord/envd/commit/0863ac5c6ea050dd5832d54c6dd90f80b5adf3d4) change python dockerfile to llb
 * [bd69c3d](https://github.com/tensorchord/envd/commit/bd69c3df326f1f213472960c099dff0c6d35e41c) feat: Support envd-server (#932)
 * [171b82f](https://github.com/tensorchord/envd/commit/171b82fcbac1da9543b9acda1911d110c322d802) fix: e2e doc test (#926)
 * [ce36545](https://github.com/tensorchord/envd/commit/ce36545d4865a6062ab4b6332a0e26e0d7030db2) feat(cli): rm image when destroying the env (#925)
 * [eecc7cf](https://github.com/tensorchord/envd/commit/eecc7cf48cf354ef3d72f17dfcf1d6f19bf90b85) feat: informative error message (#859)
 * [a3fd464](https://github.com/tensorchord/envd/commit/a3fd4649f83f95d4e349348e4367bf564447780f) chore(deps): bump github.com/onsi/ginkgo/v2 from 2.1.6 to 2.2.0 (#919)
 * [756c92d](https://github.com/tensorchord/envd/commit/756c92d48cb819149be5520948de122b687a25aa) chore(deps): bump github.com/urfave/cli/v2 from 2.16.2 to 2.16.3 (#918)
 * [b5acc08](https://github.com/tensorchord/envd/commit/b5acc08d46b6b84a0c57d19a2e1b62e247785bd4) chore(deps): bump pypa/cibuildwheel from 2.9.0 to 2.10.1 (#917)
 * [c3c0b4e](https://github.com/tensorchord/envd/commit/c3c0b4e33ab696c3863632d0d4179f8211813fcb) fix: Use macos 11 (#912)
 * [6721a88](https://github.com/tensorchord/envd/commit/6721a88aa1e2a97bfd3060e3168ab15b61f25609) chore(goreleaser): Skip homebrew (#909)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming
 * dependabot[bot]

## v0.2.0-alpha.22 (2022-09-16)

 * [4707bfe](https://github.com/tensorchord/envd/commit/4707bfeaea030a48c5601bffb87ff034e4d1b413)  fix: Fix jupyter in root (#900)
 * [71fb1ce](https://github.com/tensorchord/envd/commit/71fb1ce154ba7335de6d40e38f63d65012706c86) feat: support micromamba as an alternative to miniconda (#891)
 * [bab012c](https://github.com/tensorchord/envd/commit/bab012c5209092632d84d6cfb0c8b78fc2946523) fix: typo for git config file (#888)
 * [128f866](https://github.com/tensorchord/envd/commit/128f866f4f030cf2b10af87fe32078329e0519d8) fix(CLI): Fix build output argument and huggingface integration (#886)
 * [2e8b5d5](https://github.com/tensorchord/envd/commit/2e8b5d5d4756b5c02c6d3e846b5be093fd6394b1) fix: include update repo (#885)
 * [eb2cdd1](https://github.com/tensorchord/envd/commit/eb2cdd1a65321a1530f59190cd40560e6c31d5a3) bug: Fix detach instruction message (#882)
 * [8a02b26](https://github.com/tensorchord/envd/commit/8a02b264ac16318c5d54660ca883414ef7a15cad) refact: add envd home path func (#880)
 * [63daa5e](https://github.com/tensorchord/envd/commit/63daa5e870e5ea24c3f6e881e4932da2334688dc) chore(deps): bump github.com/spf13/viper from 1.12.0 to 1.13.0 (#875)
 * [aa53bdb](https://github.com/tensorchord/envd/commit/aa53bdb478645863dae60ca37f1ebdbc7a564c56) chore(deps): bump github.com/urfave/cli/v2 from 2.14.0 to 2.16.2 (#874)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming
 * dependabot[bot]

## v0.2.0-alpha.21 (2022-09-12)

 * [ddf3bdc](https://github.com/tensorchord/envd/commit/ddf3bdc8dd859683d7539d3c7f226b82cece40e7) chore(deps): bump actions/setup-go from 2 to 3 (#873)
 * [7220958](https://github.com/tensorchord/envd/commit/72209583b941c4eb87b282004f7c8a1ae57410ae) chore(deps): bump actions/checkout from 2 to 3 (#872)
 * [f1b3fe5](https://github.com/tensorchord/envd/commit/f1b3fe5029091cffc0a1022b665d585132c5a8d8) chore(CLI): test new release for envd-sshd (#866)

### Contributors

 * Yuedong Wu
 * dependabot[bot]

## v0.2.0-alpha.20 (2022-09-11)

 * [b38c5de](https://github.com/tensorchord/envd/commit/b38c5de82aac50fe085cd48d2111f7f5d241b6d7) chore(CLI): test new release for envd-sshd (#866)
 * [49d79fb](https://github.com/tensorchord/envd/commit/49d79fb17bee4a2baaeadd500607cba7d8426b28) fix: Update readme (#865)
 * [d7995a7](https://github.com/tensorchord/envd/commit/d7995a7171cbe48a65aad3e3b56077ffee9a625a) feat(lang): io.http download files to extra_source (#858)
 * [0d3b42f](https://github.com/tensorchord/envd/commit/0d3b42fe4f33241742986030a45431a5f068dc75) feat: Support HTTP PROXY (#857)

### Contributors

 * Ce Gao
 * Keming
 * Yuedong Wu

## v0.2.0-alpha.19 (2022-09-09)

 * [a0fbaa0](https://github.com/tensorchord/envd/commit/a0fbaa09fee056549b1d6fcd796f5b711268de2b) refact: io.mount => runtime.mount (#861)
 * [1e5e24d](https://github.com/tensorchord/envd/commit/1e5e24d1f872311a81b193443be36eaed22cc11e) bug: fix conda install with env file (#837)
 * [ecb9e26](https://github.com/tensorchord/envd/commit/ecb9e2626e65b5e1f647d7385142c10677f2d7eb) refact: unify the path env (#855)
 * [8056fda](https://github.com/tensorchord/envd/commit/8056fda28febfe6cbe64502159b302b670970517) feat: add runtime graph to image label (#815)
 * [6ad1d4c](https://github.com/tensorchord/envd/commit/6ad1d4ca2085317850e9726910bc3024b743439e) refact: apt_source, io, config mode (#853)
 * [07e2dc0](https://github.com/tensorchord/envd/commit/07e2dc0e6d4d7f3be2debf79389a792964e727b1) chore(deps): bump github.com/gliderlabs/ssh from 0.3.4 to 0.3.5 (#849)

### Contributors

 * Jinjing Zhou
 * Keming
 * dependabot[bot]
 * nullday

## v0.2.0-alpha.18 (2022-09-06)

 * [c0ba31a](https://github.com/tensorchord/envd/commit/c0ba31adc231d09600892f9445320df8eb84947b) chore(deps): bump github.com/onsi/gomega from 1.20.1 to 1.20.2 (#846)
 * [13bc9a6](https://github.com/tensorchord/envd/commit/13bc9a61c4255215661c724fb7c54f2b81642a21) chore(deps): bump github.com/urfave/cli/v2 from 2.11.2 to 2.14.0 (#845)
 * [11966f8](https://github.com/tensorchord/envd/commit/11966f8c0f997f57c45979d7857cac48ad1d1e5b) chore(deps): bump github.com/docker/go-units from 0.4.0 to 0.5.0 (#848)
 * [7ba3b3f](https://github.com/tensorchord/envd/commit/7ba3b3fa6b716efd3398a7cb3e533e915129717d) feat(cli): add msg when detach from container (#841)
 * [7fc9f34](https://github.com/tensorchord/envd/commit/7fc9f34d333b22013e847fd2c2312d18fb861068) fix: Update demo (#840)
 * [190ee76](https://github.com/tensorchord/envd/commit/190ee7635f0d9bbdee6c8b53d53a72dd8ca4e619) feat(lang): install.python_packages(local_wheels=[]) (#838)
 * [977dd47](https://github.com/tensorchord/envd/commit/977dd4725df2c16887546a48b6e3fa202a7617e2) fix: Update demo (#839)
 * [5e7c182](https://github.com/tensorchord/envd/commit/5e7c1826465511449edf6457428037cfe3afbc7e) bug: fix channels when use conda install with yaml file (#831)
 * [be02a70](https://github.com/tensorchord/envd/commit/be02a7007d65ce1227bfc61dba59de33f79c295b) fix(lang): expose host port (#832)
 * [404de31](https://github.com/tensorchord/envd/commit/404de3101cece0497084412433cf877f66cf5ee2) feat(lang): init py env by generating the bulid.envd (#827)
 * [36b1231](https://github.com/tensorchord/envd/commit/36b123142385d20fb7f7c1106c15c02f79ed4742) bug: fix permission issue when pip install from git repo (#829)
 * [1dcada4](https://github.com/tensorchord/envd/commit/1dcada4403d0c0bf8e916fc67b62c174f66df3d3) feat(build): Mount local build context into the run command (#822)
 * [f70a11c](https://github.com/tensorchord/envd/commit/f70a11c5f251b8ed0f0f42cd422b4b93efabd4a7) chore(deps): bump github.com/onsi/gomega from 1.20.0 to 1.20.1 (#821)
 * [3b440c6](https://github.com/tensorchord/envd/commit/3b440c60c59ee65444a92aa91e295be9be2125b0) chore(deps): bump github.com/moby/buildkit from 0.10.3 to 0.10.4 (#820)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming
 * dependabot[bot]

## v0.2.0-alpha.17 (2022-08-26)

 * [82fbc87](https://github.com/tensorchord/envd/commit/82fbc87ade5637fe7db8e2c0087a1555206dc1b1) doc: add include, refine others (#817)
 * [630ada1](https://github.com/tensorchord/envd/commit/630ada172bdf876c3b749329fdbe284c108051f2) feat(lang): support include other git repo for envd functions/variables (#808)
 * [5c4971b](https://github.com/tensorchord/envd/commit/5c4971b2f5fa6b32c2c247e18cf6c6a178d28f57) feat(CLI): envd env describe expose info (#801)
 * [c85766c](https://github.com/tensorchord/envd/commit/c85766cf618d55723d26427f65a385747646593d) fix: set to latest if git tag is empty (#798)

### Contributors

 * Keming
 * Zhizhen He

## v0.2.0-alpha.16 (2022-08-19)

 * [4440e22](https://github.com/tensorchord/envd/commit/4440e2246108a221585853088a766a563b2c7aad) fix: add missing expose func exposed port to oci manifest (#797)
 * [3a97375](https://github.com/tensorchord/envd/commit/3a97375383a2d135cab06665aae67f04230666e1) feat(examples): Add a streamlit mnist example (#795)
 * [7bf801b](https://github.com/tensorchord/envd/commit/7bf801bad5c5947a80e576ceb3b4bae0307fddeb) feat(example): Add streamlit hello and remove bash -c in entrypoint (#794)
 * [8e17307](https://github.com/tensorchord/envd/commit/8e173075d56cadc1485f7c526c1d59c00926c69e) built: :hammer: use latest tag when not version found for cache (#793)
 * [8225eab](https://github.com/tensorchord/envd/commit/8225eab40d1a4b37b1a6d82d300399d76cfd1320) fix: use cockroachdb errors (#790)
 * [7d293d7](https://github.com/tensorchord/envd/commit/7d293d7974cae516e7d5a9f8d514acea04a0ff13) fix: panic if daemon command is invalid (#788)
 * [4d48767](https://github.com/tensorchord/envd/commit/4d48767ce033a3d60fac74d810ea546f741bc174) feat: add runtime environments (#787)
 * [248fca3](https://github.com/tensorchord/envd/commit/248fca34fa67f6fd134e1c55eefdd931aa5d8939) doc: daemon and expose (#786)
 * [e30866f](https://github.com/tensorchord/envd/commit/e30866f07c1249c5a586558a22a18741b07e063d) feat(lang): implement expose func (#780)
 * [c49863e](https://github.com/tensorchord/envd/commit/c49863e1787b6b18c08d49035ec00e79e0020822) feat(data): Add support for managed dataset and provide shortcut for common framework (#751)
 * [7c2fed6](https://github.com/tensorchord/envd/commit/7c2fed6df565e55883456b934116314109c3837d) feat(lang): add daemon function to run daemon process in the container (#777)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Keming
 * Wei Zhang
 * nullday

## v0.2.0-alpha.15 (2022-08-16)

 * [00e8df2](https://github.com/tensorchord/envd/commit/00e8df22ba285685b11d217fdcb4cd7a52e32ba5) built: :hammer: keep v prefix in DOCKER_IMAGE_TAG env (#781)
 * [2f82fa5](https://github.com/tensorchord/envd/commit/2f82fa5ab884671ba2a5058d38d205cbd20bce1f) proposal: daemon process (#769)
 * [d06b878](https://github.com/tensorchord/envd/commit/d06b8786e2963f4bf5cc92b90fdcc80621a4bd5e) doc: update python api doc (#759)
 * [ce2e8b2](https://github.com/tensorchord/envd/commit/ce2e8b2b4e755b52892d9b3597489daae19f1dad) fix: Remove empty token arg (#772)
 * [8ab89b9](https://github.com/tensorchord/envd/commit/8ab89b967caa104ad839a83d7d9ece186eb80918) chore(deps): bump github.com/urfave/cli/v2 from 2.11.1 to 2.11.2 (#775)
 * [3a6b127](https://github.com/tensorchord/envd/commit/3a6b12772b956eb31b6ccc6c31e62efbc3feb6e3) chore(deps): bump pypa/cibuildwheel from 2.8.1 to 2.9.0 (#774)
 * [ec8cae1](https://github.com/tensorchord/envd/commit/ec8cae17b26de61bbb73026540db56505cffb2a8) fix: remove unnecessary if statement (#773)

### Contributors

 * Ce Gao
 * Keming
 * Wei Zhang
 * Zhizhen He
 * dependabot[bot]

## v0.2.0-alpha.14 (2022-08-12)

 * [b9f0af8](https://github.com/tensorchord/envd/commit/b9f0af8056129fe6c5a2e590cd428463186cac0e) fix: -path and -file bug (#766)
 * [4a359f1](https://github.com/tensorchord/envd/commit/4a359f1eb608a0680a94379b4ac52756d605be9d) fix(release): :hammer: drop go build dep for homebrew (#768)
 * [542c7cb](https://github.com/tensorchord/envd/commit/542c7cbc3f8b7658905b85e10ff47ea216b44f7e) Fix the color display in wezterm (#767)
 * [590b4c0](https://github.com/tensorchord/envd/commit/590b4c04bffab550778a3fe911db3280ffb72b09) fix(CLI): 🔨  use latest version for local build (#763)
 * [2342367](https://github.com/tensorchord/envd/commit/2342367d3edb6d6c240afffc6b2414b1e83f0413) fix(docs): :memo: fix contributing and dev links, clean tailing space (#764)
 * [b4a8519](https://github.com/tensorchord/envd/commit/b4a851920092b66f97ed6af301472af183c986a8) fix: modify jupyter's authority from hash password to token string (#762)
 * [1f0af98](https://github.com/tensorchord/envd/commit/1f0af9825eb6f1b8843feb7fbb9b12a10e1c902d) fix: use defined jupyter port (#757)
 * [381b653](https://github.com/tensorchord/envd/commit/381b653b4dd428e5b68c1876bac74653ae2b3068) docs(README): update Documentations link to https://envd.tensorchord.… (#758)
 * [ad4b9ec](https://github.com/tensorchord/envd/commit/ad4b9ec8894b08104a47defd0039fdb081b087e3) fix: do not expose ports for custom image (#754)
 * [e77b7d6](https://github.com/tensorchord/envd/commit/e77b7d61cfaeec1d62a7119373a646a7907cb933) docs(README): Correct the cmd of get Jupyter Notebook endpoint (#756)
 * [8675316](https://github.com/tensorchord/envd/commit/86753162812f07ad030a6aecbcaa7f57839aae54) fix: ParseFromStr, add unittest (#755)
 * [382aa2d](https://github.com/tensorchord/envd/commit/382aa2db5a6bdd26d1a9f96f775cd749f9927f4e) bug(CLI): fix short alias confusion (#752)
 * [7625cd2](https://github.com/tensorchord/envd/commit/7625cd28658ad02dd05f523446985152e8eb887b) feat: Avoid gid in base image cache (#749)
 * [2e39182](https://github.com/tensorchord/envd/commit/2e39182582bf779e304fa6b7878c1428be938694) feat(CLI): :sparkles: add --force args for init to overwrite build.envd (#748)
 * [d241898](https://github.com/tensorchord/envd/commit/d2418984c7e67cc1a5184fc2894e8439e58e97ef) docs(README): One obvious way to declare supported python version (#745)
 * [b492e83](https://github.com/tensorchord/envd/commit/b492e838eb1a9bef6769806c2604d51fecfa2fff) add entrypoint in custom image (#739)
 * [4df956b](https://github.com/tensorchord/envd/commit/4df956b14ec05ac4204d8891b6791f065205b5c5) The workaround to fix the label loss (#741)
 * [dd609da](https://github.com/tensorchord/envd/commit/dd609da753b71331bfb48f997760890b47563ad4) feat: Enable build for all languages (#738)
 * [ed998ce](https://github.com/tensorchord/envd/commit/ed998cede32a4cdaa996964c6ec7ef492b4802d6) feat(lang): Add a new func runtime.command (#736)
 * [3c46efc](https://github.com/tensorchord/envd/commit/3c46efcf651af6d80010e93212ef8d2f0e26e8fa) fix: setup.py build (#735)

### Contributors

 * Bingyi Sun
 * Ce Gao
 * Keming
 * Wei Zhang
 * Zhenguo.Li
 * nullday

## v0.2.0-alpha.13 (2022-08-06)


### Contributors


## v0.2.0-alpha.12 (2022-08-06)

 * [8bef795](https://github.com/tensorchord/envd/commit/8bef795c6acd581b68492df715d67c9ed32ccc49) feat(CLI): :recycle: refactor bootstrap command to show what's envd doing (#728)
 * [8f19748](https://github.com/tensorchord/envd/commit/8f197483f9f59e0712bc0f54259331621fc05ad8) Bug: fix SIGSEGV of envd top (#726)
 * [8c83cad](https://github.com/tensorchord/envd/commit/8c83cada88d3cb0fff0eb6c19aa880c6c163cf17) feat(base-image): Move conda to llb and cache it (#724)
 * [5587c3f](https://github.com/tensorchord/envd/commit/5587c3fa6f70e465df5ba14f3cec2f99339f7bcb) feature: add envd top commands (#718)
 * [f22b283](https://github.com/tensorchord/envd/commit/f22b28328c6bd1b59bd870ef6b3a49bf0281375f) bug: allow multiple run command and use bash -c  (#720)
 * [da8feb7](https://github.com/tensorchord/envd/commit/da8feb77175448a4ce2bb2b0815a8a9f62c347b2) bug: fix notebook entry when setting conda channel (#719)
 * [116271c](https://github.com/tensorchord/envd/commit/116271c93323de18eb33fc834c971c9ca7e85b12) example: add a dgl GAT example (#714)
 * [84ffea3](https://github.com/tensorchord/envd/commit/84ffea335115d94162896123831ad5459521e2ca) feat(lang): Add proposal for expose (#568)
 * [ace70e9](https://github.com/tensorchord/envd/commit/ace70e95f4bc54cf1d58585b62d1976344d79728) feat: support io.mount (#708)
 * [a1e0395](https://github.com/tensorchord/envd/commit/a1e039595539d966916f59daacc0613b3b06bc28) feat(lang): Add default conda pkg cache (#705)
 * [c515f3c](https://github.com/tensorchord/envd/commit/c515f3cf62b2450eb36dc659974fd8875cf56f48) feat(CLI): Add category and refine help text (#707)
 * [54f412a](https://github.com/tensorchord/envd/commit/54f412a21eb47060439e3749c34d57bbfce26ba4) feat(CLI): Support run command (#701)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * Wei Zhang
 * nullday

## v0.2.0-alpha.11 (2022-07-29)

 * [62bd2dc](https://github.com/tensorchord/envd/commit/62bd2dcf37be2c8b811b34ea6d05189250dbdd1f) bug: Fix version prefix (#696)

### Contributors

 * Jinjing Zhou

## v0.2.0-alpha.10 (2022-07-29)

 * [4dffb5d](https://github.com/tensorchord/envd/commit/4dffb5d94eaab981959fa56cfd7e5d540776f1c7) bug: Fix git tag version (#692)

### Contributors

 * Jinjing Zhou

## v0.2.0-alpha.9 (2022-07-29)

 * [6daa661](https://github.com/tensorchord/envd/commit/6daa661bdcbaf3bdcbad1e63e34039e54ffe088d) bug: Add release dependency (#689)
 * [1f4b468](https://github.com/tensorchord/envd/commit/1f4b46864a03a449127ac540687571943e965102) bug: Use git version by default when build ssh (#688)

### Contributors

 * Jinjing Zhou

## v0.2.0-alpha.8 (2022-07-29)

 * [728a421](https://github.com/tensorchord/envd/commit/728a421840ba10c6ba78c836a22dc7af140d370d) feat: add starship as the prompt manager (#681)
 * [6cbc53f](https://github.com/tensorchord/envd/commit/6cbc53f54c7b95c374dbc93a242e78f7be5dbd86) fix: Add environment variable PATH in run (#680)
 * [cb8ed6b](https://github.com/tensorchord/envd/commit/cb8ed6b5c6ae0a4ea77eb853bfccb450137b6b19) feat(lang): Add io.copy (#675)
 * [b3ee633](https://github.com/tensorchord/envd/commit/b3ee6334a54e7d89b3f2631631c112c1d5ecac6a) fix: Fix lint issues in conda env yaml feature PR (#679)
 * [a796a12](https://github.com/tensorchord/envd/commit/a796a129d83e5c707bb9d0c1c32295208b136495) feat(lang): Support env.yaml in conda_packages (#674)
 * [af3e78d](https://github.com/tensorchord/envd/commit/af3e78d6733f53f5d31eda045725c0af599aefa6) feature:  add label ai.tensorchord.envd.build.manifestBytecodeHash to image for cache robust (#661)
 * [3044002](https://github.com/tensorchord/envd/commit/30440022266b908adbb790e59c6696b36bf92b28) chore(deps): bump github.com/onsi/gomega from 1.19.0 to 1.20.0 (#657)
 * [400fd7f](https://github.com/tensorchord/envd/commit/400fd7f8af0a04d1d0aa5fd7f047fea59172a7ef) chore(deps): bump github.com/urfave/cli/v2 from 2.11.0 to 2.11.1 (#655)
 * [30e8caf](https://github.com/tensorchord/envd/commit/30e8cafc80b04b1bfa6e4d9d2ba273928c8b3919) chore(deps): bump actions/upload-artifact from 2 to 3 (#654)
 * [04360d6](https://github.com/tensorchord/envd/commit/04360d63ba984a1a58faaf300d0cee64269d442a) chore(deps): bump pypa/cibuildwheel from 2.8.0 to 2.8.1 (#653)
 * [de6c59e](https://github.com/tensorchord/envd/commit/de6c59e730da4520bccc747bb941bb12ced913d7) chore(deps): bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 (#656)
 * [62091a4](https://github.com/tensorchord/envd/commit/62091a460b5502d24fa57a04ad19cf273c832832) fix(build): Fix image config (#651)
 * [4c1df28](https://github.com/tensorchord/envd/commit/4c1df286df7c3d5c4a99b4bbc5958f86d08e0738) fix: context create with 'use' (#652)
 * [6f05072](https://github.com/tensorchord/envd/commit/6f05072f1f2a72b55d8e2bac468707219ceb30b4) feat(CLI): Support cache (#648)
 * [42e7531](https://github.com/tensorchord/envd/commit/42e75312a84e590c4c3b1d2f9de4ec7a6a716bb3) remove xdg, use $HOME/.config and $HOME/.cache (#641)

### Contributors

 * Ce Gao
 * Keming
 * dependabot[bot]
 * nullday
 * wyq

## v0.2.0-alpha.7 (2022-07-20)

 * [67a1f34](https://github.com/tensorchord/envd/commit/67a1f340f45966161dbac2afb9083baeac948674) feat(lang): Remove conda from custom base image (#626)
 * [890119d](https://github.com/tensorchord/envd/commit/890119d119e7becf9f39f392aa053ae1e70c77c0) fix: check manifest and image update in new gateway buildfunc (#624)
 * [c8471db](https://github.com/tensorchord/envd/commit/c8471db20f0a7e48c2c8346b3ab88637445200f5) support buildkit TCP socket (#599)
 * [ef8a90d](https://github.com/tensorchord/envd/commit/ef8a90df3bc1a3009381eb3fd10767f468fcefc2) feat: Refactor with Builder.Options (#615)
 * [2a88ad1](https://github.com/tensorchord/envd/commit/2a88ad120b2a24b5095883a1e18de55189ff643f) chore(deps): bump github.com/urfave/cli/v2 from 2.10.3 to 2.11.0 (#610)

### Contributors

 * Ce Gao
 * Keming
 * dependabot[bot]
 * nullday

## v0.2.0-alpha.6 (2022-07-15)

 * [18abe90](https://github.com/tensorchord/envd/commit/18abe90072835534f75b392b5f1ef6dbbf0bbeb5) feat(builder): Abstract BuildFunc to use gateway client (#606)
 * [178b8da](https://github.com/tensorchord/envd/commit/178b8dafdd4357688a911bfff103c397b08410a9) feat(WSL): Add ssh config entry to Windows ssh config if using WSL (#604)
 * [ceb07f5](https://github.com/tensorchord/envd/commit/ceb07f5fc43ced6da9ada711eb8c127d14afa969) fix: set conda as the only python provider (#602)
 * [f1dd546](https://github.com/tensorchord/envd/commit/f1dd546fe5598f11145a5a9c354a1e4878614ed3) fix: pre-create conda package cache directory (#600)
 * [9b3fbe3](https://github.com/tensorchord/envd/commit/9b3fbe3c9c91c167196be3764de4f53b1a074489) feat(lang): Support image in base (#595)
 * [b467279](https://github.com/tensorchord/envd/commit/b46727981968f8722feeff1023e0d85fa7bdc162) fix: Fix error handling issue (#597)
 * [fa041a8](https://github.com/tensorchord/envd/commit/fa041a849dde7f61f82fa2afdf105e45efc888ef) fix: Pre-mkdir the .cache directory of user envd (#592)
 * [54dfc52](https://github.com/tensorchord/envd/commit/54dfc52f4adc5cbce5de806dc6186e39adf21e0d) bug: fix missing function in example mnist (#589)
 * [00249bb](https://github.com/tensorchord/envd/commit/00249bbd2330088dc584015353837917bd557503) Use DefaultText in up.go (#587)
 * [302e449](https://github.com/tensorchord/envd/commit/302e4490992d5f74bf5a9e08293cfb42a89e0367) chore(deps): bump pypa/cibuildwheel from 2.7.0 to 2.8.0 (#583)

### Contributors

 * Ce Gao
 * Guangyang Li
 * Jinjing Zhou
 * dependabot[bot]
 * nullday

## v0.2.0-alpha.5 (2022-07-08)

 * [6cfc0f1](https://github.com/tensorchord/envd/commit/6cfc0f16224095605dd85fb38b9cf406fbb65118) feat: Support for build image update when exec build or up again (#570)
 * [8f89e4b](https://github.com/tensorchord/envd/commit/8f89e4be3d154728824c67f13086eb727f545400) Fix: image tag normalized to docker spec (#573)
 * [3fe3757](https://github.com/tensorchord/envd/commit/3fe375769487a4eaf224a6db1437c840002c7a15) fix: add -c for every single conda channel (#569)
 * [49fa961](https://github.com/tensorchord/envd/commit/49fa961111492d9a0599e9697bf2361321e9417d) fix: add auto start buildkit container (#563)
 * [4fa5ec7](https://github.com/tensorchord/envd/commit/4fa5ec7b520964a74009b397ae9755ae96193305) bug: Fix github action (#566)
 * [93027bd](https://github.com/tensorchord/envd/commit/93027bd669ddfe052c9abcd2c2679f547389b1ff) fix: py cmd exit code (#564)
 * [707d5e8](https://github.com/tensorchord/envd/commit/707d5e8ca880a7968ef0304fab5dd28fb05a1610) feat: replace IsCreated with Exists for Client interface from package docker (#558)
 * [f71cd7f](https://github.com/tensorchord/envd/commit/f71cd7f4891fab5e1be588e9b991ab4942e02d57) feat(CLI): Unify CLI style about env and image (#550)

### Contributors

 * Jinjing Zhou
 * Keming
 * nullday
 * xing0821
 * zhyon404

## v0.2.0-alpha.4 (2022-07-05)

 * [6e9e44d](https://github.com/tensorchord/envd/commit/6e9e44dfadf13c3899707beda09d92b4f907e24d) feat: Support specify build target (#497)
 * [e443784](https://github.com/tensorchord/envd/commit/e44378470ddd029e3f2c94c93e00b0399e89b772) feat(lang): Support RStudio server (#503)
 * [89eb6e8](https://github.com/tensorchord/envd/commit/89eb6e8b5bdf795f2f1b145b75b6d87745284d71) chore(deps): bump github.com/stretchr/testify from 1.7.5 to 1.8.0 (#540)
 * [74b27e9](https://github.com/tensorchord/envd/commit/74b27e9d0039c2558b3ebd152260f0c163fc7d37) chore(deps): bump dependabot/fetch-metadata from 1.3.1 to 1.3.3 (#539)

### Contributors

 * Ce Gao
 * Jinjing Zhou
 * dependabot[bot]

## v0.2.0-alpha.3 (2022-07-01)

 * [dbac24d](https://github.com/tensorchord/envd/commit/dbac24d7931832d0fed12e931e544d31a557626d) feat(docker): Add entrypoint and ports in image config (#533)
 * [60f85f5](https://github.com/tensorchord/envd/commit/60f85f53c62597bc47f9d4328bb071b9795e0474) fix(README): Update coverage (#536)
 * [8276d7d](https://github.com/tensorchord/envd/commit/8276d7da92658a2d885a00c5039a87fb78cb2376) feat(CLI): Add push (#531)
 * [956fd73](https://github.com/tensorchord/envd/commit/956fd730aa04f1dd3cf78f699a0068436ae1d2c2) feat: Add notice for users without permission to docker daemon (#535)

### Contributors

 * Ce Gao
 * nullday

## v0.2.0-alpha.2 (2022-06-30)

 * [b987fc6](https://github.com/tensorchord/envd/commit/b987fc6841f62e4f9a29633682c6481aea6de227) fix: uid corrupted when run envd by root user (#522)
 * [310b42d](https://github.com/tensorchord/envd/commit/310b42dd450c582fb421ac183cccd7f82446156b) fix: Add talk with us in README (#526)
 * [04e8444](https://github.com/tensorchord/envd/commit/04e84440e2a7239edd4ff7d4f7e374edcb0d950e) Fix Julia multiple pkg installation bug (#521)
 * [523fb40](https://github.com/tensorchord/envd/commit/523fb400742ed5b539f44232d9eedad8eaefd13c) feat(CLI): Support context (#512)
 * [fde448a](https://github.com/tensorchord/envd/commit/fde448aff613306cb5ff3d763c06f05de12ec338) feat: add envd init (#514)
 * [8722564](https://github.com/tensorchord/envd/commit/8722564b12f0c96c58fca5d912c7c9c2e57c77a6) enhancement(CLI): Use upper case in CLI description (#515)

### Contributors

 * Aaron Sun
 * Ce Gao
 * Jinjing Zhou
 * Yunchuan Zheng
 * nullday

## v0.2.0-alpha.1 (2022-06-25)


### Contributors


## v0.2.0 (2022-09-16)

 * [4707bfe](https://github.com/tensorchord/envd/commit/4707bfeaea030a48c5601bffb87ff034e4d1b413)  fix: Fix jupyter in root (#900)
 * [71fb1ce](https://github.com/tensorchord/envd/commit/71fb1ce154ba7335de6d40e38f63d65012706c86) feat: support micromamba as an alternative to miniconda (#891)
 * [bab012c](https://github.com/tensorchord/envd/commit/bab012c5209092632d84d6cfb0c8b78fc2946523) fix: typo for git config file (#888)
 * [128f866](https://github.com/tensorchord/envd/commit/128f866f4f030cf2b10af87fe32078329e0519d8) fix(CLI): Fix build output argument and huggingface integration (#886)
 * [2e8b5d5](https://github.com/tensorchord/envd/commit/2e8b5d5d4756b5c02c6d3e846b5be093fd6394b1) fix: include update repo (#885)
 * [eb2cdd1](https://github.com/tensorchord/envd/commit/eb2cdd1a65321a1530f59190cd40560e6c31d5a3) bug: Fix detach instruction message (#882)
 * [8a02b26](https://github.com/tensorchord/envd/commit/8a02b264ac16318c5d54660ca883414ef7a15cad) refact: add envd home path func (#880)
 * [63daa5e](https://github.com/tensorchord/envd/commit/63daa5e870e5ea24c3f6e881e4932da2334688dc) chore(deps): bump github.com/spf13/viper from 1.12.0 to 1.13.0 (#875)
 * [aa53bdb](https://github.com/tensorchord/envd/commit/aa53bdb478645863dae60ca37f1ebdbc7a564c56) chore(deps): bump github.com/urfave/cli/v2 from 2.14.0 to 2.16.2 (#874)
 * [ddf3bdc](https://github.com/tensorchord/envd/commit/ddf3bdc8dd859683d7539d3c7f226b82cece40e7) chore(deps): bump actions/setup-go from 2 to 3 (#873)
 * [7220958](https://github.com/tensorchord/envd/commit/72209583b941c4eb87b282004f7c8a1ae57410ae) chore(deps): bump actions/checkout from 2 to 3 (#872)
 * [f1b3fe5](https://github.com/tensorchord/envd/commit/f1b3fe5029091cffc0a1022b665d585132c5a8d8) chore(CLI): test new release for envd-sshd (#866)
 * [49d79fb](https://github.com/tensorchord/envd/commit/49d79fb17bee4a2baaeadd500607cba7d8426b28) fix: Update readme (#865)
 * [d7995a7](https://github.com/tensorchord/envd/commit/d7995a7171cbe48a65aad3e3b56077ffee9a625a) feat(lang): io.http download files to extra_source (#858)
 * [0d3b42f](https://github.com/tensorchord/envd/commit/0d3b42fe4f33241742986030a45431a5f068dc75) feat: Support HTTP PROXY (#857)
 * [a0fbaa0](https://github.com/tensorchord/envd/commit/a0fbaa09fee056549b1d6fcd796f5b711268de2b) refact: io.mount => runtime.mount (#861)
 * [1e5e24d](https://github.com/tensorchord/envd/commit/1e5e24d1f872311a81b193443be36eaed22cc11e) bug: fix conda install with env file (#837)
 * [ecb9e26](https://github.com/tensorchord/envd/commit/ecb9e2626e65b5e1f647d7385142c10677f2d7eb) refact: unify the path env (#855)
 * [8056fda](https://github.com/tensorchord/envd/commit/8056fda28febfe6cbe64502159b302b670970517) feat: add runtime graph to image label (#815)
 * [6ad1d4c](https://github.com/tensorchord/envd/commit/6ad1d4ca2085317850e9726910bc3024b743439e) refact: apt_source, io, config mode (#853)
 * [07e2dc0](https://github.com/tensorchord/envd/commit/07e2dc0e6d4d7f3be2debf79389a792964e727b1) chore(deps): bump github.com/gliderlabs/ssh from 0.3.4 to 0.3.5 (#849)
 * [c0ba31a](https://github.com/tensorchord/envd/commit/c0ba31adc231d09600892f9445320df8eb84947b) chore(deps): bump github.com/onsi/gomega from 1.20.1 to 1.20.2 (#846)
 * [13bc9a6](https://github.com/tensorchord/envd/commit/13bc9a61c4255215661c724fb7c54f2b81642a21) chore(deps): bump github.com/urfave/cli/v2 from 2.11.2 to 2.14.0 (#845)
 * [11966f8](https://github.com/tensorchord/envd/commit/11966f8c0f997f57c45979d7857cac48ad1d1e5b) chore(deps): bump github.com/docker/go-units from 0.4.0 to 0.5.0 (#848)
 * [7ba3b3f](https://github.com/tensorchord/envd/commit/7ba3b3fa6b716efd3398a7cb3e533e915129717d) feat(cli): add msg when detach from container (#841)
 * [7fc9f34](https://github.com/tensorchord/envd/commit/7fc9f34d333b22013e847fd2c2312d18fb861068) fix: Update demo (#840)
 * [190ee76](https://github.com/tensorchord/envd/commit/190ee7635f0d9bbdee6c8b53d53a72dd8ca4e619) feat(lang): install.python_packages(local_wheels=[]) (#838)
 * [977dd47](https://github.com/tensorchord/envd/commit/977dd4725df2c16887546a48b6e3fa202a7617e2) fix: Update demo (#839)
 * [5e7c182](https://github.com/tensorchord/envd/commit/5e7c1826465511449edf6457428037cfe3afbc7e) bug: fix channels when use conda install with yaml file (#831)
 * [be02a70](https://github.com/tensorchord/envd/commit/be02a7007d65ce1227bfc61dba59de33f79c295b) fix(lang): expose host port (#832)
 * [404de31](https://github.com/tensorchord/envd/commit/404de3101cece0497084412433cf877f66cf5ee2) feat(lang): init py env by generating the bulid.envd (#827)
 * [36b1231](https://github.com/tensorchord/envd/commit/36b123142385d20fb7f7c1106c15c02f79ed4742) bug: fix permission issue when pip install from git repo (#829)
 * [1dcada4](https://github.com/tensorchord/envd/commit/1dcada4403d0c0bf8e916fc67b62c174f66df3d3) feat(build): Mount local build context into the run command (#822)
 * [f70a11c](https://github.com/tensorchord/envd/commit/f70a11c5f251b8ed0f0f42cd422b4b93efabd4a7) chore(deps): bump github.com/onsi/gomega from 1.20.0 to 1.20.1 (#821)
 * [3b440c6](https://github.com/tensorchord/envd/commit/3b440c60c59ee65444a92aa91e295be9be2125b0) chore(deps): bump github.com/moby/buildkit from 0.10.3 to 0.10.4 (#820)
 * [82fbc87](https://github.com/tensorchord/envd/commit/82fbc87ade5637fe7db8e2c0087a1555206dc1b1) doc: add include, refine others (#817)
 * [630ada1](https://github.com/tensorchord/envd/commit/630ada172bdf876c3b749329fdbe284c108051f2) feat(lang): support include other git repo for envd functions/variables (#808)
 * [5c4971b](https://github.com/tensorchord/envd/commit/5c4971b2f5fa6b32c2c247e18cf6c6a178d28f57) feat(CLI): envd env describe expose info (#801)
 * [c85766c](https://github.com/tensorchord/envd/commit/c85766cf618d55723d26427f65a385747646593d) fix: set to latest if git tag is empty (#798)
 * [4440e22](https://github.com/tensorchord/envd/commit/4440e2246108a221585853088a766a563b2c7aad) fix: add missing expose func exposed port to oci manifest (#797)
 * [3a97375](https://github.com/tensorchord/envd/commit/3a97375383a2d135cab06665aae67f04230666e1) feat(examples): Add a streamlit mnist example (#795)
 * [7bf801b](https://github.com/tensorchord/envd/commit/7bf801bad5c5947a80e576ceb3b4bae0307fddeb) feat(example): Add streamlit hello and remove bash -c in entrypoint (#794)
 * [8e17307](https://github.com/tensorchord/envd/commit/8e173075d56cadc1485f7c526c1d59c00926c69e) built: :hammer: use latest tag when not version found for cache (#793)
 * [8225eab](https://github.com/tensorchord/envd/commit/8225eab40d1a4b37b1a6d82d300399d76cfd1320) fix: use cockroachdb errors (#790)
 * [7d293d7](https://github.com/tensorchord/envd/commit/7d293d7974cae516e7d5a9f8d514acea04a0ff13) fix: panic if daemon command is invalid (#788)
 * [4d48767](https://github.com/tensorchord/envd/commit/4d48767ce033a3d60fac74d810ea546f741bc174) feat: add runtime environments (#787)
 * [248fca3](https://github.com/tensorchord/envd/commit/248fca34fa67f6fd134e1c55eefdd931aa5d8939) doc: daemon and expose (#786)
 * [e30866f](https://github.com/tensorchord/envd/commit/e30866f07c1249c5a586558a22a18741b07e063d) feat(lang): implement expose func (#780)
 * [c49863e](https://github.com/tensorchord/envd/commit/c49863e1787b6b18c08d49035ec00e79e0020822) feat(data): Add support for managed dataset and provide shortcut for common framework (#751)
 * [7c2fed6](https://github.com/tensorchord/envd/commit/7c2fed6df565e55883456b934116314109c3837d) feat(lang): add daemon function to run daemon process in the container (#777)
 * [00e8df2](https://github.com/tensorchord/envd/commit/00e8df22ba285685b11d217fdcb4cd7a52e32ba5) built: :hammer: keep v prefix in DOCKER_IMAGE_TAG env (#781)
 * [2f82fa5](https://github.com/tensorchord/envd/commit/2f82fa5ab884671ba2a5058d38d205cbd20bce1f) proposal: daemon process (#769)
 * [d06b878](https://github.com/tensorchord/envd/commit/d06b8786e2963f4bf5cc92b90fdcc80621a4bd5e) doc: update python api doc (#759)
 * [ce2e8b2](https://github.com/tensorchord/envd/commit/ce2e8b2b4e755b52892d9b3597489daae19f1dad) fix: Remove empty token arg (#772)
 * [8ab89b9](https://github.com/tensorchord/envd/commit/8ab89b967caa104ad839a83d7d9ece186eb80918) chore(deps): bump github.com/urfave/cli/v2 from 2.11.1 to 2.11.2 (#775)
 * [3a6b127](https://github.com/tensorchord/envd/commit/3a6b12772b956eb31b6ccc6c31e62efbc3feb6e3) chore(deps): bump pypa/cibuildwheel from 2.8.1 to 2.9.0 (#774)
 * [ec8cae1](https://github.com/tensorchord/envd/commit/ec8cae17b26de61bbb73026540db56505cffb2a8) fix: remove unnecessary if statement (#773)
 * [b9f0af8](https://github.com/tensorchord/envd/commit/b9f0af8056129fe6c5a2e590cd428463186cac0e) fix: -path and -file bug (#766)
 * [4a359f1](https://github.com/tensorchord/envd/commit/4a359f1eb608a0680a94379b4ac52756d605be9d) fix(release): :hammer: drop go build dep for homebrew (#768)
 * [542c7cb](https://github.com/tensorchord/envd/commit/542c7cbc3f8b7658905b85e10ff47ea216b44f7e) Fix the color display in wezterm (#767)
 * [590b4c0](https://github.com/tensorchord/envd/commit/590b4c04bffab550778a3fe911db3280ffb72b09) fix(CLI): 🔨  use latest version for local build (#763)
 * [2342367](https://github.com/tensorchord/envd/commit/2342367d3edb6d6c240afffc6b2414b1e83f0413) fix(docs): :memo: fix contributing and dev links, clean tailing space (#764)
 * [b4a8519](https://github.com/tensorchord/envd/commit/b4a851920092b66f97ed6af301472af183c986a8) fix: modify jupyter's authority from hash password to token string (#762)
 * [1f0af98](https://github.com/tensorchord/envd/commit/1f0af9825eb6f1b8843feb7fbb9b12a10e1c902d) fix: use defined jupyter port (#757)
 * [381b653](https://github.com/tensorchord/envd/commit/381b653b4dd428e5b68c1876bac74653ae2b3068) docs(README): update Documentations link to https://envd.tensorchord.… (#758)
 * [ad4b9ec](https://github.com/tensorchord/envd/commit/ad4b9ec8894b08104a47defd0039fdb081b087e3) fix: do not expose ports for custom image (#754)
 * [e77b7d6](https://github.com/tensorchord/envd/commit/e77b7d61cfaeec1d62a7119373a646a7907cb933) docs(README): Correct the cmd of get Jupyter Notebook endpoint (#756)
 * [8675316](https://github.com/tensorchord/envd/commit/86753162812f07ad030a6aecbcaa7f57839aae54) fix: ParseFromStr, add unittest (#755)
 * [382aa2d](https://github.com/tensorchord/envd/commit/382aa2db5a6bdd26d1a9f96f775cd749f9927f4e) bug(CLI): fix short alias confusion (#752)
 * [7625cd2](https://github.com/tensorchord/envd/commit/7625cd28658ad02dd05f523446985152e8eb887b) feat: Avoid gid in base image cache (#749)
 * [2e39182](https://github.com/tensorchord/envd/commit/2e39182582bf779e304fa6b7878c1428be938694) feat(CLI): :sparkles: add --force args for init to overwrite build.envd (#748)
 * [d241898](https://github.com/tensorchord/envd/commit/d2418984c7e67cc1a5184fc2894e8439e58e97ef) docs(README): One obvious way to declare supported python version (#745)
 * [b492e83](https://github.com/tensorchord/envd/commit/b492e838eb1a9bef6769806c2604d51fecfa2fff) add entrypoint in custom image (#739)
 * [4df956b](https://github.com/tensorchord/envd/commit/4df956b14ec05ac4204d8891b6791f065205b5c5) The workaround to fix the label loss (#741)
 * [dd609da](https://github.com/tensorchord/envd/commit/dd609da753b71331bfb48f997760890b47563ad4) feat: Enable build for all languages (#738)
 * [ed998ce](https://github.com/tensorchord/envd/commit/ed998cede32a4cdaa996964c6ec7ef492b4802d6) feat(lang): Add a new func runtime.command (#736)
 * [3c46efc](https://github.com/tensorchord/envd/commit/3c46efcf651af6d80010e93212ef8d2f0e26e8fa) fix: setup.py build (#735)
 * [8bef795](https://github.com/tensorchord/envd/commit/8bef795c6acd581b68492df715d67c9ed32ccc49) feat(CLI): :recycle: refactor bootstrap command to show what's envd doing (#728)
 * [8f19748](https://github.com/tensorchord/envd/commit/8f197483f9f59e0712bc0f54259331621fc05ad8) Bug: fix SIGSEGV of envd top (#726)
 * [8c83cad](https://github.com/tensorchord/envd/commit/8c83cada88d3cb0fff0eb6c19aa880c6c163cf17) feat(base-image): Move conda to llb and cache it (#724)
 * [5587c3f](https://github.com/tensorchord/envd/commit/5587c3fa6f70e465df5ba14f3cec2f99339f7bcb) feature: add envd top commands (#718)
 * [f22b283](https://github.com/tensorchord/envd/commit/f22b28328c6bd1b59bd870ef6b3a49bf0281375f) bug: allow multiple run command and use bash -c  (#720)
 * [da8feb7](https://github.com/tensorchord/envd/commit/da8feb77175448a4ce2bb2b0815a8a9f62c347b2) bug: fix notebook entry when setting conda channel (#719)
 * [116271c](https://github.com/tensorchord/envd/commit/116271c93323de18eb33fc834c971c9ca7e85b12) example: add a dgl GAT example (#714)
 * [84ffea3](https://github.com/tensorchord/envd/commit/84ffea335115d94162896123831ad5459521e2ca) feat(lang): Add proposal for expose (#568)
 * [ace70e9](https://github.com/tensorchord/envd/commit/ace70e95f4bc54cf1d58585b62d1976344d79728) feat: support io.mount (#708)
 * [a1e0395](https://github.com/tensorchord/envd/commit/a1e039595539d966916f59daacc0613b3b06bc28) feat(lang): Add default conda pkg cache (#705)
 * [c515f3c](https://github.com/tensorchord/envd/commit/c515f3cf62b2450eb36dc659974fd8875cf56f48) feat(CLI): Add category and refine help text (#707)
 * [54f412a](https://github.com/tensorchord/envd/commit/54f412a21eb47060439e3749c34d57bbfce26ba4) feat(CLI): Support run command (#701)
 * [62bd2dc](https://github.com/tensorchord/envd/commit/62bd2dcf37be2c8b811b34ea6d05189250dbdd1f) bug: Fix version prefix (#696)
 * [4dffb5d](https://github.com/tensorchord/envd/commit/4dffb5d94eaab981959fa56cfd7e5d540776f1c7) bug: Fix git tag version (#692)
 * [6daa661](https://github.com/tensorchord/envd/commit/6daa661bdcbaf3bdcbad1e63e34039e54ffe088d) bug: Add release dependency (#689)
 * [1f4b468](https://github.com/tensorchord/envd/commit/1f4b46864a03a449127ac540687571943e965102) bug: Use git version by default when build ssh (#688)
 * [728a421](https://github.com/tensorchord/envd/commit/728a421840ba10c6ba78c836a22dc7af140d370d) feat: add starship as the prompt manager (#681)
 * [6cbc53f](https://github.com/tensorchord/envd/commit/6cbc53f54c7b95c374dbc93a242e78f7be5dbd86) fix: Add environment variable PATH in run (#680)
 * [cb8ed6b](https://github.com/tensorchord/envd/commit/cb8ed6b5c6ae0a4ea77eb853bfccb450137b6b19) feat(lang): Add io.copy (#675)
 * [b3ee633](https://github.com/tensorchord/envd/commit/b3ee6334a54e7d89b3f2631631c112c1d5ecac6a) fix: Fix lint issues in conda env yaml feature PR (#679)
 * [a796a12](https://github.com/tensorchord/envd/commit/a796a129d83e5c707bb9d0c1c32295208b136495) feat(lang): Support env.yaml in conda_packages (#674)
 * [af3e78d](https://github.com/tensorchord/envd/commit/af3e78d6733f53f5d31eda045725c0af599aefa6) feature:  add label ai.tensorchord.envd.build.manifestBytecodeHash to image for cache robust (#661)
 * [3044002](https://github.com/tensorchord/envd/commit/30440022266b908adbb790e59c6696b36bf92b28) chore(deps): bump github.com/onsi/gomega from 1.19.0 to 1.20.0 (#657)
 * [400fd7f](https://github.com/tensorchord/envd/commit/400fd7f8af0a04d1d0aa5fd7f047fea59172a7ef) chore(deps): bump github.com/urfave/cli/v2 from 2.11.0 to 2.11.1 (#655)
 * [30e8caf](https://github.com/tensorchord/envd/commit/30e8cafc80b04b1bfa6e4d9d2ba273928c8b3919) chore(deps): bump actions/upload-artifact from 2 to 3 (#654)
 * [04360d6](https://github.com/tensorchord/envd/commit/04360d63ba984a1a58faaf300d0cee64269d442a) chore(deps): bump pypa/cibuildwheel from 2.8.0 to 2.8.1 (#653)
 * [de6c59e](https://github.com/tensorchord/envd/commit/de6c59e730da4520bccc747bb941bb12ced913d7) chore(deps): bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 (#656)
 * [62091a4](https://github.com/tensorchord/envd/commit/62091a460b5502d24fa57a04ad19cf273c832832) fix(build): Fix image config (#651)
 * [4c1df28](https://github.com/tensorchord/envd/commit/4c1df286df7c3d5c4a99b4bbc5958f86d08e0738) fix: context create with 'use' (#652)
 * [6f05072](https://github.com/tensorchord/envd/commit/6f05072f1f2a72b55d8e2bac468707219ceb30b4) feat(CLI): Support cache (#648)
 * [42e7531](https://github.com/tensorchord/envd/commit/42e75312a84e590c4c3b1d2f9de4ec7a6a716bb3) remove xdg, use $HOME/.config and $HOME/.cache (#641)
 * [67a1f34](https://github.com/tensorchord/envd/commit/67a1f340f45966161dbac2afb9083baeac948674) feat(lang): Remove conda from custom base image (#626)
 * [890119d](https://github.com/tensorchord/envd/commit/890119d119e7becf9f39f392aa053ae1e70c77c0) fix: check manifest and image update in new gateway buildfunc (#624)
 * [c8471db](https://github.com/tensorchord/envd/commit/c8471db20f0a7e48c2c8346b3ab88637445200f5) support buildkit TCP socket (#599)
 * [ef8a90d](https://github.com/tensorchord/envd/commit/ef8a90df3bc1a3009381eb3fd10767f468fcefc2) feat: Refactor with Builder.Options (#615)
 * [2a88ad1](https://github.com/tensorchord/envd/commit/2a88ad120b2a24b5095883a1e18de55189ff643f) chore(deps): bump github.com/urfave/cli/v2 from 2.10.3 to 2.11.0 (#610)
 * [18abe90](https://github.com/tensorchord/envd/commit/18abe90072835534f75b392b5f1ef6dbbf0bbeb5) feat(builder): Abstract BuildFunc to use gateway client (#606)
 * [178b8da](https://github.com/tensorchord/envd/commit/178b8dafdd4357688a911bfff103c397b08410a9) feat(WSL): Add ssh config entry to Windows ssh config if using WSL (#604)
 * [ceb07f5](https://github.com/tensorchord/envd/commit/ceb07f5fc43ced6da9ada711eb8c127d14afa969) fix: set conda as the only python provider (#602)
 * [f1dd546](https://github.com/tensorchord/envd/commit/f1dd546fe5598f11145a5a9c354a1e4878614ed3) fix: pre-create conda package cache directory (#600)
 * [9b3fbe3](https://github.com/tensorchord/envd/commit/9b3fbe3c9c91c167196be3764de4f53b1a074489) feat(lang): Support image in base (#595)
 * [b467279](https://github.com/tensorchord/envd/commit/b46727981968f8722feeff1023e0d85fa7bdc162) fix: Fix error handling issue (#597)
 * [fa041a8](https://github.com/tensorchord/envd/commit/fa041a849dde7f61f82fa2afdf105e45efc888ef) fix: Pre-mkdir the .cache directory of user envd (#592)
 * [54dfc52](https://github.com/tensorchord/envd/commit/54dfc52f4adc5cbce5de806dc6186e39adf21e0d) bug: fix missing function in example mnist (#589)
 * [00249bb](https://github.com/tensorchord/envd/commit/00249bbd2330088dc584015353837917bd557503) Use DefaultText in up.go (#587)
 * [302e449](https://github.com/tensorchord/envd/commit/302e4490992d5f74bf5a9e08293cfb42a89e0367) chore(deps): bump pypa/cibuildwheel from 2.7.0 to 2.8.0 (#583)
 * [6cfc0f1](https://github.com/tensorchord/envd/commit/6cfc0f16224095605dd85fb38b9cf406fbb65118) feat: Support for build image update when exec build or up again (#570)
 * [8f89e4b](https://github.com/tensorchord/envd/commit/8f89e4be3d154728824c67f13086eb727f545400) Fix: image tag normalized to docker spec (#573)
 * [3fe3757](https://github.com/tensorchord/envd/commit/3fe375769487a4eaf224a6db1437c840002c7a15) fix: add -c for every single conda channel (#569)
 * [49fa961](https://github.com/tensorchord/envd/commit/49fa961111492d9a0599e9697bf2361321e9417d) fix: add auto start buildkit container (#563)
 * [4fa5ec7](https://github.com/tensorchord/envd/commit/4fa5ec7b520964a74009b397ae9755ae96193305) bug: Fix github action (#566)
 * [93027bd](https://github.com/tensorchord/envd/commit/93027bd669ddfe052c9abcd2c2679f547389b1ff) fix: py cmd exit code (#564)
 * [707d5e8](https://github.com/tensorchord/envd/commit/707d5e8ca880a7968ef0304fab5dd28fb05a1610) feat: replace IsCreated with Exists for Client interface from package docker (#558)
 * [f71cd7f](https://github.com/tensorchord/envd/commit/f71cd7f4891fab5e1be588e9b991ab4942e02d57) feat(CLI): Unify CLI style about env and image (#550)
 * [6e9e44d](https://github.com/tensorchord/envd/commit/6e9e44dfadf13c3899707beda09d92b4f907e24d) feat: Support specify build target (#497)
 * [e443784](https://github.com/tensorchord/envd/commit/e44378470ddd029e3f2c94c93e00b0399e89b772) feat(lang): Support RStudio server (#503)
 * [89eb6e8](https://github.com/tensorchord/envd/commit/89eb6e8b5bdf795f2f1b145b75b6d87745284d71) chore(deps): bump github.com/stretchr/testify from 1.7.5 to 1.8.0 (#540)
 * [74b27e9](https://github.com/tensorchord/envd/commit/74b27e9d0039c2558b3ebd152260f0c163fc7d37) chore(deps): bump dependabot/fetch-metadata from 1.3.1 to 1.3.3 (#539)
 * [dbac24d](https://github.com/tensorchord/envd/commit/dbac24d7931832d0fed12e931e544d31a557626d) feat(docker): Add entrypoint and ports in image config (#533)
 * [60f85f5](https://github.com/tensorchord/envd/commit/60f85f53c62597bc47f9d4328bb071b9795e0474) fix(README): Update coverage (#536)
 * [8276d7d](https://github.com/tensorchord/envd/commit/8276d7da92658a2d885a00c5039a87fb78cb2376) feat(CLI): Add push (#531)
 * [956fd73](https://github.com/tensorchord/envd/commit/956fd730aa04f1dd3cf78f699a0068436ae1d2c2) feat: Add notice for users without permission to docker daemon (#535)
 * [b987fc6](https://github.com/tensorchord/envd/commit/b987fc6841f62e4f9a29633682c6481aea6de227) fix: uid corrupted when run envd by root user (#522)
 * [310b42d](https://github.com/tensorchord/envd/commit/310b42dd450c582fb421ac183cccd7f82446156b) fix: Add talk with us in README (#526)
 * [04e8444](https://github.com/tensorchord/envd/commit/04e84440e2a7239edd4ff7d4f7e374edcb0d950e) Fix Julia multiple pkg installation bug (#521)
 * [523fb40](https://github.com/tensorchord/envd/commit/523fb400742ed5b539f44232d9eedad8eaefd13c) feat(CLI): Support context (#512)
 * [fde448a](https://github.com/tensorchord/envd/commit/fde448aff613306cb5ff3d763c06f05de12ec338) feat: add envd init (#514)
 * [8722564](https://github.com/tensorchord/envd/commit/8722564b12f0c96c58fca5d912c7c9c2e57c77a6) enhancement(CLI): Use upper case in CLI description (#515)
 * [0ae8df9](https://github.com/tensorchord/envd/commit/0ae8df9e9eb4ec0c755920a7a3697cbab12b22e3) chore(deps): bump github.com/stretchr/testify from 1.7.2 to 1.7.5 (#505)
 * [1cd5a27](https://github.com/tensorchord/envd/commit/1cd5a27027dd696276fa75e1d52a02858b7e7de8) chore(deps): bump github.com/urfave/cli/v2 from 2.8.1 to 2.10.3 (#504)
 * [ca6435d](https://github.com/tensorchord/envd/commit/ca6435ddfd29f5be2e85a1039da50c4db2fea03f) feat(lang): Support julia (#495)
 * [220d874](https://github.com/tensorchord/envd/commit/220d874712d7f7120f16337324d23b440163e3f7) feat(ssh): Config ssh key permanently and globally (#487)
 * [7127365](https://github.com/tensorchord/envd/commit/7127365e47c6a2923ce93ecf0164df49810b13fc) feat(lang): Support R language (#491)
 * [3f86086](https://github.com/tensorchord/envd/commit/3f860868426e8a96c3daf25962f79d964f57a8b1) feat(lang): Support python requirements.txt in python_packages (#484)
 * [b31f4cd](https://github.com/tensorchord/envd/commit/b31f4cd5e98d1982eb92d71759a9ecd7df36c72d) fix: correct comments and unify receiver names and simplify some code (#477)
 * [35c6b76](https://github.com/tensorchord/envd/commit/35c6b7606803485eeaae4f771b4a3a60060f853c) feat #246: envd up a GPU image without GPUs  (#474)
 * [ff43aa3](https://github.com/tensorchord/envd/commit/ff43aa35236171d3a497f6556bcd3e7ab30c6164) Update build.go (#472)
 * [229821a](https://github.com/tensorchord/envd/commit/229821a0634c20c1b05d4f8c81e4dbf7b81d1e35) feat(CLI): Add '--detail' for detail version information (#283) (#471)
 * [6bd547e](https://github.com/tensorchord/envd/commit/6bd547ea3f3132d1b82634560de265644e931637) feat(CLI): Add go API to support LSP (#302)
 * [45fcbe3](https://github.com/tensorchord/envd/commit/45fcbe39171085414c1f08b579d238760b869a2a) fix typos (#468)
 * [f2ff5b7](https://github.com/tensorchord/envd/commit/f2ff5b749af6c25e695676e40cd16e6accb21325) feat(CLI): Prune cache (#464)
 * [5d21bd3](https://github.com/tensorchord/envd/commit/5d21bd383ead10d410daa7e67e9b23b30c830cce) chore(Makefile): add default goal as `build` (#465)
 * [df4a395](https://github.com/tensorchord/envd/commit/df4a395ad48d3f45c77c2ec7de419d74561d7c7e) feat: Print out container info when wait timeout (#460)
 * [657c2cb](https://github.com/tensorchord/envd/commit/657c2cb1226020eb754510ffa9876d353d4baefe) feat: Add base image for R language (#457)
 * [3c0afb3](https://github.com/tensorchord/envd/commit/3c0afb3b9bbf35a5ce211f3b19d5bf226f1028ae) fix: enable release with Homebrew only on stable versions (#455)

### Contributors

 * Aaron Sun
 * Aka.Fido
 * Bingyi Sun
 * Ce Gao
 * Guangyang Li
 * Gui-Yue
 * Haiker Sun
 * Jinjing Zhou
 * Keming
 * Wei Zhang
 * Yuan Tang
 * Yuchen Cheng
 * Yuedong Wu
 * Yunchuan Zheng
 * Zhenguo.Li
 * Zhenzhen Zhao
 * Zhizhen He
 * dependabot[bot]
 * kenwoodjw
 * nullday
 * wyq
 * xing0821
 * zhyon404

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


### Contributors


## v0.1.0 (2022-06-18)

 * [3abef45](https://github.com/tensorchord/envd/commit/3abef452fb45bdcdbe4291caeae1ebd1a12589e4) fix: Fix the bug about uid (#447)
 * [eff6ffa](https://github.com/tensorchord/envd/commit/eff6ffac3dd6d7f0ffd0313dc4eec06eb753d4de) fix: Fix typo (#445)
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
 * [8531491](https://github.com/tensorchord/envd/commit/853149189d88e6ecf6a5a924e3aa19d8f7993f7e) fix: Fix default ssh shell (#411)
 * [5f3b16b](https://github.com/tensorchord/envd/commit/5f3b16bf579fe372ec2063bf3ce8904a4e5d2e2b) feat: Support configuring CRAN mirror for R environment (#405)
 * [8e27e99](https://github.com/tensorchord/envd/commit/8e27e9962588bd1b14268eaed6f6e067a9c5d908) fix: Only configure conda for Python environment (#406)
 * [6160899](https://github.com/tensorchord/envd/commit/6160899a10ec941f178dd821bee42eb615b932e3) feat(cli): support `envd build --output` (#402)
 * [d3fda6d](https://github.com/tensorchord/envd/commit/d3fda6db2783c93f4f5ea7954627500431059558) fix: Hack the gid (#399)
 * [e478c1a](https://github.com/tensorchord/envd/commit/e478c1a51191ab74f5c733d79106520032240e7f) bug: Fix source is released twice for macos and linux (#394)
 * [a55ad88](https://github.com/tensorchord/envd/commit/a55ad8808e5a23dba628f5a55c335603771083e9) feat(lang): Set default user to current (#390)
 * [df5bde3](https://github.com/tensorchord/envd/commit/df5bde376bcdd7db84ccd50b132cf05c3546deb9) feat(release): Support Homebrew in goreleaser (#389)
 * [7c10b71](https://github.com/tensorchord/envd/commit/7c10b71266c1b4ed250dc1a850603e930f531cc2) fix: cannot assign requested address (#386)
 * [d48e3ab](https://github.com/tensorchord/envd/commit/d48e3abf0b2443e1740dd0c2ddd5ebcf81e6fa6f) fix: Output error details when debug flag is enabled (#385)
 * [e2c8adb](https://github.com/tensorchord/envd/commit/e2c8adb41eb394d7fed56d5494f1ed5fc0832356) fix: use python3 explicitly to avoid type hints error (#379)
 * [fc2afe9](https://github.com/tensorchord/envd/commit/fc2afe9b6fe2bf9e9430039a0270de45981ba5ac) fix: add classifiers in setup.py (#380)
 * [33bdd7a](https://github.com/tensorchord/envd/commit/33bdd7a3313d18a9e781db98b033f5a7ceffe58b) doc: Add universe api doc (#374)
 * [3b3945a](https://github.com/tensorchord/envd/commit/3b3945aea9297fc69a8f9c787ef27812f1d0efb9) fix: Add v before tags (#371)
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
 * [12cf334](https://github.com/tensorchord/envd/commit/12cf3345b6c09106508271dd84bd41bc03ceedbc) fix: Fix twine (#301)
 * [f42e162](https://github.com/tensorchord/envd/commit/f42e1625fd11332411b821931d0664494bfc1927) fix: Instal twine (#300)
 * [7720529](https://github.com/tensorchord/envd/commit/7720529500b58d854c09c817485f0edc2a1198dc) feat(lang): Support config.conda_channel and install.conda_packages (#293)
 * [452f3dc](https://github.com/tensorchord/envd/commit/452f3dc8033d1d2163aed7bc61bcd8d54ad81aec) feat: add pypi sdist (#298)
 * [cfe65fe](https://github.com/tensorchord/envd/commit/cfe65fe690f050a8e759063c3d6a3f71aa051f05) fix: py27 subprocess (#296)
 * [5cf52ec](https://github.com/tensorchord/envd/commit/5cf52ec7f03f94394aa17140986cb85545dfe942) fix: Update readme about installation (#295)
 * [3bf2710](https://github.com/tensorchord/envd/commit/3bf27107eb3bf01916ec703b9ae697dd87a92ad7) action: Add pypi release pipeline (#277)
 * [35e6e1b](https://github.com/tensorchord/envd/commit/35e6e1baa1af2c8b25ba47f88949f049b996d5d9) workflow: Enable macOS in CI (#287)
 * [2bc13df](https://github.com/tensorchord/envd/commit/2bc13df54f4ce9150db0519986344e781c3e5f32) bug: fix version without tag (#288)
 * [ef3886c](https://github.com/tensorchord/envd/commit/ef3886c7730540104737e207970afcb0b3876c2a) Revert "workflow: enable macOS in CI (#280)" (#286)
 * [02f83aa](https://github.com/tensorchord/envd/commit/02f83aa6bb7d1944988ab2e62e328d6cfcd3ff77) workflow: enable macOS in CI (#280)
 * [2cd9a0e](https://github.com/tensorchord/envd/commit/2cd9a0e3fa0191d801e8c8c3d6e4a64244c58861) fix: Update contributing (#284)
 * [f166cf8](https://github.com/tensorchord/envd/commit/f166cf8059cf658c59c51f155402ddf5eef9a922) feat: Support destroy environment by name (#281)
 * [9c237e4](https://github.com/tensorchord/envd/commit/9c237e4cad996be5baf1116d8794708bf598c893) fix: Bump version and fix base image (#279)
 * [e048fc0](https://github.com/tensorchord/envd/commit/e048fc06c3f1b5dd2d0f69aded389a67a4608ced) fix: Use 127.0.0.1 instead of containerIP in ssh (#276)
 * [ae16402](https://github.com/tensorchord/envd/commit/ae16402016f27505d100d820d6a2aeba1aa9838a) fix: Hard code OS (#270)
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

 * Aaron Sun
 * Ce Gao
 * Jinjing Zhou
 * Jun
 * Keming
 * Kevin Su
 * Ling Jin
 * Manjusaka
 * Xu Jin
 * Yuan Tang
 * Yuchen Cheng
 * Zhenzhen Zhao
 * Zhizhen He
 * dependabot[bot]
 * kenwoodjw

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


