<h1 align="center">GTP - Go one-Time Password</h1>

Simple ⏱Time-based 🔑OTP Client

![toproject](https://img.shields.io/badge/toyproject-doing-green) ![gonewbie](https://img.shields.io/badge/golang-newbie-blue) ![build](https://img.shields.io/github/workflow/status/jonnung/gtp/1)  


## What is GTP?
GTP is a time-based OTP [CLI](https://en.wikipedia.org/wiki/Command-line_interface) client.  
time-based OTP is an algorithm implementation defined in the [RFC 6238](https://tools.ietf.org/html/rfc6238).
If you received a secret text from [multi-factor authentication](http://en.wikipedia.org/wiki/Multi-factor_authentication) system, GTP is store the secret text on your computer($HOME directory).  
The secret text is secure data, but GTP expects your computer to be used only you.

![Alt text](./gtp_usage.svg)

## Installation
```
$ go get github.com/jonnung/gtp
```

or download compiled binary file for multiple platform. See [release page](https://github.com/jonnung/gtp/releases).


## Usage
```
$ gtp [{number}|list|add|remove|clear]
```

- `{number}`: Show time based one-time password by specified secret
- `list`: All registered OTP secrets
- `add`: Add new OTP secret
- `remove`: Remove the specified secret
- `clear`: Clear all secrets
