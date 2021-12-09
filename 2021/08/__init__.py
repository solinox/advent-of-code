def parse(f):
    return [list(map(lambda x: [set(xx) for xx in x.split()], l.split(' | '))) for l in f.readlines()]

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
    for num in input:
        match len(num):
            case 2: one = num
            case 4: four = num
            case 3: seven = num
    n = ''
    for num in output:
        match (len(num), len(num&one), len(num&seven), len(num&four)):
            case 2, _, _, _: n += '1'
            case 3, _, _, _: n += '7'
            case 4, _, _, _: n += '4'
            case 7, _, _, _: n += '8'
            case 6, 2, 3, 3: n += '0'
            case 5, 1, 2, 2: n += '2'
            case 5, 2, 3, 3: n += '3'
            case 5, 1, 2, 3: n += '5'
            case 6, 1, 2, 3: n += '6'
            case 6, 2, 3, 4: n += '9'
    return int(n)

def two(data):
    return sum([solve(l[0], l[1]) for l in data])
