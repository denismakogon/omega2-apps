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

import requests

from fdk.application import decorators


@decorators.fn_app
class Application(object):

    def __init__(self, *args, **kwargs):
        pass

    @decorators.fn(fn_type="sync")
    def square(self, x, y, *args, **kwargs):
        return x * y

    @decorators.fn(fn_type="sync", dependencies={
        "requests": requests
    })
    def request(self, *args, **kwargs):
        cached_dependencies = kwargs["dependencies"]
        requests = cached_dependencies.get("requests")
        r = requests.get('https://api.github.com/events')
        r.raise_for_status()
        return r.text


if __name__ == "__main__":

    app = Application(config={})

    # res, err = app.square(10, 20)
    # if err:
    #     raise err
    # print(res)

    res, err = app.request()
    if err:
        raise err
    print(res)
