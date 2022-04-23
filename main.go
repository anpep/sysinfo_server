// sysinfo_server -- Provides access to system information and measurements
// Copyright (c) 2022 Ángel Pérez <ap@anpep.co>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type sysInfoParameter struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type response struct {
	Ok    bool              `json:"ok"`
	Error string            `json:"error,omitempty"`
	Param *sysInfoParameter `json:"param,omitempty"`
}

func writeResponse(res response, w http.ResponseWriter, encodeAsJson bool) (bytesWritten int, err error) {
	if encodeAsJson {
		jsonResponse, err := json.Marshal(res)
		if err != nil {
			return 0, err
		}

		return fmt.Fprintf(w, "%s\n", string(jsonResponse))
	}

	if res.Error != "" {
		// Write error string
		return fmt.Fprintf(w, "%s\n", res.Error)
	}

	if res.Param != nil {
		// Write parameter value
		return fmt.Fprintf(w, "%v\n", res.Param.Value)
	}

	return 0, nil
}

func getSysInfoParameter(paramName string) (*sysInfoParameter, error) {
	switch paramName {
	case "version":
		// Application version
		return &sysInfoParameter{Name: paramName, Value: "1.0.0"}, nil
	case "duration":
		// Boot time
		var cmd = exec.Command("systemd-analyze", "time")
		var cmdOut bytes.Buffer

		cmd.Stdout = &cmdOut
		err := cmd.Run()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not execute %s: %s\n", cmd.Path, err.Error())
			return nil, errors.New("could not measure boot duration")
		}

		regExp := regexp.MustCompile("=\\s*([\\d.]+)s")
		bootDuration, err := strconv.ParseFloat(regExp.FindStringSubmatch(cmdOut.String())[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not parse boot duration: %s\n", err.Error())
			return nil, errors.New("could not measure boot duration")
		}

		return &sysInfoParameter{Name: paramName, Value: bootDuration}, nil
	default:
		return nil, errors.New("no such parameter")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var paramName = strings.TrimPrefix(r.URL.Path, "/")
	var isJsonRequested = strings.HasSuffix(paramName, ".json")

	if isJsonRequested {
		paramName = strings.TrimSuffix(paramName, ".json")
	}

	param, err := getSysInfoParameter(paramName)
	if err != nil {
		w.WriteHeader(404)
		_, err := writeResponse(response{Ok: false, Error: err.Error(), Param: nil}, w, isJsonRequested)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write response: %s\n", err.Error())
			return
		}
		return
	}

	_, err = writeResponse(response{Ok: true, Error: "", Param: param}, w, isJsonRequested)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not write response: %s\n", err.Error())
		return
	}
}

func main() {
	http.Handle("/", http.HandlerFunc(handler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not serve: %s\n", err.Error())
		return
	}
}
