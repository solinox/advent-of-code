from statistics import median, mean

def parse(f):
    return [int(x) for x in f.read().split(',')]

def fuel_cost_1(data, n):
    return sum([abs(x-n) for x in data])

def step_sum(d):
    return d * (d+1) // 2

def fuel_cost_2(data, n):
    return sum([step_sum(abs(x-n)) for x in data])

def one(data):
    # answer seems to be the median for test and actual input
    return fuel_cost_1(data, int(median(data)))
    
def two(data):
    # answer is close to the mean
    k = int(mean(data))
    costs = [fuel_cost_2(data, x) for x in range(k-2, k+2, 1)]
    return min(costs)
