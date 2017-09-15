import asyncio
import collections

import json
import os
import sys

import aiopg

MAIN_SELECT = {
    "q": "SELECT main_emotion, COUNT(id) AS count FROM emotions GROUP BY main_emotion",
    "name": "main"
}

ALT_SELECT = {
    "q": "SELECT alt_emotion, COUNT(id) AS count FROM emotions GROUP BY alt_emotion",
    "name": "alt"
}


async def select_votes(pg_dns):
    final = {}
    async with aiopg.create_pool(pg_dns) as pool:
        async with pool.acquire() as conn:
            async with conn.cursor() as cur:
                for o in [MAIN_SELECT, ALT_SELECT]:
                    await cur.execute(o["q"])
                    result = collections.defaultdict(int)
                    async for row in cur:
                        emotion, count = row
                        result[emotion] = count
                    full_result = dict(result)
                    final[o["name"]] = full_result
                return final


if __name__ == "__main__":
    if not os.isatty(sys.stdin.fileno()):
        pg_host = os.environ.get('PG_HOST'.lower())
        pg_port = os.environ.get('PG_PORT'.lower())
        pg_db = os.environ.get('PG_DB'.lower())
        pg_user = os.environ.get('PG_USER'.lower())
        pg_pswd = os.environ.get('PG_PSWD'.lower())
        pg_dns = (
            'dbname={database} user={user} password={passwd} host={host}'
            .format(host=pg_host, database=pg_db, user=pg_user, passwd=pg_pswd))

        loop = asyncio.get_event_loop()
        result = loop.run_until_complete(select_votes(pg_dns))
        print(json.dumps(result))
