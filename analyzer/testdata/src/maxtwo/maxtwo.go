package maxtwo

type options struct{ Value string }
type service struct{ inner store }
type store struct{}

var opts options

func target2(first, second string) string                { return first + second }
func target3(first, second, third string) string         { return first + second + third }
func target4(first, second, third, fourth string) string { return first + second + third + fourth }
func sideEffect() string                                 { return "side" }
func normalize(value string) string                      { return value }

func extraBefore(value string) string { // want `function extraBefore is a passthrough to target2`
	return target2("extra", value)
}

func extraBetween(first, second string) string { // want `function extraBetween is a passthrough to target3`
	return target3(first, "extra", second)
}

func extraAfter(value string) string { // want `function extraAfter is a passthrough to target2`
	return target2(value, "extra")
}

func exactlyTwo(value string) string { // want `function exactlyTwo is a passthrough to target3`
	return target3(oneString(), value, sideEffect())
}

func selectorExtra(value string) string { // want `function selectorExtra is a passthrough to target2`
	return target2(value, opts.Value)
}

func conversionExtra(value string) string { // want `function conversionExtra is a passthrough to target2`
	return target2(value, string([]byte{'x'}))
}

func callExtra(value string) string { // want `function callExtra is a passthrough to target2`
	return target2(value, sideEffect())
}

func directAndTransformed(value string) string { // want `function directAndTransformed is a passthrough to target2`
	return target2(value, normalize(value))
}

func (receiver service) explicitReceiver(value string) string { // want `function explicitReceiver is a passthrough to receiver.inner.WithReceiver`
	return receiver.inner.WithReceiver(receiver, value)
}

func (store) WithReceiver(_ service, value string) string { return value }

func zeroWithExtras() string { // want `function zeroWithExtras is a passthrough to target2`
	return target2("one", "two")
}

func tooMany(value string) string {
	return target4(value, "one", "two", "three")
}

func omittedEvenWithinLimit(first, second string) string {
	return target2(first, "extra")
}

func oneString() string { return "one" }
