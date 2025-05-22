package scheduler

import (
	// "encoding/json"
	"fmt"
)

func PrintNodeList(nodes map[string]*Node) {
	fmt.Printf("LIST=====================\n")
	for _, node := range nodes {
		fmt.Printf("%s\n", node.name)
		fmt.Printf("%d\n", node.score)
		// bytes, err := json.MarshalIndent(node.node, "", "   ")
		// if err != nil {
		// 	return
		// }

		// fmt.Println(string(bytes))
	}
	fmt.Printf("END_LIST=====================\n")
}

func NodeInPrimaryList(name string) bool {
	_, node := PrimaryNodeList[name]
	return node
}

func NodeInSecondaryList(name string) bool {
	_, node := SecondaryNodeList[name]
	return node
}

func NodeExist(name string) bool {
	return NodeInPrimaryList(name) || NodeInSecondaryList(name)
}

func DeleteNode(name string) {
	delete(PrimaryNodeList, name)
	delete(SecondaryNodeList, name)
}
