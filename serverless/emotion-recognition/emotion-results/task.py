import asyncio
import collections

import json
import os
import sys

import aiopg

import fdk

from fdk.http import response


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


async def select_votes(context, data=None, loop=None):
    print("Entering coroutine\n", file=sys.stderr, flush=True)
    pg_host = os.environ.get('pg_host')
    pg_port = os.environ.get('pg_port')
    pg_db = os.environ.get('pg_db')
    pg_user = os.environ.get('pg_user')
    pg_pswd = os.environ.get('pg_pswd')
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
    print("Establishing connection\n", file=sys.stderr, flush=True)
    final = {}
    async with aiopg.create_pool(pg_dns) as pool:
        print("pool created\n", file=sys.stderr, flush=True)
        async with pool.acquire() as conn:
            print("connection acquired\n", file=sys.stderr, flush=True)
            async with conn.cursor() as cur:
                print("cursor created\n", file=sys.stderr, flush=True)
                await cur.execute(CREATE)
                for o in [MAIN_SELECT, ALT_SELECT]:
                    await cur.execute(o["q"])
                    result = collections.defaultdict(int)
                    async for row in cur:
                        emotion, count = row
                        result[emotion] = count
                    full_result = dict(result)
                    final[o["name"]] = full_result
            print("stats created\n", file=sys.stderr, flush=True)
            return response.RawResponse(context.version, 200, "OK", http_headers={
                "Content-Type": "application/json; charset=utf-8",
            }, response_data=json.dumps(final))


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    fdk.handle(select_votes, loop=loop)
