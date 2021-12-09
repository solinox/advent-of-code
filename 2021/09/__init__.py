def parse(f):
    return [[int(c) for c in x.strip()] for x in f]

def is_low(data, x, y, mx, my, ):
    pt = data[y][x]
    if pt == 9:
        return False
    if x > 0 and pt >= data[y][x-1]:
        return False
    if x < mx and pt >= data[y][x+1]:
        return False
    if y > 0 and pt >= data[y-1][x]:
        return False
    if y < my and pt >= data[y+1][x]:
        return False
    return True

def low_points(data):
    maxy = len(data)
    lows, lows_locs = [], []
    for y in range(maxy):
        maxx = len(data[y])
        for x in range(maxx):
            if is_low(data, x, y, maxx-1, maxy-1):
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
        if data[pt[1]][pt[0]] != 9:
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
