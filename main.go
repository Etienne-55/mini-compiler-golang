package main

import(
	"mini_go_compiler/lexical"
)
func main() {

	input := `int x = 10;
int w = 1.1;
float y = 3.14;
z == x;
x = x++;
/* comentário de múltiplas linhas
	aslkd;laskd;laskd;askd;alskd;
	sadjlkjdlaskjdlaskjdlaskjdlaskjd
*/
if (x >= y) {
    print(x);
}`

	scanner := lexical.NewScanner(input)
	tokens := scanner.ScanAll()

	lexical.PrintTokensGrouped(tokens)
}

