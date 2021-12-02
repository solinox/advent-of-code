from functools import lru_cache


def parse(f):
    # using tuples so input can be hashable for lru_cache
    data = [x.split(' ') for x in f]
    data = tuple((step[0], int(step[1])) for step in data)
    return data

@lru_cache()
def dive(data):
    x = 0
    y = 0
    aim = 0
    for dir, dist in data:
        if dir == 'forward':
            x += dist
            y += dist*aim
        elif dir == 'down':
            aim += dist
        elif dir == 'up':
            aim -= dist
    return x, y, aim

def one(data):
    x, y, a = dive(data)
    return x*a

def two(data):
    x, y, a = dive(data)
    return x*y