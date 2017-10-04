import jinja2
import requests
import os
import sys
import asyncio

from hotfn.http import response
from hotfn.http import worker

loader = jinja2.FileSystemLoader('./index.html')
env = jinja2.Environment(loader=loader)
template = env.get_template('')

fn_app = os.environ.get("FN_APP_NAME")
recorder = "{0}/r/{1}/results".format(os.environ.get("FN_API_URL"), fn_app)


async def build_view(context, data=None, loop=None):
    print("entering coroutine\n", file=sys.stderr, flush=True)
    resp = requests.get(recorder)
    resp.raise_for_status()
    print("stats received\n", file=sys.stderr, flush=True)
    data = resp.json()
    main_emotions = data['main']
    alt_emotions = data['alt']
    main, alt = [], []
    total = sum(list(main_emotions.values()))
    for emotion, count in main_emotions.items():
        main.append(dict(emotion=emotion,
                         stat=float("{:.2f}".format(float(count / total) * 100)),
                         times=count))
    for emotion, count in alt_emotions.items():
        alt.append(dict(emotion=emotion,
                        stat=float("{:.2f}".format(float(count / total) * 100)),
                        times=count))
    print("final stats assembled\n", file=sys.stderr, flush=True)
    render_context = {
        "main_emotions": main,
        "alt_emotions": alt,
        "total": total
    }
    return response.RawResponse(context.version, 200, "OK", http_headers={
        "Content-Type": "text/html",
    }, response_data=template.render(render_context))


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    worker.run(build_view, loop=loop)
