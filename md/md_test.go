package md

import (
	"testing"

	"github.com/go-test/deep"
)

const s = `**test**

__bro__, that is ***cringe!***

****test****

` + "```" + `go
println('a')
` + "```" + `

> cringe lmao
> kek

>lol
>   cringe

__**test**__

\*lol\*`

const results = `[::b]test[::-]

[::u]bro[::-], that is [::ib]cringe![::-]

[::ib][::i]test[::-][::-]

[grey]┃[-] println([#af0000]'a'[-])

[#789922]
> cringe lmao
> kek[-]

>lol[#789922]
>   cringe[-]

[::u][::b]test[::-][::-]

\*lol\*`

func TestParse(t *testing.T) {
	html := Parse(s)
	if diff := deep.Equal(html, results); diff != nil {
		t.Fatal(diff)
	}
}
