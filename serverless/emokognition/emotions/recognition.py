import numpy as np  # WTF? Is it makes code work?
import os

import tflearn
from tflearn.layers.core import input_data, dropout, fully_connected
from tflearn.layers.conv import conv_2d, max_pool_2d
from tflearn.layers.estimator import regression

from emotions import dataset_loader
from emotions import constants


class EmotionRecognition(object):

    def __init__(self):
        self.model = None
        self.network = None
        self.dataset = dataset_loader.DatasetLoader()

    def build_network(self):
        data = input_data(shape=[None, constants.SIZE_FACE, constants.SIZE_FACE, 1])
        dd_data = conv_2d(data, 64, 5, activation='relu')
        max_pool_dd = max_pool_2d(dd_data, 3, strides=2)
        conv_dd = conv_2d(max_pool_dd, 64, 5, activation='relu')
        max_pool_conv_dd = max_pool_2d(conv_dd, 3, strides=2)
        conv_dd_max_pool_conv_dd = conv_2d(max_pool_conv_dd, 128, 4, activation='relu')
        conv_dd_max_pool_conv_dd_dropout = dropout(conv_dd_max_pool_conv_dd, 0.3)
        fully_connected_conv_dd_max_pool_conv_dd_dropout = fully_connected(
            conv_dd_max_pool_conv_dd_dropout, 3072, activation='relu')
        a = fully_connected(
            fully_connected_conv_dd_max_pool_conv_dd_dropout, len(constants.EMOTIONS), activation='softmax')
        self.network = regression(a, optimizer='momentum')
        self.model = tflearn.DNN(self.network,
                                 checkpoint_path=constants.SAVE_DIRECTORY + '/emotion_recognition',
                                 max_checkpoints=1,
                                 tensorboard_verbose=2)

    def load_saved_dataset(self):
        self.dataset.load_from_save()
        print('[+] Dataset found and loaded')

    def start_training(self):
        self.load_saved_dataset()
        self.build_network()
        print('[+] Training network')
        self.model.fit(
          self.dataset.images, self.dataset.labels,
          n_epoch=100,
          batch_size=50,
          shuffle=True,
          show_metric=True,
          snapshot_step=200,
          snapshot_epoch=True,
          run_id='emotion_recognition'
        )

    def predict(self, image):
        image = image.reshape([-1, constants.SIZE_FACE, constants.SIZE_FACE, 1])
        return self.model.predict(image)

    def save_model(self):
        self.model.save(os.path.join(constants.SAVE_DIRECTORY, constants.SAVE_MODEL_FILENAME))
        print('[+] Model trained and saved at ' + constants.SAVE_MODEL_FILENAME)

    def load_model_from_external_file(self, path):
        if os.path.isfile(path):
            self.model.load(path)
