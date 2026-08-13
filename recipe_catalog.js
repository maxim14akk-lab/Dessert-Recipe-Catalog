// recipe_catalog.js — JavaScript версия

const fs = require('fs');
const readline = require('readline');

const DATA_FILE = 'recipes.json';

class Recipe {
    constructor(id, name, category, ingredients, instructions, prepTime, difficulty, rating = 0, ratingsCount = 0) {
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
}

class RecipeCatalog {
    constructor() {
        this.recipes = [];
        this.load();
    }

    load() {
        if (fs.existsSync(DATA_FILE)) {
            try {
                const data = JSON.parse(fs.readFileSync(DATA_FILE, 'utf8'));
                this.recipes = data.map(r => new Recipe(r.id, r.name, r.category, r.ingredients, r.instructions, r.prepTime, r.difficulty, r.rating, r.ratingsCount));
            } catch {
                this.recipes = [];
            }
        }
    }

    save() {
        fs.writeFileSync(DATA_FILE, JSON.stringify(this.recipes, null, 2));
    }

    add(name, category, ingredients, instructions, prepTime, difficulty) {
        const id = this.recipes.length + 1;
        const recipe = new Recipe(id, name, category, ingredients, instructions, prepTime, difficulty);
        this.recipes.push(recipe);
        this.save();
        return id;
    }

    listAll() {
        if (this.recipes.length === 0) {
            console.log('\x1b[33mКаталог пуст.\x1b[0m');
            return;
        }
        console.log('\x1b[36m' + 'ID'.padEnd(4) + 'Название'.padEnd(25) + 'Категория'.padEnd(12) + 'Время'.padEnd(8) + 'Сложность'.padEnd(10) + 'Рейтинг'.padEnd(8) + '\x1b[0m');
        console.log('-'.repeat(75));
        for (const r of this.recipes) {
            const ratingStr = r.ratingsCount > 0 ? `${r.rating.toFixed(1)}★` : '—';
            const diffColor = r.difficulty === 'лёгкая' ? '\x1b[32m' : r.difficulty === 'средняя' ? '\x1b[33m' : '\x1b[31m';
            console.log(`${String(r.id).padEnd(4)} ${r.name.padEnd(25)} ${r.category.padEnd(12)} ${String(r.prepTime).padEnd(8)} ${diffColor}${r.difficulty.padEnd(10)}\x1b[0m ${ratingStr.padEnd(8)}`);
        }
    }

    search(query) {
        query = query.toLowerCase();
        const results = this.recipes.filter(r =>
            r.name.toLowerCase().includes(query) ||
            r.ingredients.some(i => i.toLowerCase().includes(query))
        );
        if (results.length === 0) {
            console.log('\x1b[33mНичего не найдено.\x1b[0m');
            return;
        }
        for (const r of results) {
            console.log(`${r.id}: ${r.name} | ${r.category} | ${r.prepTime} мин | ${r.difficulty} | ${r.rating.toFixed(1)}★`);
        }
    }

    filterBy(field, value) {
        let results = [];
        switch (field) {
            case 'category':
                results = this.recipes.filter(r => r.category.toLowerCase() === value.toLowerCase());
                break;
            case 'difficulty':
                results = this.recipes.filter(r => r.difficulty.toLowerCase() === value.toLowerCase());
                break;
            case 'prep_time':
                results = this.recipes.filter(r => r.prepTime <= parseInt(value));
                break;
            default:
                results = [];
        }
        if (results.length === 0) {
            console.log('\x1b[33mНет рецептов, соответствующих фильтру.\x1b[0m');
            return;
        }
        for (const r of results) {
            console.log(`${r.id}: ${r.name} | ${r.category} | ${r.prepTime} мин | ${r.difficulty}`);
        }
    }

    sortBy(field, reverse) {
        switch (field) {
            case 'name':
                this.recipes.sort((a, b) => reverse ? b.name.localeCompare(a.name) : a.name.localeCompare(b.name));
                break;
            case 'prep_time':
                this.recipes.sort((a, b) => reverse ? b.prepTime - a.prepTime : a.prepTime - b.prepTime);
                break;
            case 'rating':
                this.recipes.sort((a, b) => reverse ? b.rating - a.rating : a.rating - b.rating);
                break;
            default:
                console.log('\x1b[31mНеверное поле для сортировки.\x1b[0m');
                return;
        }
        this.listAll();
    }

    delete(id) {
        const index = this.recipes.findIndex(r => r.id === id);
        if (index !== -1) {
            this.recipes.splice(index, 1);
            this.save();
            return true;
        }
        return false;
    }

    edit(id, field, value) {
        const recipe = this.recipes.find(r => r.id === id);
        if (!recipe) return false;
        switch (field) {
            case 'name': recipe.name = value; break;
            case 'category': recipe.category = value; break;
            case 'ingredients': recipe.ingredients = value.split(',').map(s => s.trim()); break;
            case 'instructions': recipe.instructions = value; break;
            case 'prep_time': recipe.prepTime = parseInt(value); break;
            case 'difficulty': recipe.difficulty = value; break;
            default: return false;
        }
        this.save();
        return true;
    }

    rate(id, rating) {
        const recipe = this.recipes.find(r => r.id === id);
        if (!recipe) return false;
        const total = recipe.rating * recipe.ratingsCount + rating;
        recipe.ratingsCount++;
        recipe.rating = total / recipe.ratingsCount;
        this.save();
        return true;
    }

    stats() {
        if (this.recipes.length === 0) {
            console.log('Нет данных.');
            return;
        }
        const total = this.recipes.length;
        const categories = {};
        const difficulties = {};
        let totalRating = 0;
        let ratedCount = 0;
        for (const r of this.recipes) {
            categories[r.category] = (categories[r.category] || 0) + 1;
            difficulties[r.difficulty] = (difficulties[r.difficulty] || 0) + 1;
            if (r.ratingsCount > 0) {
                totalRating += r.rating;
                ratedCount++;
            }
        }
        const avgRating = ratedCount > 0 ? totalRating / ratedCount : 0;
        console.log('\x1b[36m📊 Статистика:\x1b[0m');
        console.log(`  Всего рецептов: ${total}`);
        console.log(`  Средний рейтинг: ${avgRating.toFixed(1)}★ (из ${ratedCount} оценённых)`);
        console.log('  По категориям:');
        for (const [c, count] of Object.entries(categories).sort((a, b) => b[1] - a[1])) {
            console.log(`    ${c}: ${count}`);
        }
        console.log('  По сложности:');
        for (const [d, count] of Object.entries(difficulties)) {
            console.log(`    ${d}: ${count}`);
        }
    }

    exportJSON(filename = 'recipes_export.json') {
        fs.writeFileSync(filename, JSON.stringify(this.recipes, null, 2));
        console.log(`\x1b[32m💾 Экспорт JSON: ${filename}\x1b[0m`);
    }

    exportCSV(filename = 'recipes_export.csv') {
        if (this.recipes.length === 0) {
            console.log('\x1b[33mНет данных для экспорта.\x1b[0m');
            return;
        }
        let csv = 'ID,Название,Категория,Ингредиенты,Инструкция,Время,Сложность,Рейтинг\n';
        for (const r of this.recipes) {
            csv += `${r.id},${r.name},${r.category},"${r.ingredients.join(', ')}","${r.instructions}",${r.prepTime},${r.difficulty},${r.rating.toFixed(1)}\n`;
        }
        fs.writeFileSync(filename, csv);
        console.log(`\x1b[32m💾 Экспорт CSV: ${filename}\x1b[0m`);
    }
}

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const catalog = new RecipeCatalog();

function ask(question) {
    return new Promise(resolve => rl.question(question, resolve));
}

async function main() {
    while (true) {
        console.log('\x1b[36m\n🍰 Dessert Recipe Catalog (JavaScript)\x1b[0m');
        console.log('1. Добавить рецепт');
        console.log('2. Показать все рецепты');
        console.log('3. Поиск рецептов');
        console.log('4. Фильтрация');
        console.log('5. Сортировка');
        console.log('6. Удалить рецепт');
        console.log('7. Редактировать рецепт');
        console.log('8. Оценить рецепт');
        console.log('9. Статистика');
        console.log('10. Экспорт');
        console.log('11. Выход');
        const choice = await ask('Выберите действие: ');
        switch (choice.trim()) {
            case '1': {
                const name = await ask('Название: ');
                const category = await ask('Категория (торт, пирожное, печенье, мороженое, другое): ');
                const ingredientsStr = await ask('Ингредиенты (через запятую): ');
                const ingredients = ingredientsStr.split(',').map(s => s.trim());
                const instructions = await ask('Инструкция: ');
                const prepTime = parseInt(await ask('Время приготовления (мин): '));
                let difficulty = await ask('Сложность (лёгкая/средняя/сложная): ');
                difficulty = difficulty.trim().toLowerCase();
                if (!['лёгкая', 'средняя', 'сложная'].includes(difficulty)) difficulty = 'средняя';
                const id = catalog.add(name, category, ingredients, instructions, prepTime, difficulty);
                console.log(`\x1b[32m✅ Рецепт добавлен (ID: ${id})\x1b[0m`);
                break;
            }
            case '2': catalog.listAll(); break;
            case '3': {
                const query = await ask('Введите запрос (название или ингредиент): ');
                catalog.search(query);
                break;
            }
            case '4': {
                console.log('Фильтровать по: category, difficulty, prep_time');
                const field = await ask('Поле: ');
                const value = await ask('Значение: ');
                catalog.filterBy(field.trim().toLowerCase(), value.trim());
                break;
            }
            case '5': {
                console.log('Сортировать по: name, prep_time, rating');
                const field = await ask('Поле: ');
                const rev = (await ask('По убыванию? (y/n): ')).toLowerCase() === 'y';
                catalog.sortBy(field.trim().toLowerCase(), rev);
                break;
            }
            case '6': {
                catalog.listAll();
                const id = parseInt(await ask('Введите ID для удаления: '));
                if (catalog.delete(id)) {
                    console.log('\x1b[32m✅ Рецепт удалён.\x1b[0m');
                } else {
                    console.log('\x1b[31m❌ Рецепт не найден.\x1b[0m');
                }
                break;
            }
            case '7': {
                catalog.listAll();
                const id = parseInt(await ask('Введите ID для редактирования: '));
                const field = await ask('Какое поле редактировать (name, category, ingredients, instructions, prep_time, difficulty): ');
                const value = await ask('Новое значение: ');
                if (catalog.edit(id, field.trim().toLowerCase(), value.trim())) {
                    console.log('\x1b[32m✅ Рецепт обновлён.\x1b[0m');
                } else {
                    console.log('\x1b[31m❌ Не удалось обновить.\x1b[0m');
                }
                break;
            }
            case '8': {
                catalog.listAll();
                const id = parseInt(await ask('Введите ID для оценки: '));
                const rating = parseInt(await ask('Оценка (1-5): '));
                if (rating >= 1 && rating <= 5) {
                    if (catalog.rate(id, rating)) {
                        console.log('\x1b[32m✅ Оценка добавлена.\x1b[0m');
                    } else {
                        console.log('\x1b[31m❌ Рецепт не найден.\x1b[0m');
                    }
                } else {
                    console.log('\x1b[31m❌ Оценка должна быть от 1 до 5.\x1b[0m');
                }
                break;
            }
            case '9': catalog.stats(); break;
            case '10': {
                console.log('1. Экспорт в JSON');
                console.log('2. Экспорт в CSV');
                const sub = await ask('Выберите формат: ');
                if (sub === '1') catalog.exportJSON();
                else if (sub === '2') catalog.exportCSV();
                else console.log('\x1b[31mНеверный выбор.\x1b[0m');
                break;
            }
            case '11':
                console.log('До свидания!');
                rl.close();
                return;
            default:
                console.log('\x1b[31mНеверный выбор.\x1b[0m');
        }
    }
}

main().catch(console.error);
