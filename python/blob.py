import argparse
import json

from textblob import TextBlob
from textblob_de import TextBlobDE

# CC coordinating conjunction
# CD cardinal digit
# DT determiner
# EX existential there (like: “there is” … think of it like “there exists”)
# FW foreign word
# IN preposition/subordinating conjunction
# JJ adjective ‘big’
# JJR adjective, comparative ‘bigger’
# JJS adjective, superlative ‘biggest’
# LS list marker 1)
# MD modal could, will
# NN noun, singular ‘desk’
# NNS noun plural ‘desks’
# NNP proper noun, singular ‘Harrison’
# NNPS proper noun, plural ‘Americans’
# PDT predeterminer ‘all the kids’
# POS possessive ending parent‘s
# PRP personal pronoun I, he, she
# PRP$ possessive pronoun my, his, hers
# RB adverb very, silently,
# RBR adverb, comparative better
# RBS adverb, superlative best
# RP particle give up
# TO to go ‘to‘ the store.
# UH interjection errrrrrrrm
# VB verb, base form take
# VBD verb, past tense took
# VBG verb, gerund/present participle taking
# VBN verb, past participle taken
# VBP verb, sing. present, non-3d take
# VBZ verb, 3rd person sing. present takes
# WDT wh-determiner which
# WP wh-pronoun who, what
# WP$ possessive wh-pronoun whose
# WRB wh-abverb where, when

def main(args):
    sentences = args.text.split('. ')

    if args.source == "de":
        blob = [TextBlobDE(s).tags for s in sentences]
    else:
        blob = [TextBlob(s).tags for s in sentences]

    print(json.dumps({
        "tags": blob
    }, ensure_ascii=False))


def init():
    parser = argparse.ArgumentParser(description='Translate blobber')
    parser.add_argument('-t', '--text', action='store', dest='text', type=str, default='')
    parser.add_argument('-s', '--source', action='store', dest='source', type=str, default='')
    args = parser.parse_args()

    try:
        main(args)
    except Exception as err:
        print(err)


if __name__ == '__main__':
    init()
