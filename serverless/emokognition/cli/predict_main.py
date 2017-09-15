import cv2
import json
import numpy as np
import ssl
import sys
import os

from urllib import request

from emotions import constants
from emotions import recognition
from emotions import utils

os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
network = recognition.EmotionRecognition()
network.build_network()
network.load_model_from_external_file("/code/cli/face_recognition_model.tflearn")

if __name__ == "__main__":
    data = json.loads(sys.stdin.read())
    emotion_dict = {}
    ctx = None
    if "https" in data["image_url"]:
        ctx = ssl.create_default_context()
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE

    url_response = request.urlopen(data["image_url"], context=ctx)
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
    msg = "I does seem like person on the image is {} or {}".format(s[0][0], s[1][0])
    print(json.dumps(dict(s)))
