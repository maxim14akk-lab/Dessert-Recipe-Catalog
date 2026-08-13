// recipe_catalog.rs — Rust версия

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::io::{self, Write};

#[derive(Serialize, Deserialize, Clone)]
struct Recipe {
    id: usize,
    name: String,
    category: String,
    ingredients: Vec<String>,
    instructions: String,
    prep_time: u32,
    difficulty: String,
    rating: f64,
    ratings_count: usize,
}

impl Recipe {
    fn add_rating(&mut self, rating: u32) {
        let total = self.rating * self.ratings_count as f64 + rating as f64;
        self.ratings_count += 1;
        self.rating = total / self.ratings_count as f64;
    }
}

struct Catalog {
    recipes: Vec<Recipe>,
    file: String,
}

impl Catalog {
    fn new(file: &str) -> Self {
        let mut c = Catalog { recipes: Vec::new(), file: file.to_string() };
        c.load();
        c
    }

    fn load(&mut self) {
        if let Ok(data) = fs::read_to_string(&self.file) {
            if let Ok(recipes) = serde_json::from_str(&data) {
                self.recipes = recipes;
                return;
            }
        }
        self.recipes = Vec::new();
    }

    fn save(&self) {
        let data = serde_json::to_string_pretty(&self.recipes).unwrap();
        fs::write(&self.file, data).unwrap();
    }

    fn add(&mut self, name: String, category: String, ingredients: Vec<String>, instructions: String, prep_time: u32, difficulty: String) -> usize {
        let id = self.recipes.len() + 1;
        self.recipes.push(Recipe {
            id,
            name,
            category,
            ingredients,
            instructions,
            prep_time,
            difficulty,
            rating: 0.0,
            ratings_count: 0,
        });
        self.save();
        id
    }

    fn list_all(&self) {
        if self.recipes.is_empty() {
            println!("\x1b[33mКаталог пуст.\x1b[0m");
            return;
        }
        println!("\x1b[36m{:<4} {:<25} {:<12} {:<8} {:<10} {:<8}\x1b[0m", "ID", "Название", "Категория", "Время", "Сложность", "Рейтинг");
        println!("{}", "-".repeat(75));
        for r in &self.recipes {
            let rating_str = if r.ratings_count > 0 { format!("{:.1}★", r.rating) } else { "—".to_string() };
            let diff_color = if r.difficulty == "лёгкая" { "\x1b[32m" } else if r.difficulty == "средняя" { "\x1b[33m" } else { "\x1b[31m" };
            println!("{:<4} {:<25} {:<12} {:<8} {}{:<10}\x1b[0m {:<8}", r.id, r.name, r.category, r.prep_time, diff_color, r.difficulty, rating_str);
        }
    }

    fn search(&self, query: &str) {
        let query = query.to_lowercase();
        let results: Vec<&Recipe> = self.recipes.iter()
            .filter(|r| r.name.to_lowercase().contains(&query) ||
                   r.ingredients.iter().any(|i| i.to_lowercase().contains(&query)))
            .collect();
        if results.is_empty() {
            println!("\x1b[33mНичего не найдено.\x1b[0m");
            return;
        }
        for r in results {
            println!("{}: {} | {} | {} мин | {} | {:.1}★", r.id, r.name, r.category, r.prep_time, r.difficulty, r.rating);
        }
    }

    fn filter_by(&self, field: &str, value: &str) {
        let results: Vec<&Recipe> = match field {
            "category" => self.recipes.iter().filter(|r| r.category.eq_ignore_ascii_case(value)).collect(),
            "difficulty" => self.recipes.iter().filter(|r| r.difficulty.eq_ignore_ascii_case(value)).collect(),
            "prep_time" => {
                if let Ok(v) = value.parse::<u32>() {
                    self.recipes.iter().filter(|r| r.prep_time <= v).collect()
                } else {
                    Vec::new()
                }
            }
            _ => Vec::new(),
        };
        if results.is_empty() {
            println!("\x1b[33mНет рецептов, соответствующих фильтру.\x1b[0m");
            return;
        }
        for r in results {
            println!("{}: {} | {} | {} мин | {}", r.id, r.name, r.category, r.prep_time, r.difficulty);
        }
    }

    fn sort_by(&mut self, field: &str, reverse: bool) {
        match field {
            "name" => {
                if reverse {
                    self.recipes.sort_by(|a, b| b.name.cmp(&a.name));
                } else {
                    self.recipes.sort_by(|a, b| a.name.cmp(&b.name));
                }
            }
            "prep_time" => {
                if reverse {
                    self.recipes.sort_by(|a, b| b.prep_time.cmp(&a.prep_time));
                } else {
                    self.recipes.sort_by(|a, b| a.prep_time.cmp(&b.prep_time));
                }
            }
            "rating" => {
                if reverse {
                    self.recipes.sort_by(|a, b| b.rating.partial_cmp(&a.rating).unwrap());
                } else {
                    self.recipes.sort_by(|a, b| a.rating.partial_cmp(&b.rating).unwrap());
                }
            }
            _ => {
                println!("\x1b[31mНеверное поле для сортировки.\x1b[0m");
                return;
            }
        }
        self.list_all();
    }

    fn delete(&mut self, id: usize) -> bool {
        let pos = self.recipes.iter().position(|r| r.id == id);
        if let Some(idx) = pos {
            self.recipes.remove(idx);
            self.save();
            true
        } else {
            false
        }
    }

    fn edit(&mut self, id: usize, field: &str, value: &str) -> bool {
        for r in &mut self.recipes {
            if r.id == id {
                match field {
                    "name" => r.name = value.to_string(),
                    "category" => r.category = value.to_string(),
                    "ingredients" => r.ingredients = value.split(',').map(|s| s.trim().to_string()).collect(),
                    "instructions" => r.instructions = value.to_string(),
                    "prep_time" => {
                        if let Ok(v) = value.parse() {
                            r.prep_time = v;
                        } else {
                            return false;
                        }
                    }
                    "difficulty" => r.difficulty = value.to_string(),
                    _ => return false,
                }
                self.save();
                return true;
            }
        }
        false
    }

    fn rate(&mut self, id: usize, rating: u32) -> bool {
        for r in &mut self.recipes {
            if r.id == id {
                r.add_rating(rating);
                self.save();
                return true;
            }
        }
        false
    }

    fn stats(&self) {
        if self.recipes.is_empty() {
            println!("Нет данных.");
            return;
        }
        let mut categories = HashMap::new();
        let mut difficulties = HashMap::new();
        let mut total_rating = 0.0;
        let mut rated_count = 0;
        for r in &self.recipes {
            *categories.entry(&r.category).or_insert(0) += 1;
            *difficulties.entry(&r.difficulty).or_insert(0) += 1;
            if r.ratings_count > 0 {
                total_rating += r.rating;
                rated_count += 1;
            }
        }
        let avg_rating = if rated_count > 0 { total_rating / rated_count as f64 } else { 0.0 };
        println!("\x1b[36m📊 Статистика:\x1b[0m");
        println!("  Всего рецептов: {}", self.recipes.len());
        println!("  Средний рейтинг: {:.1}★ (из {} оценённых)", avg_rating, rated_count);
        println!("  По категориям:");
        let mut categories_vec: Vec<_> = categories.iter().collect();
        categories_vec.sort_by(|a, b| b.1.cmp(a.1));
        for (c, count) in categories_vec {
            println!("    {}: {}", c, count);
        }
        println!("  По сложности:");
        for (d, count) in difficulties {
            println!("    {}: {}", d, count);
        }
    }

    fn export_json(&self, filename: &str) {
        let data = serde_json::to_string_pretty(&self.recipes).unwrap();
        fs::write(filename, data).unwrap();
        println!("\x1b[32m💾 Экспорт JSON: {}\x1b[0m", filename);
    }

    fn export_csv(&self, filename: &str) {
        if self.recipes.is_empty() {
            println!("\x1b[33mНет данных для экспорта.\x1b[0m");
            return;
        }
        let mut csv = String::from("ID,Название,Категория,Ингредиенты,Инструкция,Время,Сложность,Рейтинг\n");
        for r in &self.recipes {
            csv.push_str(&format!("{},{},{},\"{}\",\"{}\",{},{},{:.1}\n",
                r.id, r.name, r.category, r.ingredients.join(", "), r.instructions, r.prep_time, r.difficulty, r.rating));
        }
        fs::write(filename, csv).unwrap();
        println!("\x1b[32m💾 Экспорт CSV: {}\x1b[0m", filename);
    }
}

fn main() {
    let mut catalog = Catalog::new("recipes.json");
    loop {
        println!("\n\x1b[36m🍰 Dessert Recipe Catalog (Rust)\x1b[0m");
        println!("1. Добавить рецепт");
        println!("2. Показать все рецепты");
        println!("3. Поиск рецептов");
        println!("4. Фильтрация");
        println!("5. Сортировка");
        println!("6. Удалить рецепт");
        println!("7. Редактировать рецепт");
        println!("8. Оценить рецепт");
        println!("9. Статистика");
        println!("10. Экспорт");
        println!("11. Выход");
        print!("Выберите действие: ");
        io::stdout().flush().unwrap();
        let mut choice = String::new();
        io::stdin().read_line(&mut choice).unwrap();
        match choice.trim() {
            "1" => {
                print!("Название: ");
                io::stdout().flush().unwrap();
                let mut name = String::new();
                io::stdin().read_line(&mut name).unwrap();
                let name = name.trim().to_string();
                print!("Категория (торт, пирожное, печенье, мороженое, другое): ");
                io::stdout().flush().unwrap();
                let mut category = String::new();
                io::stdin().read_line(&mut category).unwrap();
                let category = category.trim().to_string();
                print!("Ингредиенты (через запятую): ");
                io::stdout().flush().unwrap();
                let mut ing_str = String::new();
                io::stdin().read_line(&mut ing_str).unwrap();
                let ingredients: Vec<String> = ing_str.split(',').map(|s| s.trim().to_string()).collect();
                print!("Инструкция: ");
                io::stdout().flush().unwrap();
                let mut instructions = String::new();
                io::stdin().read_line(&mut instructions).unwrap();
                let instructions = instructions.trim().to_string();
                print!("Время приготовления (мин): ");
                io::stdout().flush().unwrap();
                let mut time_str = String::new();
                io::stdin().read_line(&mut time_str).unwrap();
                let prep_time: u32 = time_str.trim().parse().unwrap();
                print!("Сложность (лёгкая/средняя/сложная): ");
                io::stdout().flush().unwrap();
                let mut difficulty = String::new();
                io::stdin().read_line(&mut difficulty).unwrap();
                let difficulty = difficulty.trim().to_lowercase();
                let difficulty = if difficulty == "лёгкая" || difficulty == "средняя" || difficulty == "сложная" { difficulty } else { "средняя".to_string() };
                let id = catalog.add(name, category, ingredients, instructions, prep_time, difficulty);
                println!("\x1b[32m✅ Рецепт добавлен (ID: {})\x1b[0m", id);
            }
            "2" => catalog.list_all(),
            "3" => {
                print!("Введите запрос (название или ингредиент): ");
                io::stdout().flush().unwrap();
                let mut query = String::new();
                io::stdin().read_line(&mut query).unwrap();
                catalog.search(query.trim());
            }
            "4" => {
                println!("Фильтровать по: category, difficulty, prep_time");
                print!("Поле: ");
                io::stdout().flush().unwrap();
                let mut field = String::new();
                io::stdin().read_line(&mut field).unwrap();
                let field = field.trim().to_lowercase();
                print!("Значение: ");
                io::stdout().flush().unwrap();
                let mut value = String::new();
                io::stdin().read_line(&mut value).unwrap();
                catalog.filter_by(&field, value.trim());
            }
            "5" => {
                println!("Сортировать по: name, prep_time, rating");
                print!("Поле: ");
                io::stdout().flush().unwrap();
                let mut field = String::new();
                io::stdin().read_line(&mut field).unwrap();
                let field = field.trim().to_lowercase();
                print!("По убыванию? (y/n): ");
                io::stdout().flush().unwrap();
                let mut rev_str = String::new();
                io::stdin().read_line(&mut rev_str).unwrap();
                let reverse = rev_str.trim().to_lowercase() == "y";
                catalog.sort_by(&field, reverse);
            }
            "6" => {
                catalog.list_all();
                print!("Введите ID для удаления: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                let id: usize = id_str.trim().parse().unwrap();
                if catalog.delete(id) {
                    println!("\x1b[32m✅ Рецепт удалён.\x1b[0m");
                } else {
                    println!("\x1b[31m❌ Рецепт не найден.\x1b[0m");
                }
            }
            "7" => {
                catalog.list_all();
                print!("Введите ID для редактирования: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                let id: usize = id_str.trim().parse().unwrap();
                print!("Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): ");
                io::stdout().flush().unwrap();
                let mut field = String::new();
                io::stdin().read_line(&mut field).unwrap();
                let field = field.trim().to_lowercase();
                print!("Новое значение: ");
                io::stdout().flush().unwrap();
                let mut value = String::new();
                io::stdin().read_line(&mut value).unwrap();
                if catalog.edit(id, &field, value.trim()) {
                    println!("\x1b[32m✅ Рецепт обновлён.\x1b[0m");
                } else {
                    println!("\x1b[31m❌ Не удалось обновить.\x1b[0m");
                }
            }
            "8" => {
                catalog.list_all();
                print!("Введите ID для оценки: ");
                io::stdout().flush().unwrap();
                let mut id_str = String::new();
                io::stdin().read_line(&mut id_str).unwrap();
                let id: usize = id_str.trim().parse().unwrap();
                print!("Оценка (1-5): ");
                io::stdout().flush().unwrap();
                let mut rating_str = String::new();
                io::stdin().read_line(&mut rating_str).unwrap();
                let rating: u32 = rating_str.trim().parse().unwrap();
                if (1..=5).contains(&rating) {
                    if catalog.rate(id, rating) {
                        println!("\x1b[32m✅ Оценка добавлена.\x1b[0m");
                    } else {
                        println!("\x1b[31m❌ Рецепт не найден.\x1b[0m");
                    }
                } else {
                    println!("\x1b[31m❌ Оценка должна быть от 1 до 5.\x1b[0m");
                }
            }
            "9" => catalog.stats(),
            "10" => {
                println!("1. Экспорт в JSON");
                println!("2. Экспорт в CSV");
                print!("Выберите формат: ");
                io::stdout().flush().unwrap();
                let mut sub = String::new();
                io::stdin().read_line(&mut sub).unwrap();
                let sub = sub.trim();
                if sub == "1" {
                    catalog.export_json("recipes_export.json");
                } else if sub == "2" {
                    catalog.export_csv("recipes_export.csv");
                } else {
                    println!("\x1b[31mНеверный выбор.\x1b[0m");
                }
            }
            "11" => {
                println!("До свидания!");
                break;
            }
            _ => println!("\x1b[31mНеверный выбор.\x1b[0m"),
        }
    }
}
