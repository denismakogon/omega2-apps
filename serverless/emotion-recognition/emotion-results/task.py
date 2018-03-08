# All Rights Reserved.
#
#    Licensed under the Apache License, Version 2.0 (the "License"); you may
#    not use this file except in compliance with the License. You may obtain
#    a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#    License for the specific language governing permissions and limitations
#    under the License.

import os
import sys
import psycopg2
import collections
import json
import fdk


CREATE = ("CREATE TABLE IF NOT EXISTS emotions ("
          "id SERIAL, "
          "main_emotion VARCHAR(255) NOT NULL, "
          "alt_emotion VARCHAR(255) NOT NULL)")

MAIN_SELECT = {
    "q": "SELECT main_emotion, COUNT(id) AS count FROM emotions GROUP BY main_emotion",
    "name": "main"
}

ALT_SELECT = {
    "q": "SELECT alt_emotion, COUNT(id) AS count FROM emotions GROUP BY alt_emotion",
    "name": "alt"
}


def select_votes(context, data=None, loop=None):
    print("entering function\n", file=sys.stderr, flush=True)
    pg_host = os.environ.get('postgres_host'.upper())
    pg_port = os.environ.get('postgres_port'.upper())
    pg_db = os.environ.get('postgres_db'.upper())
    pg_user = os.environ.get('postgres_user'.upper())
    pg_pswd = os.environ.get('postgres_password'.upper())
    pg_dns = (
        'dbname={database} '
        'user={user} '
        'password={passwd} '
        'host={host} '
        'port={port}'
        .format(host=pg_host,
                database=pg_db,
                user=pg_user,
                passwd=pg_pswd,
                port=pg_port)
    )
    final = {}
    print("establishing connection\n", file=sys.stderr, flush=True)
    conn = psycopg2.connect(pg_dns)
    print("connection acquired\n", file=sys.stderr, flush=True)
    cur = conn.cursor()
    print("cursor created\n", file=sys.stderr, flush=True)
    cur.execute(CREATE)
    for o in [MAIN_SELECT, ALT_SELECT]:
        cur.execute(o["q"])
        result = collections.defaultdict(int)
        rows = cur.fetchall()
        for row in rows:
            emotion, count = row
            result[emotion] = count

        full_result = dict(result)
        final[o["name"]] = full_result
    print("stats created\n", file=sys.stderr, flush=True)
    return json.dumps(final)


if __name__ == "__main__":
    fdk.handle(select_votes)
