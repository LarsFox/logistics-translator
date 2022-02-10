import argparse
import json
import re


PATH = re.compile(' +d="(.*)"')
HEIGHT = re.compile(' +height="(.*)"')
WIDTH = re.compile(' +width="(.*)"')


def main(args):
    with open('svg/' + args.input + '.svg', 'r') as target:
        lines = target.readlines()

    if args.reverse:
        lines = lines[::-1]

    with open('jsons/' + args.input + '.json') as target:
        src = json.load(target)

    parsed = []
    for line in lines:
        match = HEIGHT.match(line)
        if match:
            src['height'] = int(float(match.groups()[0]))

        match = WIDTH.match(line)
        if match:
            src['width'] = int(float(match.groups()[0]))

        match = PATH.match(line)
        if match:
            parsed.append(match.groups()[0])
            continue


    print("found {} paths, had {}".format(len(parsed), len(src['paths'])))
    src['paths'] = parsed

    with open('jsons/' + args.input + '.json', 'w') as target:
        json.dump(src, target, ensure_ascii=False, indent=2)


def init():
    parser = argparse.ArgumentParser(description='SVG to SQL converter')
    parser.add_argument('-i', '--input', action='store', dest='input', type=str, default='', help='Path to input file.')

    parser.add_argument('-r', dest='reverse', action='store_true')
    parser.set_defaults(reverse=False)


    args = parser.parse_args()

    try:
        main(args)
    except Exception as err:
        print(err)


if __name__ == '__main__':
    init()
