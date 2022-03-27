package workflow

import (
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestInsertInFuncGoEnd(t *testing.T) {
	output, err := insertGo(&config.Insert{
		Section:   "func main()",
		Placement: "end",
	}, `package main

import "fmt"

func main() {
	fmt.Printf("hello world")
}
`, `	fmt.Printf("hello again")`)

	assert.Nil(t, err)
	assert.Equal(t, `package main

import "fmt"

func main() {
	fmt.Printf("hello world")
	fmt.Printf("hello again")
}
`, output)

}

func TestInsertInFuncGoBegin(t *testing.T) {
	output, err := insertGo(&config.Insert{
		Section:   "func main()",
		Placement: "begin",
	}, `package main

import "fmt"

func main() {
	fmt.Printf("hello world")
}
`, `	fmt.Printf("hello first")`)

	assert.Nil(t, err)
	assert.Equal(t, `package main

import "fmt"

func main() {
	fmt.Printf("hello first")
	fmt.Printf("hello world")
}
`, output)

}

func TestInsertImportGoDuplicate(t *testing.T) {
	output, err := insertGo(&config.Insert{
		Section: "import",
	}, `package main

import "fmt"

func main() {
	fmt.Printf("hello world")
}
`, `"fmt"`)

	assert.Nil(t, err)
	assert.Equal(t, `package main

import (
	"fmt"
)

func main() {
	fmt.Printf("hello world")
}`, output)

}

func TestInsertImportGoNew(t *testing.T) {
	output, err := insertGo(&config.Insert{
		Section: "import",
	}, `package main

import "fmt"

func main() {
	fmt.Printf("hello world")
}
`, `"testing"`)

	assert.Nil(t, err)
	assert.Equal(t, `package main

import (
	"fmt"
	"testing"
)

func main() {
	fmt.Printf("hello world")
}`, output)

}

func TestInsertImportGoMultiple(t *testing.T) {
	output, err := insertGo(&config.Insert{
		Section: "import",
	}, `package main

import "fmt"

func main() {
	fmt.Printf("hello world")
}
`, `"testing"
"fmt"`)

	assert.Nil(t, err)
	assert.Equal(t, `package main

import (
	"fmt"
	"testing"
)

func main() {
	fmt.Printf("hello world")
}`, output)

}
