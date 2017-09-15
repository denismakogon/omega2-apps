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
