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

import os
import numpy as np

from emotions.constants import *


class DatasetLoader(object):

    def __init__(self):
        self._images = None
        self._labels = None

    def load_from_save(self):
        self._images = np.load(
            os.path.join(SAVE_DIRECTORY, SAVE_DATASET_IMAGES_FILENAME)
        ).reshape([-1, SIZE_FACE, SIZE_FACE, 1])
        self._labels = np.load(
            os.path.join(SAVE_DIRECTORY, SAVE_DATASET_LABELS_FILENAME)
        ).reshape([-1, len(EMOTIONS)])

    @property
    def images(self):
        return self._images

    @property
    def labels(self):
        return self._labels
