# Fast CORS

[![License][badge-1-img]][badge-1-link]
[![go.dev][badge-2-img]][badge-2-link]
[![Go Report Card][badge-3-img]][badge-3-link]

CORS middleware for [fasthttp][1] server. This project is under
construction and shouldn't be used in production environment yet.

This lib is inspired by [rs/cors][2] and
[AdhityaRamadhanus/fasthttpcors][3].

## To-do

- [ ] Improve documentation
- [ ] Improve tests

## Source code

The main repository is on [GitHub][source-github].

## License and contribution

All original code is licensed under the
[Mozilla Public License 2.0][license-code].

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you shall be licensed under the same license as
the material to which it is contributed, without any additional terms or
conditions.

### License history

Starting with v1.0.0, this project is licensed under MPL-2.0

### Third-party and exceptions

The `lib/cors/` folder contains code from [rs/cors][2], licensed under the MIT
License. See `lib/cors/LICENSE`.

[1]: https://github.com/valyala/fasthttp
[2]: https://github.com/rs/cors/
[3]: https://github.com/AdhityaRamadhanus/fasthttpcors

[badge-1-img]: https://img.shields.io/badge/license-MPL--2.0-green?style=flat-square
[badge-1-link]: ./LICENSE
[badge-2-img]: https://img.shields.io/badge/go.dev-reference-007d9c?style=flat-square&logo=go&logoColor=white
[badge-2-link]: https://pkg.go.dev/github.com/paemuri/fastcors
[badge-3-img]: https://goreportcard.com/badge/github.com/paemuri/fastcors?style=flat-square
[badge-3-link]: https://goreportcard.com/report/github.com/paemuri/fastcors

[source-github]: https://github.com/paemuri/fastcors

[license-code]: ./LICENSES/MPL-2.0.txt

