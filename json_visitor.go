package raymond

import (
	"encoding/json"
	"github.com/mailgun/raymond/v2/ast"
	"strings"
)

type list []interface{}

func (l *list) Add(item interface{}) {
	*l = append(*l, item)
}

func newList(item interface{}) *list {
	l := new(list)
	l.Add(item)
	return l
}

type JSONVisitor struct {
	JSON map[string]interface{}
	ctx  *handlebarsContext
}

func newJSONVisitor() *JSONVisitor {
	j := map[string]interface{}{}
	v := &JSONVisitor{JSON: j, ctx: newHandlebarsContext()}
	return v
}

func ToJSON(node ast.Node) string {
	visitor := newJSONVisitor()
	node.Accept(visitor)
	b, _ := json.Marshal(visitor.JSON)
	return string(b)
}

func (v *JSONVisitor) VisitProgram(node *ast.Program) interface{} {
	for _, n := range node.Body {
		n.Accept(v)
	}
	return v.JSON
}

func (v *JSONVisitor) VisitMustache(node *ast.MustacheStatement) interface{} {

	node.Expression.Accept(v)

	return nil
}

func (v *JSONVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	var action string
	fp := node.Expression.FieldPath()
	if fp != nil {
		action = node.Expression.HelperName()
	}
	if action == "with" || action == "each" {
		blockParamsPath := make([]string, 0)
		blockParams := make([]string, 0)
		for _, params := range node.Expression.Params {
			// Extract block params from nested nodes.
			if pe, ok := params.(*ast.PathExpression); ok {
				blockParamsPath = append(blockParamsPath, pe.Parts...)
			}
		}
		if node.Program != nil {
			if len(node.Program.BlockParams) > 0 {
				blockParams = append(node.Program.BlockParams)
			}
		}
		if action == "each" {
			blockParamsPath = append(blockParamsPath, "[0]")
		}
		if len(blockParams) > 0 {
			v.ctx.AddMemberContext(strings.Join(blockParamsPath, "."), strings.Join(blockParams, "."))
		} else {
			v.ctx.AddMemberContext(strings.Join(blockParamsPath, "."), "")
		}
		if node.Program != nil {
			node.Program.Accept(v)
		}
		if node.Inverse != nil {
			node.Inverse.Accept(v)
		}
		v.ctx.MoveUpContext()
	} else {
		for _, param := range node.Expression.Params {
			param.Accept(v)
		}
		if node.Program != nil {
			node.Program.Accept(v)
		}
		if node.Inverse != nil {
			node.Inverse.Accept(v)
		}
	}
	return nil
}

func (v *JSONVisitor) VisitPartial(node *ast.PartialStatement) interface{} {

	return nil
}

func (v *JSONVisitor) VisitContent(node *ast.ContentStatement) interface{} {

	return nil
}

func (v *JSONVisitor) VisitComment(node *ast.CommentStatement) interface{} {

	return nil
}

func (v *JSONVisitor) VisitExpression(node *ast.Expression) interface{} {
	var action string
	fp := node.FieldPath()
	if fp != nil {
		if len(fp.Parts) > 0 {
			action = node.HelperName()
			if action == "lookup" {
				if len(node.Params) > 0 {
					path, ok := node.Params[0].(*ast.PathExpression)
					if ok {
						depth := path.Depth
						tmpPath := []string{}
						for _, p := range path.Parts {
							tmpPath = append(tmpPath, p)
						}
						for _, n := range node.Params[1:] {
							pe, ok := n.(*ast.PathExpression)
							if ok {
								pe.Depth = depth
								pe.Parts = append(tmpPath, pe.Parts...)
								pe.Accept(v)
							}
						}
						return nil
					}
				}
			}
		}
	}
	node.Path.Accept(v)
	for _, n := range node.Params {
		n.Accept(v)
	}
	return nil
}

func (v *JSONVisitor) VisitSubExpression(node *ast.SubExpression) interface{} {
	node.Expression.Accept(v)

	return nil
}

func (v *JSONVisitor) VisitPath(node *ast.PathExpression) interface{} {
	if node.Data {
		data := node.Parts[len(node.Parts)-1]
		if data == "index" {
			node.Parts[len(node.Parts)-1] = "[0]"
		}
	}
	if node.Scoped {
		if strings.HasPrefix(node.Original, ".") && !strings.HasPrefix(node.Original, "..") {
			if len(node.Parts) == 0 {
				node.Parts = []string{""}
			}
		}
	}
	res := v.ctx.GetMappedContext(node.Parts, node.Depth)
	v.appendToJSON(res)
	return nil
}

func (v *JSONVisitor) VisitString(node *ast.StringLiteral) interface{} {
	return nil
}

func (v *JSONVisitor) VisitBoolean(node *ast.BooleanLiteral) interface{} {
	return nil
}

func (v *JSONVisitor) VisitNumber(node *ast.NumberLiteral) interface{} {

	return nil
}

func (v *JSONVisitor) VisitHash(node *ast.Hash) interface{} {
	return nil
}

func (v *JSONVisitor) VisitHashPair(node *ast.HashPair) interface{} {
	return nil
}

func (v *JSONVisitor) appendToJSON(templateLabels []string) {
	var tmp interface{}
	tmp = v.JSON
	for idx, name := range templateLabels {
		switch c := tmp.(type) {
		case map[string]interface{}:
			if _, ok := c[name]; !ok {
				// Peak at the next value to determine if
				// it is an array or map.
				var isArray bool
				if idx < len(templateLabels)-1 {
					if strings.HasPrefix(templateLabels[idx+1], "[") {
						isArray = true
					}
				}
				// If it is the last value.
				if idx == len(templateLabels)-1 {
					if isArray {
						//If last value and is array add a mocked string value into the array
						c[name] = newList(mockStringValue(templateLabels[idx-1]))
						//If not just add a mocked string value
					} else {
						c[name] = mockStringValue(name)
					}
				} else {
					if isArray {
						//If not last value just add array
						c[name] = new(list)
					} else {
						//If not an array it is a map.
						c[name] = map[string]interface{}{}
					}
				}
			} else {
				//if item is not set in the child path make it.
				if idx < len(templateLabels)-1 {
					//Check to see if nested content is an array.
					if strings.HasPrefix(name, "[") {
						if li, ok := c[name].(list); ok {
							//If it is check to see if this is the last item.
							if idx == len(templateLabels)-1 {
								//Always use the parent name as the name for an array value.
								li.Add(mockStringValue(templateLabels[idx-1]))
								c[name] = li
							} else {
								//Add an empty list
								c[name] = new(list)
							}
						}
					} else {
						//If it's not an array see if it is a map.
						if _, ok := c[name].(map[string]interface{}); !ok {
							//If it's not anything yet... Make it a map.
							c[name] = map[string]interface{}{}
						}

					}
				}

			}
			tmp = c[name]
		case *list:
			var isArray bool
			if idx < len(templateLabels)-1 {
				//Check to see if item is an array.
				if strings.HasPrefix(templateLabels[idx+1], "[") {
					isArray = true
				}
			}
			if idx == len(templateLabels)-1 {
				//If we are on the last value just add the name to the array as a value.
				c.Add(mockStringValue(templateLabels[idx-1]))
			} else if isArray {
				//If it is an array and it's the lst item add both the new array and it's parent name.
				if idx == len(templateLabels)-1 {
					c.Add(mockStringValue(templateLabels[idx-1]))
				} else {
					//If it is not the last item just add an array.
					c.Add(new(list))
				}

			} else {
				//If it's not an array it's gotta be a map.
				//Check to see if it's the last item.
				if idx == len(templateLabels)-1 {
					//If it is add a map and name with mocked value.
					c.Add(map[string]interface{}{name: mockStringValue(name)})
				} else {
					//If it's not the last item... Add just the map.
					c.Add(map[string]interface{}{})
				}
			}
			tmp = (*c)[0]

		}
	}
}

func mockStringValue(name string) string {
	return "test_" + name
}
