def parse(f):
    edges = [[xx.strip() for xx in x.split('-')] for x in f]
    node_map = {}
    for e in edges:
        l, r = e
        node_map[l] = node_map.get(l, []) + [r]
        node_map[r] = node_map.get(r, []) + [l]
    return node_map

def traverse(node_map, current, target, path=set(), p2=False):
    if current == target:
        return 1
    n = 0
    for node in node_map[current]:
        if node.islower() and node in path:
            if p2 and node != 'start':
                # sets p2 to False since single exception has been done
                n += traverse(node_map, node, target, path|{current}, False)
            continue
        n += traverse(node_map, node, target, path|{current}, p2)
    return n

def one(data):
    return traverse(data, 'start', 'end')

def two(data):
    return traverse(data, 'start', 'end', p2=True)
