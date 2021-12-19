from math import floor, ceil
from copy import deepcopy

class Node:
    def __init__(self, parent, val):
        self.parent = parent
        self.val = val
        self.children = []
        if parent is not None:
            parent.link(self)
    
    def __str__(self):
        if self.val is not None:
            return str(self.val)
        s = '['
        s += ','.join([str(c) for c in self.children])
        s += ']'
        return s

    def link(self, n):
        self.children.append(n)
        n.parent = self

    def add(self, other):
        root = Node(None, None)
        root.link(self)
        root.link(other)
        return root.reduce()

    def split(self):
        if self.val is not None and self.val > 9:
            l, r = floor(self.val/2), ceil(self.val/2)
            self.val = None
            Node(self, l)
            Node(self, r)
            return True
        for c in self.children:
            if c.split():
                return True
        return False

    def exploding(self, dir, n):
        if self.parent is None:
            return
        if dir == 'left':
            cc = self.parent.children
            i = -1
        else:
            cc = reversed(self.parent.children)
            i = 0
        for c in cc:
            if c == self:
                return self.parent.exploding(dir, n)
            else:
                if c.val is not None:
                    c.val += n
                    return
                cc = c
                while cc.children:
                    cc = cc.children[i]
                cc.val += n
                return

    def explode(self, d=0):
        if d >= 4 and self.children:
            l, r = self.children
            self.exploding('left', l.val)
            self.exploding('right', r.val)
            self.children = []
            self.val = 0
            return True
        for c in self.children:
            if c.explode(d+1):
                return True
        return False

    def reduce(self):
        t = True
        while t:
            t = self.explode()
            if not t:
                t = self.split()
        return self

    def mag(self):
        if self.val is not None:
            return self.val
        if len(self.children) == 1:
            return self.children[0].mag()
        l, r = self.children
        return 3*l.mag() + 2*r.mag()


def parse(f):
    def parse_line(line):
        current = None
        for c in line[:-1]:
            if c == '[':
                new = Node(current, None)
                current = new
            elif c == ']':
                current = current.parent
            elif c.isnumeric():
                new = Node(current, int(c))
        return current
    return [parse_line(x.strip()) for x in f]
    
def one(data):
    n = data[0]
    for m in data[1:]:
        n = n.add(m)
    return n.mag()

def two(data):
    # slow (7s) due to deepcopy, which prevents add iterations from affecting each other
    return max(
        max(
            deepcopy(data[i]).add(deepcopy(data[j])).mag(),
            deepcopy(data[j]).add(deepcopy(data[i])).mag()
        )
        for i in range(len(data)) for j in range(i+1, len(data))
    )
