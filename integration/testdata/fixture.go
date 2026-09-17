package fixture

func target(first string, rest ...string) string {
	return first
}

func ordinary(value string) string {
	return target(value)
}

func withExtra(value string) string {
	return target("prefix", value)
}

//nolint:passthrough // This intentional API boundary tests golangci-lint suppression.
func suppressed(value string) string {
	return target(value)
}
