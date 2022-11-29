from math import floor, ceil
from copy import deepcopy

class Scanner:
    def __init__(self):
        self.beacons = []

    def __repr__(self):
        return f'Scanner({self.beacons})'

def parse(f):
    scanners = []
    scanner = None
    for x in f:
        if x.startswith('---'):
            if scanner:
                scanners.append(scanner)
            scanner = Scanner()
        else:
            l = x.strip()
            if l:
                scanner.beacons.append(tuple([int(n) for n in l.split(',')]))
    scanners.append(scanner)
    return scanners
    
def one(data):
    return data

def two(data):
    return data
