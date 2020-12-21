package main

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/solinox/advent-of-code/2020/pkg/input"
)

func main() {
	foodsString := input.StringSlice(os.Stdin)

	t0 := time.Now()
	log.Println(part1(foodsString), time.Since(t0))

	t0 = time.Now()
	log.Println(part2(foodsString), time.Since(t0))
}

type ingredient string
type allergen string

type food struct {
	Ingredients []ingredient
	Allergens   []allergen
}

func part1(foodsString []string) int {
	foods := parseFoods(foodsString)
	allergenMap, ingredientMap := makeMaps(foods)
	allergenMap, ingredientMap = solve(allergenMap, ingredientMap)

	allergenFreeIngredients := make([]ingredient, 0)
	for ing, alls := range ingredientMap {
		if len(alls) == 0 {
			allergenFreeIngredients = append(allergenFreeIngredients, ing)
		}
	}

	count := 0
	for _, food := range foods {
		count += len(intersectIngredients(food.Ingredients, allergenFreeIngredients))
	}
	return count
}

func part2(foodsString []string) string {
	foods := parseFoods(foodsString)
	allergenMap, ingredientMap := makeMaps(foods)
	allergenMap, ingredientMap = solve(allergenMap, ingredientMap)

	allergens := make([]allergen, 0, len(allergenMap))
	for all := range allergenMap {
		allergens = append(allergens, all)
	}

	sort.Slice(allergens, func(i, j int) bool { return allergens[i] < allergens[j] })

	cdil := ""
	for i := range allergens {
		cdil += "," + string(allergenMap[allergens[i]][0])
	}
	return cdil[1:]
}

func makeMaps(foods []food) (map[allergen][]ingredient, map[ingredient][]allergen) {
	allergenMap, ingredientMap := make(map[allergen][]ingredient), make(map[ingredient][]allergen)

	for _, food := range foods {
		for _, all := range food.Allergens {
			if ings, ok := allergenMap[all]; !ok {
				allergenMap[all] = copyIngredients(food.Ingredients)
			} else {
				allergenMap[all] = intersectIngredients(ings, food.Ingredients)
			}
		}
		for _, ing := range food.Ingredients {
			if alls, ok := ingredientMap[ing]; !ok {
				ingredientMap[ing] = copyAllergens(food.Allergens)
			} else {
				ingredientMap[ing] = unionDistinctAllergens(alls, food.Allergens)
			}
		}
	}

	return allergenMap, ingredientMap
}

func solve(allergenMap map[allergen][]ingredient, ingredientMap map[ingredient][]allergen) (map[allergen][]ingredient, map[ingredient][]allergen) {
	done := false
	solvedAllergens := make([]allergen, 0)
	for !done {
		done = true
		for all, ings := range allergenMap {
			if len(ings) == 1 && !containsAllergen(solvedAllergens, all) {
				done = false
				solvedAllergens = append(solvedAllergens, all)
				ing := ings[0]
				for all2, ings2 := range allergenMap {
					if all == all2 {
						continue
					}
					allergenMap[all2] = removeIngredient(ings2, ing)
				}
				ingredientMap[ing] = []allergen{all}
				for ing2, alls := range ingredientMap {
					if ing == ing2 {
						continue
					}
					ingredientMap[ing2] = removeAllergen(alls, all)
				}
				break
			}
		}
	}

	return allergenMap, ingredientMap
}

func parseFoods(foodsString []string) []food {
	foods := make([]food, 0, len(foodsString))
	for _, foodString := range foodsString {
		fields := strings.Fields(foodString)

		ingredients := make([]ingredient, 0)
		allergens := make([]allergen, 0)

		isIngredient := true
		for _, f := range fields {
			f = strings.Trim(f, "(), ")
			if f == "contains" {
				isIngredient = false
				continue
			}
			if isIngredient {
				ingredients = append(ingredients, ingredient(f))
			} else {
				allergens = append(allergens, allergen(f))
			}
		}
		foods = append(foods, food{ingredients, allergens})
	}
	return foods
}

func copyAllergens(src []allergen) []allergen {
	dst := make([]allergen, len(src))
	copy(dst, src)
	return dst
}

func copyIngredients(src []ingredient) []ingredient {
	dst := make([]ingredient, len(src))
	copy(dst, src)
	return dst
}

func containsAllergen(slice []allergen, val allergen) bool {
	for i := range slice {
		if slice[i] == val {
			return true
		}
	}
	return false
}

func containsIngredient(slice []ingredient, val ingredient) bool {
	for i := range slice {
		if slice[i] == val {
			return true
		}
	}
	return false
}

func intersectIngredients(slice1, slice2 []ingredient) []ingredient {
	intersection := make([]ingredient, 0)
	for i := range slice1 {
		if containsIngredient(slice2, slice1[i]) {
			intersection = append(intersection, slice1[i])
		}
	}
	return intersection
}

func unionDistinctAllergens(slice1, slice2 []allergen) []allergen {
	union := make([]allergen, 0, len(slice1)+len(slice2))
	for _, v := range slice1 {
		if containsAllergen(union, v) {
			continue
		}
		union = append(union, v)
	}
	for _, v := range slice2 {
		if containsAllergen(union, v) {
			continue
		}
		union = append(union, v)
	}
	return union
}

func removeIngredient(slice []ingredient, val ingredient) []ingredient {
	index := -1
	for i := range slice {
		if slice[i] == val {
			index = i
			break
		}
	}
	if index >= 0 {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}

func removeAllergen(slice []allergen, val allergen) []allergen {
	index := -1
	for i := range slice {
		if slice[i] == val {
			index = i
			break
		}
	}
	if index >= 0 {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}
