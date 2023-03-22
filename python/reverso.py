import argparse

from reverso_api.context import ReversoContextAPI


def main(args):
    api = ReversoContextAPI(
        args.text,
        "",
        "en",
        "ru"
    )

    if args.command == "translation" or args.command == "t":
        translation(api)
    elif args.command == "example" or args.command == "e":
        example(api)


def example(api):
    for ex in api.get_examples():
        print(ex[1].text)
        return


def translation(api):
    for ex in api.get_translations():
        print(ex.translation)
        return


def init():
    parser = argparse.ArgumentParser(description='Translate with reverso')
    parser.add_argument('-t', '--text', action='store', dest='text', type=str, default='', help='Text for blob')
    parser.add_argument('-c', '--command', action='store', dest='command', type=str, default='', help='Get example or translation')
    args = parser.parse_args()

    try:
        main(args)
    except Exception as err:
        print(err)


if __name__ == '__main__':
    init()
