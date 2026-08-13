

### 1. `recipe_catalog.py` (Python)

```python
# recipe_catalog.py — Python версия

import json
import os
from datetime import datetime
from colorama import init, Fore, Style

init(autoreset=True)
DATA_FILE = "recipes.json"

class Recipe:
    def __init__(self, id, name, category, ingredients, instructions, prep_time, difficulty, rating=0):
        self.id = id
        self.name = name
        self.category = category
        self.ingredients = ingredients
        self.instructions = instructions
        self.prep_time = prep_time
        self.difficulty = difficulty
        self.rating = rating
        self.ratings_count = 0

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "category": self.category,
            "ingredients": self.ingredients,
            "instructions": self.instructions,
            "prep_time": self.prep_time,
            "difficulty": self.difficulty,
            "rating": self.rating,
            "ratings_count": self.ratings_count
        }

    @classmethod
    def from_dict(cls, data):
        recipe = cls(data["id"], data["name"], data["category"], data["ingredients"],
                     data["instructions"], data["prep_time"], data["difficulty"], data.get("rating", 0))
        recipe.ratings_count = data.get("ratings_count", 0)
        return recipe

    def add_rating(self, rating):
        total = self.rating * self.ratings_count + rating
        self.ratings_count += 1
        self.rating = total / self.ratings_count

class RecipeCatalog:
    def __init__(self):
        self.recipes = []
        self.load()

    def load(self):
        if os.path.exists(DATA_FILE):
            try:
                with open(DATA_FILE, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                    self.recipes = [Recipe.from_dict(r) for r in data]
            except:
                self.recipes = []

    def save(self):
        with open(DATA_FILE, 'w', encoding='utf-8') as f:
            json.dump([r.to_dict() for r in self.recipes], f, indent=2, ensure_ascii=False)

    def add(self, name, category, ingredients, instructions, prep_time, difficulty):
        id = len(self.recipes) + 1
        recipe = Recipe(id, name, category, ingredients, instructions, prep_time, difficulty)
        self.recipes.append(recipe)
        self.save()
        return id

    def list_all(self):
        if not self.recipes:
            print(Fore.YELLOW + "Каталог пуст.")
            return
        print(Fore.CYAN + f"{'ID':<4} {'Название':<25} {'Категория':<12} {'Время':<8} {'Сложность':<10} {'Рейтинг':<8}")
        print("-" * 75)
        for r in self.recipes:
            rating_str = f"{r.rating:.1f}★" if r.ratings_count > 0 else "—"
            diff_color = Fore.GREEN if r.difficulty == "лёгкая" else Fore.YELLOW if r.difficulty == "средняя" else Fore.RED
            print(f"{r.id:<4} {r.name[:25]:<25} {r.category[:12]:<12} {r.prep_time:<8} {diff_color}{r.difficulty[:10]:<10}{Style.RESET_ALL} {rating_str:<8}")

    def search(self, query):
        query = query.lower()
        results = [r for r in self.recipes if query in r.name.lower() or
                   any(query in ing.lower() for ing in r.ingredients)]
        if not results:
            print(Fore.YELLOW + "Ничего не найдено.")
            return
        for r in results:
            print(f"{r.id}: {r.name} | {r.category} | {r.prep_time} мин | {r.difficulty} | {r.rating:.1f}★")

    def filter_by(self, field, value):
        results = []
        for r in self.recipes:
            if field == "category" and r.category.lower() == value.lower():
                results.append(r)
            elif field == "difficulty" and r.difficulty.lower() == value.lower():
                results.append(r)
            elif field == "prep_time" and r.prep_time <= int(value):
                results.append(r)
        if not results:
            print(Fore.YELLOW + "Нет рецептов, соответствующих фильтру.")
            return
        for r in results:
            print(f"{r.id}: {r.name} | {r.category} | {r.prep_time} мин | {r.difficulty}")

    def sort_by(self, field, reverse=False):
        if field == "name":
            self.recipes.sort(key=lambda r: r.name, reverse=reverse)
        elif field == "prep_time":
            self.recipes.sort(key=lambda r: r.prep_time, reverse=reverse)
        elif field == "rating":
            self.recipes.sort(key=lambda r: r.rating, reverse=reverse)
        else:
            print(Fore.RED + "Неверное поле для сортировки.")
            return
        self.list_all()

    def delete(self, id):
        for i, r in enumerate(self.recipes):
            if r.id == id:
                del self.recipes[i]
                self.save()
                return True
        return False

    def edit(self, id, field, value):
        for r in self.recipes:
            if r.id == id:
                if field == "name":
                    r.name = value
                elif field == "category":
                    r.category = value
                elif field == "ingredients":
                    r.ingredients = [i.strip() for i in value.split(',')]
                elif field == "instructions":
                    r.instructions = value
                elif field == "prep_time":
                    r.prep_time = int(value)
                elif field == "difficulty":
                    r.difficulty = value
                else:
                    return False
                self.save()
                return True
        return False

    def rate(self, id, rating):
        for r in self.recipes:
            if r.id == id:
                r.add_rating(rating)
                self.save()
                return True
        return False

    def stats(self):
        if not self.recipes:
            print("Нет данных.")
            return
        total = len(self.recipes)
        categories = {}
        difficulties = {}
        total_rating = 0
        rated_count = 0
        for r in self.recipes:
            categories[r.category] = categories.get(r.category, 0) + 1
            difficulties[r.difficulty] = difficulties.get(r.difficulty, 0) + 1
            if r.ratings_count > 0:
                total_rating += r.rating
                rated_count += 1
        avg_rating = total_rating / rated_count if rated_count > 0 else 0
        print(Fore.CYAN + "📊 Статистика:")
        print(f"  Всего рецептов: {total}")
        print(f"  Средний рейтинг: {avg_rating:.1f}★ (из {rated_count} оценённых)")
        print("  По категориям:")
        for c, count in sorted(categories.items(), key=lambda x: -x[1]):
            print(f"    {c}: {count}")
        print("  По сложности:")
        for d, count in difficulties.items():
            print(f"    {d}: {count}")

    def export_json(self, filename="recipes_export.json"):
        data = [r.to_dict() for r in self.recipes]
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print(Fore.GREEN + f"💾 Экспорт JSON: {filename}")

    def export_csv(self, filename="recipes_export.csv"):
        if not self.recipes:
            print(Fore.YELLOW + "Нет данных для экспорта.")
            return
        import csv
        with open(filename, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            writer.writerow(["ID", "Название", "Категория", "Ингредиенты", "Инструкция", "Время", "Сложность", "Рейтинг"])
            for r in self.recipes:
                writer.writerow([r.id, r.name, r.category, ", ".join(r.ingredients), r.instructions, r.prep_time, r.difficulty, f"{r.rating:.1f}"])
        print(Fore.GREEN + f"💾 Экспорт CSV: {filename}")

def main():
    catalog = RecipeCatalog()
    while True:
        print(Fore.CYAN + "\n🍰 Dessert Recipe Catalog (Python)")
        print("1. Добавить рецепт")
        print("2. Показать все рецепты")
        print("3. Поиск рецептов")
        print("4. Фильтрация")
        print("5. Сортировка")
        print("6. Удалить рецепт")
        print("7. Редактировать рецепт")
        print("8. Оценить рецепт")
        print("9. Статистика")
        print("10. Экспорт")
        print("11. Выход")
        choice = input("Выберите действие: ").strip()
        if choice == "1":
            name = input("Название: ")
            category = input("Категория (торт, пирожное, печенье, мороженое, другое): ").strip()
            ingredients = [i.strip() for i in input("Ингредиенты (через запятую): ").split(',')]
            instructions = input("Инструкция: ")
            prep_time = int(input("Время приготовления (мин): "))
            difficulty = input("Сложность (лёгкая/средняя/сложная): ").strip().lower()
            if difficulty not in ["лёгкая", "средняя", "сложная"]:
                difficulty = "средняя"
            id = catalog.add(name, category, ingredients, instructions, prep_time, difficulty)
            print(Fore.GREEN + f"✅ Рецепт добавлен (ID: {id})")
        elif choice == "2":
            catalog.list_all()
        elif choice == "3":
            query = input("Введите запрос (название или ингредиент): ")
            catalog.search(query)
        elif choice == "4":
            print("Фильтровать по: category, difficulty, prep_time")
            field = input("Поле: ").lower()
            value = input("Значение: ")
            if field == "prep_time":
                value = int(value)
            catalog.filter_by(field, value)
        elif choice == "5":
            print("Сортировать по: name, prep_time, rating")
            field = input("Поле: ").lower()
            rev = input("По убыванию? (y/n): ").lower() == 'y'
            catalog.sort_by(field, rev)
        elif choice == "6":
            catalog.list_all()
            id = int(input("Введите ID для удаления: "))
            if catalog.delete(id):
                print(Fore.GREEN + "✅ Рецепт удалён.")
            else:
                print(Fore.RED + "❌ Рецепт не найден.")
        elif choice == "7":
            catalog.list_all()
            id = int(input("Введите ID для редактирования: "))
            field = input("Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): ").lower()
            value = input("Новое значение: ")
            if catalog.edit(id, field, value):
                print(Fore.GREEN + "✅ Рецепт обновлён.")
            else:
                print(Fore.RED + "❌ Не удалось обновить.")
        elif choice == "8":
            catalog.list_all()
            id = int(input("Введите ID для оценки: "))
            rating = int(input("Оценка (1-5): "))
            if 1 <= rating <= 5:
                if catalog.rate(id, rating):
                    print(Fore.GREEN + "✅ Оценка добавлена.")
                else:
                    print(Fore.RED + "❌ Рецепт не найден.")
            else:
                print(Fore.RED + "❌ Оценка должна быть от 1 до 5.")
        elif choice == "9":
            catalog.stats()
        elif choice == "10":
            print("1. Экспорт в JSON")
            print("2. Экспорт в CSV")
            sub = input("Выберите формат: ")
            if sub == "1":
                catalog.export_json()
            elif sub == "2":
                catalog.export_csv()
            else:
                print(Fore.RED + "Неверный выбор.")
        elif choice == "11":
            print("До свидания!")
            break
        else:
            print(Fore.RED + "Неверный выбор.")

if __name__ == "__main__":
    main()
