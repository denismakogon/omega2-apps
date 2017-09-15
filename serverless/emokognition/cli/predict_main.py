import cv2
import json
import numpy as np
import ssl
import sys
import os
import requests

from urllib import request
from emotions import constants
from emotions import recognition
from emotions import utils

os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
network = recognition.EmotionRecognition()
network.build_network()
network.load_model_from_external_file("/code/cli/face_recognition_model")

if __name__ == "__main__":
    data = json.loads(sys.stdin.read())
    emotion_dict = {}
    ctx = None
    if "https" in data["media_url"]:
        ctx = ssl.create_default_context()
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE

    url_response = request.urlopen(data["media_url"], context=ctx)
    img = cv2.imdecode(
            np.array(bytearray(url_response.read()), dtype=np.uint8),
            cv2.COLOR_GRAY2BGR
    )
    frame = utils.format_image_for_prediction(img)
    if frame is None:
        print("Unable to detect face.")
        exit(1)
    result = network.predict(frame)
    for index, emotion in enumerate(constants.EMOTIONS):
        emotion_dict[emotion] = result[0][index]

    s = [(k, str(emotion_dict[k])) for k in sorted(emotion_dict, key=emotion_dict.get, reverse=True)]
    sys.stderr.write(json.dumps(dict(s)))
    main_emotion, _ = s[0]
    alt_emotion, _ = s[1]
    recorder = "{}/r/emokognition/recorder".format(data.get("api_url"))
    try:
        requests.post(recorder, json={
            "alt_emotion": alt_emotion,
            "main_emotion": main_emotion,
        })
        print("OK")
    except Exception as ex:
        sys.stderr.write(str(ex))
