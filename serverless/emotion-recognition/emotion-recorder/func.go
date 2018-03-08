package main

import (
	"context"
	"encoding/json"
	"github.com/denismakogon/omega2-apps/serverless/twitter-daemon/api"
	"github.com/fnproject/fdk-go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

var emotionsTable = `CREATE TABLE IF NOT EXISTS emotions (id serial NOT NULL, main_emotion VARCHAR(255) NOT NULL, alt_emotion VARCHAR(255) NOT NULL)`

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	pgConf := new(api.PostgresConfig)
	err := pgConf.FromEnv()

	if err != nil {
		fdk.WriteStatus(out, 500)
		out.Write([]byte(err.Error()))
		return
	}
	pg_dns := pgConf.DNS()
	db, err := sqlx.Open("postgres", pg_dns)
	if err != nil {
		fdk.WriteStatus(out, 500)
		out.Write([]byte(err.Error()))
		return
	}
	defer db.Close()

	_, err = db.Exec(emotionsTable)
	if err != nil {
		fdk.WriteStatus(out, 500)
		out.Write([]byte(err.Error()))
		return
	}

	var payload api.RequestPayload
	err = json.NewDecoder(in).Decode(&payload)
	if err != nil {
		fdk.WriteStatus(out, 500)
		out.Write([]byte(err.Error()))
		return
	}

	q := db.Rebind("INSERT INTO emotions (main_emotion, alt_emotion) VALUES (?, ?);")
	_, err = db.Exec(q, payload.MainEmotion, payload.AltEmotion)
	if err != nil {
		fdk.WriteStatus(out, 500)
		out.Write([]byte(err.Error()))
		return
	}

	out.Write([]byte("OK"))
}
