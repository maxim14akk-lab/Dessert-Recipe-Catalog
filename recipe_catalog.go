// recipe_catalog.go — Go версия

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Recipe struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
	PrepTime     int      `json:"prep_time"`
	Difficulty   string   `json:"difficulty"`
	Rating       float64  `json:"rating"`
	RatingsCount int      `json:"ratings_count"`
}

type Catalog struct {
	Recipes []Recipe `json:"recipes"`
	file    string
}

func NewCatalog(file string) *Catalog {
	c := &Catalog{file: file}
	c.load()
	return c
}

func (c *Catalog) load() {
	data, err := os.ReadFile(c.file)
	if err != nil {
		c.Recipes = []Recipe{}
		return
	}
	json.Unmarshal(data, &c.Recipes)
}

func (c *Catalog) save() {
	data, _ := json.MarshalIndent(c.Recipes, "", "  ")
	os.WriteFile(c.file, data, 0644)
}

func (c *Catalog) add(name, category string, ingredients []string, instructions string, prepTime int, difficulty string) int {
	id := len(c.Recipes) + 1
	c.Recipes = append(c.Recipes, Recipe{
		ID:           id,
		Name:         name,
		Category:     category,
		Ingredients:  ingredients,
		Instructions: instructions,
		PrepTime:     prepTime,
		Difficulty:   difficulty,
		Rating:       0,
		RatingsCount: 0,
	})
	c.save()
	return id
}

func (c *Catalog) listAll() {
	if len(c.Recipes) == 0 {
		fmt.Println("\u001B[33mКаталог пуст.\u001B[0m")
		return
	}
	fmt.Printf("\u001B[36m%-4s %-25s %-12s %-8s %-10s %-8s\u001B[0m\n", "ID", "Название", "Категория", "Время", "Сложность", "Рейтинг")
	fmt.Println(strings.Repeat("-", 75))
	for _, r := range c.Recipes {
		ratingStr := "—"
		if r.RatingsCount > 0 {
			ratingStr = fmt.Sprintf("%.1f★", r.Rating)
		}
		diffColor := "\u001B[32m"
		if r.Difficulty == "средняя" {
			diffColor = "\u001B[33m"
		} else if r.Difficulty == "сложная" {
			diffColor = "\u001B[31m"
		}
		fmt.Printf("%-4d %-25s %-12s %-8d %s%-10s\u001B[0m %-8s\n", r.ID, r.Name, r.Category, r.PrepTime, diffColor, r.Difficulty, ratingStr)
	}
}

func (c *Catalog) search(query string) {
	query = strings.ToLower(query)
	results := []Recipe{}
	for _, r := range c.Recipes {
		if strings.Contains(strings.ToLower(r.Name), query) {
			results = append(results, r)
		} else {
			for _, ing := range r.Ingredients {
				if strings.Contains(strings.ToLower(ing), query) {
					results = append(results, r)
					break
				}
			}
		}
	}
	if len(results) == 0 {
		fmt.Println("\u001B[33mНичего не найдено.\u001B[0m")
		return
	}
	for _, r := range results {
		fmt.Printf("%d: %s | %s | %d мин | %s | %.1f★\n", r.ID, r.Name, r.Category, r.PrepTime, r.Difficulty, r.Rating)
	}
}

func (c *Catalog) filterBy(field, value string) {
	results := []Recipe{}
	for _, r := range c.Recipes {
		switch field {
		case "category":
			if strings.EqualFold(r.Category, value) {
				results = append(results, r)
			}
		case "difficulty":
			if strings.EqualFold(r.Difficulty, value) {
				results = append(results, r)
			}
		case "prep_time":
			if v, err := strconv.Atoi(value); err == nil && r.PrepTime <= v {
				results = append(results, r)
			}
		}
	}
	if len(results) == 0 {
		fmt.Println("\u001B[33mНет рецептов, соответствующих фильтру.\u001B[0m")
		return
	}
	for _, r := range results {
		fmt.Printf("%d: %s | %s | %d мин | %s\n", r.ID, r.Name, r.Category, r.PrepTime, r.Difficulty)
	}
}

func (c *Catalog) sortBy(field string, reverse bool) {
	switch field {
	case "name":
		if reverse {
			SortBy(c.Recipes, func(i, j int) bool { return c.Recipes[i].Name > c.Recipes[j].Name })
		} else {
			SortBy(c.Recipes, func(i, j int) bool { return c.Recipes[i].Name < c.Recipes[j].Name })
		}
	case "prep_time":
		if reverse {
			SortBy(c.Recipes, func(i, j int) bool { return c.Recipes[i].PrepTime > c.Recipes[j].PrepTime })
		} else {
			SortBy(c.Recipes, func(i, j int) bool { return c.Recipes[i].PrepTime < c.Recipes[j].PrepTime })
		}
	case "rating":
		if reverse {
			SortBy(c.Recipes, func(i, j int) bool { return c.Recipes[i].Rating > c.Recipes[j].Rating })
		} else {
			SortBy(c.Recipes, func(i, j int) bool { return c.Recipes[i].Rating < c.Recipes[j].Rating })
		}
	default:
		fmt.Println("\u001B[31mНеверное поле для сортировки.\u001B[0m")
		return
	}
	c.listAll()
}

func SortBy(recipes []Recipe, less func(i, j int) bool) {
	for i := 0; i < len(recipes)-1; i++ {
		for j := i + 1; j < len(recipes); j++ {
			if less(i, j) {
				recipes[i], recipes[j] = recipes[j], recipes[i]
			}
		}
	}
}

func (c *Catalog) delete(id int) bool {
	for i, r := range c.Recipes {
		if r.ID == id {
			c.Recipes = append(c.Recipes[:i], c.Recipes[i+1:]...)
			c.save()
			return true
		}
	}
	return false
}

func (c *Catalog) edit(id int, field, value string) bool {
	for i, r := range c.Recipes {
		if r.ID == id {
			switch field {
			case "name":
				c.Recipes[i].Name = value
			case "category":
				c.Recipes[i].Category = value
			case "ingredients":
				c.Recipes[i].Ingredients = strings.Split(value, ",")
				for j := range c.Recipes[i].Ingredients {
					c.Recipes[i].Ingredients[j] = strings.TrimSpace(c.Recipes[i].Ingredients[j])
				}
			case "instructions":
				c.Recipes[i].Instructions = value
			case "prep_time":
				if v, err := strconv.Atoi(value); err == nil {
					c.Recipes[i].PrepTime = v
				} else {
					return false
				}
			case "difficulty":
				c.Recipes[i].Difficulty = value
			default:
				return false
			}
			c.save()
			return true
		}
	}
	return false
}

func (c *Catalog) rate(id int, rating int) bool {
	for i, r := range c.Recipes {
		if r.ID == id {
			total := r.Rating*float64(r.RatingsCount) + float64(rating)
			c.Recipes[i].RatingsCount++
			c.Recipes[i].Rating = total / float64(c.Recipes[i].RatingsCount)
			c.save()
			return true
		}
	}
	return false
}

func (c *Catalog) stats() {
	if len(c.Recipes) == 0 {
		fmt.Println("Нет данных.")
		return
	}
	total := len(c.Recipes)
	categories := make(map[string]int)
	difficulties := make(map[string]int)
	totalRating := 0.0
	ratedCount := 0
	for _, r := range c.Recipes {
		categories[r.Category]++
		difficulties[r.Difficulty]++
		if r.RatingsCount > 0 {
			totalRating += r.Rating
			ratedCount++
		}
	}
	avgRating := 0.0
	if ratedCount > 0 {
		avgRating = totalRating / float64(ratedCount)
	}
	fmt.Println("\u001B[36m📊 Статистика:\u001B[0m")
	fmt.Printf("  Всего рецептов: %d\n", total)
	fmt.Printf("  Средний рейтинг: %.1f★ (из %d оценённых)\n", avgRating, ratedCount)
	fmt.Println("  По категориям:")
	for c, cnt := range categories {
		fmt.Printf("    %s: %d\n", c, cnt)
	}
	fmt.Println("  По сложности:")
	for d, cnt := range difficulties {
		fmt.Printf("    %s: %d\n", d, cnt)
	}
}

func (c *Catalog) exportJSON(filename string) {
	data, _ := json.MarshalIndent(c.Recipes, "", "  ")
	os.WriteFile(filename, data, 0644)
	fmt.Printf("\u001B[32m💾 Экспорт JSON: %s\u001B[0m\n", filename)
}

func (c *Catalog) exportCSV(filename string) {
	if len(c.Recipes) == 0 {
		fmt.Println("\u001B[33mНет данных для экспорта.\u001B[0m")
		return
	}
	file, _ := os.Create(filename)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"ID", "Название", "Категория", "Ингредиенты", "Инструкция", "Время", "Сложность", "Рейтинг"})
	for _, r := range c.Recipes {
		writer.Write([]string{
			strconv.Itoa(r.ID),
			r.Name,
			r.Category,
			strings.Join(r.Ingredients, ", "),
			r.Instructions,
			strconv.Itoa(r.PrepTime),
			r.Difficulty,
			fmt.Sprintf("%.1f", r.Rating),
		})
	}
	fmt.Printf("\u001B[32m💾 Экспорт CSV: %s\u001B[0m\n", filename)
}

func main() {
	catalog := NewCatalog("recipes.json")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n\u001B[36m🍰 Dessert Recipe Catalog (Go)\u001B[0m")
		fmt.Println("1. Добавить рецепт")
		fmt.Println("2. Показать все рецепты")
		fmt.Println("3. Поиск рецептов")
		fmt.Println("4. Фильтрация")
		fmt.Println("5. Сортировка")
		fmt.Println("6. Удалить рецепт")
		fmt.Println("7. Редактировать рецепт")
		fmt.Println("8. Оценить рецепт")
		fmt.Println("9. Статистика")
		fmt.Println("10. Экспорт")
		fmt.Println("11. Выход")
		fmt.Print("Выберите действие: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			fmt.Print("Название: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			fmt.Print("Категория (торт, пирожное, печенье, мороженое, другое): ")
			category, _ := reader.ReadString('\n')
			category = strings.TrimSpace(category)
			fmt.Print("Ингредиенты (через запятую): ")
			ingStr, _ := reader.ReadString('\n')
			ingredients := strings.Split(ingStr, ",")
			for j := range ingredients {
				ingredients[j] = strings.TrimSpace(ingredients[j])
			}
			fmt.Print("Инструкция: ")
			instructions, _ := reader.ReadString('\n')
			instructions = strings.TrimSpace(instructions)
			fmt.Print("Время приготовления (мин): ")
			timeStr, _ := reader.ReadString('\n')
			prepTime, _ := strconv.Atoi(strings.TrimSpace(timeStr))
			fmt.Print("Сложность (лёгкая/средняя/сложная): ")
			difficulty, _ := reader.ReadString('\n')
			difficulty = strings.TrimSpace(strings.ToLower(difficulty))
			if difficulty != "лёгкая" && difficulty != "средняя" && difficulty != "сложная" {
				difficulty = "средняя"
			}
			id := catalog.add(name, category, ingredients, instructions, prepTime, difficulty)
			fmt.Printf("\u001B[32m✅ Рецепт добавлен (ID: %d)\u001B[0m\n", id)
		case "2":
			catalog.listAll()
		case "3":
			fmt.Print("Введите запрос (название или ингредиент): ")
			query, _ := reader.ReadString('\n')
			query = strings.TrimSpace(query)
			catalog.search(query)
		case "4":
			fmt.Println("Фильтровать по: category, difficulty, prep_time")
			fmt.Print("Поле: ")
			field, _ := reader.ReadString('\n')
			field = strings.TrimSpace(strings.ToLower(field))
			fmt.Print("Значение: ")
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			catalog.filterBy(field, value)
		case "5":
			fmt.Println("Сортировать по: name, prep_time, rating")
			fmt.Print("Поле: ")
			field, _ := reader.ReadString('\n')
			field = strings.TrimSpace(strings.ToLower(field))
			fmt.Print("По убыванию? (y/n): ")
			revStr, _ := reader.ReadString('\n')
			reverse := strings.TrimSpace(strings.ToLower(revStr)) == "y"
			catalog.sortBy(field, reverse)
		case "6":
			catalog.listAll()
			fmt.Print("Введите ID для удаления: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			if catalog.delete(id) {
				fmt.Println("\u001B[32m✅ Рецепт удалён.\u001B[0m")
			} else {
				fmt.Println("\u001B[31m❌ Рецепт не найден.\u001B[0m")
			}
		case "7":
			catalog.listAll()
			fmt.Print("Введите ID для редактирования: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): ")
			field, _ := reader.ReadString('\n')
			field = strings.TrimSpace(strings.ToLower(field))
			fmt.Print("Новое значение: ")
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)
			if catalog.edit(id, field, value) {
				fmt.Println("\u001B[32m✅ Рецепт обновлён.\u001B[0m")
			} else {
				fmt.Println("\u001B[31m❌ Не удалось обновить.\u001B[0m")
			}
		case "8":
			catalog.listAll()
			fmt.Print("Введите ID для оценки: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			fmt.Print("Оценка (1-5): ")
			ratingStr, _ := reader.ReadString('\n')
			rating, _ := strconv.Atoi(strings.TrimSpace(ratingStr))
			if rating >= 1 && rating <= 5 {
				if catalog.rate(id, rating) {
					fmt.Println("\u001B[32m✅ Оценка добавлена.\u001B[0m")
				} else {
					fmt.Println("\u001B[31m❌ Рецепт не найден.\u001B[0m")
				}
			} else {
				fmt.Println("\u001B[31m❌ Оценка должна быть от 1 до 5.\u001B[0m")
			}
		case "9":
			catalog.stats()
		case "10":
			fmt.Println("1. Экспорт в JSON")
			fmt.Println("2. Экспорт в CSV")
			fmt.Print("Выберите формат: ")
			sub, _ := reader.ReadString('\n')
			sub = strings.TrimSpace(sub)
			if sub == "1" {
				catalog.exportJSON("recipes_export.json")
			} else if sub == "2" {
				catalog.exportCSV("recipes_export.csv")
			} else {
				fmt.Println("\u001B[31mНеверный выбор.\u001B[0m")
			}
		case "11":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("\u001B[31mНеверный выбор.\u001B[0m")
		}
	}
}
