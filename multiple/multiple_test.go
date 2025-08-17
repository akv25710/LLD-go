package multiple

import (
	"context"
	"testing"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func TestMultiple(t *testing.T) {
	manager := NewManager()
	manager.Start(context.Background(), 5)
}
