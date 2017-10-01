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

import numpy as np
import pandas as pd

from PIL import Image

from emotions import constants
from emotions import utils


def emotion_to_vec(x):
    d = np.zeros(len(constants.EMOTIONS))
    d[x] = 1.0
    return d


def data_to_image(numpy_array):
    image_from_data = np.array(
        Image.fromarray(
            np.fromstring(
                str(numpy_array),
                dtype=np.uint8,
                sep=' ').reshape(
                (constants.SIZE_FACE, constants.SIZE_FACE))
        ).convert('RGB')
    )
    copy_of_image = image_from_data[:, :, ::-1].copy()
    data_image = utils.format_image_for_learning(copy_of_image)
    return data_image


if __name__ == "__main__":
    FILE_PATH = './data/fer2013.csv'
    data = pd.read_csv(FILE_PATH)
    labels = []
    images = []
    index = 1
    total = data.shape[0]
    for index, row in data.iterrows():
        emotion = emotion_to_vec(row['emotion'])
        image = data_to_image(row['pixels'])
        if image is not None:
            labels.append(emotion)
            images.append(image)
        index += 1
        print("Progress: {}/{} {:.2f}%".format(index, total, index * 100.0 / total))

    np.save('./data/data.npy', images)
    np.save('./data/labels.npy', labels)
