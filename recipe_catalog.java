// recipe_catalog.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.util.*;

class Recipe {
    int id;
    String name;
    String category;
    List<String> ingredients;
    String instructions;
    int prepTime;
    String difficulty;
    double rating;
    int ratingsCount;

    Recipe(int id, String name, String category, List<String> ingredients, String instructions, int prepTime, String difficulty) {
        this.id = id;
        this.name = name;
        this.category = category;
        this.ingredients = ingredients;
        this.instructions = instructions;
        this.prepTime = prepTime;
        this.difficulty = difficulty;
        this.rating = 0;
        this.ratingsCount = 0;
    }

    Recipe(int id, String name, String category, List<String> ingredients, String instructions, int prepTime, String difficulty, double rating, int ratingsCount) {
        this.id = id;
        this.name = name;
        this.category = category;
        this.ingredients = ingredients;
        this.instructions = instructions;
        this.prepTime = prepTime;
        this.difficulty = difficulty;
        this.rating = rating;
        this.ratingsCount = ratingsCount;
    }

    String toJson() {
        return String.format("{\"id\":%d,\"name\":\"%s\",\"category\":\"%s\",\"ingredients\":%s,\"instructions\":\"%s\",\"prep_time\":%d,\"difficulty\":\"%s\",\"rating\":%.1f,\"ratings_count\":%d}",
                id, name, category, ingredientsToString(), instructions, prepTime, difficulty, rating, ratingsCount);
    }

    String ingredientsToString() {
        StringBuilder sb = new StringBuilder("[");
        for (int i = 0; i < ingredients.size(); i++) {
            sb.append("\"").append(ingredients.get(i)).append("\"");
            if (i < ingredients.size() - 1) sb.append(",");
        }
        sb.append("]");
        return sb.toString();
    }

    void addRating(int rating) {
        double total = this.rating * ratingsCount + rating;
        ratingsCount++;
        this.rating = total / ratingsCount;
    }
}

public class recipe_catalog {
    private static List<Recipe> recipes = new ArrayList<>();
    private static final String DATA_FILE = "recipes.json";
    private static Scanner scanner = new Scanner(System.in);

    public static void main(String[] args) {
        load();
        while (true) {
            System.out.println("\n\u001B[36m🍰 Dessert Recipe Catalog (Java)\u001B[0m");
            System.out.println("1. Добавить рецепт");
            System.out.println("2. Показать все рецепты");
            System.out.println("3. Поиск рецептов");
            System.out.println("4. Фильтрация");
            System.out.println("5. Сортировка");
            System.out.println("6. Удалить рецепт");
            System.out.println("7. Редактировать рецепт");
            System.out.println("8. Оценить рецепт");
            System.out.println("9. Статистика");
            System.out.println("10. Экспорт");
            System.out.println("11. Выход");
            System.out.print("Выберите действие: ");
            String choice = scanner.nextLine().trim();
            switch (choice) {
                case "1": addRecipe(); break;
                case "2": listAll(); break;
                case "3": searchRecipes(); break;
                case "4": filterRecipes(); break;
                case "5": sortRecipes(); break;
                case "6": deleteRecipe(); break;
                case "7": editRecipe(); break;
                case "8": rateRecipe(); break;
                case "9": showStats(); break;
                case "10": exportRecipes(); break;
                case "11": System.out.println("До свидания!"); return;
                default: System.out.println("\u001B[31mНеверный выбор.\u001B[0m");
            }
        }
    }

    private static void load() {
        try {
            String content = new String(Files.readAllBytes(Paths.get(DATA_FILE)));
            // Упрощённо: если файл есть, парсим, иначе пустой список
            recipes = new ArrayList<>();
        } catch (IOException e) {
            recipes = new ArrayList<>();
        }
    }

    private static void save() {
        try {
            StringBuilder sb = new StringBuilder("[");
            for (int i = 0; i < recipes.size(); i++) {
                sb.append(recipes.get(i).toJson());
                if (i < recipes.size() - 1) sb.append(",");
            }
            sb.append("]");
            Files.write(Paths.get(DATA_FILE), sb.toString().getBytes());
        } catch (IOException e) {
            System.out.println("Ошибка сохранения.");
        }
    }

    private static void addRecipe() {
        System.out.print("Название: ");
        String name = scanner.nextLine().trim();
        System.out.print("Категория (торт, пирожное, печенье, мороженое, другое): ");
        String category = scanner.nextLine().trim();
        System.out.print("Ингредиенты (через запятую): ");
        String[] ingArr = scanner.nextLine().split(",");
        List<String> ingredients = new ArrayList<>();
        for (String s : ingArr) ingredients.add(s.trim());
        System.out.print("Инструкция: ");
        String instructions = scanner.nextLine().trim();
        System.out.print("Время приготовления (мин): ");
        int prepTime = Integer.parseInt(scanner.nextLine().trim());
        System.out.print("Сложность (лёгкая/средняя/сложная): ");
        String difficulty = scanner.nextLine().trim().toLowerCase();
        if (!difficulty.equals("лёгкая") && !difficulty.equals("средняя") && !difficulty.equals("сложная")) {
            difficulty = "средняя";
        }
        int id = recipes.size() + 1;
        recipes.add(new Recipe(id, name, category, ingredients, instructions, prepTime, difficulty));
        save();
        System.out.println("\u001B[32m✅ Рецепт добавлен (ID: " + id + ")\u001B[0m");
    }

    private static void listAll() {
        if (recipes.isEmpty()) {
            System.out.println("\u001B[33mКаталог пуст.\u001B[0m");
            return;
        }
        System.out.printf("\u001B[36m%-4s %-25s %-12s %-8s %-10s %-8s\u001B[0m\n", "ID", "Название", "Категория", "Время", "Сложность", "Рейтинг");
        System.out.println("-".repeat(75));
        for (Recipe r : recipes) {
            String ratingStr = r.ratingsCount > 0 ? String.format("%.1f★", r.rating) : "—";
            String diffColor = r.difficulty.equals("лёгкая") ? "\u001B[32m" : r.difficulty.equals("средняя") ? "\u001B[33m" : "\u001B[31m";
            System.out.printf("%-4d %-25s %-12s %-8d %s%-10s\u001B[0m %-8s\n", r.id, r.name, r.category, r.prepTime, diffColor, r.difficulty, ratingStr);
        }
    }

    private static void searchRecipes() {
        System.out.print("Введите запрос (название или ингредиент): ");
        String query = scanner.nextLine().trim().toLowerCase();
        List<Recipe> results = new ArrayList<>();
        for (Recipe r : recipes) {
            if (r.name.toLowerCase().contains(query)) {
                results.add(r);
            } else {
                for (String ing : r.ingredients) {
                    if (ing.toLowerCase().contains(query)) {
                        results.add(r);
                        break;
                    }
                }
            }
        }
        if (results.isEmpty()) {
            System.out.println("\u001B[33mНичего не найдено.\u001B[0m");
            return;
        }
        for (Recipe r : results) {
            System.out.printf("%d: %s | %s | %d мин | %s | %.1f★\n", r.id, r.name, r.category, r.prepTime, r.difficulty, r.rating);
        }
    }

    private static void filterRecipes() {
        System.out.println("Фильтровать по: category, difficulty, prep_time");
        System.out.print("Поле: ");
        String field = scanner.nextLine().trim().toLowerCase();
        System.out.print("Значение: ");
        String value = scanner.nextLine().trim();
        List<Recipe> results = new ArrayList<>();
        for (Recipe r : recipes) {
            if (field.equals("category") && r.category.equalsIgnoreCase(value)) {
                results.add(r);
            } else if (field.equals("difficulty") && r.difficulty.equalsIgnoreCase(value)) {
                results.add(r);
            } else if (field.equals("prep_time")) {
                try {
                    int v = Integer.parseInt(value);
                    if (r.prepTime <= v) results.add(r);
                } catch (NumberFormatException e) {}
            }
        }
        if (results.isEmpty()) {
            System.out.println("\u001B[33mНет рецептов, соответствующих фильтру.\u001B[0m");
            return;
        }
        for (Recipe r : results) {
            System.out.printf("%d: %s | %s | %d мин | %s\n", r.id, r.name, r.category, r.prepTime, r.difficulty);
        }
    }

    private static void sortRecipes() {
        System.out.println("Сортировать по: name, prep_time, rating");
        System.out.print("Поле: ");
        String field = scanner.nextLine().trim().toLowerCase();
        System.out.print("По убыванию? (y/n): ");
        boolean reverse = scanner.nextLine().trim().equalsIgnoreCase("y");
        Comparator<Recipe> comp = null;
        switch (field) {
            case "name": comp = Comparator.comparing(r -> r.name); break;
            case "prep_time": comp = Comparator.comparingInt(r -> r.prepTime); break;
            case "rating": comp = Comparator.comparingDouble(r -> r.rating); break;
            default: System.out.println("\u001B[31mНеверное поле для сортировки.\u001B[0m"); return;
        }
        if (reverse) comp = comp.reversed();
        recipes.sort(comp);
        listAll();
    }

    private static void deleteRecipe() {
        listAll();
        System.out.print("Введите ID для удаления: ");
        int id = Integer.parseInt(scanner.nextLine().trim());
        Iterator<Recipe> it = recipes.iterator();
        while (it.hasNext()) {
            if (it.next().id == id) {
                it.remove();
                save();
                System.out.println("\u001B[32m✅ Рецепт удалён.\u001B[0m");
                return;
            }
        }
        System.out.println("\u001B[31m❌ Рецепт не найден.\u001B[0m");
    }

    private static void editRecipe() {
        listAll();
        System.out.print("Введите ID для редактирования: ");
        int id = Integer.parseInt(scanner.nextLine().trim());
        Recipe target = null;
        for (Recipe r : recipes) {
            if (r.id == id) {
                target = r;
                break;
            }
        }
        if (target == null) {
            System.out.println("\u001B[31m❌ Рецепт не найден.\u001B[0m");
            return;
        }
        System.out.print("Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): ");
        String field = scanner.nextLine().trim().toLowerCase();
        System.out.print("Новое значение: ");
        String value = scanner.nextLine().trim();
        boolean ok = true;
        switch (field) {
            case "name": target.name = value; break;
            case "category": target.category = value; break;
            case "ingredients":
                String[] parts = value.split(",");
                target.ingredients = new ArrayList<>();
                for (String s : parts) target.ingredients.add(s.trim());
                break;
            case "instructions": target.instructions = value; break;
            case "prep_time": target.prepTime = Integer.parseInt(value); break;
            case "difficulty": target.difficulty = value; break;
            default: ok = false;
        }
        if (ok) {
            save();
            System.out.println("\u001B[32m✅ Рецепт обновлён.\u001B[0m");
        } else {
            System.out.println("\u001B[31m❌ Не удалось обновить.\u001B[0m");
        }
    }

    private static void rateRecipe() {
        listAll();
        System.out.print("Введите ID для оценки: ");
        int id = Integer.parseInt(scanner.nextLine().trim());
        System.out.print("Оценка (1-5): ");
        int rating = Integer.parseInt(scanner.nextLine().trim());
        if (rating >= 1 && rating <= 5) {
            for (Recipe r : recipes) {
                if (r.id == id) {
                    r.addRating(rating);
                    save();
                    System.out.println("\u001B[32m✅ Оценка добавлена.\u001B[0m");
                    return;
                }
            }
            System.out.println("\u001B[31m❌ Рецепт не найден.\u001B[0m");
        } else {
            System.out.println("\u001B[31m❌ Оценка должна быть от 1 до 5.\u001B[0m");
        }
    }

    private static void showStats() {
        if (recipes.isEmpty()) {
            System.out.println("Нет данных.");
            return;
        }
        Map<String, Integer> categories = new HashMap<>();
        Map<String, Integer> difficulties = new HashMap<>();
        double totalRating = 0;
        int ratedCount = 0;
        for (Recipe r : recipes) {
            categories.put(r.category, categories.getOrDefault(r.category, 0) + 1);
            difficulties.put(r.difficulty, difficulties.getOrDefault(r.difficulty, 0) + 1);
            if (r.ratingsCount > 0) {
                totalRating += r.rating;
                ratedCount++;
            }
        }
        double avgRating = ratedCount > 0 ? totalRating / ratedCount : 0;
        System.out.println("\u001B[36m📊 Статистика:\u001B[0m");
        System.out.printf("  Всего рецептов: %d\n", recipes.size());
        System.out.printf("  Средний рейтинг: %.1f★ (из %d оценённых)\n", avgRating, ratedCount);
        System.out.println("  По категориям:");
        for (Map.Entry<String, Integer> e : categories.entrySet()) {
            System.out.printf("    %s: %d\n", e.getKey(), e.getValue());
        }
        System.out.println("  По сложности:");
        for (Map.Entry<String, Integer> e : difficulties.entrySet()) {
            System.out.printf("    %s: %d\n", e.getKey(), e.getValue());
        }
    }

    private static void exportRecipes() {
        System.out.println("1. Экспорт в JSON");
        System.out.println("2. Экспорт в CSV");
        System.out.print("Выберите формат: ");
        String sub = scanner.nextLine().trim();
        if (sub.equals("1")) {
            try {
                StringBuilder sb = new StringBuilder("[");
                for (int i = 0; i < recipes.size(); i++) {
                    sb.append(recipes.get(i).toJson());
                    if (i < recipes.size() - 1) sb.append(",");
                }
                sb.append("]");
                Files.write(Paths.get("recipes_export.json"), sb.toString().getBytes());
                System.out.println("\u001B[32m💾 Экспорт JSON: recipes_export.json\u001B[0m");
            } catch (IOException e) {
                System.out.println("Ошибка экспорта.");
            }
        } else if (sub.equals("2")) {
            try (FileWriter fw = new FileWriter("recipes_export.csv")) {
                fw.write("ID,Название,Категория,Ингредиенты,Инструкция,Время,Сложность,Рейтинг\n");
                for (Recipe r : recipes) {
                    fw.write(r.id + "," + r.name + "," + r.category + ",\"" + String.join(", ", r.ingredients) + "\",\"" + r.instructions + "\"," + r.prepTime + "," + r.difficulty + "," + String.format("%.1f", r.rating) + "\n");
                }
                System.out.println("\u001B[32m💾 Экспорт CSV: recipes_export.csv\u001B[0m");
            } catch (IOException e) {
                System.out.println("Ошибка экспорта.");
            }
        } else {
            System.out.println("\u001B[31mНеверный выбор.\u001B[0m");
        }
    }
}
