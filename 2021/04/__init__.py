from dataclasses import dataclass, field

@dataclass
class BingoBoard:
    cells: list
    hit: list = field(default_factory=list)

    def _wins(self, i):
        # check row
        base = 5*(i//5)
        if all([x in self.hit for x in self.cells[base:base+5]]):
            return True
        # check column
        if all([x in self.hit for x in self.cells[i%5:25:5]]):
            return True

    def score(self, num):
        not_hit = [x for x in self.cells if x not in self.hit]
        return sum(not_hit)*num

    def play(self, num):
        if num in self.cells:
            i = self.cells.index(num)
            self.hit.append(num)
            if self._wins(i):
                return self.score(num)
        

def parse(f):
    nums = [int(x) for x in f.readline().split(',')]
    f.readline() # blank line
    boards = [BingoBoard([int(x) for x in board.split()]) for board in f.read().split('\n\n')]
    return (nums, boards)


def one(data):
    (nums, boards) = data
    for v in nums:
        for board in boards:
            score = board.play(v)
            if score is not None:
                return score
    
def two(data):
    (nums, boards) = data
    for v in nums:
        won_boards = []
        for board in boards:
            score = board.play(v)
            if score is not None:
                won_boards.append(board)
        if len(won_boards) == len(boards):
            return won_boards[-1].score(v)
        boards = [b for b in boards if b not in won_boards]
