from collections import Counter

def parse(f):
    return [list(map(lambda x: [Counter(xx) for xx in x.split()], l.split(' | '))) for l in f.readlines()]

def one(data):
    # shortcut, just get number of output words with lengths 2,3,4,7
    # which correspond to numbers 1, 7, 4, 8 respectively
    lengths = [2,3,4,7]
    return len([x for l in data for x in l[1] if len(x) in lengths])

# matching rules. All numbers can be found due to their length
# and count of matching wires they have with numbers 1, 4, and 7
# 0 (len 6) needs to match 2/2 of 1, 3/3 of 7, 3/4 of 4
# 1 (len 2) is only len=2
# 2 (len 5) needs to match 1/2 of 1, 2/3 of 7, 2/4 of 4
# 3 (len 5) needs to match 2/2 of 1, 3/3 of 7, 3/4 of 4
# 4 (len 4) is only len=4
# 5 (len 5) needs to match 1/2 of 1, 2/3 of 7, 3/4 of 4
# 6 (len 6) needs to match 1/2 of 1, 2/3 of 7, 3/4 of 4
# 7 (len 3) is only len=3 
# 8 (len 7) is only len=7
# 9 (len 6) needs to match 2/2 of 1, 3/3 of 7, 4/4 of 4
def solve(input, output):
    one = next(x for x in input if len(x) == 2)
    four = next(x for x in input if len(x) == 4)
    seven = next(x for x in input if len(x) == 3)
    eight = next(x for x in input if len(x) == 7)
    number_rules = {
        '0': lambda x: len(x) == 6 and len(x&one) == 2 and len(x&seven) == 3 and len(x&four) == 3,
        '1': lambda x: x == one,
        '2': lambda x: len(x) == 5 and len(x&one) == 1 and len(x&seven) == 2 and len(x&four) == 2,
        '3': lambda x: len(x) == 5 and len(x&one) == 2 and len(x&seven) == 3 and len(x&four) == 3,
        '4': lambda x: x == four,
        '5': lambda x: len(x) == 5 and len(x&one) == 1 and len(x&seven) == 2 and len(x&four) == 3,
        '6': lambda x: len(x) == 6 and len(x&one) == 1 and len(x&seven) == 2 and len(x&four) == 3,
        '7': lambda x: x == seven,
        '8': lambda x: x == eight,
        '9': lambda x: len(x) == 6 and len(x&one) == 2 and len(x&seven) == 3 and len(x&four) == 4,
    }
    in_numbers = [next(n for n,f in number_rules.items() if f(x)) for x in input]
    return int(''.join([in_numbers[input.index(x)] for x in output]))

def two(data):
    return sum([solve(l[0], l[1]) for l in data])
