def parse(f):
    return [[int(c) for c in x.strip()] for x in f]

def is_low(data, pt):
    x, y = pt
    v = data[y][x]
    if v == 9:
        return False
    for pt in [(x-1, y), (x+1, y), (x, y-1), (x, y+1)]:
        xx, yy = pt
        if yy < 0 or yy >= len(data) or xx < 0 or xx >= len(data[yy]):
            continue
        if v > data[yy][xx]:
            return False
    return True

def low_points(data):
    lows, lows_locs = [], []
    for y in range(len(data)):
        for x in range(len(data[y])):
            if is_low(data, (x, y)):
                lows.append(data[y][x])
                lows_locs.append((x, y))
    return lows, lows_locs

def one(data):
    lows, _ = low_points(data)
    return sum(lows) + len(lows)

def explore_basin(data, explored):
    x, y = explored[-1]
    for pt in [(x-1, y), (x+1, y), (x, y-1), (x, y+1)]:
        if pt in explored:
            continue
        xx, yy = pt
        if yy < 0 or yy >= len(data) or xx < 0 or xx >= len(data[yy]):
            continue
        if data[yy][xx] != 9:
            explored.append(pt)
            explored = explore_basin(data, explored)
    return explored


def two(data):
    _, lows_locs = low_points(data)
    basins = []
    for pt in lows_locs:
        basins.append(explore_basin(data, [pt]))
    basin_lengths = [len(b) for b in basins]
    basin_lengths.sort(reverse=True)
    return basin_lengths[0] * basin_lengths[1] * basin_lengths[2]
