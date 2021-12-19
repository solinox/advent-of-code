def parse(f):
    fields = f.read().split()
    x0, x1 = map(lambda x: int(x.strip(',')), fields[2].split('=')[-1].split('..'))
    y0, y1 = map(int, fields[3].split('=')[-1].split('..'))
    return [(min(x0, x1), max(y0, y1)), (max(x0, x1), min(y0, y1))]

def one(data):
    (x0, y0), (x1, y1) = data
    for y in range(500, 0, -1):
        for x in range(x1):
            max_height = 0
            px, py, dx, dy = 0, 0, x, y
            while px < x1 and py > y1:
                px += dx
                py += dy
                if py > max_height:
                    max_height = py
                if dx != 0:
                    dx = dx-1
                dy -= 1
                if x0 <= px <= x1 and y0 >= py >= y1:
                    return max_height

def two(data):
    (x0, y0), (x1, y1) = data
    valid_count = 0
    for y in range(y1, 215): # answer to part1 had y=214 as the max
        for x in range(x1+1):
            px, py, dx, dy = 0, 0, x, y
            while px < x1 and py > y1:
                px += dx
                py += dy
                if dx != 0:
                    dx = dx-1
                dy -= 1
                if x0 <= px <= x1 and y0 >= py >= y1:
                    valid_count += 1
                    break
    return valid_count