import heapq

def parse(f):
    grid = [[int(b) for b in x.strip()] for x in f]
    cavern = {(x, y): grid[y][x] for y in range(len(grid)) for x in range(len(grid[y]))}
    return cavern

def astar(cavern, start, end):
    open = [(0, start)]
    came_from = {}
    g_score = {start: 0}
    while len(open) > 0:
        _, current = heapq.heappop(open)
        if current == end:
            return g_score[current]
        neighbors = [(0, 1), (0, -1), (-1, 0), (1, 0)]
        for dn in neighbors:
            n = (current[0] + dn[0], current[1] + dn[1])
            if n not in cavern:
                continue
            tentative_g = g_score[current] + cavern[n]
            if n not in g_score or tentative_g < g_score[n]:
                came_from[n] = current
                g_score[n] = tentative_g
                f_score = tentative_g - 1
                if n not in open:
                    heapq.heappush(open, (f_score, n))


def one(data):
    return astar(data, (0, 0), (99, 99))

def expand(cavern):
    new_cavern = cavern.copy()
    for pt, v in cavern.items():
        for x in range(5):
            for y in range(5):
                new_pt = (pt[0]+(100*x), pt[1]+(100*y))
                if new_pt in new_cavern:
                    continue
                new_v = (v + x + y)
                if new_v > 9:
                    new_v -= 9
                new_cavern[new_pt] = new_v
    return new_cavern

def two(data):
    data = expand(data)
    return astar(data, (0, 0), (499, 499))
