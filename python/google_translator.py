import argparse

from deep_translator import GoogleTranslator


def main(args):
    translator = GoogleTranslator(source=args.source, target=args.destination)
    print(translator.translate(args.text))


def init():
    parser = argparse.ArgumentParser(description='Translate blobber')
    parser.add_argument('-t', '--text', action='store', dest='text', type=str, default='')
    parser.add_argument('-s', '--source', action='store', dest='source', type=str, default='')
    parser.add_argument('-d', '--destination', action='store', dest='destination', type=str, default='')
    args = parser.parse_args()

    try:
        main(args)
    except Exception as err:
        print(err)


if __name__ == '__main__':
    init()
