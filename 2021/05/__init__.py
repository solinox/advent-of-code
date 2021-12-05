from dataclasses import dataclass

@dataclass
class Point:
    x: int
    y: int

    def __hash__(self):
        return hash((self.x, self.y))
    
class Vector:
    def __init__(self, pts):
        start, end = pts
        startx, starty = start.split(',')
        endx, endy = end.split(',')
        self.start = Point(int(startx), int(starty))
        self.end = Point(int(endx), int(endy))

def parse(f):
    return [Vector(l.split(' -> ')) for l in f.readlines()]

def build_grid(data, part2 = False):
    grid = {}
    for vec in data:
        dx = vec.end.x - vec.start.x
        dy = vec.end.y - vec.start.y
        if not part2 and dx != 0 and dy != 0:
            # skip diagonals for part1
            continue
        mx = dx/abs(dx) if dx != 0 else 0
        my = dy/abs(dy) if dy != 0 else 0
        for i in range(max(abs(dx), abs(dy))+1):
            pt = Point(vec.start.x + i*mx, vec.start.y + i*my)
            grid[pt] = grid.get(pt, 0) + 1
    return grid

def one(data):
    grid = build_grid(data, False)
    return len({k: v for k,v in grid.items() if v > 1})


def two(data):
    grid = build_grid(data, True)
    return len({k: v for k,v in grid.items() if v > 1})
