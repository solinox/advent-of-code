package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/solinox/advent-of-code/2022/util"
)

//go:embed input.txt
var input string

type FileInfo struct {
	Name  string
	Size  int
	Dir   bool
	Files []*FileInfo
}

type FileTree map[string]*FileInfo

// RandomWalk iterates through every element of the FileTree, but in no specific order
func (t FileTree) RandomWalk(fn func(*FileInfo)) {
	for _, v := range t {
		fn(v)
	}
}

// Walk iterates a DFS walk through elements contained in the starting directory, including the start
func (f *FileInfo) Walk(fn func(*FileInfo)) {
	if f != nil {
		fn(f)
	}
	for _, v := range f.Files {
		v.Walk(fn)
	}
}

// TotalSize is the size of the file, or the size of all inner files if f is a directory
func (f *FileInfo) TotalSize() int {
	s := 0
	f.Walk(func(ff *FileInfo) {
		s += ff.Size
	})
	return s
}

func main() {
	tree := buildTree(input)
	util.RunTimed(part1, tree)
	util.RunTimed(part2, tree)
}

func buildTree(in string) FileTree {
	t := make(FileTree)
	s := bufio.NewScanner(strings.NewReader(in))
	currentPath := ""
	var currentDir *FileInfo
	for s.Scan() {
		line := s.Text()
		if line[0:4] == "$ cd" {
			currentPath = filepath.Join(currentPath, line[5:])
			f := &FileInfo{Name: currentPath, Size: 0, Dir: true}
			if currentDir != nil && t[currentPath] == nil {
				currentDir.Files = append(currentDir.Files, f)
			}
			if t[currentPath] == nil {
				t[currentPath] = f
			}
			currentDir = t[currentPath]
		} else if line[0] == '$' {
			continue // ls
		} else if line[0:3] == "dir" {
			f := &FileInfo{Dir: true, Size: 0, Name: filepath.Join(currentPath, line[4:])}
			if _, ok := t[f.Name]; !ok {
				t[f.Name] = f
			}
			currentDir.Files = append(currentDir.Files, t[f.Name])
		} else {
			f := &FileInfo{Dir: false}
			fmt.Sscanf(line, "%d %s", &f.Size, &f.Name)
			f.Name = filepath.Join(currentPath, f.Name)
			t[f.Name] = f
			currentDir.Files = append(currentDir.Files, f)
		}
	}
	return t
}

func part1(t FileTree) int {
	sum := 0
	t.RandomWalk(func(f *FileInfo) {
		if !f.Dir {
			return
		}
		dirSum := f.TotalSize()
		if dirSum <= 100000 {
			// fmt.Printf("Dir %s has eligible size %d\n", f.Name, dirSum)
			sum += dirSum
		}
	})
	return sum
}

func part2(t FileTree) int {
	const (
		totalDisk        = 70000000
		unusedDiskNeeded = 30000000
	)
	unusedDisk := totalDisk - t["/"].TotalSize()
	threshold := unusedDiskNeeded - unusedDisk
	min := totalDisk
	t.RandomWalk(func(f *FileInfo) {
		if !f.Dir {
			return
		}
		dirSum := f.TotalSize()
		if dirSum >= threshold && dirSum < min {
			min = dirSum
		}

	})
	return min
}
