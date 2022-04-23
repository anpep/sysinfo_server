# sysinfo_server
> Provides access to system information and measurements

![version: 1.0.0](https://img.shields.io/badge/version-1.0.0-blue.svg)
![license: GNU GPL v2](https://img.shields.io/badge/license-GNU_GPL_v2-brightgreen.svg)

# What's this?
`sysinfo_server` is a simple REST API application written in Go, which retrieves system information and parameters.
Currently, only system bootup duration is implemented.

This application was written for the technical assessment at Canonical.

# Build instructions
You will require Go 1.13 or newer. This application was built and tested with Go version 1.18.1.

```shell
# Fetch the application sources
$ git clone https://github.com/anpep/sysinfo_server && cd sysinfo_server
# Build the application
$ go build
```

# Usage
The application will provide a REST API listening at port `8080`. These are the implemented routes:
- GET `/<param_name>`
  - Returns the plaintext value of the system parameter with name `param_name`. If an error occurs, the error message will be provided with an error status code.
- GET `/<param_name>.json`
  - Returns the value of the system parameter with name `param_name` as a JSON response with the following schema:
    - `ok: boolean` - Whether or not the request succeeded. If `false`, the `error` field must be present in the response body.
    - `error: string?` - If present, a message describing an error condition.
    - `param: {name: string, value: any}` - If present, the data of the system information parameter.
      - `name: string` - Name of the parameter. Must be equal to the `param_name` value in the request URI.
      - `value: any` - The value of the parameter.

## Sample execution output
```shell
# Start the server
$ ./sysinfo_server &
[1] 7335

# No parameter
$ curl http://localhost:8080/
no such parameter
$ curl http://localhost:8080/index.json
{"ok":false,"error":"no such parameter"}

# Version parameter
$ curl http://localhost:8080/version
1.0.0
$ curl http://localhost:8080/version.json
{"ok":true,"param":{"name":"version","value":"1.0.0"}}

# Boot duration parameter
$ curl http://localhost:8080/duration
9.566
$ curl http://localhost:8080/duration.json
{"ok":true,"param":{"name":"duration","value":9.566}}
```

## License
`sysinfo_server` is licensed under the GNU General Public License v2.

```
sysinfo_server -- Provides access to system information and measurements
Copyright (c) 2022 Ángel Pérez <ap@anpep.co>

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
```