from functools import lru_cache

def parse(f):
    return [int(x) for x in f.read().strip().split(',')]

@lru_cache
def fish(start_timer, days):
    d = {start_timer: 1}
    for i in range(days, 0, -1):
        new = d.get(0, 0)
        if new > 0:
            d[9] = new
            d[7] = d.get(7, 0) + new
            del d[0]
        d = {k-1:v for k,v in d.items()}
    return sum(d.values())

def one(data):
    return sum([fish(x, 80) for x in data])


def two(data):
    return sum([fish(x, 256) for x in data])
