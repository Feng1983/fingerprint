package main
import(
    "fmt"
    sb "github.com/dropbox/godropbox/database/sqlbuilder"
)

func Example() {
    t1 := sb.NewTable(                                                                                                                                   
        "parent_prefix",
        sb.IntColumn("ns_id", false),
        sb.IntColumn("hash", false),
        sb.StrColumn(
            "prefix",
            sb.UTF8,
            sb.UTF8CaseInsensitive,
            false))

    t2 := sb.NewTable(
        "sfj",
        sb.IntColumn("ns_id", false),
        sb.IntColumn("sjid", false),
        sb.StrColumn(
            "filename",
            sb.UTF8,
            sb.UTF8CaseInsensitive,
            false))
    fmt.Println(t1.Name())
    for _, element := range t2.Columns(){
        fmt.Printf("%v",element)
    }
    fmt.Println()
    ns_id1 := t1.C("ns_id")
    prefix := t1.C("prefix")
    ns_id2 := t2.C("ns_id")
    sjid := t2.C("sjid")
    filename := t2.C("filename")

    in := []int32{1, 2, 3}
    join := t2.LeftJoinOn(t1, sb.Eq(ns_id1, ns_id2))
	q := join.Select(ns_id2, sjid, prefix, filename).Where(
        sb.And(sb.EqL(ns_id2, 123), sb.In(sjid, in)))
    fmt.Println(q.String("shard1"))
    fmt.Println(q.String("bbf"))
}
type Element struct {
    Key int
    Value interface{}
}
type Heap struct{
    size int
    max  int
    items []*Element
}
func main(){
    Example()
    fmt.Println("run...")
    h := &Heap{size:0,max:10,items:make([]*Element,10,10)}
    //for _,d:= range h.items{
    //  fmt.Println(d.Key,d.Value)
    //}
    s:= h.size
    h.size++
    fmt.Println(s,h.size)
}              
