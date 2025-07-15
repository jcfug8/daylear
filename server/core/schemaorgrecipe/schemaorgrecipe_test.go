package schemaorgrecipe

import (
	"testing"
)

func TestParseIngredient(t *testing.T) {
	tests := []struct {
		input   string
		amount1 float64
		unit1   string
		conj    string
		amount2 float64
		unit2   string
		name    string
	}{
		{"1 ½ cups sugar", 1.5, "cups", "", 0, "", "sugar"},
		{"1 and 1/2 cup sugar", 1.5, "cup", "", 0, "", "sugar"},
		{"1 + 1/2 cup sugar", 1.5, "cup", "", 0, "", "sugar"},
		{"1 cup sugar", 1, "cup", "", 0, "", "sugar"},
		{"1 cup and 2 tablespoons sugar", 1, "cup", "&&", 2, "tablespoons", "sugar"},
		{"1 cup or 100 grams flour", 1, "cup", "||", 100, "grams", "flour"},
		{"1 to 2 cups milk", 1, "cups", "--", 2, "cups", "milk"},
		{"2 tablespoons + 1 teaspoon vanilla", 2, "tablespoons", "&&", 1, "teaspoon", "vanilla"},
		{"3 eggs", 3, "eggs", "", 0, "", ""},
		{"1 package of yeast", 1, "package", "", 0, "", "of yeast"},
		{"1 to 2 packages of yeast", 1, "packages", "--", 2, "packages", "of yeast"},
		{"1 1/2 cups sugar", 1.5, "cups", "", 0, "", "sugar"},
		{"1 1/2 cups and 1/4 cup sugar", 1.5, "cups", "&&", 0.25, "cup", "sugar"},
		{"1 1/2 cups or 1/4 cup sugar", 1.5, "cups", "||", 0.25, "cup", "sugar"},
		{"1 1/2 cups to 1/4 cup sugar", 1.5, "cups", "--", 0.25, "cup", "sugar"},
		{"1 1/2 cups + 1/4 cup sugar", 1.5, "cups", "&&", 0.25, "cup", "sugar"},
		{"1 1/2 cups - 2 1/4 cups sugar", 1.5, "cups", "--", 2.25, "cups", "sugar"},
		{"1 1/2 - 2 1/4 cups sugar", 1.5, "cups", "--", 2.25, "cups", "sugar"},
		{"Salt to taste", 0, "", "", 0, "", "Salt to taste"},
		{"Salt and pepper to taste", 0, "", "", 0, "", "Salt and pepper to taste"},
	}

	for _, tt := range tests {
		amount1, unit1, conj, amount2, unit2, name := ParseIngredient(tt.input)
		if amount1 != tt.amount1 || unit1 != tt.unit1 || conj != tt.conj || amount2 != tt.amount2 || unit2 != tt.unit2 || name != tt.name {
			t.Errorf("ParseIngredient(%q) = \nhave: \t%v, \t%q, \t%q, \t%v, \t%q, \t%q \nwant: \t%v, \t%q, \t%q, \t%v, \t%q, \t%q", tt.input, amount1, unit1, conj, amount2, unit2, name, tt.amount1, tt.unit1, tt.conj, tt.amount2, tt.unit2, tt.name)
		}
	}
}
