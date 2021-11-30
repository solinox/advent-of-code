from datetime import datetime
import importlib
import sys

def run(part: int, fn, data):
    t0 = datetime.now()
    v = fn(data)
    t1 = datetime.now()
    dt = t1 - t0
    print(f'--- Part {part} --- ({dt})')
    print(v)

if __name__ == '__main__':
    day = f'0{sys.argv[1]}'[-2:]
    module = importlib.import_module(day)
    print(f'Running day {day}')
    with open(f'{day}/input.txt') as f:
        parse = getattr(module, 'parse', lambda x: x.read())
        data = parse(f)
        p1 = getattr(module, 'one', None)
        if p1:
            run(1, p1, data)
        p2 = getattr(module, 'two', None)
        if p2:
            run(2, p2, data)