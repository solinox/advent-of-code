import functools

def parse(f):
    return [int(x) for x in f]

def one(data):
    count = 0
    for i in range(1, len(data)):
        if data[i] > data[i-1]:
            count += 1
    return count

def two(data):
    count = 0
    for i in range(1, len(data)-2):
        if sum(data[i:i+3]) > sum(data[i-1:i+2]):
            count += 1
    return count