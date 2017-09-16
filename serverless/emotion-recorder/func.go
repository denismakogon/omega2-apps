package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var emotionsTable = `CREATE TABLE IF NOT EXISTS emotions (id serial NOT NULL, main_emotion VARCHAR(255) NOT NULL, alt_emotion VARCHAR(255) NOT NULL)`

func writeBadResponse(buf *bytes.Buffer, resp *http.Response, errMsg string) {
	resp.StatusCode = 500
	resp.Status = http.StatusText(resp.StatusCode)
	fmt.Fprintln(buf, errMsg)
	fmt.Fprintf(os.Stderr, errMsg)
}

func main() {
	res := http.Response{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		StatusCode: 200,
		Status:     "OK",
	}
	var buf bytes.Buffer
	pgConf := new(api.PostgresConfig)
	err := pgConf.FromEnv()
	if err != nil {
		writeBadResponse(&buf, &res,
			fmt.Sprintf("Unable to setup PG struct from env. Error: %s", err.Error()))
		res.Body = ioutil.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Write(os.Stdout)
		return
	}
	pg_dns := pgConf.DNS()
	db, err := sqlx.Open("postgres", pg_dns)
	if err != nil {
		writeBadResponse(&buf, &res,
			fmt.Sprintf("Unable to talk to PG by DNS %s, error: %s", pg_dns, err.Error()))
		res.Body = ioutil.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Write(os.Stdout)
		return
	} else {
		defer db.Close()
		for {
			res := http.Response{
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				StatusCode: 200,
				Status:     "OK",
			}
			r := bufio.NewReader(os.Stdin)
			req, err := http.ReadRequest(r)
			var buf bytes.Buffer
			if err != nil {
				writeBadResponse(&buf, &res,
					fmt.Sprintf("Unable to read request from STDIN, "+
						"it might be empty. Error: %v", err.Error()))
			} else {
				l, err := strconv.Atoi(req.Header.Get("Content-Length"))
				if err != nil {
					writeBadResponse(&buf, &res,
						fmt.Sprintf("Failed to convert content length to int: %s", err.Error()))
				} else {
					fmt.Fprintf(os.Stderr, "Income content len: %v\n", l)
					p := make([]byte, l)
					_, err = r.Read(p)
					if err != nil {
						writeBadResponse(&buf, &res,
							fmt.Sprintf("Unable to read request data, error: %s", err.Error()))
					} else {
						payload := &api.RequestPayload{}
						err = json.Unmarshal(p, payload)
						if err != nil {
							if e, ok := err.(*json.SyntaxError); ok {
								writeBadResponse(&buf, &res,
									fmt.Sprintf("syntax error at byte offset %d", e.Offset))
							} else {
								writeBadResponse(&buf, &res,
									fmt.Sprintf("Unable to decode input object, error: %s", err.Error()))
							}
						} else {
							if err != nil {
								writeBadResponse(&buf, &res,
									fmt.Sprintf("Unable to talk to PG, error: %s", err.Error()))
							} else {
								_, err = db.Exec(emotionsTable)
								if err != nil {
									writeBadResponse(&buf, &res, fmt.Sprintf("Unable to create table, error: %s", err.Error()))
								} else {
									q := db.Rebind("INSERT INTO emotions (main_emotion, alt_emotion) VALUES (?, ?);")
									_, err = db.Exec(q, payload.MainEmotion, payload.AltEmotion)
									if err != nil {
										writeBadResponse(&buf, &res, err.Error())
									} else {
										fmt.Fprint(&buf, "OK\n")
									}
								}
							}
						}
					}
				}
			}
			res.Body = ioutil.NopCloser(&buf)
			res.ContentLength = int64(buf.Len())
			res.Write(os.Stdout)
		}
	}
}
