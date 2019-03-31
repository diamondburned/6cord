package md

import (
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func TestParse(t *testing.T) {
	var testSuite = []string{
		"asdasd\n\nasd\nasdasd\n\n",
		"*test* **strong** ~~no~~",
		"____**__test__**",
		"```js\nconsole.log(\"Your mom gay\");```",
		">be me\n>is retarded\nokay homo",
		"`just normal code`",
		"https://google.com",
		"### expert mode!",
		"| lol | retard |\n| - | - |\n| | |",
		"![wtf](https://google.com)",
		"[ur mom](https://google.com)",
		"- that's\n- bullshit",
		"3. ur mom\n4. gay",
	}

	var result = `asdasd
asd
asdasd
[::i]test[::-] [::b]strong[::-] [::s]no[::-]
__[::u]**[::-]test__**

[grey]┃[-] console.log([#af0000]"Your mom gay"[-]);
[grey]┃[-] 

[green]>be me[-]
[green]>is retarded[-]
okay homo
[:#4f4f4f:]just normal code[:-:]
https://google.com
### expert mode!
| lol | retard |
| - | - |
| | |
![wtf](https://google.com)
[ur mom](https://google.com)
- that's
- bullshit
3. ur mom
4. gay`

	if p := Parse(strings.Join(testSuite, "\n\n")); p != result {
		println(p)
		println()
		t.Errorf("Test failed---\n%v", diff.LineDiff(p, result))
	}
}
