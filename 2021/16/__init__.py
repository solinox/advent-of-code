from math import prod

def parse(f):
    return next(x.strip() for x in f)

def decode(transmission):
    bits = ''
    for b in transmission:
        bits += ('0000' + bin(int(b, 16))[2:])[-4:]
    return bits

def solve(transmission):
    bits = decode(transmission)
    v, n, _ = parse_packet(bits)
    return v, n

def parse_literal_number(bits):
    cont = True
    lit_bits = ''
    read = 0
    while cont:
        cont = bool(int(bits[0], 2))
        lit_bits += bits[1:5]
        read += 5
        bits = bits[5:]
    return int(lit_bits, 2), read

def parse_packet(packet):
    version = int(packet[:3], 2)
    type_id = int(packet[3:6], 2)
    if type_id != 4:
        length_type_id = packet[6]
        subpacket_values = []
        if length_type_id == '0':
            subpacket_len = int(packet[7:22], 2)
            read = 0
            while read < subpacket_len:
                v, n, r = parse_packet(packet[22+read:])
                read += r
                version += v
                subpacket_values.append(n)
            read += 22
        elif length_type_id == '1':
            num_subpackets = int(packet[7:18], 2)
            read = 0
            for _ in range(num_subpackets):
                v, n, r = parse_packet(packet[18+read:])
                read += r
                version += v
                subpacket_values.append(n)
            read += 18
        
    match type_id:
        case 0: n = sum(subpacket_values)
        case 1: n = prod(subpacket_values)
        case 2: n = min(subpacket_values)
        case 3: n = max(subpacket_values)
        case 4:
            n, read = parse_literal_number(packet[6:])
            read += 6
        case 5: n = 1 if subpacket_values[0] > subpacket_values[1] else 0
        case 6: n = 1 if subpacket_values[0] < subpacket_values[1] else 0
        case 7: n = 1 if subpacket_values[0] == subpacket_values[1] else 0

    return version, n, read

def one(data):
    v, n = solve(data)
    return v

def two(data):
    _, n = solve(data)
    return n