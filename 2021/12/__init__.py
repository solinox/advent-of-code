def parse(f):
    edges = [[xx.strip() for xx in x.split('-')] for x in f]
    node_map = {}
    for e in edges:
        l, r = e
        node_map[l] = node_map.get(l, []) + [r]
        node_map[r] = node_map.get(r, []) + [l]
    return node_map

def traverse(node_map, current, target, path=[], p2=False, p2_exception=None):
    travelled = path + [current]
    if current == target:
        return 1
    n = 0
    for node in node_map[current]:
        if node.lower() == node and node in travelled:
            if p2 and p2_exception is None and node not in ['start', 'end']:
                n += traverse(node_map, node, target, travelled, p2, node)
            continue
        n += traverse(node_map, node, target, travelled, p2, p2_exception)
    return n

def one(data):
    return traverse(data, 'start', 'end')

def two(data):
    return traverse(data, 'start', 'end', p2=True)
