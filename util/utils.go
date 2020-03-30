package util

func PackageArgs(args []string) [][]byte {
	var byteSlice [][]byte
	for i := 0; i < len(args); i++ {
		byteSlice = append(byteSlice, []byte(args[i]))
	}
	return byteSlice
}
