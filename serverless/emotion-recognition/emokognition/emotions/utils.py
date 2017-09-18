import numpy as np

import cv2

from emotions import constants

cascade_classifier = cv2.CascadeClassifier(
    "/usr/local/share/OpenCV/haarcascades/haarcascade_frontalface_default.xml")


def format_image_for_learning(frame):
    if len(frame.shape) > 2 and frame.shape[2] == 3:
        loaded_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
    else:
        loaded_frame = cv2.cvtColor(frame, cv2.IMREAD_GRAYSCALE)

    gray_border = np.zeros((150, 150), np.uint8)
    gray_border[:, :] = 200
    face_border_lower = int((150 / 2 - constants.SIZE_FACE / 2))
    face_border_upper = int((150 + constants.SIZE_FACE) / 2)
    gray_border[
        face_border_lower:face_border_upper,
        face_border_lower:face_border_upper
    ] = loaded_frame
    faces = cascade_classifier.detectMultiScale(gray_border, scaleFactor=1.3, minNeighbors=5)
    if not len(faces) > 0:
        return
    max_area_face = faces[0]
    for face in faces:
        if face[2] * face[3] > max_area_face[2] * max_area_face[3]:
            max_area_face = face
    frame_gray_border = gray_border[
            max_area_face[1]:(max_area_face[1] + max_area_face[2]),
            max_area_face[0]:(max_area_face[0] + max_area_face[3])
    ]
    try:
        image_resize = cv2.resize(frame_gray_border,
                                  (constants.SIZE_FACE, constants.SIZE_FACE),
                                  interpolation=cv2.INTER_CUBIC) / 255.
    except Exception as ex:
        print("[+] Problem during resize, err: ", str(ex))
        return None

    return image_resize


def format_image_for_prediction(image):
    if len(image.shape) > 2 and image.shape[2] == 3:
        image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    else:
        image = cv2.imdecode(image, cv2.IMREAD_GRAYSCALE)

    attempts = 100
    faces = []

    for i in range(1, attempts + 1):
        faces = cascade_classifier.detectMultiScale(
            image,
            scaleFactor=1.0 + 0.001 * i,
            minNeighbors=8,
            minSize=(55, 55),
            flags=cv2.CASCADE_SCALE_IMAGE,
        )
        if len(faces) > 0:
            break

    if not len(faces) > 0:
        return None
    max_area_face = faces[0]
    for face in faces:
        if face[2] * face[3] > max_area_face[2] * max_area_face[3]:
                max_area_face = face
    face = max_area_face
    image = image[face[1]:(face[1] + face[2]), face[0]:(face[0] + face[3])]
    try:
        image = cv2.resize(image, (48, 48), interpolation=cv2.INTER_CUBIC) / 255.
    except Exception as ex:
        print("[+] Problem during resize: ", str(ex))
        return None

    return image
