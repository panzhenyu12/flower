package model

// type TreeNode interface {
// 	GetID() string
// 	GetSuperiorID() string
// 	AddChild(child interface{}) error
// 	GetChild() interface{}
// }

// func GetTree(r *TreeNode, ns []*TreeNode) error {
// 	root := *r
// 	for _, n := range ns {
// 		node := *n
// 		if root.GetID() == node.GetSuperiorID() {
// 			root.AddChild(n)
// 		}
// 	}
// 	return nil
// }

func has(v1 OrgModel, vs []*OrgModel) bool {
	var has bool
	has = false
	for _, v2 := range vs {
		if v2.GetID() == v2.GetSuperiorID() {
			continue
		}
		if v1.GetID() == v2.GetSuperiorID() {
			has = true
			break
		}
	}
	return has

}

func MakeTree(vs []*OrgModel, node *OrgModel) {
	childs := findChild(node, vs)
	for _, child := range childs {
		node.AddChild(child)
		if has(*child, vs) {
			MakeTree(vs, child)
		}
	}
}

func findChild(b *OrgModel, ns []*OrgModel) (ret []*OrgModel) {
	for _, n := range ns {
		if n.GetID() == n.GetSuperiorID() {
			continue
		}
		if b.GetID() == n.GetSuperiorID() {
			ret = append(ret, n)
		}
	}
	return
}
