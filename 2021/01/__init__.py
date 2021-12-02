def parse(f):
    return [int(x) for x in f]

def one(data):
    count = 0
    for i, v in enumerate(data[1:]):
        if data[i] < v:
            count += 1
    return count

def two(data):
    count = 0
    for i, v in enumerate(data[:-3]):
        if data[i+3] > v:
            count += 1
    return count