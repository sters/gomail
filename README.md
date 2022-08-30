# Gomail

[![Build Status](https://travis-ci.org/Shopify/gomail.svg?branch=master)](https://travis-ci.org/Shopify/gomail) [![codecov](https://codecov.io/gh/Shopify/gomail/branch/master/graph/badge.svg)](https://codecov.io/gh/Shopify/gomail) [![Documentation](https://pkg.go.dev/github.com/Shopify/gomail?status.svg)](https://pkg.go.dev/github.com/Shopify/gomail)

This is an actively maintained fork of [Gomail][1] and includes fixes and
improvements for a number of outstanding issues. The current progress is
as follows:

- [x] Timeouts and retries can be specified outside of the 10 second default.
- [x] Proxying is supported through specifying a custom [NetDialTimeout][2].
- [ ] Filenames are properly encoded for non-ASCII characters.
- [ ] Email addresses are properly encoded for non-ASCII characters.
- [ ] Embedded files and attachments are tested for their existence.
- [ ] An `io.Reader` can be supplied when embedding and attaching files.

[1]: https://github.com/go-gomail/gomail
[2]: https://pkg.go.dev/github.com/Shopify/gomail#NetDialTimeout


## Introduction

Gomail is a simple and efficient package to send emails. It is well tested and
documented.

Gomail can only send emails using an SMTP server. But the API is flexible and it
is easy to implement other methods for sending emails using a local Postfix, an
API, etc.

It requires Go 1.18 or newer.


## Features

Gomail supports:

- Attachments
- Embedded images
- HTML and text templates
- Automatic encoding of special characters
- SSL and TLS
- Sending multiple emails with the same SMTP connection


## Documentation

https://pkg.go.dev/github.com/Shopify/gomail


## Examples

See the [examples in the documentation](https://pkg.go.dev/github.com/Shopify/gomail#example-package).


## FAQ

### x509: certificate signed by unknown authority

If you get this error it means the certificate used by the SMTP server is not
considered valid by the client running Gomail. As a quick workaround you can
bypass the verification of the server's certificate chain and host name by using
`SetTLSConfig`:

```go
package main

import (
	"crypto/tls"

	"github.com/Shopify/gomail"
)

func main() {
	d := mail.NewDialer("smtp.example.com", 587, "user", "123456")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send emails using d.
}
```

Note, however, that this is insecure and should not be used in production.

## Contribute

Contributions are more than welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for
more info.


## Change log

See [CHANGELOG.md](CHANGELOG.md).


## License

[MIT](LICENSE)
