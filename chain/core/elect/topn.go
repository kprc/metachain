package elect

import "errors"


type TopNNode interface {
	Cmp(v TopNNode) int
}

type TopN struct {
	n int
	vs []TopNNode
}

func NewTopN(n int) *TopN {
	return &TopN{n: n}
}

func (tn *TopN)Insert(v TopNNode) (old TopNNode,err error)  {
	if len(tn.vs) == 0{
		tn.vs = append(tn.vs, v)
		return nil,nil
	}

	idx:=-1

	for i:=0;i<len(tn.vs);i++{
		if tn.vs[i].Cmp(v) <= 0{
			idx = i
			break
		}
	}

	if idx == -1{
		if len(tn.vs) < tn.n{
			tn.vs = append(tn.vs, v)
			return nil,nil
		}else{
			return nil,errors.New("not a top n node")
		}
	}

	vs := make([]TopNNode,0)

	for i:=0;i<len(tn.vs);i++{
		if i == idx {
			vs = append(vs, v)
		}
		vs = append(vs, tn.vs[i])
	}

	tn.vs = vs

	if len(tn.vs) > tn.n{
		old = tn.vs[tn.n]
		tn.vs = tn.vs[:tn.n]
	}

	return
}

func (tn *TopN)GetCount() int  {
	return len(tn.vs)
}

