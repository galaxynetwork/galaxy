package util

func GenStringWithLength(length int) string {
	bz := []byte{}
	for i := 0; i < length; i++ {
		bz = append(bz, byte(i))
	}
	return string(bz[:])
}
