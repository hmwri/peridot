package object

//Env Enviroment
type Env struct {
	envs map[string]Object
	out  *Env
}

//NewEnv make new Env struct
func NewEnv() *Env {
	e := make(map[string]Object)
	return &Env{envs: e}
}

//AddEnv Create new enclosed environment
func AddEnv(out *Env) *Env {
	in := NewEnv()
	in.out = out
	return in
}

//GetEnv Get env
func (e *Env) GetEnv(name string) (Object, bool) {
	obj, found := e.envs[name]
	if !found && e.out != nil {
		obj, found = e.out.GetEnv(name)
	}
	return obj, found
}

//SetEnv Set env
func (e *Env) SetEnv(name string, value Object) Object {
	e.envs[name] = value
	return value
}
