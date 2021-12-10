from functools import reduce


def parse(f):
    return [x.strip() for x in f]

def check(line):
    chunks = []
    for c in line:
        match c:
            case '(': chunks.append(')')
            case '{': chunks.append('}')
            case '[': chunks.append(']')
            case '<': chunks.append('>')
            case _:
                expected = chunks.pop()
                if c != expected:
                    return c, None
    return None, chunks

def one(data):
    score_map = {')': 3, ']': 57, '}': 1197, '>': 25137, None: 0}
    return sum([score_map[check(x)[0]] for x in data])

def two(data):
    score_map = {')': 1, ']': 2, '}': 3, '>': 4}
    lines = [check(x)[1] for x in data]
    scores = [reduce(lambda a,b: a*5 + score_map[b], reversed(x), 0) for x in lines if x is not None]
    scores.sort()
    return scores[len(scores)//2]
