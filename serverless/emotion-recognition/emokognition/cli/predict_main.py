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

import cv2
import numpy as np
import ssl
import sys
import os
import requests
import json

import fdk

from urllib import request
from emotions import constants
from emotions import recognition
from emotions import utils

os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
network = recognition.EmotionRecognition()
network.build_network()
network.load_model_from_external_file("/code/cli/face_recognition_model")


@fdk.coerce_input_to_content_type
def handler(context, data=None, loop=None):
    # NOTE: this really depends on content type,
    # i assume here that request body can be:
    # - plain/text
    # - application/json
    if isinstance(data, str):
        data = json.loads(data)
    if isinstance(data, dict):
        pass

    emotion_dict = {}
    ctx = None
    if "https" in data["media_url"]:
        ctx = ssl.create_default_context()
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE

    print("attempting to read image from the media URL", file=sys.stderr, flush=True)
    url_response = request.urlopen(data["media_url"], context=ctx)
    print("done reading image from the media URL: {0}".format(data["media_url"]),
          file=sys.stderr, flush=True)
    img = cv2.imdecode(
            np.array(bytearray(url_response.read()), dtype=np.uint8),
            cv2.COLOR_GRAY2BGR
    )
    frame = utils.format_image_for_prediction(img)
    if frame is None:
        print("Unable to detect face.", file=sys.stderr, flush=True)
        return
    result = network.predict(frame)
    for index, emotion in enumerate(constants.EMOTIONS):
        emotion_dict[emotion] = result[0][index]

    s = [(k, str(emotion_dict[k])) for k in
         sorted(emotion_dict, key=emotion_dict.get, reverse=True)]
    sys.stderr.write(json.dumps(dict(s)))
    main_emotion, _ = s[0]
    alt_emotion, _ = s[1]

    print("done with predictions, results: {0}, {1}"
          .format(main_emotion, alt_emotion), file=sys.stderr, flush=True)
    fn_app = os.environ.get("FN_APP_NAME")
    recorder = "{}/r/{}/recorder".format(os.environ.get("FN_API_URL"), fn_app)
    try:
        print("attempting to send prediction results "
              "to the next function", file=sys.stderr, flush=True)
        requests.post(recorder, json={
            "alt_emotion": alt_emotion,
            "main_emotion": main_emotion,
        })
        return "OK"
    except Exception as ex:
        print(str(ex), file=sys.stderr, flush=True)
        raise ex


if __name__ == "__main__":
    fdk.handle(handler)
