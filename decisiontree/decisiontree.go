package decisiontree

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"log"
	"sort"

	"gopkg.in/mgo.v2/bson"
)

var RootID = "5ffb8cb8a70f9e1daf3b24dd"

type TreeNodes struct {
	// TODO: Implement
	Nodes []Node
}

type Node struct {
	ID        bson.ObjectId          `bson:"_id"   json:"id"`
	CatID     bson.ObjectId          `bson:"cat_id,omitempty"   json:"cat_id"`
	CatName   string                 `bson:"cat_name,omitempty"   json:"cat_name"`
	Name      string                 `bson:"name,omitempty" form:"name" json:"name"`
	Key       string                 `bson:"key"  form:"key" json:"key"`
	Desc      string                 `bson:"desc"  form:"desc" json:"desc"`
	Type      string                 `bson:"type"  form:"type" json:"type"`
	IsSymptom bool                   `bson:"is_symptom"  form:"is_symptom" json:"is_symptom"`
	Lkey      string                 `bson:"lkey"  form:"lkey" json:"lkey"`
	Value     interface{}            `bson:"value"  form:"value" json:"value"`
	Operator  string                 `bson:"operator"  form:"operator" json:"operator"`
	Order     int                    `bson:"order"  form:"order" json:"order"`
	Content   interface{}            `bson:"content"  form:"content" json:"content"`
	Headers   map[string]interface{} `bson:"headers"  form:"headers" json:"headers"`
	Parent    ParentNode             `bson:"parent,omitempty" form:"parent" json:"parent"`
}
type ParentNode struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Name     string        `bson:"name,omitempty" form:"name" json:"name"`
	Key      string        `json:"key"  form:"key" json:"key"`
	Value    interface{}   `json:"value"  form:"value" json:"value"`
	Operator string        `json:"operator"  form:"operator" json:"operator"`
}

func TreeTypeCheck(val interface{}) interface{} {

	iAreaId, ok := val.(float64)
	iAreaId2, ok2 := val.(int)
	if ok {
		return float64(iAreaId)
	} else if ok2 {
		return float64(iAreaId2)
	}
	return val
}

func TreeLoader() (*Tree, error) {

	file, _ := ioutil.ReadFile("treedb.json")

	lst := TreeNodes{}
	err := json.Unmarshal([]byte(file), &lst)

	// lst, err := Load(2000)
	var trees []Tree
	ids := map[bson.ObjectId]int{}

	for i, l := range lst.Nodes {
		ids[l.ID] = (i + 1)
		t := Tree{
			ID:       ids[l.ID],
			Name:     l.Name,
			ParentID: ids[l.Parent.ID],
			Value:    TreeTypeCheck(l.Value),
			Key:      l.Key,
			Operator: l.Operator,
			Order:    l.Order,
			Content:  l.Content,
			Headers:  l.Headers,
		}
		trees = append(trees, t)
	}

	if err != nil {
		log.Println(lst)
	}
	return CreateTree(trees), nil
}

//mongo db codes end

// AddNode Add a new Node (leaf) to the Tree
func (t *Tree) AddNode(node *Tree) {
	node.parent = t
	t.Nodes = append(t.Nodes, node)
	sort.Sort(byOrder(t.Nodes))
}

// GetChild get the nodes child of this one
func (t *Tree) GetChild() []*Tree {
	return t.Nodes
}

// GetParent get the parent node of this one
func (t *Tree) GetParent() *Tree {
	return t.parent
}

// WithContext returns a tree with a context
func (t *Tree) WithContext(ctx context.Context) *Tree {
	t.ctx = ctx
	return t
}

// Context returns the context
func (t *Tree) Context() context.Context {
	return t.ctx
}

// Next evaluate which will be the next Node according to the jsonRequest
func (t *Tree) Next(jsonRequest map[string]interface{}, config *TreeOptions) (*Tree, error) {
	var jsonValue interface{}
	var oldName string
	for _, n := range t.Nodes {

		if oldName != n.Key {
			jsonValue = jsonRequest[n.Key]
			oldName = n.Key
		}

		selected, err := compare(jsonRequest, jsonValue, n, config)
		if config.StopIfConvertingError == true && err != nil {
			return n, err
		}

		if selected != nil {
			if t.ctx != nil {
				t.ctx = contextValue(t.ctx, selected.ID, selected.Key, jsonValue, selected.Operator, selected.Value)
			}

			if config.context != nil {
				config.context = contextValue(config.context, selected.ID, selected.Key, jsonValue, selected.Operator, selected.Value)
			}
			return selected, nil
		}
	}

	return nil, nil
}

// LoadTree gets a json on build the Tree related

// CreateTree attach the nodes to the Tree
func CreateTree(data []Tree) *Tree {
	temp := make(map[int]*Tree)
	var root *Tree
	for i := range data {
		leaf := &data[i]
		temp[leaf.ID] = leaf
		if leaf.ParentID == 0 {
			root = leaf
		}
	}

	for _, v := range temp {
		if v.ParentID != 0 {
			temp[v.ParentID].AddNode(v)
		}
	}

	return root
}

// ResolveJSON calculate which will be the selected node according to the jsonRequest
func (t *Tree) ResolveJSON(jsonRequest []byte, options ...func(t *TreeOptions)) (*Tree, error) {
	var request map[string]interface{}
	err := json.Unmarshal(jsonRequest, &request)
	if err != nil {
		return nil, err
	}

	return t.Resolve(request, options...)
}

// Resolve calculate which will be the selected node according to the map request
func (t *Tree) Resolve(request map[string]interface{}, options ...func(t *TreeOptions)) (*Tree, error) {
	config := &TreeOptions{}

	for _, option := range options {
		option(config)
	}

	if len(config.Operators) > 0 {
		for k := range config.Operators {
			if isExistingOperator(k) {
				config.OverrideExistingOperator = true
				break
			}
		}
	}

	return t.resolve(request, config)
}

// ResolveJSONWithContext calculate which will be the selected node according to the jsonRequest
func (t *Tree) ResolveJSONWithContext(ctx context.Context, jsonRequest []byte, options ...func(t *TreeOptions)) (*Tree, context.Context, error) {
	var request map[string]interface{}
	err := json.Unmarshal(jsonRequest, &request)
	if err != nil {
		return nil, ctx, err
	}

	return t.ResolveWithContext(ctx, request, options...)
}

// ResolveWithContext calculate which will be the selected node according to the map request
func (t *Tree) ResolveWithContext(ctx context.Context, request map[string]interface{}, options ...func(t *TreeOptions)) (*Tree, context.Context, error) {
	config := &TreeOptions{
		context: ctx,
	}

	for _, option := range options {
		option(config)
	}

	if len(config.Operators) > 0 {
		for k := range config.Operators {
			if isExistingOperator(k) {
				config.OverrideExistingOperator = true
				break
			}
		}
	}

	result, err := t.resolve(request, config)
	return result, config.context, err
}

func (t *Tree) resolve(request map[string]interface{}, config *TreeOptions) (*Tree, error) {
	temp, err := t.Next(request, config)
	if err != nil {
		return t, err
	}

	if temp == nil {
		return t, err
	}
	return temp.resolve(request, config)
}

// Operator represent a function that will evaluate a node
type Operator func(requests map[string]interface{}, node *Tree) (*Tree, error)

// TreeOptions allow to extend the comparator
type TreeOptions struct {
	StopIfConvertingError    bool
	Operators                map[string]Operator
	OverrideExistingOperator bool
	context                  context.Context
}

// Tree represent a Tree
type Tree struct {
	Nodes  []*Tree
	parent *Tree

	ctx context.Context

	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	ParentID int                    `json:"parent_id"`
	Value    interface{}            `json:"value"`
	Operator string                 `json:"operator"`
	Key      string                 `json:"key"`
	Order    int                    `json:"order"`
	Content  interface{}            `json:"content"`
	Headers  map[string]interface{} `json:"headers"`
}

type byOrder []*Tree

func (o byOrder) Len() int      { return len(o) }
func (o byOrder) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o byOrder) Less(i, j int) bool {
	// fallback is always the last
	if s, ok := o[i].Value.(string); ok {
		if s == FallbackType {
			return false
		}
	}

	if s, ok := o[j].Value.(string); ok {
		if s == FallbackType {
			return true
		}
	}

	return o[i].Order < o[j].Order
}
