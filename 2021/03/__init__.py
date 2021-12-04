def parse(f):
    return [x.strip() for x in f.readlines()]

def _to_int(bits):
    return int(''.join(bits), 2)

def one(data):
    transposed = zip(*data)
    transposed_bits = [(('1', '0') if x.count('1') > len(x)//2 else ('0', '1')) for x in transposed]
    gamma_bits, epsilon_bits = tuple(zip(*transposed_bits))
    gamma = _to_int(gamma_bits)
    epsilon = _to_int(epsilon_bits)
    return gamma * epsilon
    
def get_bit_for_column(lst, i, v):
    tsl = [x[i] for x in lst]
    return v[0] if tsl.count('1') >= len(tsl)/2 else v[1]

def oxyco(data, v):
    i = 0
    lst = [x for x in data]
    while len(lst) > 1 and i < len(lst[0]):
        bit = get_bit_for_column(lst, i, v)
        lst = [x for x in lst if x[i] == bit]
        i += 1
    return lst[0]

def two(data):
    oxygen = _to_int(oxyco(data, ('1', '0')))
    co2 = _to_int(oxyco(data, ('0', '1')))
    return oxygen*co2

