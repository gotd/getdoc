package href

import "strconv"

func superscript(i int) string {
	var out []rune
	for _, r := range strconv.Itoa(i) {
		var superscript rune
		switch r {
		case '0':
			superscript = '⁰'
		case '1':
			superscript = '¹'
		case '2':
			superscript = '²'
		case '3':
			superscript = '³'
		case '4':
			superscript = '⁴'
		case '5':
			superscript = '⁵'
		case '6':
			superscript = '⁶'
		case '7':
			superscript = '⁷'
		case '8':
			superscript = '⁸'
		case '9':
			superscript = '⁹'
		default:
			panic(r)
		}

		out = append(out, superscript)
	}

	return string(out)
}
