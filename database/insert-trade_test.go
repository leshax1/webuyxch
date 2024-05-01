package database

import "testing"

func TestAdd(t *testing.T) {
	cases := []string{"1", "2", "3", "4", "5"}

	for _, c := range cases {
		t.Run("Test"+c, func(t *testing.T) {

		})
	}

	/*t.Run("Fail", func(t *testing.T) {
		t.Errorf("Error")
	})*/
}
