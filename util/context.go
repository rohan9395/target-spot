package util

import (
	"github.com/Jeffail/gabs"
)

func ContextGet(jsonParsed gabs.Container)(string,map[string]*gabs.Container){
	contextArray,_ := jsonParsed.S("queryResult","outputContexts").Children()
	child := contextArray[0]
	contextName,_ :=child.Path("name").Data().(string)
	contextMap,_ := child.S("parameters").ChildrenMap()
	return contextName,contextMap
}

func ContextSet(jsonBuild gabs.Container,lifespanCount string,contextName string,contextMap map[string]*gabs.Container)(gabs.Container){
	innerElement:= gabs.New()
	innerElement.Set(contextName,"name")
	innerElement.Set(lifespanCount,"lifespanCount")
	parameters:=gabs.New()
	for k,v := range contextMap{
		parameters.Set(v.Data().(string),k)
	}
	innerElement.Set(parameters.Data(),"parameters")
	jsonBuild.ArrayAppend(innerElement.Data(),"outputContexts")

	return jsonBuild
}