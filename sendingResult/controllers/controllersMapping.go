package controllers

type MapInfo struct {
	MethodName    string
	ControllerRef interface{}
}

//ControllersMapping is a mapping for functions invocation through reflection
var ControllersMapping map[int]MapInfo = map[int]MapInfo{
	1300: MapInfo{"HandleSendOrder", &AgentsController{}},
	1301: MapInfo{"HandleSendCancel", &AgentsController{}},
}
