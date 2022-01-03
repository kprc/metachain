package elect

import "testing"

type MyInt int
func (mi *MyInt)Cmp(v TopNNode) int {
	d:=v.(*MyInt)

	if *mi > *d{
		return 1
	}

	return -1
}

func NewMyInt(n int) *MyInt  {
	a:=MyInt(n)

	return &a
}

func TestTopN_Insert(t *testing.T) {

	topn:=NewTopN(5)


	topn.Insert(NewMyInt(100))
	topn.Insert(NewMyInt(99))
	t.Log("current nodes count",topn.GetCount())
	topn.Insert(NewMyInt(10))
	topn.Insert(NewMyInt(120))
	topn.Insert(NewMyInt(30))
	topn.Insert(NewMyInt(40))
	old,_:=topn.Insert(NewMyInt(50))
	if old!=nil{
		t.Log("old node",*(old.(*MyInt)))
	}

	for i:=0;i<len(topn.vs);i++{
		t.Log(*(topn.vs[i].(*MyInt)))
	}

}
