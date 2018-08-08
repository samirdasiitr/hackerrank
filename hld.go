package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}
	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type graphType struct {
	adjList map[int64]map[int64]struct{}
	n       int64
	depth   map[int64]int
	size    map[int64]int64
	heavy   map[int64]int64
	parent  map[int64]int64
	roots   map[int64]int64
	treePos map[int64]int64
	tree    []int64
}

func newGraph(n int64) *graphType {
	graph := &graphType{
		adjList: make(map[int64]map[int64]struct{}),
		depth:   make(map[int64]int),
		heavy:   make(map[int64]int64),
		parent:  make(map[int64]int64),
		size:    make(map[int64]int64),
		roots:   make(map[int64]int64),
		treePos: make(map[int64]int64),
		n:       n,
	}

	for ii := int64(0); ii < n; ii++ {
		graph.heavy[ii] = -1
	}
	graph.parent[0] = -1
	graph.depth[0] = 0
	return graph
}

func (gh *graphType) addEdge(u, v int64) {
	if _, ok := gh.adjList[u]; !ok {
		gh.adjList[u] = make(map[int64]struct{})
	}
	if _, ok := gh.adjList[v]; !ok {
		gh.adjList[v] = make(map[int64]struct{})
	}
	gh.adjList[u][v] = struct{}{}
	gh.adjList[v][u] = struct{}{}
}

func (gh *graphType) doDFS(curr, prev int64, depth int) int64 {
	gh.depth[curr] = depth
	gh.parent[curr] = prev

	var maxSubTree int64
	var size int64 = 1
	for next := range gh.adjList[curr] {
		if next != gh.parent[curr] {
			gh.parent[next] = curr
			gh.depth[next] = gh.depth[curr] + 1
			subTreeSize := gh.doDFS(next, curr, depth+1)
			if subTreeSize > maxSubTree {
				gh.heavy[curr] = next
				maxSubTree = subTreeSize
			}
			size += subTreeSize
		}
	}
	gh.size[curr] = size
	return size
}

func (gh *graphType) doHLD(curr int64) {
	var currentPos int64 = 0

	for ii := int64(0); ii < gh.n; ii++ {
		if gh.parent[ii] == -1 || gh.heavy[gh.parent[ii]] != ii {
			for jj := ii; jj != -1; jj = gh.heavy[jj] {
				gh.roots[jj] = ii
				gh.treePos[jj] = currentPos
				currentPos++
			}
		}
	}
}

func (gh *graphType) makeSegmentTree() {
	gh.tree = make([]int64, gh.n*2+1)
	gh.tree[0] = 0
}

func (gh *graphType) update(u, val int64) {
	//for (t[p += n] = value; p > 1; p >>= 1) t[p>>1] = t[p] + t[p^1];
	p := gh.n + gh.treePos[u]
	gh.tree[p] = val
	for p > 1 {
		gh.tree[p>>1] = gh.tree[p^1]
		if gh.tree[p] > gh.tree[p^1] {
			gh.tree[p>>1] = gh.tree[p]
		}
		p >>= 1
	}
}

func (gh *graphType) processPath(u, v int64) int64 {
	root := gh.roots
	depth := gh.depth
	maxValue := int64(0)

	//	fmt.Printf("%d: %d, %d: %d\n", u, root[u], v, root[v])
	for root[u] != root[v] {
		//		fmt.Printf("depth %d, %d\n", depth[root[u]], depth[root[v]])
		if depth[root[u]] > depth[root[v]] {
			// swap
			//			fmt.Printf("swapping: %d %d\n", u, v)
			u, v = v, u
		}
		//		fmt.Printf("maxValue before: %d\n", maxValue)
		gh.query(gh.treePos[root[v]], gh.treePos[v]+1, &maxValue)
		//		fmt.Printf("maxValue after: %d\n", maxValue)
		// move up the chain
		//		fmt.Printf("moving to chain: %d %d\n", gh.parent[root[v]], root[gh.parent[root[v]]])
		v = gh.parent[root[v]]
	}

	if depth[u] > depth[v] {
		u, v = v, u
	}
	gh.query(gh.treePos[u], gh.treePos[v]+1, &maxValue)
	return maxValue
}

func (gh *graphType) query(u, v int64, res *int64) {
	//	fmt.Printf("q: %d %d\n", u, v)
	l := gh.n + u
	r := gh.n + v
	for l < r {
		if (l & 1) == 1 {
			if *res < gh.tree[l] {
				*res = gh.tree[l]
			}
			l++
		}
		if (r & 1) == 1 {
			r--
			if *res < gh.tree[r] {
				*res = gh.tree[r]
			}
		}
		l >>= 1
		r >>= 1
	}
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	//stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	//checkError(err)
	//defer stdout.Close()

	//writer := bufio.NewWriterSize(stdout, 1024 * 1024)

	nq := strings.Split(readLine(reader), " ")

	nNodes, err := strconv.ParseInt(nq[0], 10, 64)
	checkError(err)
	nQueries, err := strconv.ParseInt(nq[1], 10, 64)
	checkError(err)

	//fmt.Printf("nNodes: %d, nQueries: %d\n", nNodes, nQueries)
	graph := newGraph(nNodes)

	for ii := 0; ii < int(nNodes)-1; ii++ {
		uv := strings.Split(readLine(reader), " ")
		u, err := strconv.ParseInt(uv[0], 10, 64)
		checkError(err)
		v, err := strconv.ParseInt(uv[1], 10, 64)
		checkError(err)
		graph.addEdge(u, v)
	}

	graph.doDFS(0, -1, 0)
	graph.doHLD(0)
	graph.adjList = nil
	graph.makeSegmentTree()

	/*
		for nn := int64(0); nn < graph.n; nn++ {
			//fmt.Printf("%d:%d ", nn, graph.treePos[nn])
		}
		//fmt.Printf("\n")
			for ii := range graph.tree {
				fmt.Printf("%d ", graph.tree[ii])
			}
			fmt.Printf("\n")
	*/
	for ii := 0; ii < int(nQueries); ii++ {
		qs := strings.Split(readLine(reader), " ")
		//fmt.Printf("qs: %s\n", qs)
		t, err := strconv.ParseInt(qs[0], 10, 64)
		checkError(err)
		u, err := strconv.ParseInt(qs[1], 10, 64)
		checkError(err)
		v, err := strconv.ParseInt(qs[2], 10, 64)
		checkError(err)
		if t == 1 {
			graph.update(u, v)
			/*
				for ii := range graph.tree {
					if int64(ii) == graph.n {
						fmt.Printf(":")
					}
					fmt.Printf("%d ", graph.tree[ii])
				}
				fmt.Printf("\n")
			*/
		}

		if t == 2 {
			fmt.Printf("%d\n", graph.processPath(u, v))
		}
	}
}
