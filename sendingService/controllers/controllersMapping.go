package controllers

//MapInfo is a structure with information for controller methods invocation through reflection
type MapInfo struct {
	MethodName    string
	ControllerRef interface{}
}

//ControllersMapping a map for registering "controllers"
var ControllersMapping map[int]MapInfo = map[int]MapInfo{
	1300: MapInfo{"SendOrder", &AgentsController{}},
	1301: MapInfo{"CancelOrder", &AgentsController{}},
}
