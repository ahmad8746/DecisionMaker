package main

import (
	dt "DecisionMaker/decisiontree"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labstack/echo"
)

func LoadTree(c echo.Context) error {
	file, _ := ioutil.ReadFile("treedb.json")

	data := dt.TreeNodes{}
	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(200, data)

}
func ResolveingTree(c echo.Context) error {
	/**
	 * @api {post} /tree/resolve resolve Decision tree
	 * @apiName  resolve Decision tree
	 * @apiGroup Decision tree
	 * @apiVersion 1.0.0
	 * @apiDescription resolve Decision tree.
	 *
	 * @apiParam {interface} data[]

	 *

	 * @apiSuccess {String}   msg http code 200
	 * @apiSuccess {String}   answere
	 * @apiSuccess {String}   children[]
	 * @apiSuccess {String}   path
	 */
	data := make(map[string]interface{})
	if err := c.Bind(&data); err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// t1, _ := dt.LoadTree(jsonTree)
	t, _ := dt.TreeLoader()
	// log.Println(t1)
	log.Println(t)
	t.WithContext(context.Background())
	var vv []*dt.Tree
	v, _ := t.Resolve(data)

	//nodess := v.Nodes
	var answere string
	fmt.Println(answere)

	if v.Operator == "" {
		answere = fmt.Sprintf("%v", v.Name)
		//	mp = dt.FindDisease(answere)

	} else {
		answere = fmt.Sprintf(" %v %v opt: %v value:%v", v.Name, v.Key, v.Operator, v.Value)
	}

	vv = v.GetChild()
	var x []string
	for i, j := range vv {
		d := fmt.Sprintf("%v- %v %v opt: %v value:%v", i, j.Name, j.Key, j.Operator, j.Value)
		x = append(x, d)
		fmt.Println(i, "name:", j.Name, " key:", j.Key, "opt:", j.Operator, j.Value)
	}

	return c.JSON(200, echo.Map{
		"msg":       200,
		"questions": v.Nodes,
		"answer":    v.Name,
	})

}

func ResolveTree(data map[string]interface{}) {

	// t1, _ := dt.LoadTree(jsonTree)
	t, _ := dt.TreeLoader()
	// log.Println(t1)
	log.Println(t)
	t.WithContext(context.Background())
	var vv []*dt.Tree
	v, _ := t.Resolve(data)

	//nodess := v.Nodes
	var answere string
	fmt.Println(answere)

	if v.Operator == "" {
		answere = fmt.Sprintf("%v", v.Name)
		//	mp = dt.FindDisease(answere)

	} else {
		answere = fmt.Sprintf(" %v %v opt: %v value:%v", v.Name, v.Key, v.Operator, v.Value)
	}

	vv = v.GetChild()
	var x []string
	for i, j := range vv {
		d := fmt.Sprintf("%v- %v %v opt: %v value:%v", i, j.Name, j.Key, j.Operator, j.Value)
		x = append(x, d)
		fmt.Println(i, "name:", j.Name, " key:", j.Key, "opt:", j.Operator, j.Value)
	}
	fmt.Println(v.Name)
	fmt.Println(v.Nodes)

	// return c.JSON(200, echo.Map{
	// 	"msg":       200,
	// 	"questions": v.Nodes,
	// 	"answer":    v.Name,
	// })

}
func main() {
	e := echo.New()

	e.GET("/", LoadTree)
	e.POST("/tree/resolve", ResolveingTree)
	e.Logger.Fatal(e.Start(":1323"))

}
