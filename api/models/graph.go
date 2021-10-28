package models

type Graph map[string][]string

func (g Graph) GetTree(current string) Graph {
	graph := Graph{}
	queue := []string{current}
	for i := 0; i < len(queue); i++ {
		k := queue[i]
		for _, v := range g[k] {
			queue = append(queue, v)
			graph[k] = append(graph[k], v)
		}
	}
    return graph
}
