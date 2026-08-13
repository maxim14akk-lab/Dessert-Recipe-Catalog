// recipe_catalog.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;

class Recipe {
    public int Id { get; set; }
    public string Name { get; set; }
    public string Category { get; set; }
    public List<string> Ingredients { get; set; }
    public string Instructions { get; set; }
    public int PrepTime { get; set; }
    public string Difficulty { get; set; }
    public double Rating { get; set; }
    public int RatingsCount { get; set; }

    public void AddRating(int rating) {
        double total = Rating * RatingsCount + rating;
        RatingsCount++;
        Rating = total / RatingsCount;
    }
}

class RecipeCatalog {
    private List<Recipe> recipes = new List<Recipe>();
    private const string DataFile = "recipes.json";

    public RecipeCatalog() {
        Load();
    }

    private void Load() {
        if (File.Exists(DataFile)) {
            try {
                string json = File.ReadAllText(DataFile);
                recipes = JsonSerializer.Deserialize<List<Recipe>>(json) ?? new List<Recipe>();
            } catch {
                recipes = new List<Recipe>();
            }
        }
    }

    private void Save() {
        string json = JsonSerializer.Serialize(recipes, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(DataFile, json);
    }

    public int Add(string name, string category, List<string> ingredients, string instructions, int prepTime, string difficulty) {
        int id = recipes.Count + 1;
        recipes.Add(new Recipe {
            Id = id,
            Name = name,
            Category = category,
            Ingredients = ingredients,
            Instructions = instructions,
            PrepTime = prepTime,
            Difficulty = difficulty,
            Rating = 0,
            RatingsCount = 0
        });
        Save();
        return id;
    }

    public void ListAll() {
        if (recipes.Count == 0) {
            Console.WriteLine("\u001B[33mКаталог пуст.\u001B[0m");
            return;
        }
        Console.WriteLine($"\u001B[36m{"ID",-4} {"Название",-25} {"Категория",-12} {"Время",-8} {"Сложность",-10} {"Рейтинг",-8}\u001B[0m");
        Console.WriteLine(new string('-', 75));
        foreach (var r in recipes) {
            string ratingStr = r.RatingsCount > 0 ? $"{r.Rating:F1}★" : "—";
            string diffColor = r.Difficulty == "лёгкая" ? "\u001B[32m" : r.Difficulty == "средняя" ? "\u001B[33m" : "\u001B[31m";
            Console.WriteLine($"{r.Id,-4} {r.Name,-25} {r.Category,-12} {r.PrepTime,-8} {diffColor}{r.Difficulty,-10}\u001B[0m {ratingStr,-8}");
        }
    }

    public void Search(string query) {
        query = query.ToLower();
        var results = recipes.Where(r =>
            r.Name.ToLower().Contains(query) ||
            r.Ingredients.Any(i => i.ToLower().Contains(query))
        ).ToList();
        if (results.Count == 0) {
            Console.WriteLine("\u001B[33mНичего не найдено.\u001B[0m");
            return;
        }
        foreach (var r in results) {
            Console.WriteLine($"{r.Id}: {r.Name} | {r.Category} | {r.PrepTime} мин | {r.Difficulty} | {r.Rating:F1}★");
        }
    }

    public void FilterBy(string field, string value) {
        List<Recipe> results = new List<Recipe>();
        switch (field.ToLower()) {
            case "category":
                results = recipes.Where(r => r.Category.Equals(value, StringComparison.OrdinalIgnoreCase)).ToList();
                break;
            case "difficulty":
                results = recipes.Where(r => r.Difficulty.Equals(value, StringComparison.OrdinalIgnoreCase)).ToList();
                break;
            case "prep_time":
                if (int.TryParse(value, out int v)) {
                    results = recipes.Where(r => r.PrepTime <= v).ToList();
                }
                break;
        }
        if (results.Count == 0) {
            Console.WriteLine("\u001B[33mНет рецептов, соответствующих фильтру.\u001B[0m");
            return;
        }
        foreach (var r in results) {
            Console.WriteLine($"{r.Id}: {r.Name} | {r.Category} | {r.PrepTime} мин | {r.Difficulty}");
        }
    }

    public void SortBy(string field, bool reverse) {
        switch (field.ToLower()) {
            case "name":
                recipes = reverse ? recipes.OrderByDescending(r => r.Name).ToList() : recipes.OrderBy(r => r.Name).ToList();
                break;
            case "prep_time":
                recipes = reverse ? recipes.OrderByDescending(r => r.PrepTime).ToList() : recipes.OrderBy(r => r.PrepTime).ToList();
                break;
            case "rating":
                recipes = reverse ? recipes.OrderByDescending(r => r.Rating).ToList() : recipes.OrderBy(r => r.Rating).ToList();
                break;
            default:
                Console.WriteLine("\u001B[31mНеверное поле для сортировки.\u001B[0m");
                return;
        }
        ListAll();
    }

    public bool Delete(int id) {
        var recipe = recipes.FirstOrDefault(r => r.Id == id);
        if (recipe != null) {
            recipes.Remove(recipe);
            Save();
            return true;
        }
        return false;
    }

    public bool Edit(int id, string field, string value) {
        var recipe = recipes.FirstOrDefault(r => r.Id == id);
        if (recipe == null) return false;
        switch (field.ToLower()) {
            case "name": recipe.Name = value; break;
            case "category": recipe.Category = value; break;
            case "ingredients": recipe.Ingredients = value.Split(',').Select(s => s.Trim()).ToList(); break;
            case "instructions": recipe.Instructions = value; break;
            case "prep_time": if (int.TryParse(value, out int t)) recipe.PrepTime = t; else return false; break;
            case "difficulty": recipe.Difficulty = value; break;
            default: return false;
        }
        Save();
        return true;
    }

    public bool Rate(int id, int rating) {
        var recipe = recipes.FirstOrDefault(r => r.Id == id);
        if (recipe == null) return false;
        recipe.AddRating(rating);
        Save();
        return true;
    }

    public void Stats() {
        if (recipes.Count == 0) {
            Console.WriteLine("Нет данных.");
            return;
        }
        var categories = recipes.GroupBy(r => r.Category).ToDictionary(g => g.Key, g => g.Count());
        var difficulties = recipes.GroupBy(r => r.Difficulty).ToDictionary(g => g.Key, g => g.Count());
        double totalRating = recipes.Where(r => r.RatingsCount > 0).Sum(r => r.Rating);
        int ratedCount = recipes.Count(r => r.RatingsCount > 0);
        double avgRating = ratedCount > 0 ? totalRating / ratedCount : 0;
        Console.WriteLine("\u001B[36m📊 Статистика:\u001B[0m");
        Console.WriteLine($"  Всего рецептов: {recipes.Count}");
        Console.WriteLine($"  Средний рейтинг: {avgRating:F1}★ (из {ratedCount} оценённых)");
        Console.WriteLine("  По категориям:");
        foreach (var kv in categories.OrderByDescending(kv => kv.Value)) {
            Console.WriteLine($"    {kv.Key}: {kv.Value}");
        }
        Console.WriteLine("  По сложности:");
        foreach (var kv in difficulties) {
            Console.WriteLine($"    {kv.Key}: {kv.Value}");
        }
    }

    public void ExportJSON(string filename = "recipes_export.json") {
        string json = JsonSerializer.Serialize(recipes, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(filename, json);
        Console.WriteLine($"\u001B[32m💾 Экспорт JSON: {filename}\u001B[0m");
    }

    public void ExportCSV(string filename = "recipes_export.csv") {
        if (recipes.Count == 0) {
            Console.WriteLine("\u001B[33mНет данных для экспорта.\u001B[0m");
            return;
        }
        using var writer = new StreamWriter(filename);
        writer.WriteLine("ID,Название,Категория,Ингредиенты,Инструкция,Время,Сложность,Рейтинг");
        foreach (var r in recipes) {
            writer.WriteLine($"{r.Id},{r.Name},{r.Category},\"{string.Join(", ", r.Ingredients)}\",\"{r.Instructions}\",{r.PrepTime},{r.Difficulty},{r.Rating:F1}");
        }
        Console.WriteLine($"\u001B[32m💾 Экспорт CSV: {filename}\u001B[0m");
    }

    public static void Main() {
        var catalog = new RecipeCatalog();
        while (true) {
            Console.WriteLine("\n\u001B[36m🍰 Dessert Recipe Catalog (C#)\u001B[0m");
            Console.WriteLine("1. Добавить рецепт");
            Console.WriteLine("2. Показать все рецепты");
            Console.WriteLine("3. Поиск рецептов");
            Console.WriteLine("4. Фильтрация");
            Console.WriteLine("5. Сортировка");
            Console.WriteLine("6. Удалить рецепт");
            Console.WriteLine("7. Редактировать рецепт");
            Console.WriteLine("8. Оценить рецепт");
            Console.WriteLine("9. Статистика");
            Console.WriteLine("10. Экспорт");
            Console.WriteLine("11. Выход");
            Console.Write("Выберите действие: ");
            string choice = Console.ReadLine().Trim();
            switch (choice) {
                case "1":
                    Console.Write("Название: ");
                    string name = Console.ReadLine().Trim();
                    Console.Write("Категория (торт, пирожное, печенье, мороженое, другое): ");
                    string category = Console.ReadLine().Trim();
                    Console.Write("Ингредиенты (через запятую): ");
                    var ingredients = Console.ReadLine().Split(',').Select(s => s.Trim()).ToList();
                    Console.Write("Инструкция: ");
                    string instructions = Console.ReadLine().Trim();
                    Console.Write("Время приготовления (мин): ");
                    int prepTime = int.Parse(Console.ReadLine().Trim());
                    Console.Write("Сложность (лёгкая/средняя/сложная): ");
                    string difficulty = Console.ReadLine().Trim().ToLower();
                    if (difficulty != "лёгкая" && difficulty != "средняя" && difficulty != "сложная") difficulty = "средняя";
                    int id = catalog.Add(name, category, ingredients, instructions, prepTime, difficulty);
                    Console.WriteLine($"\u001B[32m✅ Рецепт добавлен (ID: {id})\u001B[0m");
                    break;
                case "2": catalog.ListAll(); break;
                case "3":
                    Console.Write("Введите запрос (название или ингредиент): ");
                    string query = Console.ReadLine().Trim();
                    catalog.Search(query);
                    break;
                case "4":
                    Console.WriteLine("Фильтровать по: category, difficulty, prep_time");
                    Console.Write("Поле: ");
                    string field = Console.ReadLine().Trim().ToLower();
                    Console.Write("Значение: ");
                    string value = Console.ReadLine().Trim();
                    catalog.FilterBy(field, value);
                    break;
                case "5":
                    Console.WriteLine("Сортировать по: name, prep_time, rating");
                    Console.Write("Поле: ");
                    string sortField = Console.ReadLine().Trim().ToLower();
                    Console.Write("По убыванию? (y/n): ");
                    bool reverse = Console.ReadLine().Trim().ToLower() == "y";
                    catalog.SortBy(sortField, reverse);
                    break;
                case "6":
                    catalog.ListAll();
                    Console.Write("Введите ID для удаления: ");
                    int delId = int.Parse(Console.ReadLine().Trim());
                    if (catalog.Delete(delId)) {
                        Console.WriteLine("\u001B[32m✅ Рецепт удалён.\u001B[0m");
                    } else {
                        Console.WriteLine("\u001B[31m❌ Рецепт не найден.\u001B[0m");
                    }
                    break;
                case "7":
                    catalog.ListAll();
                    Console.Write("Введите ID для редактирования: ");
                    int editId = int.Parse(Console.ReadLine().Trim());
                    Console.Write("Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): ");
                    string editField = Console.ReadLine().Trim().ToLower();
                    Console.Write("Новое значение: ");
                    string editValue = Console.ReadLine().Trim();
                    if (catalog.Edit(editId, editField, editValue)) {
                        Console.WriteLine("\u001B[32m✅ Рецепт обновлён.\u001B[0m");
                    } else {
                        Console.WriteLine("\u001B[31m❌ Не удалось обновить.\u001B[0m");
                    }
                    break;
                case "8":
                    catalog.ListAll();
                    Console.Write("Введите ID для оценки: ");
                    int rateId = int.Parse(Console.ReadLine().Trim());
                    Console.Write("Оценка (1-5): ");
                    int rating = int.Parse(Console.ReadLine().Trim());
                    if (rating >= 1 && rating <= 5) {
                        if (catalog.Rate(rateId, rating)) {
                            Console.WriteLine("\u001B[32m✅ Оценка добавлена.\u001B[0m");
                        } else {
                            Console.WriteLine("\u001B[31m❌ Рецепт не найден.\u001B[0m");
                        }
                    } else {
                        Console.WriteLine("\u001B[31m❌ Оценка должна быть от 1 до 5.\u001B[0m");
                    }
                    break;
                case "9": catalog.Stats(); break;
                case "10":
                    Console.WriteLine("1. Экспорт в JSON");
                    Console.WriteLine("2. Экспорт в CSV");
                    Console.Write("Выберите формат: ");
                    string sub = Console.ReadLine().Trim();
                    if (sub == "1") catalog.ExportJSON();
                    else if (sub == "2") catalog.ExportCSV();
                    else Console.WriteLine("\u001B[31mНеверный выбор.\u001B[0m");
                    break;
                case "11": Console.WriteLine("До свидания!"); return;
                default: Console.WriteLine("\u001B[31mНеверный выбор.\u001B[0m"); break;
            }
        }
    }
}
