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

import jinja2
import requests
import os
import sys

import fdk

from fdk import response

loader = jinja2.FileSystemLoader('./index.html')
env = jinja2.Environment(loader=loader)
template = env.get_template('')

fn_app = os.environ.get("FN_APP_NAME")
recorder = "{0}/r/{1}/results".format(os.environ.get("FN_API_URL"), fn_app)


def build_view(context, data=None, loop=None):
    print("entering the function \n", file=sys.stderr, flush=True)
    resp = requests.get(recorder, timeout=200)
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
    headers = {
        "Content-Type": "text/html",
    }
    return response.RawResponse(
        context,
        status_code=200,
        headers=headers,
        response_data=template.render(render_context)
    )


if __name__ == "__main__":
    fdk.handle(build_view)
