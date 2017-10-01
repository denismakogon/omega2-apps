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

import sys

from emotions.emotions import recognition


def show_usage():
    print('[!] Usage: python dnn_train_main.py')
    print('\t dnn_train_main.py train \t Trains and saves model with saved dataset')


if __name__ == "__main__":
    if len(sys.argv) <= 1:
        show_usage()
        exit()

    network = recognition.EmotionRecognition()
    if sys.argv[1] == 'train':
        network.start_training()
        network.save_model()
    elif sys.argv[1] == 'help':
        show_usage()
