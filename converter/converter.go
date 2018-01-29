package converter

import (
	"errors"
	"fmt"

	dirule "github.com/jcantonio/di-rule"
)

func GetRule(ruleMap map[string]interface{}) (dirule.Rule, error) {
	var rule dirule.Rule
	name := ruleMap["name"].(string)
	entity := ruleMap["entity"].(string)
	actions := []dirule.Action{}

	conditionMap := ruleMap["condition"]

	condition, err := getCondition(conditionMap.(map[string]interface{}))
	if err != nil {
		return rule, err
	}

	rule = dirule.Rule{
		Name:      name,
		Entity:    entity,
		Actions:   actions,
		Condition: condition,
	}

	return rule, nil
}

func getCondition(conditionMap map[string]interface{}) (dirule.Condition, error) {
	operation := conditionMap["op"]

	if operation == nil {
		return nil, errors.New("No op found")
	}

	// Logical Condition
	switch operation {
	case "or", "and":
		condition := &dirule.LogicalCondition{
			Operator: operation.(string)}

		subconditions := conditionMap["conditions"].([]interface{})
		for _, subcondition := range subconditions {
			subconditionMap := subcondition.(map[string]interface{})
			subcondition, err := getCondition(subconditionMap)
			if err != nil {
				return nil, err
			}
			condition.Add(subcondition)
		}
		return condition, nil
	}

	// Value Condition
	path := conditionMap["path"]
	value := conditionMap["value"]

	switch v := value.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		fmt.Printf("Integer: %v", v)
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		fmt.Printf("Float64: %v", v)
	case string:
		condition := &dirule.StringComparatorCondition{
			Path:     path.(string),
			Operator: operation.(string),
			Value:    value.(string),
		}
		fmt.Printf("String: %v", v)
		return condition, nil
	default:
		// And here I'm feeling dumb. ;)
		fmt.Printf("I don't know, ask stackoverflow.")
	}

	return nil, errors.New("Type not handled yet")
}