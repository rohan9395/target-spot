package util

import (
	"github.com/Jeffail/gabs"
	"strings"
)

func ContextGet(jsonParsed gabs.Container)(string,map[string]*gabs.Container){
	contextArray,_ := jsonParsed.S("queryResult","outputContexts").Children()
	for _,child := range(contextArray){
		context,_ := child.Path("name").Data().(string)
		segments :=strings.Split(context,"/")
		contextName := segments[len(segments) -1]
		if contextName == "target-assist"{
			contextMap,_ := child.S("parameters").ChildrenMap()
			return context,contextMap
		}
	}
	return "",nil
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