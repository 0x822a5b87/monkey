package object

var BuiltIns = []*BuiltIn{
	{
		Name: builtInFnNameLen,
		BuiltInFn: func(objs ...Object) Object {
			if len(objs) != 1 {
				return newWrongArgumentSizeError(len(objs), 1)
			}

			obj := objs[0]
			builtInLen, ok := obj.(Len)
			if !ok {
				return newWrongArgumentTypeError(builtInFnNameLen, obj.Type())
			}
			return &Integer{Value: builtInLen.Len().Value}
		},
	},

	{
		Name: builtInFnNameFirst,
		BuiltInFn: func(objs ...Object) Object {
			if len(objs) != 1 {
				return newWrongArgumentSizeError(len(objs), 1)
			}

			obj := objs[0]
			builtInFirst, ok := obj.(List)
			if !ok {
				return newWrongArgumentTypeError(builtInFnNameFirst, obj.Type())
			}
			return builtInFirst.First()
		},
	},

	{
		Name: builtInFnNameLast,
		BuiltInFn: func(objs ...Object) Object {
			if len(objs) != 1 {
				return newWrongArgumentSizeError(len(objs), 1)
			}

			obj := objs[0]
			buildInLast, ok := obj.(List)
			if !ok {
				return newWrongArgumentTypeError(builtInFnNameLast, obj.Type())
			}
			return buildInLast.Last()
		},
	},

	{
		Name: builtInFnNameRest,
		BuiltInFn: func(objs ...Object) Object {
			if len(objs) != 1 {
				return newWrongArgumentSizeError(len(objs), 1)
			}

			obj := objs[0]
			rest, ok := obj.(Rest)
			if !ok {
				return newWrongArgumentTypeError(builtInFnNameRest, obj.Type())
			}
			return rest.Rest()
		},
	},

	{
		Name: builtInFnNamePush,
		BuiltInFn: func(objs ...Object) Object {
			if len(objs) != 2 {
				return newWrongArgumentSizeError(len(objs), 2)
			}

			obj := objs[0]
			push, ok := obj.(Push)
			if !ok {
				return newWrongArgumentTypeError(builtInFnNamePush, obj.Type())
			}
			return push.Push(objs[1])
		},
	},
}
