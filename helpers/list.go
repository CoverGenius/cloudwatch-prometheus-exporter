package helpers

func GetLastStringElement(list []string) *string {
	last_element := list[len(list)-1]
	return &last_element
}
