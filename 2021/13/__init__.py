def parse(f):
    sections = f.read().split('\n\n')
    dots = set(tuple(map(int, p.split(','))) for p in sections[0].split('\n'))
    folds = [p.split()[-1].split('=') for p in sections[1].split('\n')]
    return (dots, folds)

def fold_dir(dots: set, dir, mirror):
    folded = set()
    for dot in dots:
        val = dot[0] if dir == 'x' else dot[1]
        if val > mirror:
            folded_dot = (mirror + mirror - dot[0], dot[1]) if dir == 'x' else (dot[0], mirror + mirror - dot[1])
            folded.add(folded_dot)
        else:
            folded.add(dot)
    return folded

def fold(dots, folds):
    for f in folds:
        dots = fold_dir(dots, f[0], int(f[1]))
    return dots

def one(data):
    dots, folds = data
    folds = [folds[0]]
    dots = fold(dots, folds)
    return len(dots)

def display(dots):
    maxx = max(dot[0] for dot in dots) + 1
    maxy = max(dot[1] for dot in dots) + 1
    text = [[' ' for i in range(maxx)] for i in range(maxy)]
    for dot in dots:
        text[dot[1]][dot[0]] = '#'
    return '\n'.join(map(''.join, text))

def two(data):
    dots, folds = data
    dots = fold(dots, folds)
    return display(dots)
