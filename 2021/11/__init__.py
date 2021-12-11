from dataclasses import dataclass

@dataclass
class Octopus:
    pos: tuple[int, int]
    energy: int

def parse(f):
    return [[Octopus((x, y), int(v)) for x, v in enumerate(l.strip())] for y, l in enumerate(f)]

def neighbors(pos):
    x, y = pos
    return [(xx, yy) for yy in range(y-1, y+2) for xx in range(x-1, x+2) if xx >=0 and xx < 10 and yy >= 0 and yy < 10 and not (xx == x and yy == y)]

def cascade(octos, flash):
    if not flash:
        return
    new = {}
    for pos in flash.keys():
        for neighbor_pt in neighbors(pos):
            octo = octos[neighbor_pt]
            if octo.energy > 9:
                continue
            octo.energy += 1
            if octo.energy > 9:
                new[octo.pos] = octo
    cascade(octos, new)

def step(octos):
    # increase all octopus energy by 1
    for octo in octos.values():
        octo.energy += 1
    
    # recursively cascade through flashing octopus and update neighbors
    cascade(octos, {k:v for k,v in octos.items() if v.energy > 9})

    # count and return number of flashing octopus, reset flashing octupus energy
    count = 0
    for octo in octos.values():
        if octo.energy > 9:
            count += 1
            octo.energy = 0
    return count
    

def one(data):
    octos = {data[y][x].pos: data[y][x] for y in range(len(data)) for x in range(len(data[y]))}
    return sum([step(octos) for _ in range(100)])

def two(data):
    octos = {data[y][x].pos: data[y][x] for y in range(len(data)) for x in range(len(data[y]))}
    steps = 0
    while step(octos) != 100:
        steps += 1
    return steps + 1
