<p align="center">
  <img height="125" alt="Fiber Utils" src="https://raw.githubusercontent.com/gofiber/utils/master/.github/logo-dark.svg#gh-dark-mode-only" />
  <img height="125" alt="Fiber Utils" src="https://raw.githubusercontent.com/gofiber/utils/master/.github/logo.svg#gh-light-mode-only" />
</p>

# Fiber Utils

![Release](https://img.shields.io/github/release/gofiber/utils.svg)
![Test](https://github.com/gofiber/utils/workflows/Test/badge.svg)
![Codecov](https://img.shields.io/codecov/c/github/gofiber/utils?token=3Cr92CwaPQ&style=flat-square&logo=codecov&label=codecov)
![Linter](https://github.com/gofiber/utils/actions/workflows/lint.yml/badge.svg)
[![Benchmarks](https://img.shields.io/badge/%F0%9F%93%8A%20benchmarks-charts-00ACD7.svg)](https://gofiber.github.io/utils/benchmarks/)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)

A collection of common functions for [Fiber](https://github.com/gofiber/fiber) with better performance, fewer allocations, and fewer dependencies.

## Benchmarks

Per-commit benchmark charts: <https://gofiber.github.io/utils/benchmarks/>

Each table below starts with the `go test` command that produced it. Maintainers refresh the
numbers with `make benchupdate`, or by commenting `/bench-readme` (arm64) or
`/bench-readme-amd64` on a pull request.

Environment:
goos: darwin
goarch: arm64
pkg: github.com/gofiber/utils/v2
cpu: Apple M2 Pro

```text
// go test ./... -benchmem -run=^$ -bench=Benchmark_ -count=1
// (swar's Benchmark_Load8_Fusion, an advisory codegen guard, is deliberately omitted)

# Case Conversion
Benchmark_ToLowerBytes/empty/fiber-12                              599341221    1.970  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/fiber/unsafe-12                       613673802    1.951  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/empty/default-12                            363300402    3.316  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber-12                           350302147    3.434  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/fiber/unsafe-12                    250488015    4.803  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get/default-12                          65883986    18.03  ns/op     8  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber-12                      95234625    12.53  ns/op     3  B/op   1  allocs/op
Benchmark_ToLowerBytes/http-get-upper/fiber/unsafe-12              242418936    4.802  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/http-get-upper/default-12                    86405012    13.85  ns/op     8  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber-12           52263289    22.74  ns/op    48  B/op   1  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/fiber/unsafe-12   131711175    9.101  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/header-content-type-mixed/default-12         25312339    47.13  ns/op    48  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber-12                        120148472    10.02  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/fiber/unsafe-12                 100000000    10.57  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-lower/default-12                       16746064    71.43  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber-12                         46270634    26.82  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-upper/fiber/unsafe-12                 100000000    10.71  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-upper/default-12                       15734505    76.92  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber-12                         46607823    26.06  ns/op    64  B/op   1  allocs/op
Benchmark_ToLowerBytes/large-mixed/fiber/unsafe-12                 100000000    10.86  ns/op     0  B/op   0  allocs/op
Benchmark_ToLowerBytes/large-mixed/default-12                       16291084    73.76  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/empty/fiber-12                              607678140    1.976  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/fiber/unsafe-12                       611928006    1.981  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/empty/default-12                            362632894    3.309  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/fiber-12                            95358867    12.55  ns/op     3  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get/fiber/unsafe-12                    249490471    4.905  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get/default-12                          79248020    14.63  ns/op     8  B/op   1  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber-12                     365516280    3.330  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/fiber/unsafe-12              249981619    4.803  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/http-get-upper/default-12                    64343162    19.01  ns/op     8  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber-12           52506741    22.69  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/fiber/unsafe-12   130478352    9.350  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/header-content-type-mixed/default-12         22415607    54.44  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber-12                         44765679    26.92  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-lower/fiber/unsafe-12                 100000000    13.68  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-lower/default-12                       14464340    110.7  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber-12                        100000000    11.98  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/fiber/unsafe-12                  83808488    13.94  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-upper/default-12                       13142096    87.35  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber-12                         33038661    32.75  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpperBytes/large-mixed/fiber/unsafe-12                 100000000    10.93  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpperBytes/large-mixed/default-12                       15034422    98.27  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/empty/fiber-12                                   595298628    2.309  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/empty/fiber/unsafe-12                            595273404    2.021  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/empty/default-12                                 636980288    1.883  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/fiber-12                                 97025577    12.40  ns/op     3  B/op   1  allocs/op
Benchmark_ToUpper/http-get/fiber/unsafe-12                         244711635    4.857  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get/default-12                               59790234    19.54  ns/op     8  B/op   1  allocs/op
Benchmark_ToUpper/http-get-upper/fiber-12                          383218539    3.056  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/fiber/unsafe-12                   247851331    4.864  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/http-get-upper/default-12                        223044249    5.182  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber-12                52628502    24.13  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpper/header-content-type-mixed/fiber/unsafe-12        127726341    9.233  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/header-content-type-mixed/default-12              11549816    104.0  ns/op    48  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber-12                              42409070    26.04  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/large-lower/fiber/unsafe-12                      100000000    11.00  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-lower/default-12                             7051911    169.0  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/large-upper/fiber-12                             123837582    9.590  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/fiber/unsafe-12                      100000000    10.83  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-upper/default-12                            19115356    61.51  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/fiber-12                              45872143    27.77  ns/op    64  B/op   1  allocs/op
Benchmark_ToUpper/large-mixed/fiber/unsafe-12                      100000000    11.00  ns/op     0  B/op   0  allocs/op
Benchmark_ToUpper/large-mixed/default-12                             5927710    203.1  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/empty/fiber-12                                   613625554    1.952  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/empty/fiber/unsafe-12                            572670608    1.989  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/empty/default-12                                 644769141    1.863  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber-12                                395302272    3.034  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get/fiber/unsafe-12                         252231990    4.813  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get/default-12                              284202345    4.219  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/fiber-12                          100000000    12.08  ns/op     3  B/op   1  allocs/op
Benchmark_ToLower/http-get-upper/fiber/unsafe-12                   252859507    4.816  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/http-get-upper/default-12                         60415860    19.54  ns/op     8  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber-12                53446660    22.49  ns/op    48  B/op   1  allocs/op
Benchmark_ToLower/header-content-type-mixed/fiber/unsafe-12        131152108    9.153  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/header-content-type-mixed/default-12              18039788    66.79  ns/op    48  B/op   1  allocs/op
Benchmark_ToLower/large-lower/fiber-12                             120227584    9.982  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/fiber/unsafe-12                      100000000    10.68  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-lower/default-12                            24070928    48.79  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/fiber-12                              44699118    25.92  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/large-upper/fiber/unsafe-12                      100000000    10.74  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-upper/default-12                             7075200    170.0  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber-12                              46426453    26.01  ns/op    64  B/op   1  allocs/op
Benchmark_ToLower/large-mixed/fiber/unsafe-12                      100000000    10.74  ns/op     0  B/op   0  allocs/op
Benchmark_ToLower/large-mixed/default-12                             5682778    204.3  ns/op    64  B/op   1  allocs/op

# Add Trailing Slash
Benchmark_AddTrailingSlashBytes/empty-12                          1000000000   0.5857  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/slash-only-12                     1000000000   0.8694  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/short-no-slash-12                  100000000    11.65  ns/op     4  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/short-with-slash-12               1000000000   0.8867  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashBytes/path-no-slash-12                    93543546    13.36  ns/op    16  B/op   1  allocs/op
Benchmark_AddTrailingSlashBytes/path-with-slash-12                1000000000   0.9011  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/empty-12                         1000000000   0.2926  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/slash-only-12                    1000000000   0.4482  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/short-no-slash-12                 100000000    11.37  ns/op     4  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/short-with-slash-12              1000000000   0.4502  ns/op     0  B/op   0  allocs/op
Benchmark_AddTrailingSlashString/path-no-slash-12                   82427726    14.41  ns/op    16  B/op   1  allocs/op
Benchmark_AddTrailingSlashString/path-with-slash-12               1000000000   0.4529  ns/op     0  B/op   0  allocs/op

# EqualFold
Benchmark_EqualFoldBytes/fiber-12                                   67368261    18.48  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-12                                 17774803    65.91  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold/fiber-12                                        80501566    14.75  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold/default-12                                      18314365    65.96  ns/op     0  B/op   0  allocs/op

# Search
// amd64 note: IndexAny2/IndexAny3 dispatch to package simd for 32+ byte
// inputs; these arm64 rows keep the portable SWAR path. For amd64 numbers
// see the Benchmark_Memchr2/Benchmark_Memchr3 rows in the simd section
// below (the kernels these functions dispatch to).
Benchmark_IndexAny2/7B/swar-12                                     282163779    4.297  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/7B/scalar-12                                   314553402    3.826  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/7B/stdlib-indexany-12                           50085909    23.97  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/8B/swar-12                                     522780714    2.289  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/8B/scalar-12                                   271494799    4.206  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/8B/stdlib-indexany-12                           44831599    26.78  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/16B/swar-12                                    393245086    3.047  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/16B/scalar-12                                  168476984    7.117  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/16B/stdlib-indexany-12                         100000000    10.03  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/32B/swar-12                                    252199587    4.802  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/32B/scalar-12                                   84105762    14.27  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/32B/stdlib-indexany-12                          70291048    17.37  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/64B/swar-12                                    144335962    8.290  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/64B/scalar-12                                   42694686    28.60  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/64B/stdlib-indexany-12                          37615540    32.23  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/512B/swar-12                                    20231538    59.71  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/512B/scalar-12                                   5071586    233.6  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny2/512B/stdlib-indexany-12                          4448119    274.2  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/7B/swar-12                                     221889592    5.397  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/7B/scalar-12                                   240859849    4.977  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/8B/swar-12                                     445742876    2.689  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/8B/scalar-12                                   211981612    5.620  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/16B/swar-12                                    326475568    3.739  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/16B/scalar-12                                  100000000    10.31  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/32B/swar-12                                    196365920    6.112  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/32B/scalar-12                                   60021255    19.81  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/64B/swar-12                                    100000000    10.98  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/64B/scalar-12                                   31180405    38.97  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/512B/swar-12                                    15466316    79.22  ns/op     0  B/op   0  allocs/op
Benchmark_IndexAny3/512B/scalar-12                                   3911382    307.1  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-8B/swar-12                                332497509    3.634  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-8B/scalar-12                              491493494    2.408  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-32B/swar-12                               149769112    7.928  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-32B/scalar-12                              44041342    27.13  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-32B/swar-12                                138732252    8.758  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-32B/scalar-12                               33871633    30.73  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-tail/swar-12                               100000000    12.34  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/hit-tail/scalar-12                              16138168    75.39  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-64B/swar-12                               100000000    11.37  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-64B/scalar-12                              17902330    67.78  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-512B/swar-12                               17640939    67.65  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/miss-512B/scalar-12                              2211582    544.4  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-hit-64B/swar-12                           100000000    11.58  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-hit-64B/scalar-12                          23906326    51.33  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-miss-64B/swar-12                          100000000    11.33  ns/op     0  B/op   0  allocs/op
Benchmark_IndexFold/long-miss-64B/scalar-12                         23553600    50.27  ns/op     0  B/op   0  allocs/op

# ASCII
// amd64 note: IsASCII dispatches to package simd for 32+ byte inputs;
// these arm64 rows keep the portable SWAR path. See the simd section
// below for amd64 numbers.
Benchmark_IsASCII/7B/swar-12                                       238699563    5.001  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/7B/scalar-12                                     292669958    4.064  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/swar-12                                       439961722    2.687  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/scalar-12                                     274732192    4.364  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/16B/swar-12                                      367179319    3.255  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/16B/scalar-12                                    179331987    6.706  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/swar-12                                      368345228    3.271  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/scalar-12                                    100000000    11.62  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/swar-12                                      265505815    4.486  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/scalar-12                                     57864440    21.33  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/swar-12                                      47664915    25.14  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/scalar-12                                     7081790    169.5  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/7B/swar-12                              142234436    8.424  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/7B/scalar-12                            142394628    8.174  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/8B/swar-12                              401721262    3.099  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/8B/scalar-12                            138047568    8.654  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/16B/swar-12                             244592613    5.059  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/16B/scalar-12                            75413396    15.94  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/32B/swar-12                             138857455    8.875  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/32B/scalar-12                            39539425    30.18  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/64B/swar-12                              75193269    16.46  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/64B/scalar-12                            20898807    57.81  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/512B/swar-12                              9634008    124.9  ns/op     0  B/op   0  allocs/op
Benchmark_IndexNonQuotable/512B/scalar-12                            2619062    460.6  ns/op     0  B/op   0  allocs/op

# Trim
Benchmark_TrimRight/fiber-12                                       587360972    2.043  ns/op     0  B/op   0  allocs/op
Benchmark_TrimRight/default-12                                     417142952    2.887  ns/op     0  B/op   0  allocs/op
Benchmark_TrimRightBytes/fiber-12                                  590828253    2.067  ns/op     0  B/op   0  allocs/op
Benchmark_TrimRightBytes/default-12                                374548395    3.236  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeft/fiber-12                                        590822072    2.030  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeft/default-12                                      408184784    2.967  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/fiber-12                                   592875732    2.033  ns/op     0  B/op   0  allocs/op
Benchmark_TrimLeftBytes/default-12                                 405921434    2.967  ns/op     0  B/op   0  allocs/op
Benchmark_Trim/fiber-12                                            317625432    3.758  ns/op     0  B/op   0  allocs/op
Benchmark_Trim/default-12                                          250173384    4.735  ns/op     0  B/op   0  allocs/op
Benchmark_Trim/default.trimspace-12                                241883282    5.022  ns/op     0  B/op   0  allocs/op
Benchmark_TrimBytes/fiber-12                                       308885211    3.858  ns/op     0  B/op   0  allocs/op
Benchmark_TrimBytes/default-12                                     254064188    4.802  ns/op     0  B/op   0  allocs/op
Benchmark_TrimBytes/default.trimspace-12                           224231778    5.330  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/empty-12                                1000000000   0.3017  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/empty-12                               583064172    2.074  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/spaces-12                                387354229    3.059  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/spaces-12                              320059388    3.651  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/ascii-word-12                            345078337    3.535  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/ascii-word-12                          270603267    4.465  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-bearer-12                    291494109    4.141  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/auth-header-bearer-12                  224087944    5.329  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/auth-header-basic-12                     390765345    3.047  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/auth-header-basic-12                   258943531    4.718  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/accept-encoding-12                       344369745    3.490  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/accept-encoding-12                     277203638    4.347  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/content-type-12                          344263705    3.506  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/content-type-12                        276213367    4.343  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/x-forwarded-for-12                       460119685    2.653  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/x-forwarded-for-12                     296468840    4.058  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/query-params-12                          330053276    3.520  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/query-params-12                        275688495    4.352  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/ascii-long-12                            275511898    4.342  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/ascii-long-12                          206143018    5.807  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/no-trim-12                              1000000000   0.5914  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/no-trim-12                             451842741    2.634  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/fiber/mixed-whitespace-12                      296083767    4.054  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpace/default/mixed-whitespace-12                    224007750    5.294  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/empty-12                           1000000000   0.2914  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/empty-12                          519488766    2.352  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/spaces-12                           410959255    2.966  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/spaces-12                         315602949    3.827  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-word-12                       341217429    3.519  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-word-12                     247799318    4.667  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-bearer-12               295130076    4.182  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-bearer-12             209162235    5.685  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/auth-header-basic-12                393366796    3.053  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/auth-header-basic-12              240882009    4.930  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/accept-encoding-12                  343673244    3.512  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/accept-encoding-12                250962042    4.642  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/content-type-12                     344382758    3.516  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/content-type-12                   250622116    4.693  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/x-forwarded-for-12                  459435145    2.621  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/x-forwarded-for-12                276595178    4.359  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/query-params-12                     341144192    3.526  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/query-params-12                   255843174    4.873  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/ascii-long-12                       270183530    4.445  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/ascii-long-12                     189740618    6.335  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/no-trim-12                         1000000000   0.5969  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/no-trim-12                        406072189    2.941  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/fiber/mixed-whitespace-12                 289382116    4.057  ns/op     0  B/op   0  allocs/op
Benchmark_TrimSpaceBytes/default/mixed-whitespace-12               218559462    5.493  ns/op     0  B/op   0  allocs/op

# Convert
Benchmark_ConvertToBytes/fiber-12                                  320332576    3.787  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int-12                                          499429904    2.382  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int8-12                                         614406423    1.965  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int16-12                                        486644316    2.463  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int32-12                                        495227163    2.422  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/int64-12                                        502362236    2.429  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint-12                                         561693333    2.101  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint8-12                                        611239028    1.965  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint16-12                                       573383422    2.127  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint32-12                                       575723191    2.093  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/uint64-12                                       575901750    2.089  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/string-12                                       613909396    1.959  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/[]uint8-12                                       96832432    12.11  ns/op    16  B/op   1  allocs/op
Benchmark_ToString/bool-12                                         608058547    1.990  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/float32-12                                       23195016    47.11  ns/op     4  B/op   1  allocs/op
Benchmark_ToString/float64-12                                       17204126    70.24  ns/op     4  B/op   1  allocs/op
Benchmark_ToString/time.Time-12                                     13103935    91.85  ns/op    24  B/op   1  allocs/op
Benchmark_ToString/time.Time#01-12                                  13198113    91.22  ns/op    24  B/op   1  allocs/op
Benchmark_ToString/[]string-12                                      65187865    18.35  ns/op    16  B/op   1  allocs/op
Benchmark_ToString/[]int-12                                         54417859    21.53  ns/op     8  B/op   1  allocs/op
Benchmark_ToString/[2]int-12                                        32691942    36.57  ns/op    16  B/op   1  allocs/op
Benchmark_ToString/[][]int-12                                        9382754    128.8  ns/op   112  B/op   6  allocs/op
Benchmark_ToString/[]interface_{}-12                                10584483    114.5  ns/op    72  B/op   3  allocs/op
Benchmark_ToString/utils.MyStringer-12                             365979882    3.277  ns/op     0  B/op   0  allocs/op
Benchmark_ToString/utils.CustomType-12                              11649330    102.5  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/int-12                             1000000000   0.2619  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int8-12                            1000000000   0.2202  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int16-12                           1000000000   0.2708  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int32-12                           1000000000   0.2658  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/int64-12                           1000000000   0.2658  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint-12                            1000000000   0.2281  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint8-12                           1000000000   0.2223  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint16-12                          1000000000   0.2396  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint32-12                          1000000000   0.2273  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/uint64-12                          1000000000   0.2271  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/string-12                          1000000000   0.2117  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/[]uint8-12                          213842515    5.718  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/bool-12                            1000000000   0.2118  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/float32-12                          201081819    6.011  ns/op     4  B/op   1  allocs/op
Benchmark_ToString_concurrency/float64-12                          142000966    8.799  ns/op     4  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time-12                         71451783    18.04  ns/op    24  B/op   1  allocs/op
Benchmark_ToString_concurrency/time.Time#01-12                      69313578    17.74  ns/op    24  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]string-12                         187564670    6.455  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[]int-12                            262707805    4.498  ns/op     8  B/op   1  allocs/op
Benchmark_ToString_concurrency/[2]int-12                           129274200    8.777  ns/op    16  B/op   1  allocs/op
Benchmark_ToString_concurrency/[][]int-12                           22638411    52.65  ns/op   112  B/op   6  allocs/op
Benchmark_ToString_concurrency/[]interface_{}-12                    34120837    35.36  ns/op    72  B/op   3  allocs/op
Benchmark_ToString_concurrency/utils.MyStringer-12                1000000000   0.4259  ns/op     0  B/op   0  allocs/op
Benchmark_ToString_concurrency/utils.CustomType-12                  54806094    25.95  ns/op    16  B/op   1  allocs/op
Benchmark_UnsafeBytes/unsafe-12                                   1000000000   0.5658  ns/op     0  B/op   0  allocs/op
Benchmark_UnsafeBytes/default-12                                    88407558    13.50  ns/op    16  B/op   1  allocs/op
Benchmark_UnsafeString/unsafe-12                                  1000000000   0.3545  ns/op     0  B/op   0  allocs/op
Benchmark_UnsafeString/default-12                                  100000000    12.09  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/0-12                                            395954607    3.023  ns/op     0  B/op   0  allocs/op
Benchmark_ByteSize/1-12                                             65295756    18.15  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/500-12                                           58264801    20.92  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1024-12                                          60886531    19.80  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1126-12                                          53409986    22.69  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1048576-12                                       60969927    19.41  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1153024-12                                       53647869    22.12  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-12                                    63792180    19.05  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-12                                    54560540    22.02  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-12                                 64080169    18.85  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-12                                 52957434    22.01  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-12                              65629659    18.33  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-12                              56535770    21.73  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-12                           64569356    18.56  ns/op    16  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-12                           56481996    21.10  ns/op    16  B/op   1  allocs/op

# Format and Append
Benchmark_FormatUint/small/fiber-12                                607527520    1.966  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-12                              607521366    1.973  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-12                                64427509    19.04  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint/medium/strconv-12                              59566404    20.42  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-12                                 45395149    26.16  ns/op    24  B/op   1  allocs/op
Benchmark_FormatUint/large/strconv-12                               45620366    26.38  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/small_pos/fiber-12                             608898928    1.963  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-12                           598721480    1.976  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-12                             616958124    1.946  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-12                            87181050    13.77  ns/op     3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-12                             63105722    19.12  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-12                           57454753    20.70  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-12                             63631222    19.10  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-12                           61296942    20.17  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-12                              48604558    25.08  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-12                            45605844    26.25  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-12                              47388217    25.09  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-12                            45871998    26.10  ns/op    24  B/op   1  allocs/op
Benchmark_FormatUint32/fiber-12                                     63452196    18.63  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint32/strconv-12                                   59485942    20.08  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt32/fiber-12                                      63030724    18.98  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-12                                    61887302    19.69  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint16/fiber-12                                     80149611    14.48  ns/op     5  B/op   1  allocs/op
Benchmark_FormatUint16/strconv-12                                   77022440    15.43  ns/op     5  B/op   1  allocs/op
Benchmark_FormatInt16/fiber-12                                      79013646    15.23  ns/op     8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-12                                    79695165    15.24  ns/op     8  B/op   1  allocs/op
Benchmark_FormatUint8/fiber-12                                     599811685    1.947  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint8/strconv-12                                    84144096    14.18  ns/op     3  B/op   1  allocs/op
Benchmark_FormatInt8/fiber-12                                      630717244    1.985  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt8/strconv-12                                     83826201    14.24  ns/op     4  B/op   1  allocs/op
Benchmark_AppendUint/fiber-12                                      130817422    9.139  ns/op     0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-12                                    100000000    10.08  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/fiber-12                             369512403    3.266  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-12                           213341739    5.627  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-12                            125937435    9.531  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-12                          122490632    9.791  ns/op     0  B/op   0  allocs/op

# Token
Benchmark_GenerateSecureToken/16_bytes-12                            4565527    258.6  ns/op    24  B/op   1  allocs/op
Benchmark_GenerateSecureToken/32_bytes-12                            4306404    280.8  ns/op    48  B/op   1  allocs/op
Benchmark_TokenGenerators/UUIDv4-12                                  4019139    288.2  ns/op    64  B/op   2  allocs/op
Benchmark_TokenGenerators/SecureToken-12                             4360036    276.9  ns/op    48  B/op   1  allocs/op

# HTTP
Benchmark_GetMIME/fiber-12                                          22321272    53.90  ns/op     0  B/op   0  allocs/op
Benchmark_GetMIME/default-12                                        17492530    68.46  ns/op     0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/vendorContentType-12      125467813    9.619  ns/op     0  B/op   0  allocs/op
Benchmark_ParseVendorSpecificContentType/defaultContentType-12     404169172    2.974  ns/op     0  B/op   0  allocs/op
Benchmark_StatusMessage/fiber-12                                  1000000000   0.3349  ns/op     0  B/op   0  allocs/op
Benchmark_StatusMessage/default-12                                 460456902    2.610  ns/op     0  B/op   0  allocs/op

# IP
Benchmark_IsIPv4/fiber-4                                            56396187    21.72  ns/op     0  B/op   0  allocs/op
Benchmark_IsIPv4/default-4                                          28357074    40.38  ns/op     0  B/op   0  allocs/op
Benchmark_IsIPv6/fiber-4                                            19461662    62.01  ns/op     0  B/op   0  allocs/op
Benchmark_IsIPv6/default-4                                          13315670    89.84  ns/op     0  B/op   0  allocs/op

# Parse
Benchmark_ParseUint/fiber-12                                       236748204    5.049  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint/fiber_bytes-12                                 224650773    5.348  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint/default-12                                      88588122    13.75  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/6digit-12                            207540178    5.769  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/8digit-12                            222360716    5.403  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/12digit-12                           156132950    7.699  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/19digit-12                           135049720    8.880  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/1digit-12                            338642533    3.553  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint_DigitRuns/4digit-12                            212370570    5.672  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt/fiber-12                                        299791800    3.981  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt/fiber_bytes-12                                  288714712    4.165  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt/default-12                                       76370292    15.48  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber-12                                      200451391    5.988  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt32/fiber_bytes-12                                201629558    5.956  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt32/default-12                                     78384832    15.71  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber-12                                      220670613    5.462  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt16/fiber_bytes-12                                220175068    5.430  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt16/default-12                                    100000000    10.71  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber-12                                       277948545    4.279  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt8/fiber_bytes-12                                 264784368    4.566  ns/op     0  B/op   0  allocs/op
Benchmark_ParseInt8/default-12                                     142184748    8.373  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber-12                                     238916116    5.025  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint32/fiber_bytes-12                               222286592    5.409  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint32/default-12                                    90395479    13.33  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber-12                                     251398184    4.828  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint16/fiber_bytes-12                               248772012    4.838  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint16/default-12                                   137015901    8.788  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber-12                                      360430984    3.427  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint8/fiber_bytes-12                                327526336    3.659  ns/op     0  B/op   0  allocs/op
Benchmark_ParseUint8/default-12                                    181904056    6.438  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber-12                                     82912768    14.76  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat64/fiber_bytes-12                               95692504    12.68  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat64/default-12                                   50109265    25.05  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber-12                                     76648135    15.76  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat32/fiber_bytes-12                               86422119    13.63  ns/op     0  B/op   0  allocs/op
Benchmark_ParseFloat32/default-12                                   47198646    25.41  ns/op     0  B/op   0  allocs/op

# Time
Benchmark_CalculateTimestamp/fiber-12                             1000000000   0.3366  ns/op     0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/default-12                             36163686    32.69  ns/op     0  B/op   0  allocs/op
Benchmark_CalculateTimestamp/fiber_asserted-12                       4471299    267.8  ns/op    12  B/op   2  allocs/op
Benchmark_CalculateTimestamp/default_asserted-12                     3991147    298.7  ns/op     8  B/op   2  allocs/op

# URL
Benchmark_AppendPathSegmentsEscape/1-segment/fiber-12               84056664    13.99  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/1-segment/split-loop-12          56022625    21.97  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/3-segments/fiber-12              88724314    13.65  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/3-segments/split-loop-12         36393044    32.41  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/6-segments/fiber-12              87622688    13.82  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/6-segments/split-loop-12         23307777    51.23  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/escaping/fiber-12                37647206    31.85  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/escaping/split-loop-12           21136993    58.06  ns/op     0  B/op   0  allocs/op

# XML
Benchmark_GolangXMLEncoder-12                                         574852     2040  ns/op  4864  B/op  12  allocs/op
Benchmark_DefaultXMLEncoder-12                                        593738     2012  ns/op  4864  B/op  12  allocs/op
Benchmark_DefaultXMLDecoder-12                                        309031     3876  ns/op  2892  B/op  57  allocs/op
```

See all the benchmarks under <https://gofiber.github.io/utils/benchmarks>

## Case-insensitive matching helpers

`EqualFold`, `IndexFold`, `ContainsFold`, `HasPrefixFold`, and
`HasSuffixFold` compare ASCII case-insensitively using the SWAR word
loops: only `A`..`Z`/`a`..`z` fold, every other byte (including bytes
>= 0x80) must match exactly. `HasPrefixFold` and `HasSuffixFold` cover
the header checks Fiber otherwise spells as an allocation-prone
`ToLower` + `HasPrefix`/`HasSuffix` pair (authorization schemes such as
`Bearer`, content-type prefixes such as `multipart/form-data`, suffixes
such as `+json`), and like the other Fold helpers they take the needle
as a plain string because call sites pass constant tokens. Their
`Benchmark_HasPrefixFold`/`Benchmark_HasSuffixFold` numbers are tracked
on the [benchmark charts](https://gofiber.github.io/utils/benchmarks/)
and join the catalog above on its next regeneration.

Two words are compared through their XOR instead of being folded each:
they match when the lanes that differ do so in the case bit alone and
hold a letter, so a word that differs anywhere else is rejected in two
operations, and one in another case costs a single range test. A token
of 4 to 7 bytes, as `bearer `, `basic `, `gzip` and `+json` are, is two
overlapping 4-byte windows in one word for `EqualFold`, `HasPrefixFold`
and `HasSuffixFold` alike, and `EqualFold` packs one of 1 to 3 bytes
from its first, middle and last byte.

`IndexFold` stays linear when the haystack is built against the needle:
for a needle of nine bytes or more, once the candidates that fail their
verification have cost more than a fixed share of the bytes scanned, a
Rabin-Karp search over the folded bytes finishes the job, as in
`strings.Index`. A 1 MiB haystack of `a` searched for `AAAA...b` (63
`A`s) takes 1.7 ms where the candidate-by-candidate check took 21 ms;
benchstat over ten runs on an Intel Xeon @ 2.30GHz, and the ordinary rows
stay within 5% of what they were. `Benchmark_IndexFold/adversarial-*`
reproduces it.

`HashFold` is the hashing half of a case-insensitive lookup table whose
other half is `EqualFold`: strings `EqualFold` reports equal hash equal,
so a table of field or header names can be probed with a key in any case
without first lower-casing a copy of it. It folds each word with
`swar.ToLowerWord` as it loads it and mixes with wyhash's widening
multiply: keys under eight bytes, as most are, pack into one word and
cost one multiply, with the length in the multiplier so that keys which
pack alike, such as `a` and `aa`, stay apart; longer keys take one
multiply for each sixteen bytes before the last sixteen, and two more to
finish. It takes no seed and is no
defense against collisions planned in advance, so it suits keys of the
program's own, such as the fields of a struct, rather than keys an
attacker picks. Against `maphash.String` of a `strings.ToLower` copy,
benchstat over ten runs on an Intel Xeon @ 2.80GHz measured 4.8 ns
against 13.2 ns for a 2-byte key, 4.5 ns against 50.7 ns for a 4-byte
key and 14.7 ns against 184 ns for a 48-byte key, where the copy also
allocates;
`go test -run='^$' -bench='^Benchmark_HashFold$' -benchmem -count=10 . > hashfold.txt`
followed by `benchstat hashfold.txt` reproduces that comparison. The
`HashFold` rows in the amd64 table below, like the table's other rows,
are a single run on the machine its header names.

## HTTP dates

`AppendHTTPDate` and `FormatHTTPDate` write a time in the RFC 9110
preferred HTTP date format (`Mon, 02 Jan 2006 15:04:05 GMT`,
`net/http.TimeFormat`), byte-identical to `time.Format` with that layout
but without walking a layout string: the calendar fields come from one
civil-date conversion of `t.Unix()` (Hinnant's era arithmetic, instead of
the three `abs()` walks behind `Date`, `Clock`, and `Weekday`), and are
written into the fixed-width template with two-digit table lookups.
`ParseHTTPDate` is the reverse: canonical preferred-format input takes a
strict scalar fast path whose validated fields feed `time.Unix` directly
— the same `Time` value `time.Date` builds, minus its normalization pass —
and everything else — the obsolete RFC 850 and asctime forms, unusual
casing, non-GMT zone names, padding — falls back to `time.Parse` on the
text with the space, tab, CR, and LF around it trimmed (unlike
`net/http.ParseTime`, which does not trim), keeping its results and
its errors. Both directions are
pinned to the `time` package by a day-by-day sweep over the years
0..9999. `Date`, `Last-Modified`, and `If-Modified-Since` handling sit on
every request/response, which makes these the highest-leverage helpers in
this group for Fiber.

## URL escaping

`AppendQueryEscape`/`AppendPathEscape` and
`AppendQueryUnescape`/`AppendPathUnescape` produce byte-identical results
(and, for the unescape pair, identical `url.EscapeError` values) to their
`net/url` counterparts, append-style and without allocating when `dst`
has the room. A `dst` too small for `s` itself grows once, to the size of
the result: the bytes to escape, or the `%` to decode, are counted first,
where appending piece by piece grew a `nil` `dst` four or five times for
a 64-byte value. One that holds `s` is not counted, so that a reused
buffer with room for the result pays nothing for that, and grows as
`append` does should the escapes not fit. Nothing is written past the
result.
The escape tables are pinned to `net/url.shouldEscape` by exhaustive
per-byte tests, and bytes that need no escaping are skipped a word at a
time: eight lookups from one load, and one branch. Unescaping jumps
between escape sites with the vectorized
scans (`IndexAny2` when `+` needs rewriting, `bytes.IndexByte` otherwise)
and copies clean spans wholesale, so route parameters and query values
without escapes — the common case — cost one scan and one copy. Once an
escape comes fewer than four bytes after another, as in percent-encoded
UTF-8 or a value dense in `%20` and `+`, a byte loop decodes the rest,
without a scan and a copy for each escape, and looks for the next one
with a scan only after eight bytes without any. Decoding
never grows the input, so `dst` may be `s[:0]` on a common backing array
to unescape in place; escaping can grow the input, so there `dst` must
not alias `s`.

`AppendPathSegmentsEscape` covers a whole multi-segment path in one pass:
it escapes every segment as `AppendPathEscape` does and keeps the `/`
between them. The output equals escaping each segment and joining with
`/`, pinned per byte and by fuzzing; `net/url`'s whole-path escaping
differs, as it also leaves `;` and `,` unescaped. The catalog above
measures it against the per-segment loop a caller writes without it; the
other URL benchmarks are still in the amd64 block below.

## JSON string escaping

`AppendJSONString` appends a value as a double-quoted JSON string,
byte-identical to `encoding/json.Marshal` of the same string through Go
1.26 — including its default HTML escaping (`<`, `>`, `&`), the `\ufffd`
replacement of invalid UTF-8 bytes, and the U+2028/U+2029 escapes — without
any of `Marshal`'s reflection or allocation. From Go 1.27, `Marshal` writes
a raw U+FFFD for an invalid byte where `AppendJSONString` keeps the
six-character `\ufffd` escape; both decode to the same string, and every
other input stays byte-identical. Clean spans are located with a SWAR
scan (the same first-match-mask technique as `IndexNonQuotable`) and
copied wholesale. This is the building block for hand-rolled JSON hot
paths such as access-log lines and error bodies.

## IP address parsing

`ParseIPv4` and `ParseIPv6` parse addresses into `netip.Addr`, accepting
exactly the strings `netip.ParseAddr` accepts for the respective family
(pinned by fuzzing) while reporting failure with a `bool` instead of a
constructed error. Both take strings or byte slices, so fasthttp-style
callers skip the `string` conversion entirely; remote-address and
trusted-proxy checks parse an IP on every request in Fiber.

## Header key canonicalization

`CanonicalHeaderKey` Title-Cases an HTTP header key exactly like
`net/http.CanonicalHeaderKey`, including its return-unchanged guard for
keys containing non-token bytes. Already-canonical keys — the common case
on receive paths — are validated a word at a time and returned
as-is with zero allocations for strings and byte slices alike: in a key
of letters, digits and hyphens, as nearly all are, a letter is in the
wrong case exactly where its case bit agrees with "follows a hyphen",
which three lane masks settle for eight bytes at once. A key with any
other token byte, such as `_`, takes the table pass it always took, a
few nanoseconds later for the word test that turned it away
(`canonical-underscore` in the benchmark), and one of fewer than four
bytes takes nothing else. Keys that
do need rewriting cost one allocation, and flipping that bit is the whole
rewrite; note that for the ~40 header names
in the stdlib's interning table the stdlib returns a cached string
without allocating, so this helper's edge there is time, not allocations.

## Control-byte scanning

`IndexControl` returns the index of the first ASCII control byte — a byte
below 0x20 or DEL (0x7F), the RFC 5234 CTL set — and `IndexControlExceptTab`
is the same scan with HTAB permitted, which is the byte set an RFC 9110
field value may not contain. Both are SWAR first-match scans
(two words per branch, one overlapping word for the tail) over strings or
byte slices; the HTAB exemption is tested only in words that hold a
control byte, so clean input costs the same either way. They never match
bytes >= 0x80, so unlike
`strings.IndexFunc(s, unicode.IsControl)` they neither decode UTF-8 nor
flag the C1 range hidden inside it. Fiber spells this check by hand in
half a dozen places — header values before they are echoed, request IDs,
`Location` values, log fields, Basic auth credentials — and each copy is
one of these two calls.

## Byte-separator cuts

`CutByte` and `LastCutByte` are `strings.Cut`/`bytes.Cut` for a single
byte separator, around its first or last occurrence, generic over strings
and byte slices. They return the two parts and a found flag exactly like
the stdlib pair (including the `""`/`nil` after-part of a miss), cost one
`IndexByte`/`LastIndexByte`, and spare `[]byte` callers the `[]byte{sep}`
needle and the length dispatch inside `bytes.Index`. `CutByte` beats
`strings.Cut`; `LastCutByte` has no stdlib counterpart, so its row is
measured against a hand-written `LastIndexByte` split, which it trails by
the cost of the call — its value is the generic `[]byte` support and the
readable spelling. Header parameter splits (`;`), `key=value` pairs, and
host/port and name/extension splits are the intended call sites.

## Single-byte trimming

`Trim`, `TrimLeft` and `TrimRight` remove one repeated byte from the ends of
a string or byte slice, generic over both. Unlike `strings.Trim` and
`bytes.Trim`, the cutset is a single `byte`, not a set of characters. The
result is a subslice of the input, never a copy, and it is nil only for a
nil input, where `bytes.Trim` and `bytes.TrimLeft` return nil whenever
nothing is left. `TrimSpace` trims ASCII whitespace only and leaves Unicode
spaces alone.

## Host and port splitting

`SplitHostPort` splits `host:port`, `host%zone:port`, `[host]:port`, and
`[host%zone]:port` with exactly `net.SplitHostPort`'s acceptance rules
(pinned by fuzzing against it), but reports failure with a `bool` instead
of constructing a `*net.AddrError`, and accepts byte slices. Without
brackets, the last colon has to be the only `:`, `[` or `]` there is, so
an input of 8 to 40 bytes is one pass over its words that counts
them, where the stdlib scans it four times; bracketed, shorter and
longer inputs keep the stdlib's `IndexByte` scans, which are vectorized
and overtake the word pass on a long valid `host:port`. A malformed `Host` or
`Forwarded` value costs no allocation, and `[]byte` input needs no
string conversion first.

## List field iteration

`SplitTrimSeq` returns an iterator over the elements of a list separated
by one byte, with ASCII whitespace trimmed from each element and empty
elements skipped — how HTTP list fields such as `Accept-Encoding`,
`Connection`, and `Vary` are read (RFC 9110 Section 5.6.1: optional
whitespace around the commas, empty elements ignored). Quoted strings are
not interpreted, so a field whose elements can hold a quoted comma, such as
`Cache-Control: no-cache="Set-Cookie, Set-Cookie2"`, needs a reader that
honors the quotes. It replaces the `strings.SplitSeq` + `TrimSpace` + emptiness-check
triplet, locates each separator with a single `IndexByte`, works on
strings and byte slices alike, and yields subslices without copying. In
an ordinary caller the iterator inlines and allocates nothing, like
`strings.SplitSeq`; the benchmark ranges inside a helper function for
that reason, because calls made directly inside a `b.Loop` body are
never inlined and would show the closure allocations instead.

## Duration formatting

`AppendDuration` renders a `time.Duration` exactly as `d.String()` does —
"72h3m0.5s", "1.5ms", "0s" — into the caller's buffer, with the two-digit
table for the field digits; a value with no fraction, the common case,
never touches the fraction digits. `d.String()` allocates its result and
`fmt.Fprintf("%v", d)`, which Fiber's logger uses for the latency column,
costs several times more on top; an append-style form makes per-request
duration rendering (access logs, `Server-Timing`) allocation-free.

## Decimal formatting

`FormatUint`/`FormatInt` (and their 32- and 16-bit variants) and
`AppendUint`/`AppendInt` render integers eight digits at a time: a value
below 1e8 becomes eight zero-padded digit lanes through a fixed sequence
of multiplies, shifts, and masks (the SWAR inverse of the 8-digit parse
step in `ParseInt`), so a 64-bit value costs at most three such
conversions and two divisions, all independent of each other, instead of
one dependent divide-and-store per digit; the first significant digit is
found with a trailing-zero count. Values below 100 come from a
precomputed table, and the output is pinned to `strconv` by a dense sweep
plus every decimal and lane-group boundary. `GetMIME` in the same spirit
looks up extensions through a hash table keyed by the extension packed
into one lower-cased word with its length in the top lane, so a hit is a
multiply, a shift, and usually one slot comparison, with no allocation even
for upper-case input.

These helpers were added on an amd64 machine, so like the `simd` numbers
their benchmarks are recorded separately from the arm64 catalog above and
join it on its next regeneration. The block also carries the current
numbers for the paths reworked on the same machine (integer formatting,
`GetMIME`, `EqualFold`, the folded prefix and suffix tests, `ByteSize`,
`AppendPathSegmentsEscape`), whose arm64 rows above predate that work:

Environment:
goos: linux
goarch: amd64
pkg: github.com/gofiber/utils/v2
cpu: AMD EPYC

```text
// go test -benchmem -run=^$ -count=1 . -bench='^Benchmark_(FormatUint|FormatInt|FormatUint32|FormatInt32|FormatUint16|FormatInt16|AppendUint|AppendInt|EqualFold|EqualFoldBytes|EqualFold_Short|HasPrefixFold|HasSuffixFold|HashFold|ByteSize|GetMIME|AppendHTTPDate|ParseHTTPDate|AppendDuration|AppendQueryEscape|AppendPathSegmentsEscape|AppendQueryUnescape|AppendPathUnescape|AppendJSONString|ParseIPv4|ParseIPv6|CanonicalHeaderKey|IndexControl|CutByte|LastCutByte|SplitHostPort|SplitTrimSeq)$'
Benchmark_EqualFoldBytes/fiber-4                                    95242539    11.41  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFoldBytes/default-4                                  31083469    34.29  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold/fiber-4                                        100000000    10.23  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold/default-4                                       37118374    33.01  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/3B/fiber-4                               555012801    2.145  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/3B/default-4                             369259904    3.172  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/4B/fiber-4                               578990004    1.987  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/4B/default-4                             296558700    4.032  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/5B/fiber-4                               584135606    2.213  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/5B/default-4                             493922005    2.369  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/7B/fiber-4                               611238091    2.016  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/7B/default-4                             415740052    2.827  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/8B/fiber-4                               643355518    1.889  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/8B/default-4                             364321028    3.336  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/10B/fiber-4                              447896754    2.693  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/10B/default-4                            302122399    4.018  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/16B/fiber-4                              401625879    2.869  ns/op     0  B/op   0  allocs/op
Benchmark_EqualFold_Short/16B/default-4                            189607155    6.431  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/clean-16B/fiber-4                           515043063    2.452  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/clean-16B/fiber-except-tab-4                488922855    2.236  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/clean-16B/default-4                         149038825    8.223  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/clean-64B/fiber-4                           207569035    5.982  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/clean-64B/fiber-except-tab-4                200065816    5.961  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/clean-64B/default-4                          32232034    37.00  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/ctl-at-end-64B/fiber-4                      185200905    6.394  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/ctl-at-end-64B/fiber-except-tab-4           178918200    6.781  ns/op     0  B/op   0  allocs/op
Benchmark_IndexControl/ctl-at-end-64B/default-4                     32070586    43.63  ns/op     0  B/op   0  allocs/op
Benchmark_ByteSize/0-4                                            1000000000   0.2084  ns/op     0  B/op   0  allocs/op
Benchmark_ByteSize/1-4                                             137542143    8.351  ns/op     2  B/op   1  allocs/op
Benchmark_ByteSize/500-4                                           134812441    10.54  ns/op     4  B/op   1  allocs/op
Benchmark_ByteSize/1024-4                                          135406344    8.939  ns/op     3  B/op   1  allocs/op
Benchmark_ByteSize/1126-4                                          121401331    9.771  ns/op     5  B/op   1  allocs/op
Benchmark_ByteSize/1048576-4                                       136905632    8.886  ns/op     3  B/op   1  allocs/op
Benchmark_ByteSize/1153024-4                                       100000000    10.16  ns/op     5  B/op   1  allocs/op
Benchmark_ByteSize/1073741824-4                                    129641943    8.971  ns/op     3  B/op   1  allocs/op
Benchmark_ByteSize/1180696576-4                                    100000000    10.89  ns/op     5  B/op   1  allocs/op
Benchmark_ByteSize/1099511627776-4                                 122353975    9.391  ns/op     3  B/op   1  allocs/op
Benchmark_ByteSize/1209033293824-4                                  96642478    11.36  ns/op     5  B/op   1  allocs/op
Benchmark_ByteSize/1125899906842624-4                              126474978    9.147  ns/op     3  B/op   1  allocs/op
Benchmark_ByteSize/1238050092875776-4                              100000000    11.58  ns/op     5  B/op   1  allocs/op
Benchmark_ByteSize/1152921504606846976-4                           100000000    10.29  ns/op     3  B/op   1  allocs/op
Benchmark_ByteSize/1267763295104794624-4                           100000000    11.22  ns/op     5  B/op   1  allocs/op
Benchmark_CutByte/fiber-4                                          393649442    3.165  ns/op     0  B/op   0  allocs/op
Benchmark_CutByte/default-4                                        279499384    3.970  ns/op     0  B/op   0  allocs/op
Benchmark_CutByte/fiber-bytes-4                                    359259126    3.318  ns/op     0  B/op   0  allocs/op
Benchmark_CutByte/default-bytes-4                                  289247341    3.995  ns/op     0  B/op   0  allocs/op
Benchmark_LastCutByte/fiber-4                                      719761060    1.665  ns/op     0  B/op   0  allocs/op
Benchmark_LastCutByte/default-4                                    770568099    1.541  ns/op     0  B/op   0  allocs/op
Benchmark_AppendDuration/micros/fiber-4                            205667509    5.816  ns/op     0  B/op   0  allocs/op
Benchmark_AppendDuration/micros/default-4                           83072940    14.72  ns/op     8  B/op   1  allocs/op
Benchmark_AppendDuration/micros/fmt-4                               22109742    57.63  ns/op    16  B/op   2  allocs/op
Benchmark_AppendDuration/millis-frac/fiber-4                       146680158    8.189  ns/op     0  B/op   0  allocs/op
Benchmark_AppendDuration/millis-frac/default-4                      70176217    18.40  ns/op    16  B/op   1  allocs/op
Benchmark_AppendDuration/millis-frac/fmt-4                          20415859    58.63  ns/op    24  B/op   2  allocs/op
Benchmark_AppendDuration/hours/fiber-4                             167795612    7.209  ns/op     0  B/op   0  allocs/op
Benchmark_AppendDuration/hours/default-4                            57929433    19.89  ns/op     8  B/op   1  allocs/op
Benchmark_AppendDuration/hours/fmt-4                                19839003    59.91  ns/op    16  B/op   2  allocs/op
Benchmark_FormatUint/small/fiber-4                                1000000000    1.098  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint/small/strconv-4                              1000000000    1.218  ns/op     0  B/op   0  allocs/op
Benchmark_FormatUint/medium/fiber-4                                 82577330    14.51  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint/medium/strconv-4                               89600736    14.50  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint/large/fiber-4                                  76925103    17.43  ns/op    24  B/op   1  allocs/op
Benchmark_FormatUint/large/strconv-4                                57337281    21.09  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/small_pos/fiber-4                             1000000000    1.034  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_pos/strconv-4                           1000000000    1.154  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/fiber-4                              921751848    1.257  ns/op     0  B/op   0  allocs/op
Benchmark_FormatInt/small_neg/strconv-4                            126309010    9.483  ns/op     3  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/fiber-4                              74587729    15.48  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_pos/strconv-4                            90545918    13.90  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/fiber-4                              86877349    14.85  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/medium_neg/strconv-4                            88416847    15.03  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/fiber-4                               65433602    16.88  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_pos/strconv-4                             53400705    21.95  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/fiber-4                               77646357    16.77  ns/op    24  B/op   1  allocs/op
Benchmark_FormatInt/large_neg/strconv-4                             60803770    20.78  ns/op    24  B/op   1  allocs/op
Benchmark_FormatUint32/fiber-4                                      71134392    14.28  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint32/strconv-4                                    83261202    14.96  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt32/fiber-4                                       88697210    13.78  ns/op    16  B/op   1  allocs/op
Benchmark_FormatInt32/strconv-4                                     80188908    14.30  ns/op    16  B/op   1  allocs/op
Benchmark_FormatUint16/fiber-4                                     127857867    9.380  ns/op     5  B/op   1  allocs/op
Benchmark_FormatUint16/strconv-4                                   100000000    10.43  ns/op     5  B/op   1  allocs/op
Benchmark_FormatInt16/fiber-4                                      106312431    10.97  ns/op     8  B/op   1  allocs/op
Benchmark_FormatInt16/strconv-4                                    100000000    11.40  ns/op     8  B/op   1  allocs/op
Benchmark_AppendUint/fiber-4                                       171862742    6.993  ns/op     0  B/op   0  allocs/op
Benchmark_AppendUint/strconv-4                                     171259778    6.968  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/fiber-4                              548362210    2.293  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/small_neg/strconv-4                            386668742    3.265  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/fiber-4                             170506190    7.058  ns/op     0  B/op   0  allocs/op
Benchmark_AppendInt/medium_neg/strconv-4                           162180542    7.410  ns/op     0  B/op   0  allocs/op
Benchmark_HashFold/2B/fiber-4                                      662384479    1.846  ns/op     0  B/op   0  allocs/op
Benchmark_HashFold/2B/default-4                                    241535532    4.928  ns/op     0  B/op   0  allocs/op
Benchmark_HashFold/4B/fiber-4                                      769847120    1.544  ns/op     0  B/op   0  allocs/op
Benchmark_HashFold/4B/default-4                                     65858584    17.81  ns/op     8  B/op   1  allocs/op
Benchmark_HashFold/12B/fiber-4                                     436935361    2.619  ns/op     0  B/op   0  allocs/op
Benchmark_HashFold/12B/default-4                                    36352672    30.86  ns/op    16  B/op   1  allocs/op
Benchmark_HashFold/17B/fiber-4                                     300781833    3.976  ns/op     0  B/op   0  allocs/op
Benchmark_HashFold/17B/default-4                                    35985314    33.97  ns/op    24  B/op   1  allocs/op
Benchmark_HashFold/48B/fiber-4                                     216305008    5.534  ns/op     0  B/op   0  allocs/op
Benchmark_HashFold/48B/default-4                                    19581700    64.92  ns/op    48  B/op   1  allocs/op
Benchmark_CanonicalHeaderKey/canonical/fiber-4                     208397125    5.848  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical/default-4                    94861126    12.76  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/common-lower/fiber-4                   63354654    18.50  ns/op    16  B/op   1  allocs/op
Benchmark_CanonicalHeaderKey/common-lower/default-4                 42964386    27.67  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/custom-lower/fiber-4                   39557229    29.37  ns/op    24  B/op   1  allocs/op
Benchmark_CanonicalHeaderKey/custom-lower/default-4                 21005216    48.28  ns/op    24  B/op   1  allocs/op
Benchmark_CanonicalHeaderKey/canonical-short/fiber-4               377163668    3.192  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical-short/default-4             245415436    5.073  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical-long/fiber-4                100000000    10.82  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical-long/default-4               43972626    28.37  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical-underscore/fiber-4           91702238    13.15  ns/op     0  B/op   0  allocs/op
Benchmark_CanonicalHeaderKey/canonical-underscore/default-4         93589179    12.77  ns/op     0  B/op   0  allocs/op
Benchmark_SplitHostPort/host-port/fiber-4                          212407646    5.677  ns/op     0  B/op   0  allocs/op
Benchmark_SplitHostPort/host-port/default-4                        163514972    7.458  ns/op     0  B/op   0  allocs/op
Benchmark_SplitHostPort/ipv6-port/fiber-4                          150666772    8.119  ns/op     0  B/op   0  allocs/op
Benchmark_SplitHostPort/ipv6-port/default-4                        151922906    7.944  ns/op     0  B/op   0  allocs/op
Benchmark_SplitHostPort/missing-port/fiber-4                       245688876    4.768  ns/op     0  B/op   0  allocs/op
Benchmark_SplitHostPort/missing-port/default-4                      69299520    15.95  ns/op    32  B/op   1  allocs/op
Benchmark_SplitHostPort/long-host-port/fiber-4                     145949444    8.510  ns/op     0  B/op   0  allocs/op
Benchmark_SplitHostPort/long-host-port/default-4                   132700629    8.889  ns/op     0  B/op   0  allocs/op
Benchmark_GetMIME/fiber-4                                           58912044    18.25  ns/op     0  B/op   0  allocs/op
Benchmark_GetMIME/default-4                                         20459073    52.75  ns/op     0  B/op   0  allocs/op
Benchmark_AppendHTTPDate/fiber-4                                    82701052    14.49  ns/op     0  B/op   0  allocs/op
Benchmark_AppendHTTPDate/default-4                                  18782808    71.80  ns/op     0  B/op   0  allocs/op
Benchmark_ParseHTTPDate/fiber-4                                    124944092    9.542  ns/op     0  B/op   0  allocs/op
Benchmark_ParseHTTPDate/default-4                                   14183923    84.09  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv4/fiber-4                                        134509474    8.779  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv4/default-4                                       96120033    13.16  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/compressed/fiber-4                              60667490    18.69  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/compressed/default-4                            51461449    24.07  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/full/fiber-4                                    50006224    25.44  ns/op     0  B/op   0  allocs/op
Benchmark_ParseIPv6/full/default-4                                  29281041    34.67  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/clean-16B/fiber-4                       116575838    10.14  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/clean-16B/default-4                      18329967    66.97  ns/op    40  B/op   2  allocs/op
Benchmark_AppendJSONString/clean-64B/fiber-4                        39247297    29.34  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/clean-64B/default-4                      11358939    100.4  ns/op    96  B/op   2  allocs/op
Benchmark_AppendJSONString/escaped-64B/fiber-4                      16470195    74.15  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/escaped-64B/default-4                     9741043    127.1  ns/op    96  B/op   2  allocs/op
Benchmark_AppendJSONString/unicode-64B/fiber-4                      10718854    111.9  ns/op     0  B/op   0  allocs/op
Benchmark_AppendJSONString/unicode-64B/default-4                     7911082    149.5  ns/op   112  B/op   2  allocs/op
Benchmark_HasPrefixFold/bearer-hit/fiber-4                         543062940    2.239  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/bearer-hit/default-4                       396693123    3.057  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/bearer-miss/fiber-4                        859282924    1.355  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/bearer-miss/default-4                      481811526    2.511  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/text-hit/fiber-4                           557733916    2.165  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/text-hit/default-4                         281911194    4.259  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/multipart-hit/fiber-4                      203003868    5.928  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/multipart-hit/default-4                    100000000    10.14  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/multipart-miss/fiber-4                     807200426    1.475  ns/op     0  B/op   0  allocs/op
Benchmark_HasPrefixFold/multipart-miss/default-4                   707085358    1.748  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/gzip-hit/fiber-4                           481585503    2.498  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/gzip-hit/default-4                         284746641    4.228  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/gzip-miss/fiber-4                          783951544    1.513  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/gzip-miss/default-4                        577793062    1.939  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/json-hit/fiber-4                           524795856    2.284  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/json-hit/default-4                         282228116    4.251  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/atomxml-hit/fiber-4                        380070187    3.268  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/atomxml-hit/default-4                      217422844    5.573  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/charset-miss/fiber-4                       802579796    1.526  ns/op     0  B/op   0  allocs/op
Benchmark_HasSuffixFold/charset-miss/default-4                     641014597    1.891  ns/op     0  B/op   0  allocs/op
Benchmark_SplitTrimSeq/accept-encoding/fiber-4                     110528046    10.79  ns/op     0  B/op   0  allocs/op
Benchmark_SplitTrimSeq/accept-encoding/default-4                   100000000    12.17  ns/op     0  B/op   0  allocs/op
Benchmark_SplitTrimSeq/accept/fiber-4                               45297562    25.86  ns/op     0  B/op   0  allocs/op
Benchmark_SplitTrimSeq/accept/default-4                             44616727    26.89  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/clean-64B/fiber-4                       84951325    14.45  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/clean-64B/fiber-grow-4                  39905455    31.54  ns/op    64  B/op   1  allocs/op
Benchmark_AppendQueryEscape/clean-64B/default-4                     81626277    14.50  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/mixed-64B/fiber-4                       20833869    57.30  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/mixed-64B/fiber-grow-4                  12820696    104.0  ns/op    96  B/op   1  allocs/op
Benchmark_AppendQueryEscape/mixed-64B/default-4                     12114079    106.0  ns/op   192  B/op   2  allocs/op
Benchmark_AppendQueryEscape/words-64B/fiber-4                       26774242    40.72  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryEscape/words-64B/fiber-grow-4                  15988232    74.43  ns/op    64  B/op   1  allocs/op
Benchmark_AppendQueryEscape/words-64B/default-4                     23491057    52.78  ns/op    64  B/op   1  allocs/op
Benchmark_AppendPathSegmentsEscape/1-segment/fiber-4               167753917    7.173  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/1-segment/split-loop-4           97259446    12.25  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/3-segments/fiber-4              164162836    7.302  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/3-segments/split-loop-4          62738319    20.08  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/6-segments/fiber-4              162682024    7.289  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/6-segments/split-loop-4          38979402    32.49  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/escaping/fiber-4                 78633295    17.86  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathSegmentsEscape/escaping/split-loop-4            39662082    30.13  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/plain-64B/fiber-4                    162386299    7.303  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/plain-64B/fiber-grow-4                47873440    25.18  ns/op    64  B/op   1  allocs/op
Benchmark_AppendQueryUnescape/plain-64B/default-4                   22820241    50.43  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/escaped-64B/fiber-4                   22122288    54.84  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/escaped-64B/fiber-grow-4              16182109    77.04  ns/op    48  B/op   1  allocs/op
Benchmark_AppendQueryUnescape/escaped-64B/default-4                 11858275    104.1  ns/op    48  B/op   1  allocs/op
Benchmark_AppendQueryUnescape/utf8-62B/fiber-4                      27144454    41.86  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/utf8-62B/fiber-grow-4                 19416654    62.87  ns/op    48  B/op   1  allocs/op
Benchmark_AppendQueryUnescape/utf8-62B/default-4                    15800938    74.23  ns/op    48  B/op   1  allocs/op
Benchmark_AppendQueryUnescape/words-64B/fiber-4                     15931467    72.81  ns/op     0  B/op   0  allocs/op
Benchmark_AppendQueryUnescape/words-64B/fiber-grow-4                12764703    98.86  ns/op    64  B/op   1  allocs/op
Benchmark_AppendQueryUnescape/words-64B/default-4                    9407224    125.9  ns/op    64  B/op   1  allocs/op
Benchmark_AppendPathUnescape/plain-25B/fiber-4                     230424097    5.156  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathUnescape/plain-25B/default-4                    70821981    17.02  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathUnescape/escaped-31B/fiber-4                    72507342    16.33  ns/op     0  B/op   0  allocs/op
Benchmark_AppendPathUnescape/escaped-31B/default-4                  25323608    48.15  ns/op    24  B/op   1  allocs/op
```

## Float parsing

`ParseFloat64` and `ParseFloat32` return the float `strconv.ParseFloat`
returns, bit for bit, for an optional sign, digits with an optional point and
an optional exponent, in strings or byte slices, wherever strconv reads the
number right: with more than 800 digits ahead of the point it keeps 800 and
places the point by that many (`"2"` and 800 zeros with `E-800` gives 0.2), and
these functions return the correctly rounded float instead. The spellings strconv adds
(`inf`, `nan`, hexadecimal floats, underscores) are syntax errors, and a value
too large for the type is a range error; every error comes with a zero
result, where strconv returns ±Inf for a range error. A mantissa up to 2^53
(2^24 for `ParseFloat32`) with an exponent of at most 22 (10) either way is
one multiplication or division by an exact power of ten, which rounds
correctly (Clinger's fast path, strconv's own first step). An exponent past
which every 64-bit mantissa is zero or infinity is settled at once, and so is
a product that is an integer times a power of two (the halfway points between
floats, integers past 2^53, decimals padded with zeros), with integer
arithmetic. Everything else goes to `strconv.ParseFloat`: the text itself up
to 700 bytes, and past that an exact mantissa and exponent or, for a number
with more digits than a `uint64` holds, at most 768 digits, a final 1 for any
nonzero digit after them and an exponent, a form strconv reads correctly
however long the input was. `Benchmark_ParseFloat_Shapes` times each shape
against strconv: 2 to 3 times as fast on the fast path, and 1.3 to 1.7 times
slower on the numbers handed over to it, whose digits are read twice.

## Unsafe views

`UnsafeString` and `UnsafeBytes` convert between a string and a byte slice
without copying, so the two share memory and neither side may be modified
while the other is in use. A write through `UnsafeBytes` faults for a string
literal and changes any other string for every holder of it, a map key, a
header or a cache entry among them; a write to the bytes behind an
`UnsafeString` changes the string under whoever holds it. `CopyBytes` and
`CopyString` give a copy for the side that outlives or modifies the other.

## Filesystem walking

`Walk` is `filepath.Walk` over an `http.FileSystem`, with the same control
flow (`filepath.SkipDir` and `filepath.SkipAll` are returned as they are).
It follows symbolic links, as `http.Dir` does when it opens a file, and the
`FileInfo` it passes comes from that open file, so it never reports a link:
a link to an ancestor repeats the tree until `ELOOP`, and a link that leaves
the tree walks what it points to. A program that serves or indexes what
`Walk` finds has to keep the paths inside the root itself, for example by
resolving each one with `filepath.EvalSymlinks` under the real directory;
`http.FileSystem` has no `Lstat` to do it here.

## Size strings

`ConvertToBytes` turns `42`, `42kb`, `1.5GiB` or `130.51MB` into a byte
count: decimal units are powers of 1000, units with an `i` right after the
letter (`Ki`, `KiB`, `kibibytes`) are powers of 1024, but not the `i` of
`kilobytes` or `gigabytes`, and the result saturates at `math.MaxInt`. The
units run from `k` to `P`, and to `E` when a `b` or `i` follows it (`1.5EB`,
`1EiB`), since a bare `E` is an exponent marker. The
number may have a `+`, a point and an exponent (`+2.01k`, `1.5e3K`), with
any number of digits, and is scaled with 128-bit integer arithmetic, so the
result is the exact product rounded down (`65.35KB` is 65350) and never passes through a
`float64`. A number with more digits than a `uint64` holds is scaled from its
first digits when they and the same digits plus one agree, and otherwise
(`0.000000000931322574615478515625Gi` is exactly one byte) from all of them.
A negative or malformed size is 0, and an unknown unit counts as bytes. A
whole number, with or without a unit, is one digit loop and a multiplication.
`Benchmark_ConvertToBytes_Shapes` times each shape.

`ParseByteSize` is the strict form for configuration values, where a guess is
a bug: it takes the number, any blanks, and at most one unit of `b`, `k` to
`p` (with an `i` for powers of 1024 and a trailing `b`) and `EB` or `Ei`, in
either case, and returns a `*strconv.NumError` for anything else, such as
`10 X`, `5kilobytes`, a negative size, or a size of `math.MaxInt` or more.
The number is scaled by `ConvertToBytes`, so the two agree on every size
`ParseByteSize` accepts.

`ByteSize` goes the other way, in powers of 1024 and with one decimal
place (`1.5KB`, `1023.9KB`, `16EB`). It writes its at most eight bytes
right to left on the stack, so the string it returns is allocated at its
own length, and a size of zero allocates nothing.

## SWAR primitives

The [`swar`](swar/) package exports the SWAR (SIMD within a register)
building blocks the helpers above are composed from: `Load8`/`Store8`,
`Broadcast`, `ZeroLanes`, `MatchByteMask`/`MatchRangeMask`,
`ToLowerWord`/`ToUpperWord`, `FirstLane`/`LastLane`, and the `WordLen`,
`Ones`, `HighBits`, and `LowSeven` constants. They are exported so
downstream packages (Fiber itself, middleware) can fuse their own byte
scans — for example, finding a delimiter while classifying the bytes
before it — without re-deriving the bit tricks. The contracts (`Load8` and
`Store8` panic outside `len`, `ZeroLanes`' approximate mask, the
little-endian lane order) and runnable examples live in the package
documentation.

`Load8` and `Store8` check the offset against `len`, also for a `[]byte`
with spare capacity, where a plain reslice checks `cap` and a miscounted
offset would read or write past `len` without a panic. The check costs a
compare only where the compiler cannot prove the offset and nothing inside
a pinned group (below); the callers' benchmarks show no net change
(geomean -0.3% over 49 rows, `benchstat -count=10`).

The package benchmarks its primitives against their stdlib counterparts
the same way the helpers above do, measuring the canonical loops composed
from them: `Load8`/`Store8` vs `encoding/binary`'s little-endian
`Uint64`/`PutUint64`, a `ZeroLanes` first-match scan vs
`strings.IndexByte`, a `MatchByteMask`+`LastLane` reverse scan vs
`bytes.LastIndexByte`, a `MatchRangeMask` digit scan vs
`strings.IndexAny`, and an in-place `ToLowerWord` loop vs
`bytes.ToLower`. The composed loops are pinned to the stdlib results by a
dedicated test, and the numbers are tracked per commit on the
[benchmark charts](https://gofiber.github.io/utils/benchmarks/). Note
that `strings.IndexByte` is hand-written SIMD assembly and overtakes the
portable SWAR scan on large inputs; the SWAR primitives earn their keep
on short HTTP-sized inputs and on fused scans the stdlib has no single
function for. The package's remaining benchmark, `Benchmark_Load8_Fusion`,
is an advisory codegen guard rather than an API performance promise, so
it is deliberately not part of the catalog above.

## SIMD-accelerated search (package simd)

The [`simd`](simd/) package is the vector-width counterpart to `swar`: byte
searching and validation whose scan kernels dispatch to AVX2 assembly (four
32-byte vectors per iteration) on amd64 CPUs for inputs of `simd.MinLen`
(32) bytes or more (`MemchrPair`: 32 bytes past its offset), and fall back
to portable loops built on the `swar` primitives everywhere else, so it
builds and behaves identically on every platform.
`simd.Accelerated()` reports which mode is active. CPU capability detection
uses [`golang.org/x/sys/cpu`](https://pkg.go.dev/golang.org/x/sys/cpu),
which also verifies OS support for saving the vector register state.

The AVX2-backed kernels are the multi-needle scans (`Memchr2`, `Memchr3`),
the paired-byte scan (`MemchrPair`), the class scans
(`MemchrDigit`/`MemchrDigitAt`, `MemchrWord`/`MemchrNotWord`), the ASCII
scans (`IsASCII`, `FirstNonASCII`, `CountNonASCII`), and the substring
search (`Memmem`), which
prefilters on the needle's rarest bytes (`SelectRareBytes`, `ByteRank`)
before verifying candidates — 2-6 byte needles containing two distinct
values are scanned for their two rarest bytes at the exact relative
distance, longer or single-valued needles for the single rarest byte.
Each kernel classifies four 32-byte vectors per iteration. Where the
result is a position or a yes/no, the four masks are combined, so a
128-byte block costs a single `VPMOVMSKB` and a single branch, and the
scans that report a position re-extract the four masks in address order on
the (rare) hit path. `CountNonASCII` is the exception — a count cannot be
recovered from a combined mask, so its block loop accumulates per-lane
counts in vector byte lanes (a signed compare against zero and a byte
subtract per vector) and reduces them once at the end, while the
remaining whole vectors are extracted and population-counted; it is the
one kernel gated on `POPCNT` in addition to AVX2.
One tier is documented explicitly as unaccelerated:
`MemchrInTable`/`MemchrNotInTable` are plain
scalar loops, since an arbitrary 256-entry membership test has no cheap
vector form. There is deliberately no single-needle `Memchr`:
`bytes.IndexByte` is already vector-accelerated by the Go runtime. `Memmem`
delegates to `bytes.Index` for haystacks of up to 64 bytes, where the
stdlib's brute-force vector kernel is level with the paired prefilter
(`Benchmark_Memmem_Prefilter` forces the prefilter helpers at every size
to keep that measurable); from 65 bytes `bytes.Index` switches to a loop
anchored on the needle's first byte, and the paired scan takes over.
Needles that take the single rare-byte scan keep a 128-byte routing point,
since that scan is less selective. Both prefilters bail to `bytes.Index`
after a bounded number of failed candidate verifications, so even
adversarial inputs — a rare byte everywhere plus a long almost-matching
needle — stay within a constant of the stdlib's O(n+m) (see
`Benchmark_Memmem_Adversarial`). The
package operates on `[]byte` only; the generic top-level helpers `IsASCII`,
`IndexAny2`, and `IndexAny3` dispatch into it automatically for inputs of
`simd.MinLen`+ bytes on amd64, which is where the AVX2 kernels overtake the
SWAR word loops (-35% to -88% ns/op at 32-512B in `benchstat -count=10`
runs; shorter inputs keep the existing SWAR paths).

The kernels are adapted from the [coregex](https://github.com/coregx/coregex)
project's `simd` package (MIT License, Copyright (c) 2025 Andrey Kolkov and
contributors; see [`simd/LICENSE`](simd/LICENSE)), with the legacy-SSE
register moves in the kernels replaced by VEX-encoded ones to avoid AVX-SSE
transition stalls, the scalar tail loops replaced by overlapping vector
rescans at the buffer end (a 63-byte scan previously cost ~3.7x a 64-byte
one), the 32-byte main loops replaced by 4x unrolled 128-byte blocks whose
bound is tested once at the bottom against a precomputed limit pointer, the
`\w` classifier rewritten from three `VPMINUB`/`VPMAXUB` clamp-and-compare
range tests into two `VPSHUFB` nibble-table lookups (with the block test
reduced to one compare over the four raw class values), the digit range
test reduced from two signed compares and a `VPANDN` to one biased
`VPADDB` and one compare, `FirstNonASCII` and `CountNonASCII` given
kernels of their own (they were SWAR-only, and the count kernel's block
loop accumulates in vector byte lanes rather than extracting and
population-counting every vector's mask), the 128-byte block loops
placed at fixed offsets inside a 32-byte code window so their throughput
no longer depends on where the linker places each kernel, and the
fallbacks reworked around the stdlib and the `swar` package as described
above.

Against the pre-unrolling kernels, `benchstat -count=10` on the machine
below reports -26.6%/-49.7% ns/op for `Memchr2` at 512B/4096B,
-17.6%/-34.3% for `Memchr3`, -13.4%/-24.6% for `MemchrPair`,
-39.2%/-50.5% for `MemchrDigit`, -53.5%/-63.7% for `MemchrNotWord`,
-42.2%/-65.7% for `IsASCII`, and -13.7%/-27.1% for `Memmem`, which rides
the paired scan. Those percentages come from `benchstat` over ten runs of
each side; the single-run table below is a catalog, not a base-vs-head
comparison, so subtracting its rows from an older copy of it will not
reproduce them.

Two bands see no benefit and are worth stating plainly. Inputs below
`simd.MinLen` never reach the kernels at all. Inputs of 32-127 bytes do
reach them but skip the 128-byte block loop, so they pay its setup with
nothing to amortize it against: across that band `benchstat` reports the
class scans clearly ahead (`MemchrNotWord` -8.6% at 32B, -18.4% at 64B;
`IsASCII` -5.0%/-4.3%), most rows statistically unchanged, and two small
regressions — `MemchrDigit` +5.6% at 32B and `Memchr3` +3.8% at 64B. Read
those small-input figures against the run's own noise floor: the `default`
legs call unchanged stdlib code, so their movement bounds what a delta of
that size means. They held within 1% in the run quoted here, but a later
run on the same machine drifted up to 17% on those same rows, so anything
under ~10% at 32-64B should be taken as "no measured change" unless its
control legs are quiet.

On the portable side, pinning each unrolled group of words with a reslice
before loading at constant offsets inside it removes the per-load bounds
checks that `swar.Load8` otherwise repeats. Measured with
`benchstat -count=10` against the same fallbacks, at 32B/512B/4KiB:
`IsASCII` -28%/-60%/-63%, `CountNonASCII` -24%/-44%/-44%,
`FirstNonASCII` -23%/-37%/-43%, `Memchr2` -19%/-23%/-25%,
`MemchrDigit` -7%/-12%/-17%, `Memchr3` -8%/-13%/-18%,
`MemchrWord` -11%/-8%/-11%, `MemchrNotWord` -5%/-7%/-6%. The tier a group
adds costs one compare on inputs too short to use it, which shows up as
+3% to +10% on some of the 8B fallback rows — a few tenths of a
nanosecond, and only below the length where the group pays for itself.
`utils.IsASCII` carries the same rework, so the two copies of that scan
stay in step.

Making the counting kernel finish its own sub-vector remainder (rather
than splitting the input in Go and calling the SWAR loop for the tail)
also removed a second non-inlinable call from every whole-vector input:
`CountNonASCII` at 32B went from 11.3ns — behind its own fallback — to
6.8ns, level with it, and it leads from 64B up.

A later audit pass reworked the kernels and fallbacks once more;
`benchstat -count=10` on the machine below, base vs. head, reports:

- `MemchrDigit` -20% at 512B and -27% at 4KiB from the one-add,
  one-compare range test.
- `MemchrWord` -13%/-11% and `MemchrNotWord` -4%/-10% at 512B/4KiB from
  comparing the reduced block once instead of every vector, and their
  SWAR fallbacks -36%/-41% at 32B, -40%/-43% at 64B and -43%/-45% at
  512B: 15 ALU operations per word instead of 26, from case-folding the
  two letter ranges into one test, a one-add underscore test, and a
  single shared ASCII gate.
- `CountNonASCII` -15% at 512B and -41% at 4KiB, with 32B and 64B
  unchanged: the block loop counts in vector byte lanes (a signed compare
  against zero and a byte subtract per vector, flushed through `VPSADBW`
  every 127 blocks) instead of four port-0 `VPMOVMSKB` extractions per
  block. The remaining whole vectors are still extracted and
  population-counted, because the vector reduction is a fixed ten-cycle
  chain that one to three vectors cannot amortize — a first version that
  used the accumulators for those too was 20% slower at 32-64B.
- `MemchrPair`'s SWAR fallback -33% at 32B, -41% at 64B, -51% at 512B and
  -52% at 4KiB, and `MemchrPair` itself -30% at 32B, where the dispatch
  routes to that fallback: it now scans two words per branch through
  pinned three-index windows (no per-load capacity check or pointer
  clamp) and finishes with one overlapping word instead of a scalar
  tail. Before, at 32B, it ran no faster than a plain scalar loop.
- `Memmem` -62% at 96B: needles that take the paired prefilter are routed
  through it from 65-byte haystacks, where `bytes.Index` leaves its
  brute-force vector kernel for a first-byte-anchored loop; 64B and
  128B+ are unchanged.

The same pinned-window rework applied to the case-conversion loops in
`internal/caseconv`, which the `strings` and `bytes` packages' `ToLower`
and `ToUpper` build on, gives on the same machine -27% to -37% ns/op for
the in-place `UnsafeToLower`/`UnsafeToUpper` at 39-64B, -14% to -19% for
the scan-only path that finds nothing to convert, and -5% to -17% for the
allocating conversions, where the allocation dominates. The arm64 catalog
rows above were not regenerated for that change.

Comparing the base and head test binaries also showed the untouched
`MemchrPair` kernel ~10% slower at 4KiB, purely because the linker had
moved its entry by 32 bytes: on Skylake-derived cores the JCC erratum
mitigation keeps a macro-fused compare-and-branch that crosses or ends on
a 32-byte boundary out of the decoded-uop cache, so the block loop ran
from the legacy decoder. Every block loop is now placed by a `BLOCKALIGN`
macro at a window offset where all of its fused pairs fit (a few
once-per-call NOPs, visible as +5% on `Memchr2` at 512B and nothing at
4KiB), and an amd64 test reads each kernel's code and fails if a later
edit lets a pair straddle a boundary again.

Because the AVX2 kernels only engage on amd64, their benchmark numbers are
recorded separately from the arm64 catalog above:

Environment:
goos: linux
goarch: amd64
pkg: github.com/gofiber/utils/v2/simd
cpu: AMD EPYC

```text
// go test ./simd/ -benchmem -run=^$ -bench=Benchmark_ -count=1
Benchmark_Memchr2/8B/simd-4                                        438827100    2.751  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/8B/scalar-4                                      593404477    2.041  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/8B/default-4                                      80032118    15.64  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/32B/simd-4                                       512293941    2.395  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/32B/scalar-4                                     301949120    4.649  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/32B/default-4                                     76776531    15.63  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/64B/simd-4                                       433758777    2.715  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/64B/scalar-4                                     184668122    6.490  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/64B/default-4                                     41520069    27.35  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/512B/simd-4                                      262840983    4.351  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/512B/scalar-4                                     27931284    43.70  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/512B/default-4                                     5704527    210.4  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/4096B/simd-4                                      38343454    28.02  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/4096B/scalar-4                                     3376074    354.0  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr2/4096B/default-4                                     711501     1714  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/8B/simd-4                                        354681948    3.244  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/8B/scalar-4                                      439717762    2.660  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/8B/default-4                                      70377523    16.67  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/32B/simd-4                                       491647814    2.388  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/32B/scalar-4                                     221500863    5.426  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/32B/default-4                                     74339248    17.10  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/64B/simd-4                                       433255318    2.765  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/64B/scalar-4                                     134343738    9.015  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/64B/default-4                                     44080189    27.44  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/512B/simd-4                                      218295508    5.438  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/512B/scalar-4                                     19787991    62.38  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/512B/default-4                                     4982817    220.9  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/4096B/simd-4                                      31439382    38.11  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/4096B/scalar-4                                     2460133    488.3  ns/op     0  B/op   0  allocs/op
Benchmark_Memchr3/4096B/default-4                                     759471     1599  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/8B/simd-4                                     221210430    5.232  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/8B/swar-4                                     348585800    3.556  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/8B/scalar-4                                   529277721    2.295  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/32B/simd-4                                    202178641    5.984  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/32B/swar-4                                    253528896    4.717  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/32B/scalar-4                                   94890778    12.35  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/64B/simd-4                                    345571293    3.425  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/64B/swar-4                                    154802499    7.898  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/64B/scalar-4                                   79833866    24.55  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/512B/simd-4                                   195209959    6.152  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/512B/swar-4                                    22109486    51.10  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/512B/scalar-4                                   5816596    222.8  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/4096B/simd-4                                   33390804    36.36  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/4096B/swar-4                                    2868840    418.4  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrPair/4096B/scalar-4                                   706368     1692  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/8B/simd-4                                    428286242    3.066  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/8B/scalar-4                                  801597561    1.579  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/32B/simd-4                                   531411568    2.375  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/32B/scalar-4                                 369762055    3.232  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/64B/simd-4                                   454296792    2.509  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/64B/scalar-4                                 243953842    4.928  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/512B/simd-4                                  291673386    3.973  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/512B/scalar-4                                 34388494    36.83  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/4096B/simd-4                                  56277913    22.34  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrDigit/4096B/scalar-4                                 4410810    279.9  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/8B/simd-4                                     405158480    2.988  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/8B/scalar-4                                   539351959    2.206  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/32B/simd-4                                    449819178    2.570  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/32B/scalar-4                                  215114436    5.620  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/64B/simd-4                                    432078406    2.905  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/64B/scalar-4                                  120594146    9.878  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/512B/simd-4                                   202084159    5.930  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/512B/scalar-4                                  17016954    72.36  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/4096B/simd-4                                   27692071    44.06  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrWord/4096B/scalar-4                                  2019286    576.4  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/8B/simd-4                                  391076607    3.136  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/8B/scalar-4                                476869863    2.345  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/32B/simd-4                                 422039916    2.676  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/32B/scalar-4                               213930672    5.661  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/64B/simd-4                                 424848363    2.867  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/64B/scalar-4                               100000000    10.09  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/512B/simd-4                                204526584    6.270  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/512B/scalar-4                               15132284    77.17  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/4096B/simd-4                                26726410    43.51  ns/op     0  B/op   0  allocs/op
Benchmark_MemchrNotWord/4096B/scalar-4                               1967492    605.1  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/8B/simd-4                                         238297198    5.078  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/8B/default-4                                      312963876    3.837  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/32B/simd-4                                        223399807    5.164  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/32B/default-4                                     280149112    4.328  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/64B/simd-4                                        181056466    6.644  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/64B/default-4                                     208156309    5.839  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/96B/simd-4                                        140714542    8.129  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/96B/default-4                                      59551779    22.04  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/128B/simd-4                                       121754152    9.157  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/128B/default-4                                     46335471    25.74  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/192B/simd-4                                       141431370    8.391  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/192B/default-4                                     30501302    39.63  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/512B/simd-4                                       100000000    10.81  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/512B/default-4                                     10538218    114.6  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/4096B/simd-4                                       32581954    38.96  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem/4096B/default-4                                     1352824    888.7  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/8B/paired-4                             199809792    5.930  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/8B/paired-stdlib-4                      291923088    4.013  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/paired-4                            143105328    8.440  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/paired-stdlib-4                     269666610    4.364  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/single-4                            178805293    6.793  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/32B/single-stdlib-4                     252354752    4.979  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/paired-4                            220193251    5.567  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/paired-stdlib-4                     203060490    6.002  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/single-4                             80880027    14.56  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/64B/single-stdlib-4                     100000000    11.58  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/paired-4                            190965673    5.694  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/paired-stdlib-4                      57875270    20.70  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/single-4                             65319981    17.73  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/96B/single-stdlib-4                      76707356    15.73  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/paired-4                           200099610    6.044  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/paired-stdlib-4                     38150725    26.76  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/single-4                            48030675    25.40  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/128B/single-stdlib-4                     56638509    21.41  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/paired-4                           199293573    6.187  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/paired-stdlib-4                     29952097    40.70  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/single-4                            30438082    39.93  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/192B/single-stdlib-4                     31924005    37.64  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/paired-4                           139797565    8.697  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/paired-stdlib-4                      9185761    122.8  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/single-4                             9594715    121.1  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/512B/single-stdlib-4                      9683748    121.5  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/paired-4                           32443600    39.69  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/paired-stdlib-4                     1267646    962.3  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/single-4                            1218538    931.7  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Prefilter/4096B/single-stdlib-4                     1293177    927.2  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Adversarial/simd-4                                     1089  1059972  ns/op     0  B/op   0  allocs/op
Benchmark_Memmem_Adversarial/default-4                                  1124  1045770  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/simd-4                                        493146400    2.510  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/8B/swar-4                                       1000000000    1.086  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/simd-4                                       549037777    2.331  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/32B/swar-4                                       657224274    1.863  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/simd-4                                       537544213    2.260  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/64B/swar-4                                       563848951    2.175  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/simd-4                                      335413507    3.293  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/512B/swar-4                                      149214904    7.730  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/4096B/simd-4                                      78255848    14.58  ns/op     0  B/op   0  allocs/op
Benchmark_IsASCII/4096B/swar-4                                      22478596    54.43  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/8B/simd-4                                  416852157    2.938  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/8B/swar-4                                  673407130    1.810  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/32B/simd-4                                 520230950    2.360  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/32B/swar-4                                 606508831    1.989  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/64B/simd-4                                 484232101    2.402  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/64B/swar-4                                 432093866    2.813  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/512B/simd-4                                354636026    3.395  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/512B/swar-4                                 79015945    15.00  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/4096B/simd-4                                83377725    16.17  ns/op     0  B/op   0  allocs/op
Benchmark_FirstNonASCII/4096B/swar-4                                10899777    110.7  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/8B/simd-4                                  407292187    2.925  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/8B/swar-4                                  633446112    1.876  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/32B/simd-4                                 527741985    2.267  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/32B/swar-4                                 520401512    2.198  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/64B/simd-4                                 528453115    2.183  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/64B/swar-4                                 376214196    3.276  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/512B/simd-4                                305556792    3.753  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/512B/swar-4                                 52532768    20.09  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/4096B/simd-4                                67882563    18.61  ns/op     0  B/op   0  allocs/op
Benchmark_CountNonASCII/4096B/swar-4                                 8308380    143.3  ns/op     0  B/op   0  allocs/op
```

<!-- skip-docs -->
## ☕ Supporters

Fiber is an open-source project that runs on donations to pay the bills, e.g., our domain name, hosting, and serverless infrastructure. If you want to support Fiber, please become a [GitHub Sponsor](https://github.com/sponsors/gofiber).

<p align="center">
  <a href="https://www.coderabbit.ai/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://www.coderabbit.ai/images/logo-dark.svg">
      <img width="280" height="52" alt="CodeRabbit" src="https://www.coderabbit.ai/images/logo-orange.svg">
    </picture>
  </a>
</p>
<p align="center">
  <a href="https://blacksmith.sh/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme">
    <img width="280" height="96" alt="Blacksmith" src="https://raw.githubusercontent.com/gofiber/.github/main/assets/sponsors/blacksmith.png">
  </a>
</p>
<p align="center">
  <sub><b>Tool Sponsors</b> - supporting Fiber with free IDE licenses and AI credits</sub>
</p>
<p align="center">
  <a href="https://www.jetbrains.com/?from=gofiber" title="JetBrains - IDE licenses"><img width="36" height="36" alt="JetBrains" src="https://github.com/JetBrains.png?size=72"></a>
  &nbsp;
  <a href="https://openai.com/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme" title="OpenAI - AI credits"><img width="36" height="36" alt="OpenAI" src="https://github.com/openai.png?size=72"></a>
  &nbsp;
  <a href="https://www.anthropic.com/?utm_source=gofiber&utm_medium=sponsor&utm_content=readme" title="Anthropic - AI credits"><img width="36" height="36" alt="Anthropic" src="https://github.com/anthropics.png?size=72"></a>
</p>

<!-- sponsors -->

### 📅 Monthly Sponsors

<table>
<tr><td valign="top"><strong>🔥 Fiber Guardian</strong></td><td><a href="https://www.coderabbit.ai/?utm_source=cr_org&amp;utm_medium=github" title="@coderabbitai"><img src="https://github.com/coderabbitai.png" width="50" alt="@coderabbitai" /></a></td></tr>
<tr><td valign="top"><strong>☕ Fiber Supporter</strong></td><td><a href="https://ndole.studio" title="@NdoleStudio"><img src="https://github.com/NdoleStudio.png" width="34" alt="@NdoleStudio" /></a></td></tr>
<tr><td valign="top"><strong>🪴 Fiber Friend</strong></td><td><a href="https://www.kamaslau.com" title="@kamaslau"><img src="https://github.com/kamaslau.png" width="32" alt="@kamaslau" /></a>&nbsp;<a href="https://github.com/simonheisstpeter" title="@simonheisstpeter"><img src="https://github.com/simonheisstpeter.png" width="32" alt="@simonheisstpeter" /></a></td></tr>
</table>

### 🎁 One-time Sponsors

<table>
<tr><td valign="top"><strong>🚀 Fiber Hero</strong></td><td><a href="https://www.thanks.dev" title="@thnxdev"><img src="https://github.com/thnxdev.png" width="40" alt="@thnxdev" /></a></td></tr>
<tr><td valign="top"><strong>🪴 Fiber Friend</strong></td><td><a href="https://github.com/Gl1tchedPixzl" title="@Gl1tchedPixzl"><img src="https://github.com/Gl1tchedPixzl.png" width="26" alt="@Gl1tchedPixzl" /></a></td></tr>
</table>
<!-- sponsors -->
<!-- skip-docs -->
