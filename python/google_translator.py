import argparse

from deep_translator import GoogleTranslator


def main(args):
    translator = GoogleTranslator(source='en', target='russian')
    print(translator.translate(args.text))


def init():
    parser = argparse.ArgumentParser(description='Translate blobber')
    parser.add_argument('-t', '--text', action='store', dest='text', type=str, default='', help='Text for blob')
    args = parser.parse_args()

    try:
        main(args)
    except Exception as err:
        print(err)


if __name__ == '__main__':
    init()
