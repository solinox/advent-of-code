from collections import Counter

def parse(f):
    template = f.readline().strip()
    polymer_pairs = Counter(template[i:i+2] for i in range(len(template)-1))
    f.readline()
    rules = dict(x.strip().split(' -> ') for x in f)
    return (template, polymer_pairs, rules)

def polymer_step(pairs, letters, rules):
    for pair, count in pairs.copy().items():
        letter = rules[pair] # guaranteed to exist
        letters[letter] += count
        pairs[pair] -= count
        pairs[pair[0]+letter] += count
        pairs[letter+pair[1]] += count

def solve(data, steps):
    template, polymer_pairs, rules = data
    letters = Counter(template)
    for _ in range(steps):
        polymer_step(polymer_pairs, letters, rules)
    letters = letters.most_common()
    return letters[0][1] - letters[-1][1]

def one(data):
    return solve(data, 10)

def two(data):
    return solve(data, 40)
