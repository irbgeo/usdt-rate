package rateerr

type err struct {
	error
	args map[string]any
}

func New(e error, args ...any) *err {
	if internalError, ok := e.(*err); ok {
		internalError.addArgs(args...)
		return internalError
	}

	return &err{
		error: e,
		args:  buildArgs(args...),
	}
}

func (s *err) Args() map[string]any {
	return s.args
}

func buildArgs(args ...any) map[string]any {
	fields := make(map[string]any, len(args)/2)
	if len(args)%2 == 0 {
		for i := 0; i < len(args); i += 2 {
			fields[args[i].(string)] = args[i+1]
		}
	}

	return fields
}

func (s *err) addArgs(args ...any) {
	if len(args)%2 != 0 {
		return
	}

	for i := 0; i < len(args); i += 2 {
		s.args[args[i].(string)] = args[i+1]
	}
}
