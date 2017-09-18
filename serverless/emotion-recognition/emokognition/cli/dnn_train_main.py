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
